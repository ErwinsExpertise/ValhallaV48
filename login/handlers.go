package login

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/internal"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
)

// HandleClientPacket data
func (server *Server) HandleClientPacket(conn mnet.Client, reader mpacket.Reader) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[LOGIN] panic handling client packet from %s: %v", conn, r)
		}
	}()

	op := reader.ReadInt16()

	switch op {
	case opcode.RecvLoginRequest:
		server.handleLoginRequest(conn, reader)
	case opcode.RecvLoginWorldInfoRequest:
		server.handleWorldInfoRequest(conn, reader)
	case opcode.RecvLoginEULA:
		server.handleEULA(conn, reader)
	case opcode.RecvLoginCheckLogin:
		server.handleGoodLogin(conn, reader)
	case opcode.RecvLoginRegisterPin:
		server.handlePinRegistration(conn, reader)
	case opcode.RecvLoginLogoutWorld:
		server.handleLogoutWorld(conn, reader)
	case opcode.RecvLoginViewAllCharacters:
		server.handleViewAllCharacters(conn, reader)
	case opcode.RecvLoginUnknown14:
		server.handleLoginUnknown14(conn, reader)
	case opcode.RecvLoginWorldSelect:
		server.handleWorldSelect(conn, reader)
	case opcode.RecvLoginChannelSelect:
		server.handleChannelSelect(conn, reader)
	case opcode.RecvLoginNameCheck:
		server.handleNameCheck(conn, reader)
	case opcode.RecvLoginNewCharacter:
		server.handleNewCharacter(conn, reader)
	case opcode.RecvLoginDeleteChar:
		server.handleDeleteCharacter(conn, reader)
	case opcode.RecvLoginSelectCharacter:
		server.handleSelectCharacter(conn, reader)
	case opcode.RecvReturnToLoginScreen:
		server.handleReturnToLoginScreen(conn, reader)
	case opcode.RecvPing:
		// Thumbs Up
	default:
		log.Printf("[LOGIN] UNKNOWN CLIENT PACKET: opcode=0x%04X, data=% X", op, reader.GetBuffer())
	}
}

func (server *Server) sendWorldList(conn mnet.Client) {
	for i := range server.worlds {
		if server.worlds[i].Conn == nil {
			continue
		}
		conn.Send(packetLoginWorldListing(byte(i), server.worlds[i]))
	}
	conn.Send(packetLoginEndWorldList())
}

