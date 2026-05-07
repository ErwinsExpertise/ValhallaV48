package channel

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

type pet struct {
	name            string
	itemID          int32
	lockerSN        int64
	sn              int32
	itemDBID        int64
	level           byte
	closeness       int16
	fullness        byte
	deadDate        int64
	spawnDate       int64
	lastInteraction int64

	pos    pos
	stance byte

	spawned bool
}

func newPet(itemID, sn int32, dbID int64, lockerSN int64) *pet {
	itemInfo, err := nx.GetItem(itemID)
	if err != nil {
		log.Println(err)
	}

	return &pet{
		name:            itemInfo.Name,
		itemID:          itemID,
		lockerSN:        lockerSN,
		sn:              sn,
		itemDBID:        dbID,
		stance:          0,
		level:           1,
		closeness:       0,
		fullness:        100,
		deadDate:        petExpiryTime(),
		spawnDate:       0,
		lastInteraction: 0,
	}
}

func petExpiryTime() int64 {
	return time.Now().Add(90*24*time.Hour).UnixMilli()*10000 + 116444592000000000
}

func savePet(item *Item) error {
	// Initialize pet data if it doesn't exist
	if item.petData == nil {
		sn, _ := nx.GetCommoditySNByItemID(item.ID)
		item.petData = newPet(item.ID, sn, item.dbID, item.cashID)
		if item.expireTime != 0 && item.expireTime != neverExpire {
			item.petData.deadDate = item.expireTime
		}
	}

	if item.petData.deadDate != 0 {
		item.expireTime = item.petData.deadDate
	}

	p := item.petData

	_, err := common.DB.Exec("INSERT INTO pets (parentID, name, sn, `level`, closeness, fullness, deadDate, spawnDate, lastInteraction) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) AS new ON DUPLICATE KEY UPDATE name=new.name, `level`=new.`level`, closeness=new.closeness, fullness=new.fullness, deadDate=new.deadDate, spawnDate=new.spawnDate, lastInteraction=new.lastInteraction", item.dbID,
		p.name,
		p.sn,
		p.level,
		p.closeness,
		p.fullness,
		p.deadDate,
		p.spawnDate,
		p.lastInteraction,
	)
	return err
}

func (p *pet) updateMovement(frag movementFrag) {
	p.pos.x = frag.x
	p.pos.y = frag.y
	p.pos.foothold = frag.foothold
	p.stance = frag.stance
}

func handlePetInteraction(plr *Player, pet *pet, interactionID byte, multiplier bool) bool {
	itm, err := nx.GetItem(pet.itemID)
	if err != nil || itm.Interact == nil {
		return false
	}
	react, ok := itm.Interact[interactionID]
	if !ok {
		return false
	}

	now := time.Now().UnixMilli()
	if now < pet.lastInteraction+15_000 || pet.level < react.LevelMin || pet.level > react.LevelMax || pet.fullness < 50 {
		return false
	}

	elapsed := float64(now - pet.lastInteraction - 15_000)
	pet.lastInteraction = now
	plr.MarkDirty(DirtyPet, time.Millisecond*300)

	mult := 1.0
	if multiplier && pet.name != "" {
		mult = 1.5
	}
	successProb := float64(react.Prob) * ((elapsed/10_000.0)*0.01 + 1) * mult
	success := float64(rand.Intn(100)) < successProb
	if success {
		pet.closeness += int16(react.Inc)
		if pet.closeness < 0 {
			pet.closeness = 0
		}
		if pet.closeness > 30_000 {
			pet.closeness = 30_000
		}
		pet.level = petLevelFromCloseness(pet.closeness)
		plr.updatePet()
	}
	return success
}

var thresholds = []int16{0, 1, 100, 300, 600, 1000, 1800, 3100, 5000, 8000, 12000, 17000, 22000, 28000}

func petLevelFromCloseness(c int16) byte {
	for lvl := byte(len(thresholds) - 1); lvl > 0; lvl-- {
		if c >= thresholds[lvl] {
			return lvl
		}
	}
	return 1
}

