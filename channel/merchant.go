package channel

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

const (
	merchantStateActive byte = iota + 1
	merchantStateRetrievable
	merchantStateExpired
)

const (
	merchantEmployeeTemplateID int32 = 9030000
	merchantBankerNpcTemplate  int32 = 9030000
	merchantCashUseType        int   = 11

	merchantCheckResultSuccess  byte = 6
	merchantCheckResultOpen     byte = 7
	merchantCheckResultFredrick byte = 8
	merchantCheckResultFailed   byte = 10

	merchantReqPutItem      byte = 30
	merchantReqBuyItem      byte = 31
	merchantReqSilent27     byte = 27
	merchantReqTakeItemBack byte = 35
	merchantReqGoOut        byte = 36
	merchantReqArrangeItem  byte = 37
	merchantReqWithdrawAll  byte = 38
	merchantReqWithdrawCash byte = 40

	merchantRespUpdate         byte = 0x19
	merchantRespBuyResult      byte = 35
	merchantRespSoldItem       byte = 37
	merchantRespMoveItemToInv  byte = 38
	merchantRespArrangeResult  byte = 0x25
	merchantRespWithdrawResult byte = 0x27
	merchantRespWithdrawMoney  byte = 0x29

	storeBankReqOpenPreview byte = 23
	storeBankReqRetrieve    byte = 24
	storeBankReqClose       byte = 25

	storeBankShowDialog byte = 32
	storeBankShowFee    byte = 33
	storeBankShowStatus byte = 34
	storeBankForceClose byte = 35

	storeBankResultSuccess      byte = 27
	storeBankResultTooMuchMeso  byte = 28
	storeBankResultOnlyItem     byte = 29
	storeBankResultNoServiceFee byte = 30
	storeBankResultFullInv      byte = 31
	merchantCreateTTL                = 30 * time.Second
)

type merchantPermitState struct {
	itemID     int32
	slotID     int16
	cashID     int64
	cashSN     int32
	expiresAt  int64
	startedMap int32
	startedCh  byte
}

type merchantRoom struct {
	room

	mu sync.Mutex

	shopID        int64
	fieldInst     *fieldInstance
	ownerAccount  int32
	ownerName     string
	ownerAvatar   []byte
	title         string
	description   string
	permitItemID  int32
	permitCashID  int64
	permitCashSN  int32
	slotCount     byte
	mapID         int32
	pendingMesos  int32
	state         byte
	createdAt     int64
	expiresAt     int64
	closedAt      int64
	lastTouchedAt int64
	npcTemplateID int32
	npcSpawnID    int32
	balloonOpen   byte
	pos           pos
	items         []*shopItem
}

func newMerchantRoom(id int32, shopID int64, inst *fieldInstance, owner *Player, permit Item, title string) *merchantRoom {
	now := time.Now().UnixMilli()
	slotCount := byte(16)
	if meta, err := nx.GetItem(permit.ID); err == nil && meta.SlotMax > 0 {
		slotCount = byte(meta.SlotMax)
	}
	periodDays := merchantPermitDurationDays(permit)
	expiresAt := now + int64(periodDays)*24*int64(time.Hour/time.Millisecond)
	if periodDays <= 0 {
		expiresAt = now + 24*int64(time.Hour/time.Millisecond)
	}

	r := &merchantRoom{
		room: room{
			roomID:        id,
			ownerPlayerID: owner.ID,
			roomType:      constant.MiniRoomTypeEntrustedShop,
			players:       []*Player{},
		},
		shopID:        shopID,
		fieldInst:     inst,
		ownerAccount:  owner.accountID,
		ownerName:     owner.Name,
		ownerAvatar:   append([]byte(nil), owner.avatarLookBytes()...),
		title:         title,
		description:   title,
		permitItemID:  permit.ID,
		permitCashID:  permit.GetCashID(),
		permitCashSN:  permit.GetCashSN(),
		slotCount:     slotCount,
		mapID:         owner.mapID,
		state:         merchantStateActive,
		createdAt:     now,
		expiresAt:     expiresAt,
		lastTouchedAt: now,
		npcTemplateID: merchantEmployeeTemplateID,
		pos:           owner.pos,
		items:         []*shopItem{},
	}
	return r
}

func merchantPermitDurationDays(item Item) int32 {
	if item.cashSN != 0 {
		if c, ok := nx.GetCommodity(item.cashSN); ok && c.Period > 0 {
			return c.Period
		}
	}
	if c, ok := nx.GetCommodityByItemID(item.ID); ok && c.Period > 0 {
		return c.Period
	}
	return 1
}

func merchantExpiryDuration(item Item) time.Duration {
	period := merchantPermitDurationDays(item)
	if period <= 0 {
		period = 1
	}
	return time.Duration(period) * 24 * time.Hour
}

func (r *merchantRoom) closed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state != merchantStateActive && len(r.players) == 0
}

