package channel

import (
	"fmt"
	"time"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/mpacket"
)

const (
	opcodeSendMarriageRequest   int16 = 59
	opcodeSendMarriageResult    int16 = 60
	opcodeSendWeddingGiftResult int16 = 61
	opcodeSendPartnerTransfer   int16 = 62
	opcodeRecvRingAction        int16 = 0x7D
	opcodeRecvWeddingAction     int16 = 0x7E
	marriageResultSuccess       byte  = 11
	marriageResultEngagedNotice byte  = 36
	weddingStageLobby           byte  = 0
	weddingStageCeremony        byte  = 1
	weddingStageBlessingsClosed byte  = 2
	weddingStageMarried         byte  = 3
	weddingStageParty           byte  = 4
	weddingStageCompleted       byte  = 5
)

const (
	weddingLobbyDuration    = 10 * time.Minute
	weddingBlessingDuration = 15 * time.Minute
	weddingCeremonyDuration = 20 * time.Minute
	weddingPartyDuration    = 45 * time.Minute
	marriageDivorceCooldown = 7 * 24 * time.Hour
)

type pendingMarriageProposal struct {
	SourceID   int32
	TargetID   int32
	SourceName string
	ItemID     int32
}

var pendingMarriageProposals = make(map[int32]pendingMarriageProposal)

type weddingReservation struct {
	MarriageID    int32
	GroomID       int32
	BrideID       int32
	Cathedral     bool
	Premium       bool
	InviteItem    int32
	GuestTicket   int32
	EntryMapID    int32
	AltarMapID    int32
	ReservedAt    time.Time
	Started       bool
	Completed     bool
	Guests        map[int32]bool
	Stage         byte
	Blessings     int32
	GroomWishlist []string
	BrideWishlist []string
	GroomGifts    []Item
	BrideGifts    []Item
	GiftCounts    map[int32]int
}

var weddingReservations = make(map[int32]*weddingReservation)

func engagementOutcomeItem(itemID int32) int32 {
	switch itemID {
	case constant.ItemEngagementBoxMoonstone:
		return constant.ItemEngagementRingMoonstone
	case constant.ItemEngagementBoxStar:
		return constant.ItemEngagementRingStar
	case constant.ItemEngagementBoxGolden:
		return constant.ItemEngagementRingGolden
	case constant.ItemEngagementBoxSilver:
		return constant.ItemEngagementRingSilver
	default:
		return 0
	}
}

func weddingRingOutcome(itemID int32) int32 {
	switch itemID {
	case constant.ItemEngagementBoxMoonstone, constant.ItemEngagementRingMoonstone, constant.ItemEmptyEngagementBoxMoonstone:
		return constant.ItemWeddingRingMoonstone
	case constant.ItemEngagementBoxStar, constant.ItemEngagementRingStar, constant.ItemEmptyEngagementBoxStar:
		return constant.ItemWeddingRingStar
	case constant.ItemEngagementBoxGolden, constant.ItemEngagementRingGolden, constant.ItemEmptyEngagementBoxGolden:
		return constant.ItemWeddingRingGolden
	case constant.ItemEngagementBoxSilver, constant.ItemEngagementRingSilver, constant.ItemEmptyEngagementBoxSilver:
		return constant.ItemWeddingRingSilver
	default:
		return constant.ItemWeddingRingMoonstone
	}
}

func packetMarriageRequest(name string, playerID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodeSendMarriageRequest)
	p.WriteByte(0)
	p.WriteString(name)
	p.WriteInt32(playerID)
	return p
}

func packetMarriageResultRecord(marriageID int32, groomID int32, brideID int32, ringID int32, groomName, brideName string, wedding bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodeSendMarriageResult)
	p.WriteByte(marriageResultSuccess)
	p.WriteInt32(marriageID)
	p.WriteInt32(groomID)
	p.WriteInt32(brideID)
	if wedding {
		p.WriteInt16(3)
	} else {
		p.WriteInt16(1)
	}
	p.WriteInt32(ringID)
	p.WriteInt32(ringID)
	p.WritePaddedString(groomName, 13)
	p.WritePaddedString(brideName, 13)
	return p
}

