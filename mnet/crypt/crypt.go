package crypt

import (
	"crypto/aes"
)

const (
	encryptHeaderSize = 4
)

type Maple struct {
	mapleVersion int
	key          [16]byte
}

func New(key [4]byte, mapleVersion int) *Maple {
	var c Maple

	for i := 0; i < 4; i++ {
		copy(c.key[4*i:], key[:])
	}

	swapped := ((mapleVersion >> 8) & 0xFF) | ((mapleVersion << 8) & 0xFF00)
	c.mapleVersion = int(int16(swapped))

	return &c
}

func (c *Maple) Encrypt(p []byte, maple, aes bool) {
	c.generateHeader(p)

	if maple {
		mapleCrypt(p[encryptHeaderSize:])
	}

	if aes {
		c.aesCrypt(p[encryptHeaderSize:])
	}

	c.Shuffle()
}

func (c *Maple) Decrypt(p []byte, maple, aes bool) {
	if aes {
		c.aesCrypt(p)
	}

	c.Shuffle()

	if maple {
		mapleDecrypt(p)
	}
}

func (c *Maple) IV() []byte {
	return c.key[:]
}

func GetPacketLength(encryptedHeader []byte) int {
	packetHeader := int(encryptedHeader[0])<<24 | int(encryptedHeader[1])<<16 | int(encryptedHeader[2])<<8 | int(encryptedHeader[3])
	packetLength := ((packetHeader >> 16) ^ (packetHeader & 0xFFFF))
	packetLength = ((packetLength << 8) & 0xFF00) | ((packetLength >> 8) & 0xFF)
	return packetLength
}

func (c *Maple) checkPacket(packet []byte) bool {
	return ((packet[0]^byte(c.mapleVersion>>8))&0xFF == c.key[2]) &&
		((packet[1]^byte(c.mapleVersion))&0xFF == c.key[3])
}

var funnyBytes = [...]byte{
	0xEC, 0x3F, 0x77, 0xA4, 0x45, 0xD0, 0x71, 0xBF, 0xB7, 0x98, 0x20, 0xFC, 0x4B, 0xE9, 0xB3, 0xE1,
	0x5C, 0x22, 0xF7, 0x0C, 0x44, 0x1B, 0x81, 0xBD, 0x63, 0x8D, 0xD4, 0xC3, 0xF2, 0x10, 0x19, 0xE0,
	0xFB, 0xA1, 0x6E, 0x66, 0xEA, 0xAE, 0xD6, 0xCE, 0x06, 0x18, 0x4E, 0xEB, 0x78, 0x95, 0xDB, 0xBA,
	0xB6, 0x42, 0x7A, 0x2A, 0x83, 0x0B, 0x54, 0x67, 0x6D, 0xE8, 0x65, 0xE7, 0x2F, 0x07, 0xF3, 0xAA,
	0x27, 0x7B, 0x85, 0xB0, 0x26, 0xFD, 0x8B, 0xA9, 0xFA, 0xBE, 0xA8, 0xD7, 0xCB, 0xCC, 0x92, 0xDA,
	0xF9, 0x93, 0x60, 0x2D, 0xDD, 0xD2, 0xA2, 0x9B, 0x39, 0x5F, 0x82, 0x21, 0x4C, 0x69, 0xF8, 0x31,
	0x87, 0xEE, 0x8E, 0xAD, 0x8C, 0x6A, 0xBC, 0xB5, 0x6B, 0x59, 0x13, 0xF1, 0x04, 0x00, 0xF6, 0x5A,
	0x35, 0x79, 0x48, 0x8F, 0x15, 0xCD, 0x97, 0x57, 0x12, 0x3E, 0x37, 0xFF, 0x9D, 0x4F, 0x51, 0xF5,
	0xA3, 0x70, 0xBB, 0x14, 0x75, 0xC2, 0xB8, 0x72, 0xC0, 0xED, 0x7D, 0x68, 0xC9, 0x2E, 0x0D, 0x62,
	0x46, 0x17, 0x11, 0x4D, 0x6C, 0xC4, 0x7E, 0x53, 0xC1, 0x25, 0xC7, 0x9A, 0x1C, 0x88, 0x58, 0x2C,
	0x89, 0xDC, 0x02, 0x64, 0x40, 0x01, 0x5D, 0x38, 0xA5, 0xE2, 0xAF, 0x55, 0xD5, 0xEF, 0x1A, 0x7C,
	0xA7, 0x5B, 0xA6, 0x6F, 0x86, 0x9F, 0x73, 0xE6, 0x0A, 0xDE, 0x2B, 0x99, 0x4A, 0x47, 0x9C, 0xDF,
	0x09, 0x76, 0x9E, 0x30, 0x0E, 0xE4, 0xB2, 0x94, 0xA0, 0x3B, 0x34, 0x1D, 0x28, 0x0F, 0x36, 0xE3,
	0x23, 0xB4, 0x03, 0xD8, 0x90, 0xC8, 0x3C, 0xFE, 0x5E, 0x32, 0x24, 0x50, 0x1F, 0x3A, 0x43, 0x8A,
	0x96, 0x41, 0x74, 0xAC, 0x52, 0x33, 0xF0, 0xD9, 0x29, 0x80, 0xB1, 0x16, 0xD3, 0xAB, 0x91, 0xB9,
	0x84, 0x7F, 0x61, 0x1E, 0xCF, 0xC5, 0xD1, 0x56, 0x3D, 0xCA, 0xF4, 0x05, 0xC6, 0xE5, 0x08, 0x49,
	0x4F, 0x64, 0x69, 0x6E, 0x4D, 0x53, 0x7E, 0x46, 0x72, 0x7A}