func (r *merchantRoom) addPlayer(plr *Player) bool {
	var kicked []*Player

	r.mu.Lock()

	if r.state != merchantStateActive {
		r.mu.Unlock()
		plr.Send(packetRoomClosed())
		return false
	}
	if plr.ID != r.ownerPlayerID && r.ownerManagingLocked() {
		r.mu.Unlock()
		plr.Send(packetRoomClosed())
		return false
	}
	for _, v := range r.players {
		if v != nil && v.ID == plr.ID {
			r.mu.Unlock()
			return false
		}
	}
	if plr.ID == r.ownerPlayerID {
		kicked = r.collectVisitorsLocked()
	}
	if r.displayCountLocked() >= constant.ShopMaxPlayers {
		r.mu.Unlock()
		plr.Send(packetRoomFull())
		return false
	}

	r.players = append(r.players, plr)
	plr.Send(r.packetShowWindowLocked(plr))
	slot := r.displaySlotLocked(plr.ID)
	if slot > 0 {
		r.sendExcept(packetRoomJoin(r.roomType, slot, plr), plr)
		if r.npcSpawnID != 0 && r.fieldInst != nil {
			r.fieldInst.send(packetEmployeeMiniRoomBalloon(r))
		}
	}
	r.mu.Unlock()

	for _, visitor := range kicked {
		if visitor == nil {
			continue
		}
		visitor.Send(packetRoomLeave(0, constant.MiniRoomClosed))
		visitor.Send(packetPlayerNoChange())
	}
	return true
}

func (r *merchantRoom) removePlayer(plr *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	idx := -1
	slot := byte(0)
	for i, v := range r.players {
		if v != nil && v.ID == plr.ID {
			idx = i
			slot = r.displaySlotLocked(plr.ID)
			break
		}
	}
	if idx < 0 {
		return
	}

	r.players = append(r.players[:idx], r.players[idx+1:]...)
	plr.Send(packetRoomLeave(slot, constant.MiniRoomLeaveReason))
	r.send(packetRoomLeave(slot, constant.MiniRoomLeaveReason))
	if r.npcSpawnID != 0 && r.fieldInst != nil {
		r.fieldInst.send(packetEmployeeMiniRoomBalloon(r))
	}
}

func (r *merchantRoom) displayCountLocked() int {
	count := 1
	for _, v := range r.players {
		if v == nil {
			continue
		}
		if v.ID == r.ownerPlayerID {
			continue
		}
		count++
	}
	return count
}

func (r *merchantRoom) ownerManagingLocked() bool {
	for _, v := range r.players {
		if v != nil && v.ID == r.ownerPlayerID {
			return true
		}
	}
	return false
}

func (r *merchantRoom) ownerManaging() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ownerManagingLocked()
}

func (r *merchantRoom) collectVisitorsLocked() []*Player {
	visitors := make([]*Player, 0, len(r.players))
	kept := make([]*Player, 0, 1)
	for _, v := range r.players {
		if v == nil {
			continue
		}
		if v.ID == r.ownerPlayerID {
			kept = append(kept, v)
			continue
		}
		visitors = append(visitors, v)
	}
	r.players = kept
	return visitors
}

func (r *merchantRoom) displaySlotLocked(playerID int32) byte {
	if playerID == r.ownerPlayerID {
		return 0
	}
	var slot byte = 1
	for _, v := range r.players {
		if v == nil || v.ID == r.ownerPlayerID {
			continue
		}
		if v.ID == playerID {
			return slot
		}
		slot++
	}
	return 0
}

func (r *merchantRoom) packetShowWindowLocked(plr *Player) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelRoom)
	p.WriteByte(constant.RoomPacketShowWindow)
	p.WriteByte(r.roomType)
	p.WriteByte(0x04)
	p.WriteInt16(int16(r.displaySlotLocked(plr.ID)))
	p.WriteInt32(r.permitItemID)
	p.WriteString(r.title)

	slot := byte(1)
	for _, v := range r.players {
		if v == nil || v.ID == r.ownerPlayerID {
			continue
		}
		p.WriteByte(slot)
		v.encodeDisplayBytes(&p)
		p.WriteString(v.Name)
		slot++
	}

	p.WriteByte(constant.RoomPacketEndList)
	p.WriteInt16(0)
	p.WriteString(r.ownerName)
	if plr.ID == r.ownerPlayerID {
		p.WriteInt16(0)
		p.WriteInt16(int16(time.Since(time.UnixMilli(r.createdAt)).Minutes()))
		p.WriteByte(1)
		p.WriteByte(0)
		p.WriteInt32(r.pendingMesos)
	}
	p.WriteString(r.description)
	p.WriteByte(r.slotCount)
	if plr.ID == r.ownerPlayerID {
		p.WriteInt32(r.pendingMesos)
	} else {
		p.WriteInt32(plr.mesos)
	}
	encodeMerchantShopItems(&p, r.items)
	if len(r.items) == 0 {
		p.WriteByte(0)
	}
	return p
}

func (r *merchantRoom) chatMsg(plr *Player, msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.send(packetRoomChat(plr.Name, msg, r.displaySlotLocked(plr.ID)))
}

func (r *merchantRoom) markTouchedLocked() {
	r.lastTouchedAt = time.Now().UnixMilli()
	if _, err := common.DB.Exec("UPDATE merchant_shops SET lastTouchedAt=? WHERE id=?", r.lastTouchedAt, r.shopID); err != nil {
		log.Printf("merchant: update touch failed shop=%d err=%v", r.shopID, err)
	}
}

