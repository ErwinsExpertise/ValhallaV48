package opcode

const (
	RecvCashShopBuyItem        byte = 0x02
	RecvCashShopGiftItem       byte = 0x03
	RecvCashShopUpdateWishlist byte = 0x04
	RecvCashShopIncreaseSlots  byte = 0x05
	RecvCashShopMoveLtoS       byte = 0x0A
	RecvCashShopMoveStoL       byte = 0x0B
	RecvCashShopBuyCoupleRing  byte = 0x18
	RecvCashShopBuyPackage     byte = 0x19
	RecvCashShopGiftPackage    byte = 0x1A
	RecvCashShopBuyQuestItem   byte = 0x1B

	// SendChannelCSAction subcodes from CCashShop::OnCashItemResult.
	SendCashShopLoadLockerDone      byte = 42
	SendCashShopLoadLockerFailed    byte = 43
	SendCashShopLoadWishDone        byte = 44
	SendCashShopLoadWishFailed      byte = 45
	SendCashShopUpdateWishDone      byte = 50
	SendCashShopUpdateWishFailed    byte = 51
	SendCashShopBuyDone             byte = 52
	SendCashShopBuyFailed           byte = 53
	SendCashShopUseCouponDone       byte = 54
	SendCashShopUseGiftCouponDone   byte = 56
	SendCashShopUseCouponFailed     byte = 57
	SendCashShopGiftDone            byte = 59
	SendCashShopGiftFailed          byte = 60
	SendCashShopIncSlotCountDone    byte = 61
	SendCashShopIncSlotCountFailed  byte = 62
	SendCashShopIncTrunkCountDone   byte = 63
	SendCashShopIncTrunkCountFailed byte = 64
	SendCashShopMoveLtoSDone        byte = 65
	SendCashShopMoveLtoSFailed      byte = 66
	SendCashShopMoveStoLDone        byte = 67
	SendCashShopMoveStoLFailed      byte = 68
)