func (server *Server) handleLoginRequest(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if sess.stage != sessionStageAwaitLogin {
		_ = conn.Close()
		return
	}

	username := reader.ReadString(reader.ReadInt16())
	password := reader.ReadString(reader.ReadInt16())

	// Try to read HWID
	reader.Skip(6)
	hwidBytes := reader.ReadBytes(4)
	hwid := strings.ToUpper(hex.EncodeToString(hwidBytes))

	ip := clientRemoteIP(conn)

	var accountID int32
	var user string
	var databasePassword string
	var gender byte
	var isLogedIn bool
	var isBanned int
	var isLocked int
	var lastHwid sql.NullString
	var adminLevel int
	var eula byte
	var passwordSalt sql.NullString

	err := common.DB.QueryRow("SELECT accountID, username, password, passwordSalt, gender, isLogedIn, isBanned, isLocked, hwid, adminLevel, eula FROM accounts WHERE username=?", username).
		Scan(&accountID, &user, &databasePassword, &passwordSalt, &gender, &isLogedIn, &isBanned, &isLocked, &lastHwid, &adminLevel, &eula)

	result := constant.LoginResultSuccess

	lastHwidStr := ""
	if lastHwid.Valid {
		lastHwidStr = lastHwid.String
	}

	if server.ac != nil {
		banned, _, endEpoch, err := server.ac.IsBanned(accountID, ip, hwid)
		if err != nil {
			return
		}
		if banned {
			pkt := packetLoginBanned(endEpoch, constant.BanReasonHacking)
			conn.Send(pkt)
			return
		}
	}

	if err != nil {
		log.Println(err)
		if server.autoRegister {
			hashedPassword, passwordSaltValue, insertErr := makeStoredCredential(password)
			if insertErr != nil {
				log.Println("Failed to hash new account password", insertErr)
				result = constant.LoginResultSystemError
			} else {
				hashedPin, pinSaltValue, pinErr := makeStoredCredential(constant.AutoRegisterDefaultPIN)
				if pinErr != nil {
					log.Println("Failed to hash new account pin", pinErr)
					result = constant.LoginResultSystemError
				} else {
					res, insertErr := common.DB.Exec("INSERT INTO accounts (username, password, passwordSalt, pin, pinSalt, isLogedIn, adminLevel, isBanned, gender, dob, eula, nx, maplepoints, hwid) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
						username, hashedPassword, passwordSaltValue, hashedPin, pinSaltValue, constant.AutoRegisterDefaultIsLoggedIn,
						constant.AutoRegisterDefaultAdminLevel, constant.AutoRegisterDefaultIsBanned, constant.AutoRegisterDefaultGender,
						constant.AutoRegisterDefaultDOB, constant.AutoRegisterDefaultEULA, constant.AutoRegisterDefaultNX,
						constant.AutoRegisterDefaultMaplePoints, hwid)

					if insertErr != nil {
						log.Println("Failed to create new account", err)
						result = constant.LoginResultNotRegistered
					} else if id, err := res.LastInsertId(); err == nil {
						accountID = int32(id)
						gender = constant.AutoRegisterDefaultGender
						adminLevel = constant.AutoRegisterDefaultAdminLevel
						eula = constant.AutoRegisterDefaultEULA
						log.Println("Auto-registered new account:", username, "with ID:", accountID)
						result = constant.LoginResultSuccess
					} else {
						log.Println("Failed to get new account ID:", err)
						result = constant.LoginResultSystemError
					}
				}
			}
		} else {
			result = constant.LoginResultNotRegistered
		}
	} else if isLocked > 0 {
		result = constant.LoginResultDeletedOrBlocked
	} else if !verifyStoredCredential(password, databasePassword, passwordSalt) {
		if server.ac != nil {
			ipKey := fmt.Sprintf("ip:%s", ip)
			userKey := fmt.Sprintf("user:%s", username)
			hwidKey := fmt.Sprintf("hwid:%s", hwid)

			exceeded := server.ac.TrackFailedAuth(ipKey) || server.ac.TrackFailedAuth(hwidKey) || server.ac.TrackFailedAuth(userKey)
			if exceeded && strings.Compare(hwid, lastHwidStr) != 0 {
				// Lock Account
				err := server.ac.LockAccount(accountID)
				if err != nil {
					log.Println(err)
				}
			}
		}

		result = constant.LoginResultInvalidPassword
	} else if isLogedIn {
		active, reconcileErr := common.ReconcileAccountLoginState(accountID)
		if reconcileErr != nil {
			log.Println(reconcileErr)
			result = constant.LoginResultSystemError
		} else if active {
			result = constant.LoginResultAlreadyOnline
		}
	} else if isBanned > 0 {
		result = constant.LoginResultBanned
	} else if eula == 0 {
		result = constant.LoginResultEULA
	}
	// Banned = 2, Deleted or Blocked = 3, Invalid Password = 4, Not Registered = 5, Sys Error = 6,
	// Already online = 7, System error = 9, Too many requests = 10, Older than 20 = 11, valid login = 12, Master cannot login on this IP = 13,
	// wrong gateway korean text = 14, still processing request korean text = 15, verify email = 16, gateway english text = 17,
	// verify email = 21, eula = 23

	if result <= constant.LoginResultSuccess || result == constant.LoginResultEULA {
		sess.username = username
		sess.accountID = accountID
		conn.SetGender(gender)
		conn.SetAdminLevel(adminLevel)
		conn.SetAccountID(accountID)
		conn.SetHWID(hwid)

		if result <= constant.LoginResultSuccess {
			// Update HWID and clear failed attempts on successful login
			if hwid != "" {
				common.DB.Exec("UPDATE accounts SET hwid = ? WHERE accountID = ?", hwid, accountID)
			}
			if server.ac != nil {
				ip := ""
				if host, _, err := net.SplitHostPort(conn.String()); err == nil {
					ip = host
				} else {
					ip = conn.String()
				}
				server.ac.ClearAuth(
					fmt.Sprintf("user:%s", username),
					fmt.Sprintf("ip:%s", ip),
					fmt.Sprintf("hwid:%s", hwid),
				)
			}
		}

		if result == constant.LoginResultEULA {
			sess.stage = sessionStageAwaitEULA
		} else {
			sess.stage = sessionStageAwaitPin
		}
	}

	conn.Send(packetLoginResponse(result, accountID, gender, adminLevel > 0, username, isBanned))
}