func (r *merchantRoom) addItem(owner *Player, item Item, bundles, bundleAmount int16, price int32) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != merchantStateActive || owner == nil || owner.ID != r.ownerPlayerID {
		return false
	}
	if len(r.items) >= int(r.slotCount) || bundles <= 0 || bundleAmount <= 0 || price < 0 {
		return false
	}
	cur, err := owner.getItem(item.invID, item.slotID)
	if err != nil || cur.dbID != item.dbID || cur.ID != item.ID {
		return false
	}
	if cur.isRechargeable() {
		bundles = 1
		bundleAmount = cur.amount
	}
	need := bundles * bundleAmount
	if need <= 0 || need > cur.amount {
		return false
	}

	removed, err := owner.takeItem(cur.ID, cur.slotID, need, cur.invID)
	if err != nil {
		return false
	}
	listItem := removed
	listItem.amount = bundleAmount
	listItem.dbID = 0
	listItem.slotID = 0

	si := &shopItem{item: listItem, price: price, bundles: bundles, bundleAmount: bundleAmount, reserved: need}
	if err := merchantInsertItem(r.shopID, len(r.items), si); err != nil {
		if _, giveErr := owner.GiveItem(removed); giveErr != nil {
			log.Printf("merchant: rollback add item failed shop=%d char=%d err=%v", r.shopID, owner.ID, giveErr)
		}
		log.Printf("merchant: insert item failed shop=%d err=%v", r.shopID, err)
		return false
	}

	r.items = append(r.items, si)
	r.markTouchedLocked()
	return true
}

func (r *merchantRoom) buyItem(buyer *Player, slot byte, quantity int16) (byte, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != merchantStateActive || buyer == nil || buyer.ID == r.ownerPlayerID {
		return constant.PlayerShopNotEnoughInStock, false
	}
	if int(slot) < 0 || int(slot) >= len(r.items) {
		return constant.PlayerShopNotEnoughInStock, false
	}
	si := r.items[slot]
	if si == nil || quantity <= 0 || si.bundles < quantity {
		return constant.PlayerShopNotEnoughInStock, false
	}
	totalCost := int64(si.price) * int64(quantity)
	if totalCost < 0 || totalCost > math.MaxInt32 {
		return constant.PlayerShopPriceTooHighForTrade, false
	}
	if int64(buyer.mesos) < totalCost {
		return constant.PlayerShopBuyerNotEnoughMoney, false
	}
	realAmount := quantity * si.bundleAmount
	if realAmount <= 0 {
		return constant.PlayerShopNotEnoughInStock, false
	}
	give := si.item
	if give.isRechargeable() {
		give.amount = si.bundleAmount
	} else {
		give.amount = realAmount
	}
	if _, err := buyer.GiveItem(give); err != nil {
		return constant.PlayerShopInventoryFull, false
	}

	buyer.takeMesos(int32(totalCost))
	if int64(r.pendingMesos)+totalCost > math.MaxInt32 {
		buyer.giveMesos(int32(totalCost))
		_, _ = buyer.TakeItemSilent(give.ID, give.slotID, give.amount, give.invID)
		return constant.PlayerShopPriceTooHighForTrade, false
	}
	r.pendingMesos += int32(totalCost)
	si.bundles -= quantity
	if si.bundles <= 0 {
		r.items = append(r.items[:slot], r.items[slot+1:]...)
	}
	if err := merchantRewriteItems(r.shopID, r.items); err != nil {
		log.Printf("merchant: rewrite after sale failed shop=%d err=%v", r.shopID, err)
	}
	if _, err := common.DB.Exec("UPDATE merchant_shops SET pendingMesos=?, lastTouchedAt=? WHERE id=?", r.pendingMesos, time.Now().UnixMilli(), r.shopID); err != nil {
		log.Printf("merchant: update mesos failed shop=%d err=%v", r.shopID, err)
	}
	r.lastTouchedAt = time.Now().UnixMilli()
	return 0, true
}

func (r *merchantRoom) removeItem(owner *Player, slot byte) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if owner == nil || owner.ID != r.ownerPlayerID || int(slot) < 0 || int(slot) >= len(r.items) {
		return false
	}
	si := r.items[slot]
	if si == nil {
		return false
	}
	back := si.item
	if back.isRechargeable() {
		back.amount = si.bundleAmount
	} else {
		back.amount = si.bundles * si.bundleAmount
	}
	if _, err := owner.GiveItem(back); err != nil {
		return false
	}
	r.items = append(r.items[:slot], r.items[slot+1:]...)
	if err := merchantRewriteItems(r.shopID, r.items); err != nil {
		log.Printf("merchant: rewrite after remove failed shop=%d err=%v", r.shopID, err)
	}
	r.markTouchedLocked()
	return true
}

func (r *merchantRoom) withdrawMoney(owner *Player) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if owner == nil || owner.ID != r.ownerPlayerID {
		return false
	}
	if r.pendingMesos <= 0 {
		return true
	}
	if int64(owner.mesos)+int64(r.pendingMesos) > math.MaxInt32 {
		return false
	}
	owner.giveMesos(r.pendingMesos)
	r.pendingMesos = 0
	if _, err := common.DB.Exec("UPDATE merchant_shops SET pendingMesos=?, lastTouchedAt=? WHERE id=?", 0, time.Now().UnixMilli(), r.shopID); err != nil {
		log.Printf("merchant: withdraw mesos update failed shop=%d err=%v", r.shopID, err)
	}
	r.lastTouchedAt = time.Now().UnixMilli()
	return true
}

func (r *merchantRoom) withdrawMoneyNoPacket(owner *Player) bool {
	return r.withdrawMoney(owner)
}