func packetMarriageResultNotice(msg byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodeSendMarriageResult)
	p.WriteByte(msg)
	if msg == marriageResultEngagedNotice {
		p.WriteByte(1)
		p.WriteString("You are now engaged.")
	}
	return p
}

func packetNotifyWeddingPartnerTransfer(partnerID, mapID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodeSendPartnerTransfer)
	p.WriteInt32(mapID)
	p.WriteInt32(partnerID)
	return p
}

func lookupMarriageID(charID int32) int32 {
	var marriageID int32 = -1
	_ = common.DB.QueryRow("SELECT id FROM marriages WHERE husbandID=? OR wifeID=? LIMIT 1", charID, charID).Scan(&marriageID)
	return marriageID
}

func createMarriageRelationship(groomID, brideID int32) (int32, error) {
	res, err := common.DB.Exec("INSERT INTO marriages (husbandID, wifeID) VALUES (?, ?)", groomID, brideID)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int32(id), nil
}

func deleteMarriageRelationship(marriageID int32) error {
	_, err := common.DB.Exec("DELETE FROM marriages WHERE id=?", marriageID)
	return err
}

func (server *Server) proposeMarriage(source *Player, targetName string, itemID int32) error {
	if source.married() || source.partnerID > 0 {
		return fmt.Errorf("you are already engaged or married")
	}
	if source.underMarriageCooldown() {
		return fmt.Errorf("you must wait before marrying again")
	}
	if !source.removeItemsByID(itemID, 1, false) {
		return fmt.Errorf("you do not have the engagement box")
	}

	target, err := server.players.GetFromName(targetName)
	if err != nil {
		_ = source.GainItemByID(itemID, 1)
		return fmt.Errorf("unable to find %s on this channel", targetName)
	}
	if target.ID == source.ID {
		_ = source.GainItemByID(itemID, 1)
		return fmt.Errorf("you cannot propose to yourself")
	}
	if target.mapID != source.mapID || target.inst.id != source.inst.id {
		_ = source.GainItemByID(itemID, 1)
		return fmt.Errorf("your partner must be on the same map")
	}
	if target.partnerID > 0 || target.married() {
		_ = source.GainItemByID(itemID, 1)
		return fmt.Errorf("the other player is already engaged or married")
	}
	if target.underMarriageCooldown() {
		_ = source.GainItemByID(itemID, 1)
		return fmt.Errorf("the other player must wait before marrying again")
	}
	if target.gender == source.gender {
		_ = source.GainItemByID(itemID, 1)
		return fmt.Errorf("engagement currently requires opposite genders")
	}

	pendingMarriageProposals[target.ID] = pendingMarriageProposal{
		SourceID:   source.ID,
		TargetID:   target.ID,
		SourceName: source.Name,
		ItemID:     itemID,
	}
	target.Send(packetMarriageRequest(source.Name, source.ID))
	return nil
}

func (source *Player) GainItemByID(id int32, amount int16) error {
	it, err := CreateItemFromID(id, amount)
	if err != nil {
		return err
	}
	_, err = source.GiveItem(it)
	return err
}

