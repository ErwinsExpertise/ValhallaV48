package channel

import (
	"github.com/Hucaru/Valhalla/mpacket"
)

type movementFrag struct {
	x, y, vx, vy, foothold, duration int16
	originFh                         int16
	stance, stat, mType              byte
	equipData                        byte
	posSet                           bool
}

type movement struct {
	origX, origY      int16
	frags             []movementFrag
	stateCount        byte
	packedStateValues []byte
	minX, minY        int16
	maxX, maxY        int16
}

// values from WvsGlobal
var movementType = struct {
	normalMovement  byte
	jump            byte
	jumpKb          byte
	immediate       byte
	teleport        byte
	normalMovement2 byte
	flashJump       byte
	assaulter       byte
	immediate2      byte
	chair           byte
	equipMovement   byte
	chair2          byte
	startWings      byte
}{
	normalMovement:  0,
	jump:            1,
	jumpKb:          2,
	immediate:       3, // GM F1 teleport
	teleport:        4,
	normalMovement2: 5,
	flashJump:       6,
	assaulter:       7,
	immediate2:      8,
	chair:           9,
	equipMovement:   10,
	chair2:          11,
	startWings:      12,
}

func parseMovement(reader mpacket.Reader) (movement, movementFrag, []byte, int, bool) {
	mData := movement{}
	startRemaining := len(reader.GetRestAsBytes())

	mData.origX = reader.ReadInt16()
	mData.origY = reader.ReadInt16()

	nFrags := reader.ReadByte()
	mData.frags = make([]movementFrag, nFrags)

	final := movementFrag{x: mData.origX, y: mData.origY, posSet: false}

	for i := byte(0); i < nFrags; i++ {
		frag := movementFrag{posSet: false}
		frag.mType = reader.ReadByte()

		switch frag.mType {
		case movementType.normalMovement,
			movementType.normalMovement2:
			frag.x = reader.ReadInt16()
			frag.y = reader.ReadInt16()
			frag.vx = reader.ReadInt16()
			frag.vy = reader.ReadInt16()
			frag.foothold = reader.ReadInt16()
			frag.stance = reader.ReadByte()
			frag.duration = reader.ReadInt16()
			frag.posSet = true

		case movementType.jump,
			movementType.jumpKb,
			movementType.flashJump,
			movementType.startWings:
			frag.vx = reader.ReadInt16()
			frag.vy = reader.ReadInt16()
			frag.stance = reader.ReadByte()
			frag.duration = reader.ReadInt16()

		case movementType.immediate,
			movementType.teleport,
			movementType.assaulter,
			movementType.immediate2,
			movementType.chair,
			movementType.chair2:
			frag.x = reader.ReadInt16()
			frag.y = reader.ReadInt16()
			frag.foothold = reader.ReadInt16()
			frag.stance = reader.ReadByte()
			frag.duration = reader.ReadInt16()
			frag.posSet = true

		case movementType.equipMovement:
			frag.equipData = reader.ReadByte()

		default:
			frag.stance = reader.ReadByte()
			frag.duration = reader.ReadInt16()
		}

		if frag.posSet {
			final.x = frag.x
			final.y = frag.y
			final.foothold = frag.foothold
			final.posSet = true
		}
		final.mType = frag.mType
		final.stance = frag.stance

		mData.frags[i] = frag
	}

	mData.stateCount = reader.ReadByte()
	mData.packedStateValues = reader.ReadBytes(int((mData.stateCount + 1) / 2))
	mData.minX = reader.ReadInt16()
	mData.minY = reader.ReadInt16()
	mData.maxX = reader.ReadInt16()
	mData.maxY = reader.ReadInt16()

	consumed := startRemaining - len(reader.GetRestAsBytes())
	fragTypes := make([]byte, len(mData.frags))
	for i, frag := range mData.frags {
		fragTypes[i] = frag.mType
	}

	return mData, final, fragTypes, consumed, len(reader.GetRestAsBytes()) == 0
}