func (server *Server) handleEULA(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if sess.stage != sessionStageAwaitEULA {
		_ = conn.Close()
		return
	}

	accept := false
	if len(reader.GetRestAsBytes()) > 0 {
		accept = reader.ReadByte() != 0
	}
	if !accept {
		return
	}

	accountID := conn.GetAccountID()
	if accountID <= 0 {
		log.Println("Could not set EULA signed: invalid accountID", accountID)
		return
	}

	res, err := common.DB.Exec("UPDATE accounts SET eula=? WHERE accountID=?", 1, accountID)

	if err != nil {
		log.Println("Could not set EULA signed", err)
		return
	}

	if rows, rowsErr := res.RowsAffected(); rowsErr != nil {
		log.Println("Could not confirm EULA update", rowsErr)
	} else if rows == 0 {
		log.Println("Could not set EULA signed: no rows updated for accountID", accountID)
		return
	}

	var username string
	var gender byte
	var adminLevel int
	var isBanned int
	err = common.DB.QueryRow("SELECT username, gender, adminLevel, isBanned FROM accounts WHERE accountID=?", accountID).
		Scan(&username, &gender, &adminLevel, &isBanned)
	if err != nil {
		log.Println("Could not load account after EULA acceptance", err)
		return
	}

	conn.SetGender(gender)
	conn.SetAdminLevel(adminLevel)
	sess.stage = sessionStageAwaitPin
	conn.Send(packetLoginResponse(constant.LoginResultSuccess, accountID, gender, adminLevel > 0, username, isBanned))
}

func (server *Server) handlePinRegistration(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if sess.stage != sessionStageAwaitPinRegistration {
		_ = conn.Close()
		return
	}

	if !server.withPin {
		conn.Send(packetCancelPin())
		return
	}

	b1 := reader.ReadByte()

	if b1 == 0 { // Client canceled pin change request
		conn.Send(packetCancelPin())
		return
	}
	reader.Skip(2)

	accountID := conn.GetAccountID()
	pin := string(reader.GetRestAsBytes())
	if len(pin) != 4 {
		conn.Send(packetRegisterPin())
		return
	}
	for _, ch := range pin {
		if ch < '0' || ch > '9' {
			conn.Send(packetRegisterPin())
			return
		}
	}

	hashedPin, pinSaltValue, err := makeStoredCredential(pin)
	if err != nil {
		log.Println("handlePinRegistration failed to hash pin for accountID:", accountID, err)
		return
	}

	_, err = common.DB.Exec("UPDATE accounts SET pin=?, pinSalt=? WHERE accountID=?", hashedPin, pinSaltValue, accountID)
	if err != nil {
		log.Println("handlePinRegistration database pin update issue for accountID:", accountID, err)
	}

	sess.stage = sessionStageAwaitPin
	conn.Send(packetRequestPin())

}

func (server *Server) handleGoodLogin(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if sess.stage != sessionStageAwaitPin {
		_ = conn.Close()
		return
	}

	server.migrating[conn] = false
	accountID := conn.GetAccountID()

	if server.withPin {
		var pinDB string
		var pinSalt sql.NullString
		var authDone bool

		err := common.DB.QueryRow("SELECT pin, pinSalt FROM accounts WHERE accountID=?", accountID).
			Scan(&pinDB, &pinSalt)

		if err != nil {
			log.Println("handleCheckLogin database retrieval issue for accountID:", accountID, err)
		}

		b1 := reader.ReadByte()
		b2 := reader.ReadByte()

		if b1 == 1 && b2 == 1 { // First attempt, request for pin
			if len(pinDB) == 0 {
				sess.stage = sessionStageAwaitPinRegistration
				conn.Send(packetRegisterPin())
			} else {
				conn.Send(packetRequestPin())
			}

		} else if b1 == 1 || b1 == 2 { // Client assigned pin
			reader.Skip(6) // space padding?
			pin := string(reader.GetRestAsBytes())

			if !verifyStoredCredential(pin, pinDB, pinSalt) {
				conn.Send(packetRequestPinAfterFailure())

			} else if b1 == 2 { // Changing pin request
				sess.stage = sessionStageAwaitPinRegistration
				conn.Send(packetRegisterPin())

			} else { // Authenticated successfully
				authDone = true
			}

		} else if b1 == 0 { // Client cancels pin request
			conn.Send(packetCancelPin())
		}

		if !authDone {
			return
		}
	}

	conn.SetLogedIn(true)
	_, err := common.DB.Exec("UPDATE accounts set isLogedIn=1 WHERE accountID=?", accountID)

	if err != nil {
		log.Println("Database error with approving login of accountID", accountID, err)
	} else {
		log.Println("User", accountID, "has logged in from", conn)
	}
	sess.onlineMarked = true
	sess.stage = sessionStageAwaitWorldSelect

	server.sendWorldList(conn)
}