func (server *Server) resolveMarriageProposal(target *Player, accepted bool, proposerName string, proposerID int32) {
	proposal, ok := pendingMarriageProposals[target.ID]
	if !ok || proposal.SourceID != proposerID || proposal.SourceName != proposerName {
		target.Send(packetPlayerNoChange())
		return
	}
	delete(pendingMarriageProposals, target.ID)

	source, err := server.players.GetFromID(proposal.SourceID)
	if err != nil {
		target.Send(packetPlayerNoChange())
		return
	}

	if !accepted {
		source.Send(packetMarriageResultNotice(0))
		source.Send(packetMessageRedText(target.Name + " has declined your engagement request."))
		_ = source.GainItemByID(proposal.ItemID, 1)
		return
	}

	engagementRing := engagementOutcomeItem(proposal.ItemID)
	if engagementRing == 0 {
		source.Send(packetMarriageResultNotice(0))
		return
	}

	if err := source.GainItemByID(engagementRing, 1); err != nil {
		source.Send(packetMarriageResultNotice(0))
		return
	}
	if err := target.GainItemByID(engagementRing+1, 1); err != nil {
		_ = source.removeItemsByID(engagementRing, 1, false)
		source.Send(packetMarriageResultNotice(0))
		return
	}

	marriageID, err := createMarriageRelationship(source.ID, target.ID)
	if err != nil {
		_ = source.removeItemsByID(engagementRing, 1, false)
		_ = target.removeItemsByID(engagementRing+1, 1, false)
		source.Send(packetMarriageResultNotice(0))
		return
	}

	_ = source.setPartnerID(target.ID)
	_ = target.setPartnerID(source.ID)
	source.marriageID = marriageID
	target.marriageID = marriageID
	_ = source.setMarriageItemID(engagementRing)
	_ = target.setMarriageItemID(engagementRing + 1)

	packet := packetMarriageResultRecord(marriageID, source.ID, target.ID, weddingRingOutcome(engagementRing), source.Name, target.Name, false)
	source.Send(packet)
	target.Send(packet)
	source.Send(packetNotifyWeddingPartnerTransfer(target.ID, target.mapID))
	target.Send(packetNotifyWeddingPartnerTransfer(source.ID, source.mapID))
}

func (server *Server) breakMarriageState(plr *Player, itemID int32) {
	partnerID := plr.partnerID
	if partnerID <= 0 {
		return
	}
	marriageID := lookupMarriageID(plr.ID)
	if marriageID > 0 {
		_ = deleteMarriageRelationship(marriageID)
		delete(weddingReservations, marriageID)
	}
	if ringID := plr.firstRingIDByKind(ringKindWedding); ringID > 0 {
		_ = deleteRingRecordPair(ringID)
	}
	partner, _ := server.players.GetFromID(partnerID)
	_ = plr.setPartnerID(-1)
	_ = plr.setMarriageItemID(-1)
	plr.marriageID = -1
	plr.refreshRingRecords()
	_ = plr.setDivorceUntil(time.Now().Add(marriageDivorceCooldown).Unix())
	_ = plr.removeItemsByID(itemID, 1, false)
	plr.Send(packetNotifyWeddingPartnerTransfer(0, 0))
	if partner != nil {
		_ = partner.setPartnerID(-1)
		_ = partner.setMarriageItemID(-1)
		partner.marriageID = -1
		partner.refreshRingRecords()
		_ = partner.setDivorceUntil(time.Now().Add(marriageDivorceCooldown).Unix())
		_ = partner.removeItemsByID(itemID+1, 1, false)
		partner.Send(packetNotifyWeddingPartnerTransfer(0, 0))
	}
}

func hasAnyWeddingRing(plr *Player) bool {
	for _, id := range []int32{constant.ItemWeddingRingMoonstone, constant.ItemWeddingRingStar, constant.ItemWeddingRingGolden, constant.ItemWeddingRingSilver} {
		if plr.countItem(id) > 0 {
			return true
		}
	}
	return false
}

func weddingGuestTicket(cathedral bool) int32 {
	if cathedral {
		return constant.ItemWeddingGuestTicketCathedral
	}
	return constant.ItemWeddingGuestTicketChapel
}

func weddingInviteItem(cathedral bool) int32 {
	if cathedral {
		return constant.ItemWeddingInvitationCathedral
	}
	return constant.ItemWeddingInvitationChapel
}

func weddingEntryMap(cathedral bool) int32 {
	if cathedral {
		return 680000200
	}
	return 680000100
}

func weddingReservationTicket(cathedral, premium bool) int32 {
	if cathedral {
		if premium {
			return constant.ItemWeddingTicketCathedralPremium
		}
		return constant.ItemWeddingTicketCathedral
	}
	if premium {
		return constant.ItemWeddingTicketChapelPremium
	}
	return constant.ItemWeddingTicketChapelRegular
}

func weddingAltarMap(cathedral bool) int32 {
	if cathedral {
		return 680000210
	}
	return 680000110
}

