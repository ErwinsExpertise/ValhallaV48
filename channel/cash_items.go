package channel

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/internal"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

const (
	cashItemTypeSpeakerChannel = 12
	cashItemTypeSpeakerWorld   = 13
	cashItemTypeWeather        = 14
	cashItemTypePetNameTag     = 15
	cashItemTypeMessageBox     = 16
	cashItemTypeSendMemo       = 19
	cashItemTypeMapTransfer    = 20
	cashItemTypeAPReset        = 21
	cashItemTypeSPReset        = 22
	cashItemTypeItemNameTag    = 23
	cashItemTypePetSkill       = 26
	cashItemTypeShopScanner    = 27
	cashItemTypeChalkboard     = 28
	cashItemTypePetFood        = 30
	cashItemTypeMorph          = 31
	cashItemTypeParcel         = 34
	cashItemTypeMapleTV        = 40
	cashItemTypeMapleSoleTV    = 41
	cashItemTypeMapleLoveTV    = 42
	cashItemTypeMegaTV         = 43
	cashItemTypeMegaSoleTV     = 44
	cashItemTypeMegaLoveTV     = 45
	cashItemTypeNameChange     = 46
	cashItemTypeTransferWorld  = 50
	cashItemTypeUnsupported    = 0
	cashItemTypeExpCoupon      = 100
	cashItemTypeDropCoupon     = 101

	mapTransferResultDeleteSlot     = 2
	mapTransferResultAddSlot        = 3
	mapTransferResultInvalidField   = 5
	mapTransferResultUnavailable    = 6
	mapTransferResultTargetDead     = 7
	mapTransferResultBlocked        = 8
	mapTransferResultSameMap        = 9
	mapTransferResultLowLevel       = 11
	petNameMinLength                = 4
	petNameMaxLength                = 12
	channelMegaphoneMaxDisplayChars = 80
	worldMegaphoneMaxChars          = 60
)

func cashItemUseType(itemID int32) int {
	itemType := cashSlotItemType(itemID)
	if itemType == 0 {
		return 0
	}

	if itemType > 38 {
		switch itemType {
		case 40, 41, 42, 43, 44, 45, 46, 50:
			return itemType
		default:
			return 0
		}
	}

	if itemType < 12 {
		return 0
	}

	if itemType <= 28 {
		return itemType
	}

	if itemType == 30 || itemType == 31 || itemType == 34 {
		return itemType
	}

	return 0
}

func cashSlotItemType(itemID int32) int {
	switch itemID / 10000 {
	case 500:
		return 8
	case 501:
		return 9
	case 502:
		return 10
	case 503:
		return 11
	case 504:
		return 20
	case 505:
		suffix := itemID % 10
		if suffix == 0 {
			return 21
		}
		if suffix > 0 && suffix <= 4 {
			return 22
		}
	case 506:
		suffix := itemID % 10
		switch suffix {
		case 0:
			return 23
		case 1:
			return 24
		case 2:
			return 25
		}
	case 507:
		subType := itemID % 10000 / 1000
		switch subType {
		case 1:
			return 12
		case 2:
			return 13
		case 3, 4:
			return 13
		case 5:
			switch itemID % 10 {
			case 0:
				return 40
			case 1:
				return 41
			case 2:
				return 42
			case 3:
				return 43
			case 4:
				return 44
			case 5:
				return 45
			}
		}
	case 508:
		return 16
	case 509:
		return 19
	case 510:
		return 18
	case 512:
		return 14
	case 513:
		return 7
	case 514:
		return 4
	case 515:
		switch itemID / 1000 {
		case 5150, 5151:
			return 1
		case 5152:
			if itemID/100 == 51520 {
				return 2
			}
		case 5153:
			return 3
		case 5154:
			return 1
		}
	case 516:
		return 6
	case 517:
		if itemID%10000 == 0 {
			return 15
		}
	case 518:
		return 5
	case 519:
		return 26
	case 520:
		return 17
	case 522:
		return 29
	case 523:
		return 27
	case 524:
		return 30
	case 525:
		if itemID%1000 == 100 {
			return 32
		}
		return 33
	case 530:
		return 31
	case 533:
		return 34
	case 537:
		return 28
	case 538:
		return 35
	case 539:
		return 36
	case 540:
		switch itemID / 1000 {
		case 5400:
			return 46
		case 5401:
			return 50
		}
	case 542:
		if itemID/1000 == 5420 {
			return 47
		}
	}

	return 0
}