func (r *merchantRoom) closeAndBank(server *Server, owner *Player, reason byte) bool {
	r.mu.Lock()
	if r.state != merchantStateActive || owner == nil || owner.ID != r.ownerPlayerID {
		r.mu.Unlock()
		return false
	}
	r.state = merchantStateRetrievable
	r.closedAt = time.Now().UnixMilli()
	players := append([]*Player(nil), r.players...)
	r.players = nil
	shopID := r.shopID
	mapID := r.mapID
	wasVisible := r.npcSpawnID != 0
	r.npcSpawnID = 0
	r.mu.Unlock()

	if _, err := common.DB.Exec("UPDATE merchant_shops SET state=?, npcSpawnID=0, closedAt=?, lastTouchedAt=? WHERE id=?", merchantStateRetrievable, time.Now().UnixMilli(), time.Now().UnixMilli(), shopID); err != nil {
		log.Printf("merchant: close update failed shop=%d err=%v", shopID, err)
	}
	for _, plr := range players {
		if plr != nil {
			plr.Send(packetRoomLeave(0, reason))
			plr.Send(packetPlayerNoChange())
		}
	}
	if field, ok := server.fields[mapID]; ok {
		if inst, err := field.getInstance(r.fieldInst.id); err == nil {
			if wasVisible {
				inst.send(packetEmployeeLeaveField(r.ownerPlayerID))
			}
			_ = inst.roomPool.removeRoom(r.roomID)
		}
	}
	server.unregisterMerchant(r)
	return true
}

func (r *merchantRoom) autoCloseSoldOut(server *Server, reason byte) bool {
	r.mu.Lock()
	if r.state != merchantStateActive || len(r.items) > 0 {
		r.mu.Unlock()
		return false
	}
	r.state = merchantStateRetrievable
	r.closedAt = time.Now().UnixMilli()
	players := append([]*Player(nil), r.players...)
	r.players = nil
	shopID := r.shopID
	mapID := r.mapID
	wasVisible := r.npcSpawnID != 0
	r.npcSpawnID = 0
	r.mu.Unlock()

	if _, err := common.DB.Exec("UPDATE merchant_shops SET state=?, npcSpawnID=0, closedAt=?, lastTouchedAt=? WHERE id=?", merchantStateRetrievable, time.Now().UnixMilli(), time.Now().UnixMilli(), shopID); err != nil {
		log.Printf("merchant: sold-out close update failed shop=%d err=%v", shopID, err)
	}
	for _, plr := range players {
		if plr != nil {
			plr.Send(packetRoomLeave(0, reason))
			plr.Send(packetPlayerNoChange())
		}
	}
	if field, ok := server.fields[mapID]; ok {
		if inst, err := field.getInstance(r.fieldInst.id); err == nil {
			if wasVisible {
				inst.send(packetEmployeeLeaveField(r.ownerPlayerID))
			}
			_ = inst.roomPool.removeRoom(r.roomID)
		}
	}
	server.unregisterMerchant(r)
	return true
}

