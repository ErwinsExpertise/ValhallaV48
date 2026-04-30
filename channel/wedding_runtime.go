package channel

import (
	"time"

	"github.com/Hucaru/Valhalla/constant"
)

func (res *weddingReservation) ensureState() {
	if res.Guests == nil {
		res.Guests = make(map[int32]bool)
	}
	if res.Rewarded == nil {
		res.Rewarded = make(map[int32]bool)
	}
	if res.GiftCounts == nil {
		res.GiftCounts = make(map[int32]int)
	}
	if res.EntryMapID == 0 {
		res.EntryMapID = weddingEntryMap(res.Cathedral)
	}
	if res.AltarMapID == 0 {
		res.AltarMapID = weddingAltarMap(res.Cathedral)
	}
	if res.InviteItem == 0 {
		res.InviteItem = weddingInviteItem(res.Cathedral)
	}
	if res.GuestTicket == 0 {
		res.GuestTicket = weddingGuestTicket(res.Cathedral)
	}
	if res.ReceiptItem == 0 {
		res.ReceiptItem = weddingReceiptItem(res.Cathedral, res.Premium)
	}
	if res.ResumeMapID == 0 {
		res.ResumeMapID = res.EntryMapID
	}
	if res.ReservedAt.IsZero() {
		res.ReservedAt = time.Now()
	}
	if res.StageChangedAt.IsZero() {
		res.StageChangedAt = res.ReservedAt
	}
}

func (res *weddingReservation) participantIDs() []int32 {
	return []int32{res.GroomID, res.BrideID}
}

func weddingGiftItems(entries []weddingGiftEntry) []Item {
	items := make([]Item, 0, len(entries))
	for _, entry := range entries {
		item, err := CreateItemFromID(entry.ID, entry.Amount)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	return items
}

func touchWeddingReservation(res *weddingReservation) {
	if res == nil {
		return
	}
	res.ensureState()
}

func sameWeddingParty(plr, partner *Player) bool {
	if plr == nil || partner == nil || plr.party == nil || partner.party == nil || plr.party != partner.party {
		return false
	}
	count := 0
	hasSelf := false
	hasPartner := false
	for _, member := range plr.party.players {
		if member == nil {
			continue
		}
		count++
		if member.ID == plr.ID {
			hasSelf = true
		}
		if member.ID == partner.ID {
			hasPartner = true
		}
	}
	return count == 2 && hasSelf && hasPartner
}

func (server *Server) resumeWeddingReservation(res *weddingReservation) {
	if res == nil || res.Completed || res.DeadlineAt.IsZero() {
		return
	}
	server.scheduleWeddingStage(res, time.Until(res.DeadlineAt), func(cur *weddingReservation) {
		server.handleWeddingDeadline(cur)
	})
}

func (server *Server) handleWeddingDeadline(res *weddingReservation) {
	if res == nil || res.Completed {
		return
	}
	switch res.Stage {
	case weddingStageLobby:
		server.advanceWeddingToCeremony(res)
	case weddingStageCeremony:
		server.closeWeddingBlessings(res)
	case weddingStageBlessingsClosed:
		server.endWedding(res)
	case weddingStageMarried:
		server.startWeddingAfterPartyFromReservation(res)
	case weddingStageParty:
		if res.ResumeMapID == 680000300 {
			if res.Premium {
				res.StageChangedAt = time.Now()
				res.DeadlineAt = res.StageChangedAt.Add(weddingPartyDuration)
				res.ResumeMapID = 680000400
				touchWeddingReservation(res)
				server.warpWeddingParticipants(res, 680000400)
				server.showWeddingCountdown(res, int32(weddingPartyDuration/time.Second))
				server.resumeWeddingReservation(res)
				return
			}
			server.endWedding(res)
			return
		}
		server.endWedding(res)
	case weddingStageCompleted:
		server.cleanupWeddingReservation(res)
	}
}

func (server *Server) startWeddingAfterPartyFromReservation(res *weddingReservation) {
	if res == nil || res.Completed || !res.Finalized || res.Stage != weddingStageMarried {
		return
	}
	res.Stage = weddingStageParty
	res.StageChangedAt = time.Now()
	res.DeadlineAt = res.StageChangedAt.Add(weddingPhotoDuration())
	res.ResumeMapID = 680000300
	for _, id := range append(res.participantIDs(), keysOfGuests(res.Guests)...) {
		if plr, err := server.players.GetFromID(id); err == nil && plr.countItem(constant.ItemWeddingEntryPermission) == 0 {
			_ = plr.GainItemByID(constant.ItemWeddingEntryPermission, 1)
		}
	}
	touchWeddingReservation(res)
	server.warpWeddingParticipants(res, 680000300)
	server.showWeddingCountdown(res, int32(weddingPhotoDuration()/time.Second))
	server.resumeWeddingReservation(res)
}

func (server *Server) cleanupWeddingReservation(res *weddingReservation) {
	if res == nil {
		return
	}
	res.Completed = true
	delete(weddingReservations, res.MarriageID)
}

func (server *Server) claimWeddingExit(plr *Player) int {
	res := server.currentWeddingReservationAny(plr)
	if res == nil || !res.ExitReady || res.Completed {
		return 0
	}
	if !res.isParticipant(plr.ID) {
		return 1
	}
	if res.Rewarded[plr.ID] {
		return 1
	}
	item, err := CreateItemFromID(4031424, 1)
	if err != nil || !plr.CanReceiveItems([]Item{item}) {
		return 3
	}
	if err := plr.GainItemByID(4031424, 1); err != nil {
		return 3
	}
	res.Rewarded[plr.ID] = true
	touchWeddingReservation(res)
	return 2
}