func (server *Server) handleCashItemUse(plr *Player, reader mpacket.Reader) {
	slot := reader.ReadInt16()
	itemID := reader.ReadInt32()

	item, err := plr.getItem(constant.InventoryCash, slot)
	if err != nil {
		log.Printf("cash use invalid slot: player=%s slot=%d itemID=%d err=%v", plr.Name, slot, itemID, err)
		plr.Send(packetPlayerNoChange())
		return
	}

	if item.ID != itemID || !item.cash || item.cashID == 0 {
		log.Printf("cash use invalid item: player=%s slot=%d want=%d have=%d cash=%t cashID=%d", plr.Name, slot, itemID, item.ID, item.cash, item.cashID)
		plr.Send(packetPlayerNoChange())
		return
	}

	itemType := cashItemUseType(itemID)
	if itemType == 0 {
		switch {
		case isExpCouponItem(itemID):
			itemType = cashItemTypeExpCoupon
		case isDropCouponItem(itemID):
			itemType = cashItemTypeDropCoupon
		}
	}
	if itemType == 0 {
		log.Printf("cash use rejected type: player=%s itemID=%d rawType=%d", plr.Name, itemID, cashSlotItemType(itemID))
		plr.Send(packetPlayerNoChange())
		return
	}

	supported := map[int]bool{
		cashItemTypeSpeakerChannel: true,
		cashItemTypeSpeakerWorld:   true,
		cashItemTypeWeather:        true,
		cashItemTypePetNameTag:     true,
		cashItemTypeMessageBox:     true,
		cashItemTypeMapTransfer:    true,
		cashItemTypeAPReset:        true,
		cashItemTypeSPReset:        true,
		cashItemTypeItemNameTag:    true,
		cashItemTypeChalkboard:     true,
		cashItemTypeExpCoupon:      true,
		cashItemTypeDropCoupon:     true,
	}

	if !supported[itemType] {
		log.Printf("cash use unsupported category: player=%s itemID=%d type=%d", plr.Name, itemID, itemType)
		plr.Send(packetPlayerNoChange())
		return
	}

	var apply func() error
	switch itemType {
	case cashItemTypeSpeakerChannel:
		apply, err = server.prepareCashSpeakerChannel(plr, reader)
	case cashItemTypeSpeakerWorld:
		apply, err = server.prepareCashSpeakerWorld(plr, reader)
	case cashItemTypeWeather:
		apply, err = server.prepareCashWeather(plr, itemID, reader)
	case cashItemTypePetNameTag:
		apply, err = server.prepareCashPetNameTag(plr, reader)
	case cashItemTypeMessageBox:
		apply, err = server.prepareCashMessageBox(plr, itemID, reader)
	case cashItemTypeMapTransfer:
		apply, err = server.prepareCashTeleportRock(plr, itemID, reader)
	case cashItemTypeAPReset:
		apply, err = server.prepareCashAPReset(plr, reader)
	case cashItemTypeSPReset:
		apply, err = server.prepareCashSPReset(plr, itemID, reader)
	case cashItemTypeItemNameTag:
		apply, err = server.prepareCashItemNameTag(plr, reader)
	case cashItemTypeChalkboard:
		apply, err = server.prepareCashChalkboard(plr, reader)
	case cashItemTypeExpCoupon:
		apply, err = server.prepareCashExpCoupon(plr, itemID)
	case cashItemTypeDropCoupon:
		apply, err = server.prepareCashDropCoupon(plr, itemID)
	}

	if err != nil {
		log.Printf("cash use validation failed: player=%s itemID=%d type=%d err=%v", plr.Name, itemID, itemType, err)
		plr.Send(packetPlayerNoChange())
		return
	}

	removed, takeErr := plr.takeItem(itemID, slot, 1, constant.InventoryCash)
	if takeErr != nil {
		log.Printf("cash use consume failed: player=%s itemID=%d slot=%d err=%v", plr.Name, itemID, slot, takeErr)
		plr.Send(packetPlayerNoChange())
		return
	}

	if err := apply(); err != nil {
		log.Printf("cash use apply failed: player=%s itemID=%d type=%d err=%v", plr.Name, itemID, itemType, err)
		if _, rollbackErr := plr.GiveItem(removed); rollbackErr != nil {
			log.Printf("cash use rollback failed: player=%s itemID=%d err=%v", plr.Name, itemID, rollbackErr)
		}
		plr.Send(packetPlayerNoChange())
		return
	}

	plr.Send(packetPlayerNoChange())
}