func (server *Server) reserveWedding(plr *Player, cathedral, premium bool) error {
	if plr.partnerID <= 0 || plr.married() {
		return fmt.Errorf("you must be engaged before reserving a wedding")
	}
	marriageID := lookupMarriageID(plr.ID)
	if marriageID <= 0 {
		return fmt.Errorf("your relationship record could not be found")
	}
	if _, exists := weddingReservations[marriageID]; exists {
		return fmt.Errorf("you already have a wedding reservation")
	}
	partner, err := server.players.GetFromID(plr.partnerID)
	if err != nil {
		return fmt.Errorf("your partner must be online on this channel")
	}
	if partner.mapID != plr.mapID || partner.inst.id != plr.inst.id {
		return fmt.Errorf("your partner must be on the same map to reserve the wedding")
	}
	if hasAnyWeddingRing(plr) || hasAnyWeddingRing(partner) {
		return fmt.Errorf("one of you already has a wedding ring")
	}

	reservationTicket := weddingReservationTicket(cathedral, premium)
	inviteItem := weddingInviteItem(cathedral)
	guestTicket := weddingGuestTicket(cathedral)
	if plr.countItem(reservationTicket) < 1 {
		return fmt.Errorf("you need the correct wedding reservation ticket")
	}
	if !plr.removeItemsByID(reservationTicket, 1, false) {
		return fmt.Errorf("could not consume the wedding reservation ticket")
	}
	if err := plr.GainItemByID(inviteItem, 15); err != nil {
		_ = plr.GainItemByID(reservationTicket, 1)
		return fmt.Errorf("please make room in your ETC inventory first")
	}
	if err := partner.GainItemByID(inviteItem, 15); err != nil {
		_ = plr.removeItemsByID(inviteItem, 15, false)
		_ = plr.GainItemByID(reservationTicket, 1)
		return fmt.Errorf("your partner needs room in their ETC inventory first")
	}

	weddingReservations[marriageID] = &weddingReservation{
		MarriageID:  marriageID,
		GroomID:     plr.ID,
		BrideID:     partner.ID,
		Cathedral:   cathedral,
		Premium:     premium,
		InviteItem:  inviteItem,
		GuestTicket: guestTicket,
		EntryMapID:  weddingEntryMap(cathedral),
		AltarMapID:  weddingAltarMap(cathedral),
		ReservedAt:  time.Now(),
		Guests:      make(map[int32]bool),
		Stage:       weddingStageLobby,
		GiftCounts:  make(map[int32]int),
	}
	return nil
}

func (server *Server) startWedding(plr *Player, cathedral bool) error {
	marriageID := lookupMarriageID(plr.ID)
	res, ok := weddingReservations[marriageID]
	if !ok {
		return fmt.Errorf("you do not have a wedding reservation")
	}
	if res.Cathedral != cathedral {
		return fmt.Errorf("your reservation is for a different venue")
	}
	if res.Started {
		return fmt.Errorf("your wedding has already started")
	}
	partner, err := server.players.GetFromID(plr.partnerID)
	if err != nil {
		return fmt.Errorf("your partner must be online on this channel")
	}
	if partner.mapID != plr.mapID || partner.inst.id != plr.inst.id {
		return fmt.Errorf("your partner must be here with you to begin the ceremony")
	}
	if plr.marriageItemID <= 0 || partner.marriageItemID <= 0 {
		return fmt.Errorf("both partners must still hold their engagement rings")
	}
	res.Started = true
	res.Stage = weddingStageLobby
	res.Blessings = 0
	server.warpWeddingCouple(res, res.EntryMapID)
	server.showWeddingCountdown(res, int32(weddingLobbyDuration/time.Second))

	server.scheduleWeddingStage(res, weddingLobbyDuration, func(cur *weddingReservation) {
		if cur.Stage != weddingStageLobby {
			return
		}
		server.advanceWeddingToCeremony(cur)
	})
	server.scheduleWeddingStage(res, weddingLobbyDuration+weddingBlessingDuration, func(cur *weddingReservation) {
		if cur.Stage != weddingStageCeremony {
			return
		}
		server.closeWeddingBlessings(cur)
	})
	server.scheduleWeddingStage(res, weddingLobbyDuration+weddingCeremonyDuration, func(cur *weddingReservation) {
		if cur.Stage < weddingStageMarried {
			server.endWedding(cur)
		}
	})
	return nil
}

