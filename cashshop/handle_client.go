package cashshop

import (
	"fmt"
	"log"
	"time"

	"github.com/Hucaru/Valhalla/channel"
	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/internal"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

func (server *Server) HandleClientPacket(conn mnet.Client, reader mpacket.Reader) {
	op := reader.ReadInt16()

	switch op {
	case opcode.RecvPing:
	case opcode.RecvClientMigrate:
		server.handlePlayerConnect(conn, reader)
	case opcode.RecvCashShopCashRequest:
		server.handleCashRequest(conn, reader)
	case opcode.RecvCashShopOperation:
		server.handleCashShopOperation(conn, reader)
	case opcode.RecvChannelUserPortal:
		server.leaveCashShopToChannel(conn, reader)

	default:
		log.Println("[CASHSHOP] UNKNOWN CLIENT PACKET (", op, "):", reader)
	}
}

func (server *Server) handlePlayerConnect(conn mnet.Client, reader mpacket.Reader) {
	charID := reader.ReadInt32()
	clientIP := common.RemoteIPFromConn(conn)
	pending, err := common.ConsumePendingMigration(charID, common.MigrationTypeCashShop, 50, clientIP)
	if err != nil {
		log.Println("cashshop:playerConnect pending migration validation failed:", err)
		return
	}

	// Fetch channelID, migrationID and accountID in a single query
	var (
		migrationID byte
		channelID   int8
		accountID   int32
	)
	err = common.DB.QueryRow(
		"SELECT channelID, migrationID, accountID FROM characters WHERE ID=?",
		charID,
	).Scan(&channelID, &migrationID, &accountID)
	if err != nil {
		log.Println("playerConnect query error:", err)
		return
	}

	if migrationID != 50 || accountID != pending.AccountID {
		log.Println("cashshop:playerConnect: invalid migrationID:", migrationID)
		common.DeletePendingMigrationForCharacter(charID)
		return
	}

	conn.SetAccountID(accountID)

	var adminLevel int
	err = common.DB.QueryRow("SELECT adminLevel FROM accounts WHERE accountID=?", conn.GetAccountID()).Scan(&adminLevel)

	if err != nil {
		log.Println(err)
		return
	}

	conn.SetAdminLevel(adminLevel)

	_, err = common.DB.Exec("UPDATE characters SET migrationID=? WHERE ID=?", -1, charID)

	if err != nil {
		log.Println(err)
		return
	}

	plr := channel.LoadPlayerFromID(charID, conn)

	server.players.Add(&plr)

	// Load cash shop storage
	storage, err := server.GetOrLoadStorage(conn)
	if err != nil {
		log.Println("Failed to load cash shop storage for account", accountID, ":", err)
	}

	server.world.Send(internal.PacketChannelPlayerConnected(plr.ID, plr.Name, 0xFF, false, 0, 0))
	log.Printf("[CASHSHOP] consumed migration account=%d char=%d ip=%s", pending.AccountID, pending.CharacterID, clientIP)

	plr.Send(packetCashShopSet(&plr))

	// Send cash shop storage items to player (before wishlist and amounts, matching OpenMG order)
	if storage != nil {
		plr.Send(packetCashShopLoadLocker(storage, accountID, plr.ID))
	}

	wishlist, wishErr := loadWishList(plr.ID)
	if wishErr != nil {
		log.Printf("Failed to load cash shop wishlist for character %d: %v", plr.ID, wishErr)
		wishlist = make([]int32, 10)
	}
	plr.Send(packetCashShopWishList(wishlist, false))
	plr.Send(packetCashShopUpdateAmounts(plr.GetNX(), plr.GetMaplePoints(), 0))
}

func decodeCashShopGiftRequest(reader *mpacket.Reader) (byte, string, string, int32, error) {
	currencySel := reader.ReadByte()
	if len(reader.GetRestAsBytes()) < 2 {
		return 0, "", "", 0, fmt.Errorf("missing recipient length")
	}

	recipientLen := reader.ReadInt16()
	recipient := reader.ReadString(recipientLen)
	if recipient == "" {
		return 0, "", "", 0, fmt.Errorf("missing recipient")
	}

	message := ""
	if len(reader.GetRestAsBytes()) > 4 {
		messageLen := reader.ReadInt16()
		message = reader.ReadString(messageLen)
	}
	if len(reader.GetRestAsBytes()) < 4 {
		return 0, "", "", 0, fmt.Errorf("missing SN")
	}

	return currencySel, recipient, message, reader.ReadInt32(), nil
}

func decodeCashShopSlotIncreaseRequest(reader *mpacket.Reader) (byte, byte, int16, int32, error) {
	currencySel := reader.ReadByte()
	commodityRequest := reader.ReadByte()

	if commodityRequest == 0 {
		if len(reader.GetRestAsBytes()) < 1 {
			return 0, 0, 0, 0, fmt.Errorf("missing inventory type")
		}
		invType := reader.ReadByte()
		return currencySel, invType, 4, 4000, nil
	}

	if len(reader.GetRestAsBytes()) < 4 {
		return 0, 0, 0, 0, fmt.Errorf("missing slot increase commodity sn")
	}

	sn := reader.ReadInt32()
	commodity, ok := nx.GetCommodity(sn)
	if !ok || commodity.ItemID == 0 || commodity.OnSale == 0 || commodity.Price <= 0 {
		return 0, 0, 0, 0, fmt.Errorf("invalid slot increase commodity")
	}

	invType := byte((commodity.ItemID / 1000) % 10)
	if invType < constant.InventoryEquip || invType > constant.InventoryCash {
		return 0, 0, 0, 0, fmt.Errorf("invalid slot increase inventory type %d", invType)
	}

	return currencySel, invType, 4, commodity.Price, nil
}

func (server *Server) leaveCashShopToChannel(conn mnet.Client, reader mpacket.Reader) {
	plr, err := server.players.GetFromConn(conn)
	if err != nil || plr == nil {
		return
	}

	var prevChanID int8
	if err := common.DB.QueryRow("SELECT previousChannelID FROM characters WHERE ID=?", plr.ID).Scan(&prevChanID); err != nil {
		log.Println("Failed to fetch previousChannelID:", err)
	}

	targetChan := plr.ChannelID
	if prevChanID >= 0 && int(prevChanID) < len(server.channels) && server.channels[byte(prevChanID)].Port != 0 {
		targetChan = byte(prevChanID)
	}
	server.migrating[conn] = true

	if _, err := common.DB.Exec("UPDATE characters SET migrationID=?, inCashShop=0 WHERE ID=?", targetChan, plr.ID); err != nil {
		delete(server.migrating, conn)
		log.Println("Failed to set migrationID:", err)
		return
	}

	var worldID byte
	if err := common.DB.QueryRow("SELECT worldID FROM characters WHERE ID=?", plr.ID).Scan(&worldID); err != nil {
		delete(server.migrating, conn)
		log.Println("cashshop: failed to fetch worldID:", err)
		return
	}

	pending, err := common.CreatePendingMigration(conn.GetAccountID(), plr.ID, worldID, common.MigrationTypeChannel, int(targetChan), common.RemoteIPFromConn(conn), 30*time.Second)
	if err != nil {
		delete(server.migrating, conn)
		log.Println("cashshop: failed to create pending migration:", err)
		return
	}
	log.Printf("[CASHSHOP] created migration account=%d char=%d targetChannel=%d nonce=%s ip=%s", pending.AccountID, pending.CharacterID, pending.DestinationID, pending.Nonce, pending.ClientIP)

	var ip []byte
	var port int16

	if int(targetChan) < len(server.channels) {
		ip = server.channels[targetChan].IP
		port = server.channels[targetChan].Port
	}

	if len(ip) != 4 || port == 0 {
		log.Printf("Target channel %d missing IP/port, searching for fallback...", targetChan)

		log.Println("Sent request to world for updated channel information")
		server.world.Send(internal.PacketCashShopRequestChannelInfo())

		found := false
		for i, ch := range server.channels {
			if len(ch.IP) == 4 && ch.Port != 0 {
				targetChan = byte(i)
				ip = ch.IP
				port = ch.Port
				found = true
				log.Printf("Using fallback channel %d", targetChan)
				break
			}
		}

		if !found {
			delete(server.migrating, conn)
			log.Println("No valid fallback channels available")
			return
		}
	}

	p := mpacket.CreateWithOpcode(opcode.SendChannelChange)
	p.WriteBool(true)
	p.WriteBytes(ip)
	p.WriteInt16(port)
	conn.Send(p)
}

func (server *Server) handleCashRequest(conn mnet.Client, reader mpacket.Reader) {
	plr, err := server.players.GetFromConn(conn)
	if err != nil {
		return
	}

	plrNX := plr.GetNX()
	plrMaplePoints := plr.GetMaplePoints()

	plr.Send(packetCashShopUpdateAmounts(plrNX, plrMaplePoints, 0))
}

func (server *Server) handleCashShopOperation(conn mnet.Client, reader mpacket.Reader) {
	plr, err := server.players.GetFromConn(conn)
	if err != nil {
		return
	}

	plrNX := plr.GetNX()
	plrMaplePoints := plr.GetMaplePoints()

	sub := reader.ReadByte()
	switch sub {
	case opcode.RecvCashShopBuyItem:
		currencySel := reader.ReadByte()
		sn := reader.ReadInt32()

		commodity, ok := nx.GetCommodity(sn)
		if !ok || commodity.ItemID == 0 {
			// Unknown SN
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorOutOfStock))
			return
		}

		if commodity.OnSale == 0 || commodity.Price <= 0 {
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorOutOfStock))
			return
		}

		// Check funds
		price := commodity.Price
		switch currencySel {
		case constant.CashShopNX:
			if plrNX < price {
				plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorNotEnoughCash))
				return
			}
		case constant.CashShopMaplePoints:
			if plrMaplePoints < price {
				plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorNotEnoughCash))
				return
			}
		default:
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorUnknown))
			return
		}

		newItem, e := channel.CreateCashItemFromCommodity(commodity)
		if e != nil {
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorUnknown))
			return
		}

		// Get cash shop storage
		storage, storageErr := server.GetOrLoadStorage(conn)
		if storageErr != nil {
			log.Println("Failed to get cash shop storage:", storageErr)
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorUnknown))
			return
		}

		// Add item to storage instead of inventory
		slotIdx, added := storage.addItem(newItem, sn)
		if !added {
			log.Println("Failed to add item to cash shop storage")
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorExceededNumberOfCashItems))
			return
		}

		// Save storage
		if saveErr := storage.save(); saveErr != nil {
			log.Println("Failed to save cash shop storage:", saveErr)
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorUnknown))
			return
		}

		switch currencySel {
		case constant.CashShopNX:
			plrNX -= price
			plr.SetNX(plrNX)
		case constant.CashShopMaplePoints:
			plrMaplePoints -= price
			plr.SetMaplePoints(plrMaplePoints)
		default:
			log.Println("Unknown currency type: ", currencySel)
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorUnknown))
			return
		}

		plr.Send(packetCashShopUpdateAmounts(plrNX, plrMaplePoints, 0))

		// Send buy success packet with the specific item that was just added
		addedItem, ok := storage.getItemBySlot(int16(slotIdx + 1))
		if ok {
			plr.Send(packetCashShopBuyDone(*addedItem, conn.GetAccountID(), plr.ID))
		}

	case opcode.RecvCashShopBuyPackage, opcode.RecvCashShopGiftPackage:
		if sub == opcode.RecvCashShopGiftPackage {
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorUnknown))
			return
		}

		currencySel := reader.ReadByte()
		pkgSN := reader.ReadInt32()

		commodity, ok := nx.GetCommodity(pkgSN)
		if !ok || commodity.Price <= 0 {
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorOutOfStock))
			return
		}

		pkgItems, ok := nx.GetCashPackageEntries(pkgSN)
		if !ok || len(pkgItems) == 0 {
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorOutOfStock))
			return
		}

		price := commodity.Price
		switch currencySel {
		case constant.CashShopNX:
			if plrNX < price {
				plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorNotEnoughCash))
				return
			}
		case constant.CashShopMaplePoints:
			if plrMaplePoints < price {
				plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorNotEnoughCash))
				return
			}
		default:
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorUnknown))
			return
		}

		storage, storageErr := server.GetOrLoadStorage(conn)
		if storageErr != nil {
			log.Println("Failed to get cash shop storage:", storageErr)
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorUnknown))
			return
		}

		itemsToAdd := make([]channel.Item, 0, len(pkgItems))
		for _, entry := range pkgItems {
			var (
				item      channel.Item
				itemID    int32
				createErr error
			)

			if itCommodity, ok := nx.GetCommodity(entry); ok && itCommodity.ItemID != 0 {
				item, createErr = channel.CreateCashItemFromCommodity(itCommodity)
			} else {
				// CashPackage.img can list raw item IDs instead of SNs.
				itemID = entry
				if snByItem, ok := nx.GetCommoditySNByItemID(itemID); ok {
					if itCommodity, ok := nx.GetCommodity(snByItem); ok {
						item, createErr = channel.CreateCashItemFromCommodity(itCommodity)
					}
				}
				if item.ID == 0 && createErr == nil {
					item, createErr = channel.CreateItemFromID(itemID, 1)
					item.EnsureCashMetadata(0, 0)
				}
			}

			if createErr != nil || item.ID == 0 {
				plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorOutOfStock))
				return
			}
			itemsToAdd = append(itemsToAdd, item)
		}

		freeSlots := 0
		for i := range storage.items {
			if storage.items[i].ID == 0 {
				freeSlots++
			}
		}
		if freeSlots < len(itemsToAdd) {
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorExceededNumberOfCashItems))
			return
		}

		addedItems := make([]channel.Item, 0, len(itemsToAdd))
		for _, it := range itemsToAdd {
			slotIdx, added := storage.addItem(it, it.GetCashSN())
			if !added {
				plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorExceededNumberOfCashItems))
				return
			}
			if addedItem, ok := storage.getItemBySlot(int16(slotIdx + 1)); ok {
				addedItems = append(addedItems, *addedItem)
			}
		}

		if saveErr := storage.save(); saveErr != nil {
			log.Println("Failed to save cash shop storage:", saveErr)
			plr.Send(packetCashShopError(opcode.SendCashShopBuyFailed, constant.CashShopErrorUnknown))
			return
		}

		switch currencySel {
		case constant.CashShopNX:
			plrNX -= price
			plr.SetNX(plrNX)
		case constant.CashShopMaplePoints:
			plrMaplePoints -= price
			plr.SetMaplePoints(plrMaplePoints)
		}

		plr.Send(packetCashShopUpdateAmounts(plrNX, plrMaplePoints, 0))
		for _, it := range addedItems {
			plr.Send(packetCashShopBuyDone(it, conn.GetAccountID(), plr.ID))
		}

	case opcode.RecvCashShopGiftItem:
		currencySel, recipientName, _, sn, decodeErr := decodeCashShopGiftRequest(&reader)
		if decodeErr != nil {
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorUnknown))
			return
		}

		if currencySel != constant.CashShopNX {
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorNotEnoughCash))
			return
		}

		commodity, ok := nx.GetCommodity(sn)
		if !ok || commodity.ItemID == 0 || commodity.OnSale == 0 || commodity.Price <= 0 {
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorOutOfStock))
			return
		}

		if plrNX < commodity.Price {
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorNotEnoughCash))
			return
		}

		giftItem, createErr := channel.CreateCashItemFromCommodity(commodity)
		if createErr != nil {
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorUnknown))
			return
		}

		var recipientCharacterID int32
		var recipientAccountID int32
		var recipientGender int
		queryErr := common.DB.QueryRow(`
			SELECT ID, accountID, gender
			FROM characters
			WHERE BINARY name=? AND worldID=?`, recipientName, conn.GetWorldID()).Scan(&recipientCharacterID, &recipientAccountID, &recipientGender)
		if queryErr != nil || recipientCharacterID == 0 || recipientAccountID == 0 || recipientCharacterID == plr.ID {
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorIneligibleRecipientNameOrGender))
			return
		}
		if commodity.Gender != 2 && commodity.Gender != int32(recipientGender) {
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorIneligibleRecipientNameOrGender))
			return
		}

		if err := giftCashItemToAccount(recipientAccountID, giftItem); err != nil {
			log.Printf("giftCashItemToAccount failed sender=%d recipient=%s account=%d err=%v", plr.ID, recipientName, recipientAccountID, err)
			plr.Send(packetCashShopError(opcode.SendCashShopGiftFailed, constant.CashShopErrorExceededNumberOfCashItems))
			return
		}
		if recipient, err := server.players.GetFromID(recipientCharacterID); err == nil && recipient != nil {
			if refreshed, loadErr := LoadStorageByAccountID(recipientAccountID); loadErr != nil {
				log.Printf("cash gift locker refresh failed recipient=%d account=%d err=%v", recipientCharacterID, recipientAccountID, loadErr)
			} else {
				recipient.Conn.SetCashShopStorage(refreshed)
				recipient.Send(packetCashShopLoadLocker(refreshed, recipientAccountID, recipientCharacterID))
			}
		}

		plrNX -= commodity.Price
		plr.SetNX(plrNX)
		plr.Send(packetCashShopUpdateAmounts(plrNX, plrMaplePoints, 0))
		plr.Send(packetCashShopGiftDone(recipientName, giftItem, commodity.Price))

	case opcode.RecvCashShopUpdateWishlist:
		wishlist := make([]int32, 10)
		for i := range wishlist {
			wishlist[i] = reader.ReadInt32()
			if wishlist[i] == 0 {
				continue
			}
			if _, ok := nx.GetCommodity(wishlist[i]); !ok {
				plr.Send(packetCashShopError(opcode.SendCashShopUpdateWishFailed, constant.CashShopErrorOutOfStock))
				return
			}
		}
		if err := saveWishList(plr.ID, wishlist); err != nil {
			log.Printf("saveWishList failed characterID=%d err=%v", plr.ID, err)
			plr.Send(packetCashShopError(opcode.SendCashShopUpdateWishFailed, constant.CashShopErrorUnknown))
			return
		}
		plr.Send(packetCashShopWishList(wishlist, true))

	case opcode.RecvCashShopIncreaseSlots:
		currencySel, invType, delta, price, decodeErr := decodeCashShopSlotIncreaseRequest(&reader)
		if decodeErr != nil {
			log.Printf("cashshop slot increase decode failed player=%d err=%v", plr.ID, decodeErr)
			plr.Send(packetCashShopError(opcode.SendCashShopIncSlotCountFailed, constant.CashShopErrorUnknown))
			return
		}

		switch currencySel {
		case constant.CashShopNX:
			if plrNX < price {
				plr.Send(packetCashShopError(opcode.SendCashShopIncSlotCountFailed, constant.CashShopErrorNotEnoughCash))
				return
			}
			if err := plr.IncreaseSlotSize(invType, byte(delta)); err != nil {
				plr.Send(packetCashShopError(opcode.SendCashShopIncSlotCountFailed, constant.CashShopErrorUnknown))
				return
			}
			plrNX -= price
			plr.SetNX(plrNX)
		case constant.CashShopMaplePoints:
			if plrMaplePoints < price {
				plr.Send(packetCashShopError(opcode.SendCashShopIncSlotCountFailed, constant.CashShopErrorNotEnoughCash))
				return
			}
			if err := plr.IncreaseSlotSize(invType, byte(delta)); err != nil {
				plr.Send(packetCashShopError(opcode.SendCashShopIncSlotCountFailed, constant.CashShopErrorUnknown))
				return
			}
			plrMaplePoints -= price
			plr.SetMaplePoints(plrMaplePoints)
		default:
			log.Println("Unknown currency type: ", currencySel)
			plr.Send(packetCashShopError(opcode.SendCashShopIncSlotCountFailed, constant.CashShopErrorUnknown))
			return
		}
		plr.FlushState()

		plr.Send(packetCashShopIncreaseInv(invType, plr.GetSlotSize(invType)))
		plr.Send(packetCashShopUpdateAmounts(plrNX, plrMaplePoints, 0))

	case opcode.RecvCashShopMoveLtoS:
		cashItemID := reader.ReadInt64()
		_ = reader.ReadByte()
		_ = reader.ReadInt16()

		storage, storageErr := server.GetOrLoadStorage(conn)
		if storageErr != nil {
			plr.Send(packetCashShopError(opcode.SendCashShopMoveLtoSFailed, constant.CashShopErrorUnknown))
			return
		}

		var foundIdx = -1
		var foundItem *channel.Item
		for i := range storage.items {
			if storage.items[i].ID == 0 {
				continue
			}
			if storage.items[i].GetCashID() == cashItemID {
				foundIdx = i
				foundItem = &storage.items[i]
				break
			}
		}

		if foundIdx == -1 || foundItem == nil {
			plr.Send(packetCashShopError(opcode.SendCashShopMoveLtoSFailed, constant.CashShopErrorUnknown))
			return
		}

		removedItem, ok := storage.removeAt(foundIdx)
		if !ok {
			plr.Send(packetCashShopError(opcode.SendCashShopMoveLtoSFailed, constant.CashShopErrorUnknown))
			return
		}

		item := *removedItem
		givenItem, err := plr.GiveItem(item)
		if err != nil {
			if _, restored := storage.addItemWithCashID(item, item.GetCashSN(), item.GetCashID()); !restored {
				log.Println("CRITICAL: Restore to storage failed. Item may be lost. player:", plr.ID, "accountID:", conn.GetAccountID(), "itemID:", item.ID)
			} else {
				if saveErr := storage.save(); saveErr != nil {
					log.Println("Failed to save restored storage:", saveErr)
				}
			}
			plr.Send(packetCashShopError(opcode.SendCashShopMoveLtoSFailed, constant.CashShopErrorCheckFullInventory))
			return
		}

		if saveErr := storage.save(); saveErr != nil {
			plr.Send(packetCashShopError(opcode.SendCashShopMoveLtoSFailed, constant.CashShopErrorUnknown))
			return
		}

		plr.Send(packetCashShopMoveLtoSDone(givenItem, givenItem.GetSlotID()))

	case opcode.RecvCashShopMoveStoL:
		// Move from slot (inventory) to locker (storage)
		cashItemID := reader.ReadInt64()
		invType := reader.ReadByte()

		storage, storageErr := server.GetOrLoadStorage(conn)
		if storageErr != nil {
			log.Println("Failed to get storage:", storageErr)
			plr.Send(packetCashShopError(opcode.SendCashShopMoveStoLFailed, constant.CashShopErrorUnknown))
			return
		}

		item, itemSlot, findErr := plr.GetItemByCashID(invType, cashItemID)
		if findErr != nil {
			plr.Send(packetCashShopError(opcode.SendCashShopMoveStoLFailed, constant.CashShopErrorUnknown))
			return
		}

		expectedInvType := byte(item.ID / 1000000)
		if expectedInvType != invType {
			plr.Send(packetCashShopError(opcode.SendCashShopMoveStoLFailed, constant.CashShopErrorUnknown))
			return
		}

		takenItem, takeErr := plr.TakeItemSilent(item.ID, itemSlot, 1, invType)
		if takeErr != nil {
			plr.Send(packetCashShopError(opcode.SendCashShopMoveStoLFailed, constant.CashShopErrorUnknown))
			return
		}

		slotIdx, added := storage.addItemWithCashID(takenItem, takenItem.GetCashSN(), takenItem.GetCashID())
		if !added {
			if _, err := plr.GiveItem(takenItem); err != nil {
				log.Println("CRITICAL: Failed to return item to player after add failure:", err)
			}
			plr.Send(packetCashShopError(opcode.SendCashShopMoveStoLFailed, constant.CashShopErrorExceededNumberOfCashItems))
			return
		}

		if saveErr := storage.save(); saveErr != nil {
			log.Println("Failed to save storage:", saveErr)
			plr.Send(packetCashShopError(opcode.SendCashShopMoveStoLFailed, constant.CashShopErrorUnknown))
			return
		}

		addedItem, ok := storage.getItemBySlot(int16(slotIdx + 1))
		if ok {
			plr.Send(packetCashShopMoveStoLDone(*addedItem, conn.GetAccountID()))
		} else {
			plr.Send(packetCashShopError(opcode.SendCashShopMoveStoLFailed, constant.CashShopErrorUnknown))
		}

	default:
		log.Println("Unknown Cash Shop Packet(", sub, "): ", reader)
	}

}