func merchantInsertItem(shopID int64, order int, si *shopItem) error {
	_, err := common.DB.Exec(`INSERT INTO merchant_items(
		shopID, displayOrder, inventoryID, itemID, amount, flag, upgradeSlots, level,
		str, dex, intt, luk, hp, mp, watk, matk, wdef, mdef, accuracy, avoid, hands, speed, jump,
		expireTime, creatorName, cashID, cashSN, ringID, bundles, bundleAmount, price
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		shopID, order, si.item.invID, si.item.ID, si.item.amount, si.item.flag, si.item.upgradeSlots, si.item.scrollLevel,
		si.item.str, si.item.dex, si.item.intt, si.item.luk, si.item.hp, si.item.mp, si.item.watk, si.item.matk, si.item.wdef, si.item.mdef,
		si.item.accuracy, si.item.avoid, si.item.hands, si.item.speed, si.item.jump,
		si.item.expireTime, si.item.creatorName,
		sql.NullInt64{Int64: si.item.cashID, Valid: si.item.cashID != 0},
		sql.NullInt32{Int32: si.item.cashSN, Valid: si.item.cashSN != 0},
		sql.NullInt32{Int32: si.item.ringID, Valid: si.item.ringID > 0},
		si.bundles, si.bundleAmount, si.price,
	)
	return err
}

func merchantRewriteItems(shopID int64, items []*shopItem) error {
	tx, err := common.DB.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.Exec("DELETE FROM merchant_items WHERE shopID=?", shopID); err != nil {
		return err
	}
	for i, si := range items {
		if si == nil {
			continue
		}
		if _, err := tx.Exec(`INSERT INTO merchant_items(
			shopID, displayOrder, inventoryID, itemID, amount, flag, upgradeSlots, level,
			str, dex, intt, luk, hp, mp, watk, matk, wdef, mdef, accuracy, avoid, hands, speed, jump,
			expireTime, creatorName, cashID, cashSN, ringID, bundles, bundleAmount, price
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			shopID, i, si.item.invID, si.item.ID, si.item.amount, si.item.flag, si.item.upgradeSlots, si.item.scrollLevel,
			si.item.str, si.item.dex, si.item.intt, si.item.luk, si.item.hp, si.item.mp, si.item.watk, si.item.matk, si.item.wdef, si.item.mdef,
			si.item.accuracy, si.item.avoid, si.item.hands, si.item.speed, si.item.jump,
			si.item.expireTime, si.item.creatorName,
			sql.NullInt64{Int64: si.item.cashID, Valid: si.item.cashID != 0},
			sql.NullInt32{Int32: si.item.cashSN, Valid: si.item.cashSN != 0},
			sql.NullInt32{Int32: si.item.ringID, Valid: si.item.ringID > 0},
			si.bundles, si.bundleAmount, si.price,
		); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (server *Server) registerMerchant(r *merchantRoom) {
	if server.merchantShops == nil {
		server.merchantShops = make(map[int64]*merchantRoom)
	}
	if server.merchantByChar == nil {
		server.merchantByChar = make(map[int32]*merchantRoom)
	}
	server.merchantShops[r.shopID] = r
	server.merchantByChar[r.ownerPlayerID] = r
}

func (server *Server) unregisterMerchant(r *merchantRoom) {
	if server == nil || r == nil {
		return
	}
	delete(server.merchantShops, r.shopID)
	delete(server.merchantByChar, r.ownerPlayerID)
}

func (server *Server) hasActiveMerchant(charID int32) bool {
	if r, ok := server.merchantByChar[charID]; ok && r != nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		return r.state == merchantStateActive
	}
	return false
}

func (server *Server) createMerchant(plr *Player, description string) (*merchantRoom, error) {
	if plr == nil || plr.inst == nil {
		return nil, fmt.Errorf("player not in map")
	}
	if server.hasActiveMerchant(plr.ID) {
		return nil, fmt.Errorf("merchant already active")
	}
	if plr.pendingMerchant.expiresAt <= 0 || time.Now().UnixMilli() > plr.pendingMerchant.expiresAt {
		return nil, fmt.Errorf("merchant permit not primed")
	}
	if plr.pendingMerchant.startedMap != plr.mapID || plr.pendingMerchant.startedCh != server.id {
		return nil, fmt.Errorf("merchant permit map mismatch")
	}
	if plr.inst.fieldID < 0 || plr.inst.fieldID != plr.mapID {
		return nil, fmt.Errorf("invalid merchant field")
	}
	if plr.mapID < 910000001 || plr.mapID > 910000022 {
		return nil, fmt.Errorf("map does not allow merchants")
	}
	permit, err := plr.getItem(constant.InventoryCash, plr.pendingMerchant.slotID)
	if err != nil {
		return nil, err
	}
	if permit.ID != plr.pendingMerchant.itemID || permit.GetCashID() != plr.pendingMerchant.cashID {
		return nil, fmt.Errorf("merchant permit changed")
	}
	removed := permit
	plr.pendingMerchant = merchantPermitState{}

	res, err := common.DB.Exec(`INSERT INTO merchant_shops(
		characterID, accountID, worldID, channelID, mapID, roomID, npcSpawnID, npcTemplateID,
		ownerName, title, description, ownerAvatar, permitItemID, permitCashID, permitCashSN,
		slotCount, x, y, foothold, pendingMesos, state, createdAt, expiresAt, closedAt, lastTouchedAt
	) VALUES (?, ?, ?, ?, ?, 0, 0, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?, ?, 0, ?)`,
		plr.ID, plr.accountID, plr.worldID, server.id, plr.mapID,
		merchantEmployeeTemplateID, plr.Name, description, description, plr.avatarLookBytes(),
		removed.ID,
		sql.NullInt64{Int64: removed.cashID, Valid: removed.cashID != 0},
		sql.NullInt32{Int32: removed.cashSN, Valid: removed.cashSN != 0},
		byte(getItemSlotMax(removed.ID)), plr.pos.x, plr.pos.y, plr.pos.foothold,
		merchantStateActive, time.Now().UnixMilli(), time.Now().Add(merchantExpiryDuration(removed)).UnixMilli(), time.Now().UnixMilli(),
	)
	if err != nil {
		return nil, err
	}
	shopID, _ := res.LastInsertId()
	r := newMerchantRoom(int32(shopID), shopID, plr.inst, plr, removed, description)
	if _, err := common.DB.Exec("UPDATE merchant_shops SET roomID=?, slotCount=?, expiresAt=? WHERE id=?", r.roomID, r.slotCount, r.expiresAt, r.shopID); err != nil {
		log.Printf("merchant: finalize create failed shop=%d err=%v", r.shopID, err)
	}
	if err := plr.inst.roomPool.addRoom(r); err != nil {
		_, _ = common.DB.Exec("DELETE FROM merchant_shops WHERE id=?", r.shopID)
		return nil, err
	}
	r.addPlayer(plr)
	if _, err := common.DB.Exec("UPDATE merchant_shops SET roomID=? WHERE id=?", r.roomID, r.shopID); err != nil {
		log.Printf("merchant: update room id failed shop=%d err=%v", r.shopID, err)
	}
	server.registerMerchant(r)
	return r, nil
}

func (inst *fieldInstance) addMerchantNPC(r *merchantRoom) (int32, error) {
	spawnID, err := inst.lifePool.nextNpcID()
	if err != nil {
		return 0, err
	}
	data := nx.Life{ID: r.npcTemplateID, Type: "n", X: r.pos.x, Y: r.pos.y, FaceLeft: false, Foothold: r.pos.foothold, Rx0: r.pos.x - 40, Rx1: r.pos.x + 40}
	val := createNpcFromData(spawnID, data)
	inst.lifePool.npcs[spawnID] = &val
	inst.send(packetNpcShow(&val))
	return spawnID, nil
}

func (inst *fieldInstance) removeMerchantNPC(spawnID int32) {
	if inst == nil || spawnID == 0 {
		return
	}
	if n, ok := inst.lifePool.npcs[spawnID]; ok && n != nil {
		n.removeController()
	}
	delete(inst.lifePool.npcs, spawnID)
	inst.send(packetNpcRemove(spawnID))
}

func (server *Server) loadMerchants() {
	if server == nil {
		return
	}
	if err := server.closeLingeringChannelMerchants(); err != nil {
		log.Printf("merchant: startup close failed: %v", err)
	}
	rows, err := common.DB.Query(`SELECT id, characterID, accountID, worldID, channelID, mapID, roomID, npcSpawnID, npcTemplateID,
		ownerName, title, description, ownerAvatar, permitItemID, permitCashID, permitCashSN, slotCount,
		x, y, foothold, pendingMesos, state, createdAt, expiresAt, closedAt, lastTouchedAt
		FROM merchant_shops WHERE state=? AND channelID=? AND npcSpawnID<>0`, merchantStateActive, server.id)
	if err != nil {
		log.Printf("merchant: load failed: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &merchantRoom{room: room{roomType: constant.MiniRoomTypeEntrustedShop, players: []*Player{}}, items: []*shopItem{}}
		var permitCashID sql.NullInt64
		var permitCashSN sql.NullInt32
		var worldID byte
		var channelID byte
		var mapID int32
		if err := rows.Scan(&r.shopID, &r.ownerPlayerID, &r.ownerAccount, &worldID, &channelID, &mapID, &r.roomID, &r.npcSpawnID, &r.npcTemplateID,
			&r.ownerName, &r.title, &r.description, &r.ownerAvatar, &r.permitItemID, &permitCashID, &permitCashSN, &r.slotCount,
			&r.pos.x, &r.pos.y, &r.pos.foothold, &r.pendingMesos, &r.state, &r.createdAt, &r.expiresAt, &r.closedAt, &r.lastTouchedAt); err != nil {
			log.Printf("merchant: load scan failed: %v", err)
			continue
		}
		r.mapID = mapID
		if permitCashID.Valid {
			r.permitCashID = permitCashID.Int64
		}
		if permitCashSN.Valid {
			r.permitCashSN = permitCashSN.Int32
		}
		fieldRow, ok := server.fields[mapID]
		if !ok {
			continue
		}
		inst, err := fieldRow.getInstance(0)
		if err != nil {
			continue
		}
		r.fieldInst = inst
		items, err := merchantLoadItems(r.shopID)
		if err != nil {
			log.Printf("merchant: load items failed shop=%d err=%v", r.shopID, err)
			continue
		}
		r.items = items
		if err := inst.roomPool.addRoom(r); err != nil {
			log.Printf("merchant: room add failed shop=%d err=%v", r.shopID, err)
			continue
		}
		inst.send(packetEmployeeEnterField(r))
		server.registerMerchant(r)
	}
}

func (server *Server) closeLingeringChannelMerchants() error {
	if server == nil || common.DB == nil {
		return nil
	}

	now := time.Now().UnixMilli()
	tx, err := common.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`UPDATE merchant_shops
		SET state=?, npcSpawnID=0, closedAt=?, lastTouchedAt=?
		WHERE state=? AND channelID=?`, merchantStateRetrievable, now, now, merchantStateActive, server.id); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func merchantLoadItems(shopID int64) ([]*shopItem, error) {
	rows, err := common.DB.Query(`SELECT inventoryID, itemID, amount, flag, upgradeSlots, level, str, dex, intt, luk,
		hp, mp, watk, matk, wdef, mdef, accuracy, avoid, hands, speed, jump, expireTime, creatorName,
		cashID, cashSN, ringID, bundles, bundleAmount, price
		FROM merchant_items WHERE shopID=? ORDER BY displayOrder ASC`, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*shopItem{}
	for rows.Next() {
		var it Item
		var si shopItem
		var cashID sql.NullInt64
		var cashSN sql.NullInt32
		var ringID sql.NullInt32
		if err := rows.Scan(&it.invID, &it.ID, &it.amount, &it.flag, &it.upgradeSlots, &it.scrollLevel, &it.str, &it.dex, &it.intt, &it.luk,
			&it.hp, &it.mp, &it.watk, &it.matk, &it.wdef, &it.mdef, &it.accuracy, &it.avoid, &it.hands, &it.speed, &it.jump, &it.expireTime, &it.creatorName,
			&cashID, &cashSN, &ringID, &si.bundles, &si.bundleAmount, &si.price); err != nil {
			return nil, err
		}
		if cashID.Valid {
			it.cashID = cashID.Int64
			it.cash = true
		}
		if cashSN.Valid {
			it.cashSN = cashSN.Int32
		}
		if ringID.Valid {
			it.ringID = ringID.Int32
		}
		si.item = it
		si.reserved = si.bundles * si.bundleAmount
		items = append(items, &si)
	}
	return items, nil
}

func packetMerchantWithdrawAllResult(code byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelRoom)
	p.WriteByte(merchantRespWithdrawResult)
	p.WriteByte(code)
	return p
}

func packetMerchantArrangeResult(pendingMesos int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelRoom)
	p.WriteByte(merchantRespArrangeResult)
	p.WriteInt32(pendingMesos)
	return p
}

func packetMerchantWithdrawMoney() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelRoom)
	p.WriteByte(merchantRespWithdrawMoney)
	return p
}