func (server *Server) advanceWeddingToCeremony(res *weddingReservation) {
	if res == nil || res.Completed || !res.Started || res.Stage != weddingStageLobby {
		return
	}
	res.Stage = weddingStageCeremony
	server.warpWeddingParticipants(res, res.AltarMapID)
	server.showWeddingCountdown(res, int32(weddingBlessingDuration/time.Second))
	if groom, err := server.players.GetFromID(res.GroomID); err == nil && groom.inst != nil {
		groom.inst.send(packetMessageNotice("Wedding Assistant: The couple are heading to the altar."))
	}
}

func (server *Server) closeWeddingBlessings(res *weddingReservation) {
	if res == nil || res.Completed || res.Stage != weddingStageCeremony {
		return
	}
	res.Stage = weddingStageBlessingsClosed
	server.showWeddingCountdown(res, int32((weddingCeremonyDuration-weddingBlessingDuration)/time.Second))
	if groom, err := server.players.GetFromID(res.GroomID); err == nil {
		groom.Send(packetMessageNotice("Wedding Assistant: The blessing period is now closed."))
	}
	if bride, err := server.players.GetFromID(res.BrideID); err == nil {
		bride.Send(packetMessageNotice("Wedding Assistant: The blessing period is now closed."))
	}
}

func (server *Server) advanceWeddingCeremony(plr *Player, cathedral bool) error {
	res := server.currentWeddingReservation(plr, cathedral)
	if res == nil {
		return fmt.Errorf("your wedding reservation could not be found")
	}
	if !res.Started || res.Completed {
		return fmt.Errorf("your wedding session is not active")
	}
	if res.Stage != weddingStageLobby {
		return fmt.Errorf("the ceremony has already moved beyond the lounge")
	}
	if !res.isParticipant(plr.ID) {
		return fmt.Errorf("only the engaged couple can start the ceremony")
	}
	partner, err := server.players.GetFromID(plr.partnerID)
	if err != nil {
		return fmt.Errorf("your partner must be online on this channel")
	}
	if partner.mapID != plr.mapID || partner.inst == nil || plr.inst == nil || partner.inst.id != plr.inst.id {
		return fmt.Errorf("your partner must be with you in the lounge")
	}
	server.advanceWeddingToCeremony(res)
	return nil
}