func parseMobMovement(reader mpacket.Reader, startX, startY int16) (movement, movementFrag, []byte, int, bool) {
	mData := movement{origX: startX, origY: startY}
	startRemaining := len(reader.GetRestAsBytes())

	nFrags := reader.ReadByte()
	mData.frags = make([]movementFrag, nFrags)
	final := movementFrag{x: startX, y: startY, posSet: true}

	for i := byte(0); i < nFrags; i++ {
		frag := movementFrag{posSet: false}
		frag.mType = reader.ReadByte()

		switch frag.mType {
		case movementType.normalMovement,
			movementType.normalMovement2:
			frag.x = reader.ReadInt16()
			frag.y = reader.ReadInt16()
			frag.vx = reader.ReadInt16()
			frag.vy = reader.ReadInt16()
			frag.foothold = reader.ReadInt16()
			frag.stance = reader.ReadByte()
			frag.duration = reader.ReadInt16()
			frag.posSet = true

		case movementType.jump,
			movementType.jumpKb,
			movementType.flashJump,
			movementType.startWings:
			frag.vx = reader.ReadInt16()
			frag.vy = reader.ReadInt16()
			frag.stance = reader.ReadByte()
			frag.duration = reader.ReadInt16()

		case movementType.immediate,
			movementType.teleport,
			movementType.assaulter,
			movementType.immediate2,
			movementType.chair,
			movementType.chair2:
			frag.x = reader.ReadInt16()
			frag.y = reader.ReadInt16()
			frag.foothold = reader.ReadInt16()
			frag.stance = reader.ReadByte()
			frag.duration = reader.ReadInt16()
			frag.posSet = true

		case movementType.equipMovement:
			frag.equipData = reader.ReadByte()

		default:
			frag.stance = reader.ReadByte()
			frag.duration = reader.ReadInt16()
		}

		if frag.posSet {
			final.x = frag.x
			final.y = frag.y
			final.foothold = frag.foothold
			final.posSet = true
		}
		final.mType = frag.mType
		final.stance = frag.stance
		mData.frags[i] = frag
	}

	mData.stateCount = reader.ReadByte()
	mData.packedStateValues = reader.ReadBytes(int((mData.stateCount + 1) / 2))
	mData.minX = reader.ReadInt16()
	mData.minY = reader.ReadInt16()
	mData.maxX = reader.ReadInt16()
	mData.maxY = reader.ReadInt16()

	consumed := startRemaining - len(reader.GetRestAsBytes())
	fragTypes := make([]byte, len(mData.frags))
	for i, frag := range mData.frags {
		fragTypes[i] = frag.mType
	}

	return mData, final, fragTypes, consumed, len(reader.GetRestAsBytes()) == 0
}

func generateMovementBytes(moveData movement) mpacket.Packet {
	p := mpacket.NewPacket()

	p.WriteInt16(moveData.origX)
	p.WriteInt16(moveData.origY)

	p.WriteByte(byte(len(moveData.frags)))

	for _, frag := range moveData.frags {
		p.WriteByte(frag.mType)

		switch frag.mType {
		case movementType.normalMovement,
			movementType.normalMovement2:
			p.WriteInt16(frag.x)
			p.WriteInt16(frag.y)
			p.WriteInt16(frag.vx)
			p.WriteInt16(frag.vy)
			p.WriteInt16(frag.foothold)
			p.WriteByte(frag.stance)
			p.WriteInt16(frag.duration)

		case movementType.jump,
			movementType.jumpKb,
			movementType.flashJump,
			movementType.startWings:
			p.WriteInt16(frag.vx)
			p.WriteInt16(frag.vy)
			p.WriteByte(frag.stance)
			p.WriteInt16(frag.duration)

		case movementType.immediate,
			movementType.teleport,
			movementType.chair,
			movementType.assaulter,
			movementType.immediate2,
			movementType.chair2:
			p.WriteInt16(frag.x)
			p.WriteInt16(frag.y)
			p.WriteInt16(frag.foothold)
			p.WriteByte(frag.stance)
			p.WriteInt16(frag.duration)

		case movementType.equipMovement:
			p.WriteByte(frag.equipData)

		default:
			p.WriteByte(frag.stance)
			p.WriteInt16(frag.duration)
		}
	}

	p.WriteByte(moveData.stateCount)
	p.WriteBytes(moveData.packedStateValues)
	p.WriteInt16(moveData.minX)
	p.WriteInt16(moveData.minY)
	p.WriteInt16(moveData.maxX)
	p.WriteInt16(moveData.maxY)

	return p
}

func (data movement) validateChar(player *Player) bool {
	// Check for suspicious movement (teleport hacks)
	if len(data.frags) > 0 {
		lastFrag := data.frags[len(data.frags)-1]
		if lastFrag.posSet {
			dx := lastFrag.x - data.origX
			dy := lastFrag.y - data.origY
			if dx < 0 {
				dx = -dx
			}
			if dy < 0 {
				dy = -dy
			}
			distance := dx
			if dy > distance {
				distance = dy
			}

			// Suspicious immediate movement over 1000 pixels
			if distance > 1000 {
				for _, frag := range data.frags {
					if frag.mType == movementType.immediate && frag.mType != movementType.teleport {
						return false // Invalid immediate movement
					}
				}
			}
		}
	}

	return true
}

type mob interface {
}

// ValidateMob movement
func (data movement) validateMob(mob mob) bool {
	// run through the movement data and make sure monsters are not moving too fast

	return true
}