func (server *Server) prepareCashSpeakerChannel(plr *Player, reader mpacket.Reader) (func() error, error) {
	msg := strings.TrimSpace(reader.ReadString(reader.ReadInt16()))
	if msg == "" {
		return nil, fmt.Errorf("empty channel megaphone message")
	}
	if len(plr.Name)+3+len(msg) > channelMegaphoneMaxDisplayChars {
		return nil, fmt.Errorf("channel megaphone too long")
	}
	return func() error {
		server.players.broadcast(packetMessageBroadcastChannel(plr.Name, msg, server.id, false))
		return nil
	}, nil
}

func (server *Server) prepareCashSpeakerWorld(plr *Player, reader mpacket.Reader) (func() error, error) {
	msg := strings.TrimSpace(reader.ReadString(reader.ReadInt16()))
	if msg == "" {
		return nil, fmt.Errorf("empty world megaphone message")
	}
	if len(msg) > worldMegaphoneMaxChars {
		msg = msg[:worldMegaphoneMaxChars]
	}
	whisper := reader.ReadBool()
	return func() error {
		server.world.Send(internal.PacketChatMegaphone(plr.Name, msg, whisper))
		return nil
	}, nil
}

func (server *Server) prepareCashWeather(plr *Player, itemID int32, reader mpacket.Reader) (func() error, error) {
	msg := strings.TrimSpace(reader.ReadString(reader.ReadInt16()))
	if plr.inst == nil {
		return nil, fmt.Errorf("player not in field")
	}
	formatted := msg
	if itm, err := nx.GetItem(itemID); err == nil && itm.Msg != "" {
		formatted = fmt.Sprintf(itm.Msg, plr.Name, msg)
	}
	return func() error {
		if !plr.inst.startWeatherEffect(itemID, formatted) {
			return fmt.Errorf("weather effect rejected")
		}
		return nil
	}, nil
}

func (server *Server) prepareCashPetNameTag(plr *Player, reader mpacket.Reader) (func() error, error) {
	name := strings.TrimSpace(reader.ReadString(reader.ReadInt16()))
	if len(name) < petNameMinLength || len(name) > petNameMaxLength {
		return nil, fmt.Errorf("invalid pet name length")
	}
	if plr.pet == nil || !plr.pet.spawned {
		return nil, fmt.Errorf("no active pet")
	}
	itemInfo, err := nx.GetItem(plr.pet.itemID)
	if err != nil {
		return nil, fmt.Errorf("pet item metadata missing: %w", err)
	}
	if strings.EqualFold(itemInfo.Name, name) {
		return nil, fmt.Errorf("pet name matches default")
	}
	return func() error {
		plr.pet.name = name
		plr.MarkDirty(DirtyPet, 300*time.Millisecond)
		plr.Send(packetPlayerPetUpdate(plr.pet.lockerSN))
		if plr.inst != nil {
			plr.inst.send(packetPetNameChange(plr.ID, name))
		}
		return nil
	}, nil
}