func (server *Server) completeWedding(plr *Player, cathedral bool) error {
	marriageID := lookupMarriageID(plr.ID)
	res, ok := weddingReservations[marriageID]
	if !ok {
		return fmt.Errorf("you do not have an active wedding session")
	}
	if res.Cathedral != cathedral {
		return fmt.Errorf("your wedding is scheduled for a different venue")
	}
	if !res.Started || res.Completed {
		return fmt.Errorf("your wedding session is not active")
	}
	if res.Stage != weddingStageCeremony && res.Stage != weddingStageBlessingsClosed {
		return fmt.Errorf("the vows cannot be completed yet")
	}
	partner, err := server.players.GetFromID(plr.partnerID)
	if err != nil {
		return fmt.Errorf("your partner must be online on this channel")
	}
	if partner.mapID != plr.mapID || partner.inst.id != plr.inst.id {
		return fmt.Errorf("your partner must be here with you at the altar")
	}
	if plr.marriageItemID <= 0 || partner.marriageItemID <= 0 {
		return fmt.Errorf("both partners must still hold their engagement rings")
	}

	weddingRing := weddingRingOutcome(plr.marriageItemID)
	if weddingRing <= 0 {
		return fmt.Errorf("could not determine the wedding ring for this couple")
	}
	ownerRingID, partnerRingID, err := createPairedRingRecords(weddingRing, plr, partner)
	if err != nil {
		return fmt.Errorf("could not create wedding ring records")
	}
	if !plr.removeItemsByID(plr.marriageItemID, 1, false) || !partner.removeItemsByID(partner.marriageItemID, 1, false) {
		_ = deleteRingRecordPair(ownerRingID)
		return fmt.Errorf("both partners must carry their engagement ring items")
	}
	ownerRingItem, err := CreateItemFromID(weddingRing, 1)
	if err != nil {
		_ = deleteRingRecordPair(ownerRingID)
		return fmt.Errorf("could not create the wedding ring item")
	}
	ownerRingItem.ringID = ownerRingID
	partnerRingItem, err := CreateItemFromID(weddingRing, 1)
	if err != nil {
		_ = deleteRingRecordPair(ownerRingID)
		return fmt.Errorf("could not create the wedding ring item")
	}
	partnerRingItem.ringID = partnerRingID
	ownerGiven, err := plr.GiveItem(ownerRingItem)
	if err != nil {
		_ = deleteRingRecordPair(ownerRingID)
		return fmt.Errorf("please make room in your EQUIP inventory first")
	}
	if _, err := partner.GiveItem(partnerRingItem); err != nil {
		plr.removeItem(ownerGiven, false)
		_ = deleteRingRecordPair(ownerRingID)
		return fmt.Errorf("your partner needs room in their EQUIP inventory first")
	}
	_ = plr.setMarriageItemID(weddingRing)
	_ = partner.setMarriageItemID(weddingRing)
	plr.refreshRingRecords()
	partner.refreshRingRecords()

	packet := packetMarriageResultRecord(marriageID, plr.ID, partner.ID, weddingRing, plr.Name, partner.Name, true)
	plr.Send(packet)
	partner.Send(packet)
	res.Stage = weddingStageMarried
	server.hideWeddingCountdown(res)
	return nil
}

func (server *Server) warpWeddingParticipants(res *weddingReservation, mapID int32) {
	field, ok := server.fields[mapID]
	if !ok {
		return
	}
	inst, err := field.getInstance(0)
	if err != nil {
		return
	}
	portal, err := inst.getRandomSpawnPortal()
	if err != nil {
		return
	}
	for _, id := range append([]int32{res.GroomID, res.BrideID}, keysOfGuests(res.Guests)...) {
		if p, err := server.players.GetFromID(id); err == nil {
			_ = server.warpPlayer(p, field, portal, true)
		}
	}
}

func (server *Server) startWeddingAfterParty(plr *Player) error {
	res := server.currentWeddingReservationAny(plr)
	if res == nil {
		return fmt.Errorf("there is no active wedding session")
	}
	if !res.isParticipant(plr.ID) {
		return fmt.Errorf("only the married couple can begin the afterparty")
	}
	if res.Stage != weddingStageMarried {
		return fmt.Errorf("the vows must be completed before the afterparty can begin")
	}
	if res.Premium {
		res.Stage = weddingStageParty
		server.warpWeddingParticipants(res, 680000300)
		server.showWeddingCountdown(res, int32(weddingPartyDuration/time.Second))
		server.scheduleWeddingStage(res, weddingPartyDuration, func(cur *weddingReservation) {
			server.endWedding(cur)
		})
	} else {
		server.endWedding(res)
	}
	return nil
}

func (server *Server) enterWeddingAsGuest(plr *Player, cathedral bool) error {
	for _, res := range weddingReservations {
		if res.Cathedral != cathedral || !res.Started || res.Completed {
			continue
		}
		if !plr.removeItemsByID(res.GuestTicket, 1, false) {
			return fmt.Errorf("you do not have the correct guest ticket")
		}
		dstMap := res.EntryMapID
		if res.Stage >= weddingStageCeremony {
			dstMap = res.AltarMapID
		}
		altarField, ok := server.fields[dstMap]
		if !ok {
			return fmt.Errorf("the wedding map is unavailable right now")
		}
		inst, err := altarField.getInstance(0)
		if err != nil {
			return fmt.Errorf("the wedding instance is unavailable right now")
		}
		portal, err := inst.getRandomSpawnPortal()
		if err != nil {
			return fmt.Errorf("the wedding entry point is unavailable right now")
		}
		_ = server.warpPlayer(plr, altarField, portal, true)
		return nil
	}
	return fmt.Errorf("there is no active wedding at that venue right now")
}