func packetMerchantItemListUpdate(mesos int32, items []*shopItem) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelRoom)
	p.WriteByte(constant.RoomShopRefresh)
	p.WriteInt32(mesos)
	encodeMerchantShopItems(&p, items)
	return p
}

func encodeMerchantShopItems(p *mpacket.Packet, items []*shopItem) {
	p.WriteByte(byte(len(items)))
	for _, shopItem := range items {
		p.WriteInt16(shopItem.bundles)
		p.WriteInt16(shopItem.bundleAmount)
		p.WriteInt32(shopItem.price)
		p.Append(shopItem.item.StorageBytes())
	}
}

func packetEmployeeEnterField(r *merchantRoom) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelEmployeeEnterField)
	p.WriteInt32(r.ownerPlayerID)
	p.WriteInt32(r.permitItemID)
	p.WriteInt16(r.pos.x)
	p.WriteInt16(r.pos.y)
	p.WriteInt16(r.pos.foothold)
	p.WriteString(r.ownerName)
	p.WriteByte(r.roomType)
	p.WriteInt32(r.roomID)
	p.WriteString(r.title)
	p.WriteByte(byte(r.permitItemID % 100))
	p.WriteByte(byte(len(r.players)))
	p.WriteByte(byte(constant.ShopMaxPlayers))
	return p
}