func (server *Server) handleWorldInfoRequest(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if !sess.onlineMarked {
		_ = conn.Close()
		return
	}

	if sess.stage < sessionStageAwaitWorldSelect {
		sess.stage = sessionStageAwaitWorldSelect
	}
	server.sendWorldList(conn)
}

func (server *Server) handleWorldSelect(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if sess.stage != sessionStageAwaitWorldSelect {
		_ = conn.Close()
		return
	}

	worldID := reader.ReadByte()
	if !server.validWorld(worldID) {
		conn.Send(packetLoginReturnFromChannel())
		return
	}

	log.Printf("world %d selected", worldID)
	conn.SetWorldID(worldID)
	sess.worldID = worldID
	reader.ReadByte() // ?

	var warning, population byte = 0, 0

	if conn.GetAdminLevel() < 1 { // gms are not restricted in any capacity
		var currentPlayers int16
		var maxPlayers int16

		for _, v := range server.worlds[conn.GetWorldID()].Channels {
			currentPlayers += v.Pop
			maxPlayers += v.MaxPop
		}

		if currentPlayers >= maxPlayers {
			warning = 2
		} else if float64(currentPlayers)/float64(maxPlayers) > 0.90 { // I'm not sure if this warning is even worth it
			warning = 1
		}
	}

	conn.Send(packetLoginWorldInfo(warning, population)) // hard coded for now
	sess.stage = sessionStageAwaitChannelSelect
}

func (server *Server) handleChannelSelect(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if sess.stage != sessionStageAwaitChannelSelect {
		_ = conn.Close()
		return
	}

	selectedWorld := reader.ReadByte()   // world
	conn.SetChannelID(reader.ReadByte()) // Channel
	if selectedWorld != conn.GetWorldID() || !server.validChannel(selectedWorld, conn.GetChannelID()) {
		conn.Send(packetLoginReturnFromChannel())
		return
	}
	sess.channelID = conn.GetChannelID()

	if server.worlds[selectedWorld].Channels[conn.GetChannelID()].MaxPop == 0 {
		conn.Send(packetLoginReturnFromChannel())
		return
	}

	if selectedWorld == conn.GetWorldID() {
		characters := getCharactersFromAccountWorldID(conn.GetAccountID(), conn.GetWorldID())
		conn.Send(packetLoginDisplayCharacters(characters))
		sess.stage = sessionStageAwaitCharacterSelect
	}
}

func (server *Server) handleLogoutWorld(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if !sess.onlineMarked {
		_ = conn.Close()
		return
	}

	conn.SetWorldID(0)
	conn.SetChannelID(0)
	sess.worldID = 0
	sess.channelID = 0
	sess.stage = sessionStageAwaitWorldSelect
	conn.Send(packetLoginReturnFromChannel())
}

func (server *Server) handleViewAllCharacters(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if !sess.onlineMarked || sess.stage < sessionStageAwaitWorldSelect {
		_ = conn.Close()
		return
	}

	charactersByWorld := getCharactersFromAccountAllWorlds(conn.GetAccountID())
	if len(charactersByWorld) == 0 {
		conn.Send(packetLoginViewAllCharactersEmpty())
		return
	}

	totalCharacters := 0
	worldIDs := make([]int, 0, len(charactersByWorld))
	for worldID, characters := range charactersByWorld {
		totalCharacters += len(characters)
		worldIDs = append(worldIDs, int(worldID))
	}
	sort.Ints(worldIDs)

	conn.Send(packetLoginViewAllCharactersSummary(len(worldIDs), totalCharacters))
	for _, worldID := range worldIDs {
		conn.Send(packetLoginViewAllCharactersWorld(byte(worldID), charactersByWorld[byte(worldID)]))
	}
	sess.stage = sessionStageAwaitCharacterSelect
}