func (server *Server) inviteGuest(plr *Player, guestName string, marriageID int32, slot int16) error {
	res, ok := weddingReservations[marriageID]
	if !ok {
		return fmt.Errorf("your wedding reservation could not be found")
	}
	if res.Started {
		return fmt.Errorf("the wedding is already under way")
	}
	guest, err := server.players.GetFromName(guestName)
	if err != nil {
		return fmt.Errorf("unable to find %s on this channel", guestName)
	}
	if res.Guests[guest.ID] {
		return fmt.Errorf("that guest has already been invited")
	}
	item := plr.findEtcItemBySlot(slot)
	if item == nil || item.ID != res.InviteItem || item.amount < 1 {
		return fmt.Errorf("you need a valid wedding invitation card")
	}
	if err := guest.GainItemByID(res.GuestTicket, 1); err != nil {
		return fmt.Errorf("your guest needs room in their ETC inventory")
	}
	if _, err := plr.takeItem(res.InviteItem, slot, 1, byte(constant.InventoryEtc)); err != nil {
		_ = guest.removeItemsByID(res.GuestTicket, 1, false)
		return fmt.Errorf("could not consume the wedding invitation")
	}
	res.Guests[guest.ID] = true
	guest.Send(packetMarriageInvitation(CharacterNameByID(res.GroomID, ""), CharacterNameByID(res.BrideID, ""), res.Cathedral))
	return nil
}

func CharacterNameByID(id int32, fallback string) string {
	var name string
	if err := common.DB.QueryRow("SELECT name FROM characters WHERE id=?", id).Scan(&name); err != nil || name == "" {
		return fallback
	}
	return name
}

func packetMarriageInvitation(groom, bride string, cathedral bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodeSendMarriageResult)
	p.WriteByte(15)
	p.WriteString(groom)
	p.WriteString(bride)
	if cathedral {
		p.WriteInt16(1)
	} else {
		p.WriteInt16(2)
	}
	return p
}

func (server *Server) showWeddingCountdown(res *weddingReservation, seconds int32) {
	for _, mapID := range []int32{res.EntryMapID, res.AltarMapID, 680000300, 680000400, 680000401} {
		field, ok := server.fields[mapID]
		if !ok {
			continue
		}
		inst, err := field.getInstance(0)
		if err != nil {
			continue
		}
		inst.send(packetShowCountdown(seconds))
	}
}

func (server *Server) hideWeddingCountdown(res *weddingReservation) {
	for _, mapID := range []int32{res.EntryMapID, res.AltarMapID, 680000300, 680000400, 680000401} {
		field, ok := server.fields[mapID]
		if !ok {
			continue
		}
		inst, err := field.getInstance(0)
		if err != nil {
			continue
		}
		inst.send(packetHideCountdown())
	}
}

func (server *Server) activeWeddingByMap(mapID int32) *weddingReservation {
	for _, res := range weddingReservations {
		if !res.Started || res.Completed {
			continue
		}
		if mapID == res.EntryMapID || mapID == res.AltarMapID {
			return res
		}
	}
	return nil
}

func isWeddingMap(mapID int32) bool {
	return mapID >= 680000100 && mapID <= 680000500
}

func (res *weddingReservation) isParticipant(charID int32) bool {
	return res.GroomID == charID || res.BrideID == charID
}

func (server *Server) handleWeddingMapLeave(plr *Player, dstMapID int32) {
	for _, res := range weddingReservations {
		if !res.Started || res.Completed {
			continue
		}
		if res.isParticipant(plr.ID) && !isWeddingMap(dstMapID) {
			server.endWedding(res)
			return
		}
		if res.Guests[plr.ID] && !isWeddingMap(dstMapID) {
			delete(res.Guests, plr.ID)
		}
	}
}