func (c *Maple) Shuffle() {
	newIV := []byte{0xF2, 0x53, 0x50, 0xC6}

	for i := 0; i < 4; i++ {
		input := c.key[i]
		elina := newIV[1]
		anna := input
		moritz := funnyBytes[elina]
		moritz -= input
		newIV[0] += moritz
		moritz = newIV[2]
		moritz ^= funnyBytes[anna]
		elina -= byte(int(moritz) & 0xFF)
		newIV[1] = elina
		elina = newIV[3]
		moritz = elina
		elina -= byte(int(newIV[0]) & 0xFF)
		moritz = funnyBytes[moritz]
		moritz += input
		moritz ^= newIV[2]
		newIV[2] = moritz
		elina += funnyBytes[anna]
		newIV[3] = elina

		merry := uint32(newIV[0]) & 0xFF
		merry |= (uint32(newIV[1]) << 8) & 0xFF00
		merry |= (uint32(newIV[2]) << 16) & 0xFF0000
		merry |= (uint32(newIV[3]) << 24) & 0xFF000000
		retValue := merry >> 0x1D
		merry = merry << 3
		retValue = retValue | merry

		newIV[0] = byte(retValue & 0xFF)
		newIV[1] = byte((retValue >> 8) & 0xFF)
		newIV[2] = byte((retValue >> 16) & 0xFF)
		newIV[3] = byte((retValue >> 24) & 0xFF)
	}

	for i := byte(0); i < 4; i++ {
		copy(c.key[4*i:], newIV[:])
	}
}

func (c *Maple) generateHeader(p []byte) {
	dataLength := len(p[encryptHeaderSize:])

	iiv := int(c.key[3]) & 0xFF
	iiv |= (int(c.key[2]) << 8) & 0xFF00

	iiv ^= int(c.mapleVersion)
	mlength := ((dataLength << 8) & 0xFF00) | (dataLength >> 8)
	xoredIv := iiv ^ mlength

	p[0] = byte((iiv >> 8) & 0xFF)
	p[1] = byte(iiv & 0xFF)
	p[2] = byte((xoredIv >> 8) & 0xFF)
	p[3] = byte(xoredIv & 0xFF)
}

func rol(val byte, count int) byte {
	tmp := uint16(val) << (count % 8)
	return byte((tmp & 0xFF) | (tmp >> 8))
}

func ror(val byte, count int) byte {
	tmp := uint16(val) << 8
	tmp = tmp >> (count % 8)
	return byte((tmp & 0xFF) | (tmp >> 8))
}

func mapleDecrypt(buf []byte) {
	for j := 1; j <= 6; j++ {
		remember := byte(0)
		dataLength := byte(len(buf) & 0xFF)
		var nextRemember byte

		if j%2 == 0 {
			for i := 0; i < len(buf); i++ {
				cur := buf[i]
				cur -= 0x48
				cur = (^cur) & 0xFF
				cur = rol(cur, int(dataLength))
				nextRemember = cur
				cur ^= remember
				remember = nextRemember
				cur -= dataLength
				cur = ror(cur, 3)
				buf[i] = cur
				dataLength--
			}
		} else {
			for i := len(buf) - 1; i >= 0; i-- {
				cur := buf[i]
				cur = rol(cur, 3)
				cur ^= 0x13
				nextRemember = cur
				cur ^= remember
				remember = nextRemember
				cur -= dataLength
				cur = ror(cur, 4)
				buf[i] = cur
				dataLength--
			}
		}
	}
}

func mapleCrypt(buf []byte) {
	for j := 0; j < 6; j++ {
		remember := byte(0)
		dataLength := byte(len(buf) & 0xFF)

		if j%2 == 0 {
			for i := 0; i < len(buf); i++ {
				cur := buf[i]
				cur = rol(cur, 3)
				cur += dataLength
				cur ^= remember
				remember = cur
				cur = ror(cur, int(dataLength))
				cur = (^cur) & 0xFF
				cur += 0x48
				dataLength--
				buf[i] = cur
			}
		} else {
			for i := len(buf) - 1; i >= 0; i-- {
				cur := buf[i]
				cur = rol(cur, 4)
				cur += dataLength
				cur ^= remember
				remember = cur
				cur ^= 0x13
				cur = ror(cur, 3)
				dataLength--
				buf[i] = cur
			}
		}
	}
}

var aeskey = [32]byte{
	0x13, 0x00, 0x00, 0x00,
	0x08, 0x00, 0x00, 0x00,
	0x06, 0x00, 0x00, 0x00,
	0xB4, 0x00, 0x00, 0x00,
	0x1B, 0x00, 0x00, 0x00,
	0x0F, 0x00, 0x00, 0x00,
	0x33, 0x00, 0x00, 0x00,
	0x52, 0x00, 0x00, 0x00}

func (c *Maple) aesCrypt(buf []byte) {
	remaining := len(buf)
	llength := 0x5B0
	start := 0

	block, err := aes.NewCipher(aeskey[:])
	if err != nil {
		panic(err)
	}

	for remaining > 0 {
		myIv := multiplyBytes(c.key[:4], 4, 4)

		if remaining < llength {
			llength = remaining
		}

		for x := start; x < (start + llength); x++ {
			if (x-start)%len(myIv) == 0 {
				block.Encrypt(myIv, myIv)
			}
			buf[x] ^= myIv[(x-start)%len(myIv)]
		}

		start += llength
		remaining -= llength
		llength = 0x5B4
	}
}

func multiplyBytes(in []byte, count, mul int) []byte {
	ret := make([]byte, count*mul)
	for x := 0; x < count*mul; x++ {
		ret[x] = in[x%count]
	}
	return ret
}