func (server *Server) handleLoginUnknown14(conn mnet.Client, reader mpacket.Reader) {
	// v48 emits this around unsupported view-all/close flows; ignore instead of treating it as a fatal unknown packet.
}

func (server *Server) handleNameCheck(conn mnet.Client, reader mpacket.Reader) {
	if server.sessionFor(conn).stage != sessionStageAwaitCharacterSelect {
		_ = conn.Close()
		return
	}

	newCharName := reader.ReadString(reader.ReadInt16())

	var nameFound int
	err := common.DB.QueryRow("SELECT count(*) name FROM characters WHERE name=?", newCharName).
		Scan(&nameFound)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		// Default to name found just in-case
		conn.Send(packetLoginNameCheck(newCharName, 1))
		return
	}

	conn.Send(packetLoginNameCheck(newCharName, nameFound))
}

func (server *Server) handleNewCharacter(conn mnet.Client, reader mpacket.Reader) {
	if server.sessionFor(conn).stage != sessionStageAwaitCharacterSelect {
		_ = conn.Close()
		return
	}

	name := reader.ReadString(reader.ReadInt16())
	face := reader.ReadInt32()
	hair := reader.ReadInt32()
	hairColour := reader.ReadInt32()
	skin := reader.ReadInt32()
	top := reader.ReadInt32()
	bottom := reader.ReadInt32()
	shoes := reader.ReadInt32()
	weapon := reader.ReadInt32()

	_ = reader.ReadByte()
	_ = reader.ReadByte()
	_ = reader.ReadByte()
	_ = reader.ReadByte()

	const (
		baseStat = byte(4)
		startAP  = int16(9)
	)

	// Add str, dex, int, luk validation (check to see if client generates a constant sum)

	var counter int
	var accountCharCount int

	err := common.DB.QueryRow("SELECT count(*) FROM characters where name=? and worldID=?", name, conn.GetWorldID()).Scan(&counter)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return
	}
	if err := common.DB.QueryRow("SELECT count(*) FROM characters WHERE accountID=? AND worldID=?", conn.GetAccountID(), conn.GetWorldID()).Scan(&accountCharCount); err != nil {
		log.Println(err)
		return
	}

	allowedEyes := []int32{20000, 20001, 20002, 21000, 21001, 21002, 20100, 20401, 20402, 21700, 21201, 21002}
	allowedHair := []int32{30000, 30020, 30030, 31000, 31040, 31050}
	allowedHairColour := []int32{0, 7, 3, 2}
	allowedBottom := []int32{1060002, 1060006, 1061002, 1061008}                         // v48 missing: 1062115
	allowedTop := []int32{1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011} // v48 missing: 1042167
	allowedShoes := []int32{1072001, 1072005, 1072037, 1072038}                          // v48 missing: 1072383
	allowedWeapons := []int32{1302000, 1322005, 1312004}                                 // v48 missing: 1442079
	allowedSkinColour := []int32{0, 1, 2, 3}

	inSlice := func(val int32, s []int32) bool {
		for _, b := range s {
			if b == val {
				return true
			}
		}
		return false
	}

	valid := inSlice(face, allowedEyes) && inSlice(hair, allowedHair) && inSlice(hairColour, allowedHairColour) &&
		inSlice(bottom, allowedBottom) && inSlice(top, allowedTop) && inSlice(shoes, allowedShoes) &&
		inSlice(weapon, allowedWeapons) && inSlice(skin, allowedSkinColour) && (counter == 0) && accountCharCount < 4

	newCharacter := player{}

	if conn.GetAdminLevel() > 0 {
		name = "[GM]" + name
	} else if strings.ContainsAny(name, "[]") {
		valid = false // hacked client or packet editting
	}

	if valid {
		res, err := common.DB.Exec("INSERT INTO characters (name, accountID, worldID, face, hair, skin, gender, str, dex, intt, luk, ap) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			name, conn.GetAccountID(), conn.GetWorldID(), face, hair+hairColour, skin, conn.GetGender(), baseStat, baseStat, baseStat, baseStat, startAP)

		if err != nil {
			log.Println(err)
		}

		characterID, err := res.LastInsertId()

		if err != nil {
			log.Println(err)
		}

		char := loadPlayerFromID(int32(characterID))

		if conn.GetAdminLevel() > 0 {
			items := map[int32]int16{
				1002140: -1,  // Hat
				1032006: -4,  // Earrings
				1042003: -5,  // top
				1062007: -6,  // bottom
				1072004: -7,  // shoes
				1082002: -8,  // Gloves
				1102054: -9,  // Cape
				1092008: -10, // Shield
				1322013: -11, // weapon
			}

			for id, pos := range items {
				item := newAdminItem(id, char.name)

				if err != nil {
					log.Println(err)
					return
				}

				item.slotID = pos
				item.creatorName = name
				item.save(int32(characterID))
				char.equip = append(char.equip, item)
			}

			// TODO: Give GM all skils maxed
		} else {
			items := map[int32]int16{
				top:    -5,
				bottom: -6,
				shoes:  -7,
				weapon: -11,
			}

			for id, pos := range items {
				item := newBeginnerItem(id)

				if err != nil {
					log.Println(err)
					return
				}

				item.slotID = pos
				item.save(int32(characterID))
				char.equip = append(char.equip, item)
			}
		}

		char.save()
		newCharacter = char
	}

	conn.Send(packetLoginCreatedCharacter(valid, newCharacter))
}