func packetEmployeeLeaveField(ownerID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelEmployeeLeaveField)
	p.WriteInt32(ownerID)
	return p
}

func packetEmployeeMiniRoomBalloon(r *merchantRoom) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelEmployeeMiniRoomBalloon)
	p.WriteInt32(r.ownerPlayerID)
	p.WriteByte(r.roomType)
	p.WriteInt32(r.roomID)
	p.WriteString(r.title)
	p.WriteByte(byte(r.permitItemID % 100))
	p.WriteByte(byte(len(r.players)))
	p.WriteByte(byte(constant.ShopMaxPlayers))
	return p
}

func (r *merchantRoom) balloonOpenState() byte {
	if r.ownerManaging() {
		return 0
	}
	return 1
}

func packetEntrustedShopResult(code byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelEntrustedShopCheckResult)
	p.WriteByte(code)
	return p
}

func packetEntrustedShopOpenElsewhere(mapID int32, channelID byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelEntrustedShopCheckResult)
	p.WriteByte(merchantCheckResultOpen)
	p.WriteInt32(mapID)
	p.WriteByte(channelID)
	return p
}

func packetStoreBankShow(npcID int32, shopID int64, mesos int32, items []*shopItem) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStoreBank)
	p.WriteByte(storeBankShowDialog)
	p.WriteInt32(npcID)
	p.WriteByte(5)
	p.WriteInt16(0x007E)
	p.WriteInt32(mesos)
	for inv := byte(1); inv <= 5; inv++ {
		section := make([]Item, 0, len(items))
		for _, si := range items {
			if si == nil || si.item.invID != inv {
				continue
			}
			item := si.item
			if item.isRechargeable() {
				item.amount = si.bundleAmount
			} else {
				item.amount = si.bundles * si.bundleAmount
			}
			section = append(section, item)
		}
		p.WriteByte(byte(len(section)))
		for i := range section {
			p.WriteBytes(section[i].StorageBytes())
		}
	}
	_ = shopID
	return p
}

func packetStoreBankFee(days, fee int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStoreBank)
	p.WriteByte(storeBankShowFee)
	p.WriteInt32(days)
	p.WriteInt32(fee)
	return p
}

func packetStoreBankStatus(mapID int32, channelID byte, empty bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStoreBank)
	p.WriteByte(storeBankShowStatus)
	p.WriteInt32(0)
	if empty {
		p.WriteInt32(999999999)
		p.WriteByte(0xFE)
	} else {
		p.WriteInt32(mapID)
		p.WriteByte(channelID)
	}
	return p
}

func packetStoreBankForceClose() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStoreBank)
	p.WriteByte(storeBankForceClose)
	return p
}

func packetStoreBankResult(code byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStoreBankResult)
	p.WriteByte(code)
	return p
}

func (server *Server) primeMerchantPermit(plr *Player, item Item) error {
	if plr == nil || plr.inst == nil {
		return fmt.Errorf("player not in map")
	}
	plr.pendingMerchant = merchantPermitState{
		itemID:     item.ID,
		slotID:     item.slotID,
		cashID:     item.cashID,
		cashSN:     item.cashSN,
		expiresAt:  time.Now().Add(merchantCreateTTL).UnixMilli(),
		startedMap: plr.mapID,
		startedCh:  server.id,
	}
	return nil
}

func (server *Server) tryMerchantBanker(plr *Player) bool {
	if plr == nil {
		return false
	}
	if active, ok := server.merchantByChar[plr.ID]; ok && active != nil {
		plr.Send(packetStoreBankStatus(active.mapID, plr.ChannelID, false))
		return true
	}

	shopID, mesos, _, items, err := merchantLoadRetrievalRecord(plr.ID)
	if err != nil {
		log.Printf("merchant: load retrieval failed char=%d err=%v", plr.ID, err)
		plr.Send(packetStoreBankStatus(0, 0, true))
		return true
	}
	if shopID == 0 {
		plr.Send(packetStoreBankStatus(0, 0, true))
		return true
	}
	plr.storeBankShopID = shopID
	plr.storeBankNpcID = merchantBankerNpcTemplate
	plr.storeBankOpen = true
	plr.Send(packetStoreBankShow(plr.storeBankNpcID, shopID, mesos, items))
	return true
}

func merchantLoadRetrieval(charID int32) (int64, int32, []*shopItem, error) {
	var shopID int64
	var mesos int32
	err := common.DB.QueryRow("SELECT id, pendingMesos FROM merchant_shops WHERE characterID=? AND state<>? ORDER BY closedAt DESC, id DESC LIMIT 1", charID, merchantStateActive).Scan(&shopID, &mesos)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil, nil
		}
		return 0, 0, nil, err
	}
	items, err := merchantLoadItems(shopID)
	if err != nil {
		return 0, 0, nil, err
	}
	return shopID, mesos, items, nil
}