func packetPetAction(charID int32, op, action byte, text string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetAction)
	p.WriteInt32(charID)
	p.WriteByte(op)
	p.WriteByte(action)
	p.WriteString(text)
	p.WriteByte(0)
	return p
}

func packetPetNameChange(charID int32, name string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetNameChange)
	p.WriteInt32(charID)
	p.WriteString(name)
	return p
}

func packetPetFoodResponse(charID int32, success bool, balloon bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetInteraction)
	p.WriteInt32(charID)
	p.WriteByte(1)
	p.WriteBool(success)
	p.WriteBool(balloon)

	return p
}

func packetPetInteraction(charID int32, interactionId byte, success bool, balloon bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetInteraction)
	p.WriteInt32(charID)
	p.WriteByte(0)
	p.WriteByte(interactionId)
	p.WriteBool(success)
	p.WriteBool(balloon)

	return p
}

func packetPetMove(charID int32, move []byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetMove)
	p.WriteInt32(charID)
	p.WriteBytes(move)
	return p
}

func writePetInitData(p *mpacket.Packet, petData *pet) {
	p.WriteInt32(petData.itemID)
	p.WriteString(petData.name)
	p.WriteUint64(uint64(petData.lockerSN))
	p.WriteInt16(petData.pos.x)
	p.WriteInt16(petData.pos.y)
	p.WriteByte(petData.stance)
	p.WriteInt16(petData.pos.foothold)
	p.WriteBool(false) // bNameTag
	p.WriteBool(false) // bChatBalloon
}

func packetPetSpawn(charID int32, petData *pet) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetSpawn)
	p.WriteInt32(charID)
	writePetInitData(&p, petData)

	return p
}

func packetPetRemove(charID int32, reason byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPetRemove)
	p.WriteInt32(charID)
	p.WriteBool(false)
	p.WriteByte(reason)

	return p
}

func (p *pet) canConsumeFood(meta nx.Item) bool {
	if p == nil || meta.PetFoodInc <= 0 {
		return false
	}
	for _, petID := range meta.PetFoodPets {
		if petID == p.itemID {
			return true
		}
	}
	return false
}

func handlePetFoodUse(plr *Player, slot int16, item Item, meta nx.Item) error {
	if plr == nil || plr.pet == nil || !plr.pet.spawned {
		return fmt.Errorf("no active pet")
	}
	if !plr.pet.canConsumeFood(meta) {
		return fmt.Errorf("item %d is not valid pet food for pet %d", item.ID, plr.pet.itemID)
	}
	if _, err := plr.takeItem(item.ID, slot, 1, constant.InventoryUse); err != nil {
		return err
	}

	success := false
	oldLevel := plr.pet.level
	oldFullness := plr.pet.fullness
	oldCloseness := plr.pet.closeness
	if plr.pet.fullness < 100 {
		gain := int(meta.PetFoodInc)
		if gain > int(100-plr.pet.fullness) {
			gain = int(100 - plr.pet.fullness)
		}
		plr.pet.fullness += byte(gain)
		if plr.pet.closeness < 30000 {
			plr.pet.closeness++
		}
		success = gain > 0
	} else if plr.pet.closeness > 0 {
		plr.pet.closeness--
	}

	plr.pet.level = petLevelFromCloseness(plr.pet.closeness)
	plr.MarkDirty(DirtyPet, time.Millisecond*300)
	petStateChanged := plr.pet.fullness != oldFullness || plr.pet.closeness != oldCloseness || plr.pet.level != oldLevel
	if petStateChanged {
		plr.Send(packetPlayerPetUpdate(plr.pet.lockerSN))
		if petItem, _, err := plr.GetItemByCashID(constant.InventoryCash, plr.petCashID); err == nil {
			plr.Send(packetInventoryAddItem(petItem, true))
		}
	}
	plr.Send(packetPetFoodResponse(plr.ID, success, false))

	if success && plr.pet.level > oldLevel {
		// v48 pet level-up visuals still need live confirmation; keep the state update authoritative.
	}

	return nil
}