func (server *Server) handleDeleteCharacter(conn mnet.Client, reader mpacket.Reader) {
	if server.sessionFor(conn).stage != sessionStageAwaitCharacterSelect {
		_ = conn.Close()
		return
	}

	dob := reader.ReadInt32()
	charID := reader.ReadInt32()

	var storedDob int32
	var charCount int

	err := common.DB.QueryRow("SELECT dob FROM accounts where accountID=?", conn.GetAccountID()).Scan(&storedDob)
	err = common.DB.QueryRow("SELECT count(*) FROM characters where accountID=? AND id=?", conn.GetAccountID(), charID).Scan(&charCount)
	if err != nil {
		log.Println(err)
		return
	}

	hacking := false
	deleted := false

	if charCount != 1 {
		if server.ac != nil {
			err = server.ac.IssueBan(0, 24, "Attempted to delete character not associated with account", conn.String(), conn.GetHWID())
			if err != nil {
				log.Println(err)
			}
		}
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
		return
	}

	if dob == storedDob {
		records, err := common.DB.Query("DELETE FROM characters where id=?", charID)

		defer records.Close()

		if err != nil {
			log.Println(err)
			return
		}

		deleted = true
	}

	if deleted {
		for _, v := range server.worlds {
			if v.Conn != nil {
				v.Conn.Send(internal.PacketLoginDeletedCharacter(charID))
			}
		}
	}

	conn.Send(packetLoginDeleteCharacter(charID, deleted, hacking))
}

