package channel

import (
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/internal"
	"github.com/Hucaru/Valhalla/mpacket"
)

// TODO: login server needs to Send a deleted character event so that they can leave the party for playing players

type party struct {
	serverChannelID int32
	players         [constant.MaxPartySize]*Player
	internal.Party
}

func (d party) broadcast(p mpacket.Packet) {
	for _, v := range d.players {
		if v == nil {
			continue
		}

		v.Send(p)
	}
}

func (d *party) addExistingPlayer(plr *Player) bool {
	for i, id := range d.PlayerID {
		if id == plr.ID {
			d.players[i] = plr
			plr.party = d
			return true
		}
	}

	return false
}

func (d *party) relinkPlayers(players *Players) {
	var linked [constant.MaxPartySize]*Player

	for i, id := range d.PlayerID {
		if id == 0 {
			continue
		}

		if plr, err := players.GetFromID(id); err == nil {
			linked[i] = plr
			plr.party = d
		}
	}

	d.players = linked
}
func (d *party) addPlayer(plr *Player, index int32) {
	if plr != nil {
		d.players[index] = plr
		plr.party = d
	}

	p := packetPartyPlayerJoin(d.ID, d.Name[index], d)
	d.broadcast(p)
}

func (d *party) getPlayerIndex(plrID int32) byte {
	for i, pid := range d.PlayerID {
		if pid == plrID {
			return byte(i)
		}
	}
	return 0
}

func (d *party) removePlayer(index int32, kick bool, formerID int32, formerName string) {
	playerID := formerID
	if playerID == 0 {
		playerID = d.PlayerID[index]
	}
	name := formerName
	if name == "" {
		name = d.Name[index]
	}

	if index == 0 {
		p := packetPartyLeave(d.ID, playerID, false, kick, "", d)
		d.broadcast(p)

		for _, p := range d.players {
			if p != nil {
				p.party = nil
			}
		}
	} else {
		if d.players[index] != nil {
			d.players[index].party = nil
		}

		p := packetPartyLeave(d.ID, playerID, true, kick, name, d)
		d.broadcast(p)

		d.players[index] = nil
	}

}

func (d party) full() bool {
	for _, v := range d.players {
		if v == nil {
			return false
		}
	}

	return true
}

func (d *party) updateOnlineStatus(index int32, plr *Player) {
	d.players[index] = plr

	if plr != nil {
		plr.party = d
	}

	p := packetPartyUpdate(d.ID, d)
	d.broadcast(p)
}

func (d *party) syncPlayersHP() {
	for index := range d.players {
		plr := d.players[index]

		if plr == nil {
			continue
		}

		d.broadcast(packetPlayerHpChange(plr.ID, int32(plr.hp), int32(plr.maxHP)))
	}
}

func (d party) giveExp(playerID, amount int32, sameMap bool) {
	var mapID int32 = 0
	var instanceID int = -1

	for i, id := range d.PlayerID {
		if id == playerID {
			mapID = d.MapID[i]
			if d.players[i] != nil && d.players[i].inst != nil {
				instanceID = d.players[i].inst.id
			}
			break
		}
	}

	if sameMap {
		nPlayers := 0

		for i, id := range d.PlayerID {
			if id != playerID && d.players[i] != nil && d.MapID[i] == mapID && d.players[i].hp > 0 && d.players[i].inst != nil && d.players[i].inst.id == instanceID {
				nPlayers++
			}

			if nPlayers == 0 {
				return
			}
		}
	}

	for _, plr := range d.players {
		if plr != nil && sameMap && plr.mapID == mapID && plr.inst != nil && plr.inst.id == instanceID {
			plr.giveEXP(amount, false, true)
		}
	}
}

// - Index is 1..6 (party UI ordering), total is the total party member count (byte from packet)
// - Returns 0xFF if index > total (matching the original behavior)
func partyMemberMaskForIndex(index int, total byte) byte {
	var base int
	switch index {
	case 1:
		base = 0x40
	case 2:
		base = 0x80
	case 3:
		base = 0x100
	case 4:
		base = 0x200
	case 5:
		base = 0x400
	case 6:
		base = 0x800
	default:
		return 0xFF
	}

	if int(total) >= index {
		v := base >> uint(total)
		return byte(v & 0xFF)
	}
	return 0xFF
}