func (server *Server) prepareCashMessageBox(plr *Player, itemID int32, reader mpacket.Reader) (func() error, error) {
	text := strings.TrimSpace(reader.ReadString(reader.ReadInt16()))
	if text == "" {
		return nil, fmt.Errorf("empty message box text")
	}
	if plr.inst == nil {
		return nil, fmt.Errorf("player not in field")
	}
	return func() error {
		plr.inst.setMessageBox(plr, itemID, text)
		return nil
	}, nil
}

func (server *Server) prepareCashChalkboard(plr *Player, reader mpacket.Reader) (func() error, error) {
	text := strings.TrimSpace(reader.ReadString(reader.ReadInt16()))
	if plr.inst == nil {
		return nil, fmt.Errorf("player not in field")
	}
	if text == "" {
		return nil, fmt.Errorf("empty chalkboard text")
	}
	if plr.inst.fieldID >= 910000001 && plr.inst.fieldID <= 910000022 {
		return nil, fmt.Errorf("chalkboard blocked in free market room")
	}
	return func() error {
		plr.inst.setChalkboard(plr, text)
		return nil
	}, nil
}

func (server *Server) prepareCashTeleportRock(plr *Player, itemID int32, reader mpacket.Reader) (func() error, error) {
	mode := reader.ReadByte()
	bCanTransferContinent := itemID == constant.ItemVIPTeleportRock

	if mode == constant.TeleportToName {
		targetName := reader.ReadString(reader.ReadInt16())
		target, err := server.players.GetFromName(targetName)
		if err != nil || target == nil || target.admin() {
			plr.Send(packetMapTransferResult(mapTransferResultUnavailable, bCanTransferContinent, nil))
			return nil, fmt.Errorf("invalid target player")
		}
		if target.hp <= 0 {
			plr.Send(packetMapTransferResult(mapTransferResultTargetDead, bCanTransferContinent, nil))
			return nil, fmt.Errorf("target is dead")
		}
		if ok, reason := server.canUseTeleportRock(plr, bCanTransferContinent, target.mapID); !ok {
			plr.Send(packetMapTransferResult(reason, bCanTransferContinent, nil))
			return nil, fmt.Errorf("target map blocked: reason=%d", reason)
		}
		targetField, ok := server.fields[target.mapID]
		if !ok || target.inst == nil {
			plr.Send(packetMapTransferResult(mapTransferResultInvalidField, bCanTransferContinent, nil))
			return nil, fmt.Errorf("target field missing")
		}
		portal, err := target.inst.getRandomSpawnPortal()
		if err != nil {
			plr.Send(packetMapTransferResult(mapTransferResultInvalidField, bCanTransferContinent, nil))
			return nil, err
		}
		return func() error {
			return server.warpPlayerToInstance(plr, targetField, portal, target.inst.id, false)
		}, nil
	}

	mapID := reader.ReadInt32()
	if !teleportRockContains(plr.regTeleportRocks, mapID) && !teleportRockContains(plr.vipTeleportRocks, mapID) {
		plr.Send(packetMapTransferResult(mapTransferResultBlocked, bCanTransferContinent, nil))
		return nil, fmt.Errorf("map %d not registered", mapID)
	}
	if !bCanTransferContinent && !teleportRockContains(plr.regTeleportRocks, mapID) {
		plr.Send(packetMapTransferResult(mapTransferResultBlocked, false, nil))
		return nil, fmt.Errorf("regular teleport rock requires regular slot")
	}
	if ok, reason := server.canUseTeleportRock(plr, bCanTransferContinent, mapID); !ok {
		plr.Send(packetMapTransferResult(reason, bCanTransferContinent, nil))
		return nil, fmt.Errorf("destination blocked: reason=%d", reason)
	}
	targetField, ok := server.fields[mapID]
	if !ok {
		plr.Send(packetMapTransferResult(mapTransferResultInvalidField, bCanTransferContinent, nil))
		return nil, fmt.Errorf("field %d missing", mapID)
	}
	targetInst, err := targetField.getInstance(0)
	if err != nil {
		plr.Send(packetMapTransferResult(mapTransferResultInvalidField, bCanTransferContinent, nil))
		return nil, err
	}
	portal, err := targetInst.getRandomSpawnPortal()
	if err != nil {
		plr.Send(packetMapTransferResult(mapTransferResultInvalidField, bCanTransferContinent, nil))
		return nil, err
	}
	return func() error {
		return server.warpPlayer(plr, targetField, portal, false)
	}, nil
}

