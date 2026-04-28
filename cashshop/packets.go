package cashshop

import (
	"github.com/Hucaru/Valhalla/channel"
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/mpacket"
)

func packetCashShopSet(plr *channel.Player) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSetCashShop)
	channel.AppendCashShopCharacterData(&p, plr)

	// CCashShop::LoadData
	p.WriteByte(1)
	p.WriteString(plr.GetAccountName())

	// sub_71E3E4
	p.WriteInt32(0)
	p.WriteInt16(0)
	p.WriteByte(0)

	// Raw 0x438-byte Cash Shop data block consumed by CCashShop::LoadData.
	// Field meanings inside this block are not fully resolved yet, so preserve
	// the size and ordering with an explicit zeroed placeholder.
	p.WriteBytes(make([]byte, 0x438))

	// CCashShop::DecodeStock / CCashShop::DecodeLimitGoods
	p.WriteInt16(0)
	p.WriteInt16(0)

	// Trailing byte read by CCashShop::CCashShop after LoadData.
	p.WriteByte(0)

	return p
}

func packetCashShopUpdateAmounts(nxCredit, maplePoints, prepaidNX int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSUpdateAmounts)
	p.WriteInt32(nxCredit)
	p.WriteInt32(maplePoints)
	p.WriteInt32(prepaidNX)
	return p
}

func packetCashShopIncreaseInv(invID byte, slots int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteByte(opcode.SendCashShopIncSlotCountDone)
	p.WriteByte(invID)
	p.WriteInt16(slots)
	return p
}

func packetCashShopError(opCode, err byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteByte(opCode)
	p.WriteByte(err)

	return p
}

func packetCashShopShowBoughtItem(charID int32, cashItemSNHash int64, itemID int32, count int16, itemName string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteInt64(cashItemSNHash)
	p.WriteInt32(charID)

	for i := 0; i < 4; i++ {
		p.WriteByte(0x01)
	}

	p.WriteInt32(itemID)

	for i := 0; i < 4; i++ {
		p.WriteByte(0x01)
	}

	p.WriteInt16(count)
	p.WriteString(itemName)
	p.WriteInt64(0)
	for i := 0; i < 4; i++ {
		p.WriteByte(0x01)
	}
	return p
}

func packetCashShopShowBoughtQuestItem(position byte, itemID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteInt32(365)
	p.WriteByte(0)
	p.WriteInt16(1)
	p.WriteByte(position)
	p.WriteByte(0)
	p.WriteInt32(itemID)
	return p
}

func packetCashShopShowCouponRedeemedItem(itemID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteInt16(0x3A)
	p.WriteInt32(0)
	p.WriteInt32(1)
	p.WriteInt16(1)
	p.WriteInt16(0x1A)
	p.WriteInt32(itemID)
	p.WriteInt32(0)
	return p
}

func packetCashShopSendCSItemInventory(slotType byte, it channel.Item) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteByte(0x2F)
	p.WriteInt16(int16(slotType))
	p.WriteByte(slotType)
	p.WriteBytes(it.InventoryBytes())
	return p
}

func packetCashShopWishList(sns []int32, update bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	if update {
		p.WriteByte(opcode.SendCashShopUpdateWishDone)
	} else {
		p.WriteByte(opcode.SendCashShopLoadWishDone)
	}
	for i := 0; i < 10; i++ {
		var v int32
		if i < len(sns) {
			v = sns[i]
		}
		p.WriteInt32(v)
	}
	return p
}

func packetCashShopLoadLocker(storage *CashShopStorage, accountID, characterID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteByte(opcode.SendCashShopLoadLockerDone)

	items := storage.getAllItems()
	p.WriteInt16(int16(len(items)))
	for _, csItem := range items {
		p.WriteInt64(csItem.GetCashID())
		p.WriteInt32(accountID)
		p.WriteInt32(characterID)
		p.WriteInt32(csItem.ID)
		p.WriteInt32(csItem.GetCashSN())
		p.WriteInt16(csItem.GetAmount())
		p.WritePaddedString("", 13)
		p.WriteInt64(csItem.GetExpireTime())
		p.WriteInt64(0)
	}

	p.WriteInt16(0)
	p.WriteInt16(int16(storage.maxSlots))
	return p
}

func packetCashShopMoveLtoSDone(item channel.Item, slot int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteByte(opcode.SendCashShopMoveLtoSDone)
	p.WriteInt16(slot)
	p.WriteBytes(item.CashShopInventoryBody())
	return p
}

func packetCashShopMoveStoLDone(csItem channel.Item, accountID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteByte(opcode.SendCashShopMoveStoLDone)
	p.WriteInt64(csItem.GetCashID())
	p.WriteInt32(accountID)
	p.WriteInt32(0)
	p.WriteInt32(csItem.ID)
	p.WriteInt32(csItem.GetCashSN())
	p.WriteInt16(csItem.GetAmount())
	p.WritePaddedString("", 13)
	p.WriteInt64(csItem.GetExpireTime())
	p.WriteInt64(0)
	return p
}

func packetCashShopBuyDone(csItem channel.Item, accountID, characterID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteByte(opcode.SendCashShopBuyDone)
	p.WriteInt64(csItem.GetCashID())
	p.WriteInt32(accountID)
	p.WriteInt32(characterID)
	p.WriteInt32(csItem.ID)
	p.WriteInt32(csItem.GetCashSN())
	p.WriteInt16(csItem.GetAmount())
	p.WritePaddedString("", 13)
	p.WriteInt64(csItem.GetExpireTime())
	p.WriteInt64(0)
	return p
}

func packetCashShopWrongCoupon() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelCSAction)
	p.WriteByte(0x40)
	p.WriteByte(0x87)
	return p
}