func (server *Server) handleWeddingDisconnect(plr *Player) {
	for _, res := range weddingReservations {
		if !res.Started || res.Completed {
			continue
		}
		if res.isParticipant(plr.ID) {
			server.endWedding(res)
			return
		}
		if res.Guests[plr.ID] {
			delete(res.Guests, plr.ID)
		}
	}
}

func (server *Server) scheduleWeddingStage(res *weddingReservation, after time.Duration, fn func(*weddingReservation)) {
	time.AfterFunc(after, func() {
		server.dispatch <- func() {
			current, ok := weddingReservations[res.MarriageID]
			if !ok || current.Completed {
				return
			}
			fn(current)
		}
	})
}

func (server *Server) warpWeddingCouple(res *weddingReservation, mapID int32) {
	field, ok := server.fields[mapID]
	if !ok {
		return
	}
	inst, err := field.getInstance(0)
	if err != nil {
		return
	}
	portal, err := inst.getRandomSpawnPortal()
	if err != nil {
		return
	}
	if groom, err := server.players.GetFromID(res.GroomID); err == nil {
		_ = server.warpPlayer(groom, field, portal, true)
	}
	if bride, err := server.players.GetFromID(res.BrideID); err == nil {
		_ = server.warpPlayer(bride, field, portal, true)
	}
}

func (server *Server) endWedding(res *weddingReservation) {
	server.hideWeddingCountdown(res)
	res.Completed = true
	res.Stage = weddingStageCompleted
	delete(weddingReservations, res.MarriageID)
	exitMap := int32(680000500)
	field, ok := server.fields[exitMap]
	if !ok {
		return
	}
	inst, err := field.getInstance(0)
	if err != nil {
		return
	}
	portal, err := inst.getRandomSpawnPortal()
	if err != nil {
		return
	}
	for _, plr := range append([]int32{res.GroomID, res.BrideID}, keysOfGuests(res.Guests)...) {
		if p, err := server.players.GetFromID(plr); err == nil {
			_ = server.warpPlayer(p, field, portal, true)
		}
	}
}

func keysOfGuests(m map[int32]bool) []int32 {
	r := make([]int32, 0, len(m))
	for id := range m {
		r = append(r, id)
	}
	return r
}

func (server *Server) currentWeddingReservation(plr *Player, cathedral bool) *weddingReservation {
	marriageID := lookupMarriageID(plr.ID)
	res, ok := weddingReservations[marriageID]
	if !ok || res.Cathedral != cathedral {
		return nil
	}
	return res
}

func (server *Server) currentWeddingReservationAny(plr *Player) *weddingReservation {
	marriageID := lookupMarriageID(plr.ID)
	res, ok := weddingReservations[marriageID]
	if !ok {
		return nil
	}
	return res
}

func (res *weddingReservation) isGroom(id int32) bool {
	return res.GroomID == id
}

func (res *weddingReservation) wishlistFor(id int32) *[]string {
	if res.isGroom(id) {
		return &res.GroomWishlist
	}
	return &res.BrideWishlist
}

func (res *weddingReservation) spouseWishlistFor(id int32) *[]string {
	if res.isGroom(id) {
		return &res.BrideWishlist
	}
	return &res.GroomWishlist
}

func (res *weddingReservation) giftsFor(id int32) *[]Item {
	if res.isGroom(id) {
		return &res.GroomGifts
	}
	return &res.BrideGifts
}

func packetWeddingGiftResult(mode byte, names []string, items []Item) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodeSendWeddingGiftResult)
	p.WriteByte(mode)
	switch mode {
	case 0x09:
		p.WriteByte(byte(len(names)))
		for _, name := range names {
			p.WriteString(name)
		}
	case 0x0A, 0x0B, 0x0F:
		if mode == 0x0B {
			p.WriteByte(byte(len(names)))
			for _, name := range names {
				p.WriteString(name)
			}
		}
		p.WriteInt64(32)
		p.WriteByte(byte(len(items)))
		for _, item := range items {
			p.WriteBytes(item.StorageBytes())
		}
	}
	return p
}