func (server *Server) canUseTeleportRock(plr *Player, canTransferContinent bool, dstMapID int32) (bool, byte) {
	srcField, ok := server.fields[plr.mapID]
	if !ok || plr.inst == nil {
		return false, mapTransferResultInvalidField
	}
	if srcField.Data.ScrollDisable != 0 || (plr.mapID/1000000)%100 == 9 {
		return false, mapTransferResultInvalidField
	}
	dstField, ok := server.fields[dstMapID]
	if !ok {
		return false, mapTransferResultInvalidField
	}
	if plr.mapID == dstMapID {
		return false, mapTransferResultSameMap
	}
	if plr.level < 7 && dstMapID > 9999999 {
		return false, mapTransferResultLowLevel
	}
	if dstField.Data.ScrollDisable != 0 || dstMapID == 180000000 || (dstMapID/1000000)%100 == 9 || plr.mapID/10000 == 20009 || dstMapID/10000 == 20009 {
		return false, mapTransferResultBlocked
	}
	if !canTransferContinent && !mapsConnectedForTeleport(plr.mapID, dstMapID) {
		return false, mapTransferResultBlocked
	}
	return true, 0
}

func mapsConnectedForTeleport(srcMapID, dstMapID int32) bool {
	return srcMapID/100000000 == dstMapID/100000000
}

func teleportRockContains(rocks []int32, mapID int32) bool {
	for _, rock := range rocks {
		if rock == mapID {
			return true
		}
	}
	return false
}

func (server *Server) prepareCashAPReset(plr *Player, reader mpacket.Reader) (func() error, error) {
	statUp := reader.ReadInt32()
	statDown := reader.ReadInt32()
	if statUp == statDown {
		return nil, fmt.Errorf("identical AP reset stats")
	}

	validStat := func(stat int32) bool {
		switch stat {
		case constant.StrID, constant.DexID, constant.IntID, constant.LukID:
			return true
		default:
			return false
		}
	}
	if !validStat(statUp) || !validStat(statDown) {
		return nil, fmt.Errorf("unsupported AP reset stats up=%d down=%d", statUp, statDown)
	}

	canReset := false
	switch statDown {
	case constant.StrID:
		canReset = plr.str > 4
	case constant.DexID:
		canReset = plr.dex > 4
	case constant.IntID:
		canReset = plr.intt > 4
	case constant.LukID:
		canReset = plr.luk > 4
	}
	if !canReset {
		return nil, fmt.Errorf("stat underflow prevented")
	}

	return func() error {
		switch statDown {
		case constant.StrID:
			plr.setStr(plr.str - 1)
		case constant.DexID:
			plr.setDex(plr.dex - 1)
		case constant.IntID:
			plr.setInt(plr.intt - 1)
		case constant.LukID:
			plr.setLuk(plr.luk - 1)
		}
		switch statUp {
		case constant.StrID:
			plr.setStr(plr.str + 1)
		case constant.DexID:
			plr.setDex(plr.dex + 1)
		case constant.IntID:
			plr.setInt(plr.intt + 1)
		case constant.LukID:
			plr.setLuk(plr.luk + 1)
		}
		return nil
	}, nil
}

