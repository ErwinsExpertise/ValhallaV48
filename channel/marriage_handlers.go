package channel

import (
	"log"

	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
)

func (server *Server) playerRingAction(conn mnet.Client, reader mpacket.Reader) {
	plr, err := server.players.GetFromConn(conn)
	if err != nil {
		return
	}

	mode := reader.ReadByte()

	switch mode {
	case 0:
		targetName := reader.ReadString(reader.ReadInt16())
		itemID := reader.ReadInt32()
		if err := server.proposeMarriage(plr, targetName, itemID); err != nil {
			plr.Send(packetMessageRedText(err.Error()))
		}
	case 1:
		if plr.marriageItemID/1000000 != 4 {
			_ = plr.setMarriageItemID(-1)
		}
	case 2:
		accepted := reader.ReadByte() > 0
		name := reader.ReadString(reader.ReadInt16())
		id := reader.ReadInt32()
		server.resolveMarriageProposal(plr, accepted, name, id)
	case 3:
		itemID := reader.ReadInt32()
		server.breakMarriageState(plr, itemID)
	case 5:
		guestName := reader.ReadString(reader.ReadInt16())
		marriageID := reader.ReadInt32()
		slot := int16(reader.ReadByte())
		if err := server.inviteGuest(plr, guestName, marriageID, slot); err != nil {
			plr.Send(packetMessageRedText(err.Error()))
		}
	case 6:
		slot := int16(reader.ReadInt32())
		invitationID := reader.ReadInt32()
		item := plr.findEtcItemBySlot(slot)
		if item == nil || item.ID != invitationID {
			plr.Send(packetPlayerNoChange())
			return
		}
		for _, res := range weddingReservations {
			if (res.GuestTicket == invitationID) && res.Guests[plr.ID] {
				plr.Send(packetMarriageInvitation(CharacterNameByID(res.GroomID, ""), CharacterNameByID(res.BrideID, ""), res.Cathedral))
				break
			}
		}
	case 9:
		count := int(reader.ReadInt16())
		res := server.currentWeddingReservationAny(plr)
		if res == nil {
			break
		}
		list := res.wishlistFor(plr.ID)
		*list = (*list)[:0]
		for i := 0; i < count; i++ {
			*list = append(*list, reader.ReadString(reader.ReadInt16()))
		}
		touchWeddingReservation(res)
	default:
		log.Printf("Unhandled ring action mode=%d for %s", mode, plr.Name)
	}

	plr.Send(packetPlayerNoChange())
}

func (server *Server) playerWeddingAction(conn mnet.Client, reader mpacket.Reader) {
	plr, err := server.players.GetFromConn(conn)
	if err != nil {
		return
	}

	mode := reader.ReadByte()
	res := server.currentWeddingReservationAny(plr)
	if res == nil {
		plr.Send(packetPlayerNoChange())
		return
	}

	switch mode {
	case 6:
		slot := reader.ReadInt16()
		itemID := reader.ReadInt32()
		qty := reader.ReadInt16()
		spouseID := res.BrideID
		if !res.isGroom(plr.ID) {
			spouseID = res.GroomID
		}
		if spouseID == plr.ID {
			plr.Send(packetWeddingGiftResult(0x0E, nil, nil))
			break
		}
		if res.GiftCounts[plr.ID] >= 5 {
			plr.Send(packetWeddingGiftResult(0x0C, nil, nil))
			break
		}
		item := plr.findEtcItemBySlot(slot)
		if item == nil || item.ID != itemID || qty <= 0 || item.amount < qty {
			plr.Send(packetWeddingGiftResult(0x0E, nil, nil))
			break
		}
		gift := *item
		gift.amount = qty
		if _, err := plr.takeItem(itemID, slot, qty, byte(constant.InventoryEtc)); err != nil {
			plr.Send(packetWeddingGiftResult(0x0E, nil, nil))
			break
		}
		gifts := res.giftsFor(spouseID)
		*gifts = append(*gifts, weddingGiftEntry{ID: gift.ID, Amount: gift.amount})
		res.GiftCounts[plr.ID]++
		touchWeddingReservation(res)
		plr.Send(packetWeddingGiftResult(0x0B, *res.spouseWishlistFor(plr.ID), []Item{gift}))
	case 7:
		_ = reader.ReadByte()
		idx := int(reader.ReadByte())
		gifts := res.giftsFor(plr.ID)
		if idx < 0 || idx >= len(*gifts) {
			plr.Send(packetWeddingGiftResult(0x0E, *res.wishlistFor(plr.ID), weddingGiftItems(*gifts)))
			break
		}
		gift := (*gifts)[idx]
		if err := plr.GainItemByID(gift.ID, gift.Amount); err != nil {
			plr.Send(packetWeddingGiftResult(0x0E, *res.wishlistFor(plr.ID), weddingGiftItems(*gifts)))
			break
		}
		*gifts = append((*gifts)[:idx], (*gifts)[idx+1:]...)
		touchWeddingReservation(res)
		plr.Send(packetWeddingGiftResult(0x0F, *res.wishlistFor(plr.ID), weddingGiftItems(*gifts)))
	case 8:
		plr.Send(packetPlayerNoChange())
	default:
		log.Printf("Unhandled wedding action mode=%d for %s", mode, plr.Name)
	}
	plr.Send(packetPlayerNoChange())
}

func (server *Server) playerWeddingTalk(conn mnet.Client, reader mpacket.Reader) {
	plr, err := server.players.GetFromConn(conn)
	if err != nil {
		return
	}

	if reader.GetRestAsBytes() == nil || len(reader.GetRestAsBytes()) == 0 {
		plr.Send(packetPlayerNoChange())
		return
	}

	action := reader.ReadByte()
	res := server.activeWeddingByMap(plr.mapID)
	if res == nil {
		plr.Send(packetPlayerNoChange())
		return
	}

	isCouple := plr.ID == res.GroomID || plr.ID == res.BrideID
	switch action {
	case 1:
		if !isCouple && res.Stage == weddingStageCeremony {
			plr.Send(packetMessageNotice("High Priest John: Your blessings have been added to their love."))
		} else {
			plr.Send(packetMessageNotice("The ceremony is in progress."))
		}
	default:
		if !isCouple && res.Stage == weddingStageCeremony {
			res.Blessings++
			touchWeddingReservation(res)
			plr.Send(packetMessageNotice("High Priest John: Your blessings have been added to their love."))
		}
	}
	plr.Send(packetPlayerNoChange())
}
