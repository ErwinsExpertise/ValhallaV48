package channel

import (
	"time"

	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/constant/skill"
	"github.com/Hucaru/Valhalla/mpacket"
)

// Summon represents a Player-owned summon/puppet instance.
type summon struct {
	OwnerID int32
	SkillID int32
	Level   byte
	HP      int

	Pos       pos
	Stance    byte
	Foothold  int16
	ExpiresAt time.Time

	// Flags
	IsPuppet   bool
	SummonType int32
}

type summonState struct {
	puppet *summon
	summon *summon
}

func summonMovementType(su *summon) byte {
	if su == nil {
		return 0
	}
	if su.IsPuppet {
		return 0
	}
	switch skill.Skill(su.SkillID) {
	case skill.SilverHawk, skill.GoldenEagle, skill.SummonDragon:
		return 3
	default:
		return 0
	}
}

func packetShowSummon(ownerID int32, su *summon, animated bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSpecialMapObjectSpawn)
	p.WriteInt32(ownerID)
	p.WriteInt32(su.SkillID)
	p.WriteByte(su.Level)
	p.WriteInt16(su.Pos.x)
	p.WriteInt16(su.Pos.y)
	p.WriteByte(0)
	p.WriteByte(0)
	p.WriteByte(0)
	p.WriteByte(summonMovementType(su))
	p.WriteByte(1)
	if animated {
		p.WriteByte(0)
	} else {
		p.WriteByte(1)
	}
	return p
}

func packetRemoveSummon(ownerID int32, summonID int32, reason byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSpecialMapObjectRemove)
	p.WriteInt32(ownerID)
	p.WriteInt32(summonID)
	p.WriteByte(reason)
	return p
}

func packetSummonMove(ownerID int32, summonID int32, start pos, moveBytes mpacket.Packet) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSummonMove)
	p.WriteInt32(ownerID)
	p.WriteInt32(summonID)
	p.WriteInt16(start.x)
	p.WriteInt16(start.y)
	p.WriteBytes(moveBytes)
	return p
}

func packetSummonAttack(ownerID int32, summonID int32, anim byte, targets byte, mobDamages map[int32][]int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSummonAttack)
	p.WriteInt32(ownerID)
	p.WriteInt32(summonID)
	p.WriteByte(anim)
	p.WriteByte(targets)
	for mobID, dList := range mobDamages {
		p.WriteInt32(mobID)
		p.WriteByte(constant.SummonAttackMob)
		for _, d := range dList {
			p.WriteInt32(d)
		}
	}
	return p
}

func packetSummonDamage(ownerID int32, summonID int32, damage int32, mobID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSummonDamage)
	p.WriteInt32(ownerID)
	p.WriteInt32(summonID)
	p.WriteByte(constant.SummonTakeDamage)
	p.WriteInt32(damage)
	p.WriteInt32(mobID)
	p.WriteByte(0)
	return p
}