func (server *Server) prepareCashSPReset(plr *Player, itemID int32, reader mpacket.Reader) (func() error, error) {
	skillUp := reader.ReadInt32()
	skillDown := reader.ReadInt32()
	if skillUp == skillDown {
		return nil, fmt.Errorf("identical SP reset skills")
	}

	skillUpData, okUp := plr.skills[skillUp]
	skillDownData, okDown := plr.skills[skillDown]
	if !okUp || !okDown || skillDownData.Level == 0 {
		return nil, fmt.Errorf("missing skill records")
	}
	if skillUpData.Level >= skillUpData.Mastery {
		return nil, fmt.Errorf("target skill at mastery")
	}

	jobTier := int(itemID % 10)
	if !cashSPResetSkillAllowed(plr.job, skillUp, jobTier) || !cashSPResetSkillAllowed(plr.job, skillDown, jobTier) {
		return nil, fmt.Errorf("skill tier mismatch up=%d down=%d tier=%d", skillUp, skillDown, jobTier)
	}

	return func() error {
		skillDownData.Level--
		if skillDownData.Level == 0 {
			plr.removeSkill(skillDown)
		} else {
			plr.updateSkill(skillDownData)
		}
		skillUpData.Level++
		plr.updateSkill(skillUpData)
		return nil
	}, nil
}

func cashSPResetSkillAllowed(jobID int16, skillID int32, maxTier int) bool {
	baseSkillID := skillID / 10000
	if !validateSkillWithJob(jobID, baseSkillID) {
		return false
	}
	if baseSkillID == 0 {
		return maxTier >= 1
	}
	return skillTierForBaseSkill(jobID, baseSkillID) <= maxTier
}

func skillTierForBaseSkill(jobID int16, baseSkillID int32) int {
	if baseSkillID == 0 {
		return 1
	}
	if baseSkillID == int32(jobID) {
		return int(jobID % 10)
	}
	for tier := 4; tier >= 1; tier-- {
		candidate := int16((int(jobID)/100)*100 + tier*10)
		if candidate == jobID {
			continue
		}
		if int32(candidate) == baseSkillID {
			return tier
		}
	}
	if baseSkillID%100 == 0 {
		return 1
	}
	return int(baseSkillID % 10)
}

func (server *Server) prepareCashItemNameTag(plr *Player, reader mpacket.Reader) (func() error, error) {
	targetSlot := reader.ReadInt16()
	item, err := plr.getItem(constant.InventoryEquip, targetSlot)
	if err != nil {
		return nil, err
	}
	if item.ID/1000000 != constant.InventoryEquip || item.creatorName != "" {
		return nil, fmt.Errorf("item not taggable")
	}
	return func() error {
		item.creatorName = plr.Name
		if _, err := item.save(plr.ID); err != nil {
			return err
		}
		plr.updateItem(item)
		plr.Send(packetInventoryAddItem(item, true))
		if item.slotID < 0 && plr.inst != nil {
			plr.inst.broadcastAvatarChange(plr)
		}
		return nil
	}, nil
}

func (server *Server) prepareCashExpCoupon(plr *Player, itemID int32) (func() error, error) {
	item, err := nx.GetItem(itemID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return func() error {
		plr.activateExpCoupon(itemID, now)
		plr.Send(packetMessageNotice(couponActivationNotice(item.Name, "EXP")))
		return nil
	}, nil
}

func (server *Server) prepareCashDropCoupon(plr *Player, itemID int32) (func() error, error) {
	item, err := nx.GetItem(itemID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return func() error {
		plr.activateDropCoupon(itemID, now)
		plr.Send(packetMessageNotice(couponActivationNotice(item.Name, "Drop")))
		return nil
	}, nil
}

func couponActivationNotice(itemName, kind string) string {
	itemName = strings.TrimSpace(itemName)
	if itemName == "" {
		return fmt.Sprintf("2x %s coupon is now active.", kind)
	}
	return fmt.Sprintf("%s is now active.", itemName)
}