func merchantLoadRetrievalRecord(charID int32) (int64, int32, int64, []*shopItem, error) {
	var shopID int64
	var mesos int32
	var closedAt int64
	err := common.DB.QueryRow("SELECT id, pendingMesos, closedAt FROM merchant_shops WHERE characterID=? AND state<>? ORDER BY closedAt DESC, id DESC LIMIT 1", charID, merchantStateActive).Scan(&shopID, &mesos, &closedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, 0, nil, nil
		}
		return 0, 0, 0, nil, err
	}
	items, err := merchantLoadItems(shopID)
	if err != nil {
		return 0, 0, 0, nil, err
	}
	return shopID, mesos, closedAt, items, nil
}

func merchantStoreBankFee(closedAt int64, totalMesos int32) (int32, int32) {
	if closedAt <= 0 || totalMesos <= 0 {
		return 0, 0
	}
	days := int32(time.Since(time.UnixMilli(closedAt)) / (24 * time.Hour))
	if days < 0 {
		days = 0
	}
	fee := days * (totalMesos / 100)
	if fee > totalMesos {
		fee = totalMesos
	}
	if fee < 0 {
		fee = 0
	}
	return days, fee
}

func canReceiveMerchantItems(plr *Player, items []*shopItem) bool {
	if plr == nil {
		return false
	}
	equip := append([]Item(nil), plr.equip...)
	use := append([]Item(nil), plr.use...)
	setup := append([]Item(nil), plr.setUp...)
	etc := append([]Item(nil), plr.etc...)
	cash := append([]Item(nil), plr.cash...)
	for _, si := range items {
		if si == nil {
			continue
		}
		item := si.item
		if item.isRechargeable() {
			item.amount = si.bundleAmount
		} else {
			item.amount = si.bundles * si.bundleAmount
		}
		switch item.invID {
		case constant.InventoryEquip:
			if byte(len(equip)) >= plr.equipSlotSize {
				return false
			}
			equip = append(equip, item)
		case constant.InventoryUse:
			if !simulateReceiveStackable(&use, plr.useSlotSize, item) {
				return false
			}
		case constant.InventorySetup:
			if byte(len(setup)) >= plr.setupSlotSize {
				return false
			}
			setup = append(setup, item)
		case constant.InventoryEtc:
			if !simulateReceiveStackable(&etc, plr.etcSlotSize, item) {
				return false
			}
		case constant.InventoryCash:
			if byte(len(cash)) >= plr.cashSlotSize {
				return false
			}
			cash = append(cash, item)
		default:
			return false
		}
	}
	return true
}

func simulateReceiveStackable(inv *[]Item, slotSize byte, item Item) bool {
	if item.isRechargeable() {
		if byte(len(*inv)) >= slotSize {
			return false
		}
		*inv = append(*inv, item)
		return true
	}
	remaining := item.amount
	slotMax := getItemSlotMax(item.ID)
	for i := range *inv {
		if remaining == 0 {
			return true
		}
		if (*inv)[i].ID != item.ID || (*inv)[i].amount >= slotMax {
			continue
		}
		canAdd := slotMax - (*inv)[i].amount
		if canAdd > remaining {
			canAdd = remaining
		}
		(*inv)[i].amount += canAdd
		remaining -= canAdd
	}
	for remaining > 0 {
		if byte(len(*inv)) >= slotSize {
			return false
		}
		add := remaining
		if add > slotMax {
			add = slotMax
		}
		clone := item
		clone.amount = add
		*inv = append(*inv, clone)
		remaining -= add
	}
	return true
}

func scheduleMerchantExpiry(server *Server) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		if server == nil || server.dispatch == nil {
			continue
		}
		select {
		case server.dispatch <- func() { server.expireMerchants() }:
		default:
		}
	}
}

func (inst *fieldInstance) updateMerchantBalloons() {
	if inst == nil || inst.server == nil {
		return
	}
	for _, shop := range inst.server.merchantShops {
		if shop == nil || shop.mapID != inst.fieldID || shop.fieldInst == nil || shop.fieldInst.id != inst.id || shop.npcSpawnID == 0 {
			continue
		}
		openState := shop.balloonOpenState()
		shop.mu.Lock()
		changed := shop.balloonOpen != openState
		shop.balloonOpen = openState
		shop.mu.Unlock()
		if changed {
			inst.send(packetEmployeeMiniRoomBalloon(shop))
		}
	}
}

func (server *Server) expireMerchants() {
	if server == nil {
		return
	}
	now := time.Now().UnixMilli()
	for _, shop := range server.merchantShops {
		shop.mu.Lock()
		expired := shop.state == merchantStateActive && shop.expiresAt > 0 && shop.expiresAt <= now
		shop.mu.Unlock()
		if !expired {
			continue
		}
		shop.mu.Lock()
		shop.state = merchantStateExpired
		shop.closedAt = now
		players := append([]*Player(nil), shop.players...)
		shop.players = nil
		wasVisible := shop.npcSpawnID != 0
		shop.npcSpawnID = 0
		shop.mu.Unlock()
		if _, err := common.DB.Exec("UPDATE merchant_shops SET state=?, npcSpawnID=0, closedAt=?, lastTouchedAt=? WHERE id=?", merchantStateExpired, now, now, shop.shopID); err != nil {
			log.Printf("merchant: expire update failed shop=%d err=%v", shop.shopID, err)
		}
		for _, plr := range players {
			if plr != nil {
				plr.Send(packetRoomLeave(0, constant.MiniRoomClosed))
				plr.Send(packetPlayerNoChange())
			}
		}
		if shop.fieldInst != nil {
			if wasVisible {
				shop.fieldInst.send(packetEmployeeLeaveField(shop.ownerPlayerID))
			}
			_ = shop.fieldInst.roomPool.removeRoom(shop.roomID)
		}
		server.unregisterMerchant(shop)
	}
}