func (server *Server) handleSelectCharacter(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if sess.stage != sessionStageAwaitCharacterSelect {
		_ = conn.Close()
		return
	}

	charID := reader.ReadInt32()

	var charWorldID byte
	var channelID int8
	var migrationID int8
	var inCashShop bool

	err := common.DB.QueryRow("SELECT worldID, channelID, migrationID, inCashShop FROM characters where accountID=? AND id=?", conn.GetAccountID(), charID).Scan(&charWorldID, &channelID, &migrationID, &inCashShop)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			conn.Send(packetLoginReturnFromChannel())
			return
		}
		log.Println(err)
		if server.ac != nil {
			err = server.ac.IssueBan(0, 24, "Attempted to select character not associated with account", conn.String(), conn.GetHWID())
			if err != nil {
				log.Println(err)
			}
		}
		return
	}

	if charWorldID != conn.GetWorldID() || migrationID != -1 || inCashShop || channelID != -1 {
		conn.Send(packetLoginReturnFromChannel())
		return
	}
	if !server.validChannel(conn.GetWorldID(), conn.GetChannelID()) {
		conn.Send(packetLoginReturnFromChannel())
		return
	}

	channel := server.worlds[conn.GetWorldID()].Channels[conn.GetChannelID()]
	_, err = common.DB.Exec("UPDATE characters SET migrationID=? WHERE id=?", conn.GetChannelID(), charID)

	if err != nil {
		log.Println(err)
		return
	}

	pending, err := common.CreatePendingMigration(conn.GetAccountID(), charID, conn.GetWorldID(), common.MigrationTypeChannel, int(conn.GetChannelID()), clientRemoteIP(conn), 30*time.Second)
	if err != nil {
		log.Println("failed to create pending migration", err)
		conn.Send(packetLoginReturnFromChannel())
		return
	}

	server.migrating[conn] = true
	sess.stage = sessionStageMigrating
	sess.migrationChar = charID
	log.Printf("[LOGIN] migration created account=%d char=%d world=%d channel=%d nonce=%s ip=%s", pending.AccountID, pending.CharacterID, pending.WorldID, pending.DestinationID, pending.Nonce, pending.ClientIP)

	conn.Send(packetLoginMigrateClient(channel.IP, channel.Port, charID))
}

func (server *Server) addCharacterItem(characterID int64, itemID int32, slot int32, creatorName string) {
	_, err := common.DB.Exec("INSERT INTO items (characterID, itemID, slotNumber, creatorName) VALUES (?, ?, ?, ?)", characterID, itemID, slot, creatorName)

	if err != nil {
		log.Println(err)
		return
	}
}

func (server *Server) handleReturnToLoginScreen(conn mnet.Client, reader mpacket.Reader) {
	sess := server.sessionFor(conn)
	if sess.onlineMarked {
		sess.stage = sessionStageAwaitWorldSelect
	}
	conn.Send(packetLoginReturnFromChannel())
}

// HandleServerPacket from world
func (server *Server) HandleServerPacket(conn mnet.Server, reader mpacket.Reader) {
	switch reader.ReadByte() {
	case opcode.WorldNew:
		server.handleNewWorld(conn, reader)
	case opcode.WorldInfo:
		server.handleWorldInfo(conn, reader)
	default:
		log.Println("UNKNOWN WORLD PACKET:", reader)
	}
}

// The following logic could do with being cleaned up
func (server *Server) handleNewWorld(conn mnet.Server, reader mpacket.Reader) {
	log.Println("Server register request from", conn)
	if len(server.worlds) > 14 {
		log.Println("Rejected")
		conn.Send(mpacket.CreateInternal(opcode.WorldRequestBad))
	} else {
		name := reader.ReadString(reader.ReadInt16())

		if name == "" {
			name = constant.WORLD_NAMES[len(server.worlds)]

			registered := false
			for i, v := range server.worlds {
				if v.Conn == nil {
					server.worlds[i].Conn = conn
					name = server.worlds[i].Name

					registered = true
					break
				}
			}

			if !registered {
				server.worlds = append(server.worlds, internal.World{Conn: conn, Name: name})
			}

			p := mpacket.CreateInternal(opcode.WorldRequestOk)
			p.WriteString(name)
			conn.Send(p)

			log.Println("Registered", name)
		} else {
			registered := false
			for i, w := range server.worlds {
				if w.Name == name {
					server.worlds[i].Conn = conn
					server.worlds[i].Name = name

					p := mpacket.CreateInternal(opcode.WorldRequestOk)
					p.WriteString(name)
					conn.Send(p)

					registered = true

					break
				}
			}

			if !registered {
				server.worlds = append(server.worlds, internal.World{Conn: conn, Name: name})

				p := mpacket.CreateInternal(opcode.WorldRequestOk)
				p.WriteString(server.worlds[len(server.worlds)-1].Name)
				conn.Send(p)
			}

			log.Println("Re-registered", name)
		}
	}
}

func (server *Server) handleWorldInfo(conn mnet.Server, reader mpacket.Reader) {
	for i, v := range server.worlds {
		if v.Conn != conn {
			continue
		}

		server.worlds[i].SerialisePacket(reader)

		if v.Name == "" {
			log.Println("Registered new world", server.worlds[i].Name)
		} else {
			log.Println("Updated world info for", v.Name)
		}
	}
}