func packetPartyCreate(partyID int32, doorMap1, doorMap2 int32, point pos) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0x07)
	p.WriteInt32(partyID)

	if doorMap1 > -1 {
		p.WriteInt32(doorMap1)
		p.WriteInt32(doorMap2)
		p.WriteInt16(point.x)
		p.WriteInt16(point.y)
	} else {
		p.WriteInt32(-1)
		p.WriteInt32(-1)
		p.WriteInt16(0)
		p.WriteInt16(0)
	}

	return p
}

func packetPartyPlayerJoin(partyID int32, name string, party *party) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0x0e)
	p.WriteInt32(partyID)
	p.WriteString(name)

	updateParty(&p, party)

	return p
}

func packetPartyUpdate(partyID int32, party *party) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0x1a)
	p.WriteInt32(partyID)

	updateParty(&p, party)

	return p
}

func packetPartyLeaderChange(partyID int32, party *party) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0x06)
	p.WriteInt32(partyID)
	updateParty(&p, party)
	return p
}

func packetPartyLeave(partyID, playerID int32, keepParty, kicked bool, name string, party *party) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0x0b)
	p.WriteInt32(partyID)
	p.WriteInt32(playerID)
	p.WriteBool(keepParty)

	if keepParty {
		p.WriteBool(kicked)
		p.WriteString(name)
		updateParty(&p, party)
	}

	return p
}

func packetPartyDoorUpdate(index byte, townID, mapID int32, point pos) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0x1c)
	p.WriteByte(index)
	p.WriteInt32(townID)
	p.WriteInt32(mapID)
	p.WriteInt16(point.x)
	p.WriteInt16(point.y)

	return p
}

func packetPartyUpdateJobLevel(playerID, job, level int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0x1b)
	p.WriteInt32(playerID)
	p.WriteInt32(level)
	p.WriteInt32(job)

	return p
}

func updateParty(p *mpacket.Packet, party *party) {
	for i := 0; i < constant.MaxPartySize; i++ {
		p.WriteInt32(party.PlayerID[i])
	}

	for i := 0; i < constant.MaxPartySize; i++ {
		p.WritePaddedString(party.Name[i], 13)
	}

	for i := 0; i < constant.MaxPartySize; i++ {
		p.WriteInt32(party.Job[i])
	}

	for i := 0; i < constant.MaxPartySize; i++ {
		p.WriteInt32(party.Level[i])
	}

	for i := 0; i < constant.MaxPartySize; i++ {
		p.WriteInt32(party.ChannelID[i]) // -1 cash shop, -2 offline
	}

	for i := 0; i < constant.MaxPartySize; i++ {
		p.WriteInt32(party.MapID[i])
	}

	// Mystic door
	for i := 0; i < constant.MaxPartySize; i++ {
		p.WriteInt32(-1)
		p.WriteInt32(-1)
		p.WriteInt32(0) // x
		p.WriteInt32(0) // y
	}
}

func packetPartyCreateUnkownError() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0)

	return p
}

func packetPartyInviteNotice(partyID int32, fromName string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(0x04)
	p.WriteInt32(partyID)
	p.WriteString(fromName)

	return p
}

func packetPartyMessage(op byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(op)

	return p
}

func packetPartyAlreadyJoined() mpacket.Packet {
	return packetPartyMessage(0x08)
}

func packetPartyBeginnerCannotCreate() mpacket.Packet {
	return packetPartyMessage(0x09)
}

func packetPartyNotInParty() mpacket.Packet {
	return packetPartyMessage(0x0c)
}

func packetPartyAlreadyJoined2() mpacket.Packet {
	return packetPartyMessage(0x0f)
}

func packetPartyToJoinIsFull() mpacket.Packet {
	return packetPartyMessage(0x10)
}

func packetPartyUnableToFindPlayer() mpacket.Packet {
	return packetPartyMessage(0x11)
}

func packetPartyAdminNoCreate() mpacket.Packet {
	return packetPartyMessage(0x18)
}

func packetPartyUnableToFindPlayer2() mpacket.Packet {
	return packetPartyMessage(0x19)
}

func packetPartyMessageName(op byte, name string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPartyInfo)
	p.WriteByte(op)
	p.WriteString(name)

	return p
}

func packetPartyBlockingInvites(name string) mpacket.Packet {
	return packetPartyMessageName(0x13, name)
}

func packetPartyHasOtherRequest(name string) mpacket.Packet {
	return packetPartyMessageName(0x14, name)
}

func packetPartyRequestDenied(name string) mpacket.Packet {
	return packetPartyMessageName(0x15, name)
}
