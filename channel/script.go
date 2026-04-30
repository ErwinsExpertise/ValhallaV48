package channel

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/internal"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/nx"
	"github.com/dop251/goja"
	"github.com/fsnotify/fsnotify"
)

type scriptStore struct {
	folder   string
	scripts  map[string]*goja.Program
	dispatch chan func()
}

func createScriptStore(folder string, dispatch chan func()) *scriptStore {
	return &scriptStore{folder: folder, dispatch: dispatch, scripts: make(map[string]*goja.Program)}
}

func (s scriptStore) String() string {
	return fmt.Sprintf("%v", s.scripts)
}

func (s *scriptStore) loadScripts() error {
	err := filepath.Walk(s.folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		name, program, errComp := createScriptProgramFromFilename(path)

		if errComp == nil {
			s.scripts[name] = program
		} else {
			log.Println("Script compiling:", errComp)
		}

		return nil
	})

	return err
}

func (s *scriptStore) monitor(task func(name string, program *goja.Program)) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Println(err)
	}

	defer watcher.Close()

	err = watcher.Add(s.folder)

	if err != nil {
		log.Println(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				s.dispatch <- func() {
					log.Println("Script:", event.Name, "modified/created")
					name, program, err := createScriptProgramFromFilename(event.Name)

					if err == nil {
						s.scripts[name] = program
						task(name, program)
					} else {
						log.Println("Script compiling:", err)
					}
				}
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				s.dispatch <- func() {
					name := filepath.Base(event.Name)
					name = strings.TrimSuffix(name, filepath.Ext(name))

					if _, ok := s.scripts[name]; ok {
						log.Println("Script:", event.Name, "removed")
						task(name, nil)
						delete(s.scripts, name)
					} else {
						log.Println("Script: could not find:", name, "to delete")
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Println(err)
		}
	}
}

func createScriptProgramFromFilename(filename string) (string, *goja.Program, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return "", nil, err
	}

	program, err := goja.Compile(filename, string(data), false)

	if err != nil {
		return "", nil, err
	}

	filename = filepath.Base(filename)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	return name, program, nil
}

type npcChatNodeType int

const (
	npcYesState npcChatNodeType = iota
	npcNoState
	npcNextState
	// npcBackState we don't use this as this condition is a pop rather than insert
	npcSelectionState
	npcStringInputState
	npcNumberInputState
	npcIncorrectState
)

type npcChatStateTracker struct {
	lastPos    int
	currentPos int

	list []npcChatNodeType

	selections []int32
	selection  int

	inputs []string
	input  int

	numbers []int32
	number  int
}

func (tracker *npcChatStateTracker) addState(stateType npcChatNodeType) {
	if tracker.currentPos >= len(tracker.list) {
		tracker.list = append(tracker.list, stateType)
	} else {
		tracker.list[tracker.currentPos] = stateType
	}

	tracker.currentPos++
	tracker.lastPos = tracker.currentPos
}

func (tracker *npcChatStateTracker) performInterrupt() bool {
	if tracker.currentPos == tracker.lastPos {
		return true
	}

	tracker.currentPos++

	return false
}

func (tracker *npcChatStateTracker) getCurrentState() npcChatNodeType {
	if tracker.currentPos >= len(tracker.list) {
		return npcIncorrectState
	}

	return tracker.list[tracker.currentPos]
}

func (tracker *npcChatStateTracker) popState() {
	tracker.lastPos--
}

type scriptPlayerWrapper struct {
	plr    *Player
	server *Server
}

func (ctrl *scriptPlayerWrapper) Warp(id int32) {
	if field, ok := ctrl.server.fields[id]; ok {
		var (
			inst *fieldInstance
			err  error
		)

		if ctrl.plr.event != nil {
			inst, err = ensureFieldInstance(field, ctrl.plr.event.instanceID, &ctrl.server.rates, ctrl.server)
		} else {
			inst, err = field.getInstance(ctrl.plr.inst.id)
			if err != nil {
				inst, err = field.getInstance(0)
			}
		}

		if err != nil {
			return
		}

		portal, err := inst.getRandomSpawnPortal()

		if err != nil {
			return
		}

		ctrl.server.warpPlayer(ctrl.plr, field, portal, true)
	}
}

func (ctrl *scriptPlayerWrapper) EventRemainingTime() int32 {
	if ctrl.plr == nil || ctrl.plr.event == nil {
		return 0
	}
	return ctrl.plr.event.RemainingTime()
}

func (ctrl *scriptPlayerWrapper) WarpToPortalName(id int32, name string) {
	if field, ok := ctrl.server.fields[id]; ok {
		var (
			inst *fieldInstance
			err  error
		)

		if ctrl.plr.event != nil {
			inst, err = ensureFieldInstance(field, ctrl.plr.event.instanceID, &ctrl.server.rates, ctrl.server)
		} else {
			inst, err = field.getInstance(ctrl.plr.inst.id)
			if err != nil {
				inst, err = field.getInstance(0)
			}
		}

		if err != nil {
			return
		}

		portal, err := inst.getPortalFromName(name)

		if err != nil {
			return
		}

		ctrl.server.warpPlayer(ctrl.plr, field, portal, true)
	}
}

func (ctrl *scriptPlayerWrapper) WarpToPortalNameInInstance(id int32, name string, instanceID int) {
	if field, ok := ctrl.server.fields[id]; ok {
		inst, err := ensureFieldInstance(field, instanceID, &ctrl.server.rates, ctrl.server)
		if err != nil {
			return
		}

		portal, err := inst.getPortalFromName(name)
		if err != nil {
			return
		}

		ctrl.server.warpPlayerToInstance(ctrl.plr, field, portal, instanceID, true)
	}
}

func (ctrl *scriptPlayerWrapper) SendMessage(msg string) {
	ctrl.plr.Send(packetMessageRedText(msg))
}

func (ctrl *scriptPlayerWrapper) Mesos() int32 {
	return ctrl.plr.mesos
}

func (ctrl *scriptPlayerWrapper) GetMesos() int32 {
	return ctrl.Mesos()
}

func (ctrl *scriptPlayerWrapper) GiveMesos(amount int32) {
	ctrl.plr.giveMesos(amount)
}

func (ctrl *scriptPlayerWrapper) GainMesos(amount int32) {
	ctrl.GiveMesos(amount)
}

func (ctrl *scriptPlayerWrapper) GiveItem(id int32, amount int16) bool {
	if amount == 0 {
		return true
	}
	if amount < 0 {
		return ctrl.plr.removeItemsByID(id, int32(-amount), false)
	}

	item, err := CreateItemFromID(id, amount)

	if err != nil {
		return false
	}

	if _, err = ctrl.plr.GiveItem(item); err != nil {
		return false
	}

	return true
}

func (ctrl *scriptPlayerWrapper) GainItem(id int32, amount int16) bool {
	return ctrl.GiveItem(id, amount)
}

func (ctrl *scriptPlayerWrapper) Job() int16 {
	return ctrl.plr.job
}

func (ctrl *scriptPlayerWrapper) SetJob(id int16) {
	ctrl.plr.setJob(id)
}

func (ctrl *scriptPlayerWrapper) Level() byte {
	return ctrl.plr.level
}

func (ctrl *scriptPlayerWrapper) PartnerID() int32 {
	return ctrl.plr.partnerID
}

func (ctrl *scriptPlayerWrapper) MarriageItemID() int32 {
	return ctrl.plr.marriageItemID
}

func (ctrl *scriptPlayerWrapper) IsMarried() bool {
	return ctrl.plr.married()
}

func (ctrl *scriptPlayerWrapper) HasEngagementBox() bool {
	for _, id := range []int32{constant.ItemEngagementBoxMoonstone, constant.ItemEngagementBoxStar, constant.ItemEngagementBoxGolden, constant.ItemEngagementBoxSilver} {
		if ctrl.plr.countItem(id) > 0 {
			return true
		}
	}
	return false
}

func (ctrl *scriptPlayerWrapper) PartnerName() string {
	if ctrl.plr.partnerID <= 0 {
		return ""
	}
	if partner, err := ctrl.server.players.GetFromID(ctrl.plr.partnerID); err == nil {
		return partner.Name
	}
	var name string
	_ = common.DB.QueryRow("SELECT name FROM characters WHERE id=?", ctrl.plr.partnerID).Scan(&name)
	return name
}

func (ctrl *scriptPlayerWrapper) ReserveWedding(cathedral, premium bool) bool {
	return ctrl.server.reserveWedding(ctrl.plr, cathedral, premium) == nil
}

func (ctrl *scriptPlayerWrapper) StartWedding(cathedral bool) bool {
	return ctrl.server.startWedding(ctrl.plr, cathedral) == nil
}

func (ctrl *scriptPlayerWrapper) AdvanceWeddingCeremony(cathedral bool) bool {
	return ctrl.server.advanceWeddingCeremony(ctrl.plr, cathedral) == nil
}

func (ctrl *scriptPlayerWrapper) EnterWeddingAsGuest(cathedral bool) bool {
	return ctrl.server.enterWeddingAsGuest(ctrl.plr, cathedral) == nil
}

func (ctrl *scriptPlayerWrapper) HasWeddingReservation(cathedral bool) bool {
	return ctrl.server.currentWeddingReservation(ctrl.plr, cathedral) != nil
}

func (ctrl *scriptPlayerWrapper) WeddingStarted(cathedral bool) bool {
	res := ctrl.server.currentWeddingReservation(ctrl.plr, cathedral)
	return res != nil && res.Started
}

func (ctrl *scriptPlayerWrapper) WeddingIsPremium() bool {
	res := ctrl.server.currentWeddingReservationAny(ctrl.plr)
	return res != nil && res.Premium
}

func (ctrl *scriptPlayerWrapper) WeddingGuestTicket(cathedral bool) int32 {
	return weddingGuestTicket(cathedral)
}

func (ctrl *scriptPlayerWrapper) WeddingReservationTicket(cathedral, premium bool) int32 {
	return weddingReservationTicket(cathedral, premium)
}

func (ctrl *scriptPlayerWrapper) WeddingInviteItem(cathedral bool) int32 {
	return weddingInviteItem(cathedral)
}

func (ctrl *scriptPlayerWrapper) WeddingStage(cathedral bool) int {
	res := ctrl.server.currentWeddingReservationAny(ctrl.plr)
	if res != nil && res.Cathedral != cathedral {
		res = nil
	}
	if res == nil {
		return -1
	}
	return int(res.Stage)
}

func (ctrl *scriptPlayerWrapper) BreakMarriage(itemID int32) {
	ctrl.server.breakMarriageState(ctrl.plr, itemID)
}

func (ctrl *scriptPlayerWrapper) CompleteWedding(cathedral bool) bool {
	return ctrl.server.completeWedding(ctrl.plr, cathedral) == nil
}

func (ctrl *scriptPlayerWrapper) UnderMarriageCooldown() bool {
	return ctrl.plr.underMarriageCooldown()
}

func (ctrl *scriptPlayerWrapper) StartWeddingAfterParty() bool {
	return ctrl.server.startWeddingAfterParty(ctrl.plr) == nil
}

func (ctrl *scriptPlayerWrapper) ClaimWeddingExitReward() int {
	return ctrl.server.claimWeddingExit(ctrl.plr)
}

func (ctrl *scriptPlayerWrapper) InGuild() bool {
	return ctrl.plr.guild != nil
}

func (ctrl *scriptPlayerWrapper) GuildRank() byte {
	if ctrl.plr.guild != nil {
		for i, id := range ctrl.plr.guild.playerID {
			if id == ctrl.plr.ID {
				return ctrl.plr.guild.ranks[i]
			}
		}
	}

	return 0
}

func (ctrl *scriptPlayerWrapper) InParty() bool {
	return ctrl.plr.party != nil
}

func (ctrl *scriptPlayerWrapper) PartyQuestActive() bool {
	if ctrl == nil || ctrl.plr == nil || ctrl.plr.party == nil || ctrl.server == nil {
		return false
	}
	_, ok := ctrl.server.events[ctrl.plr.party.ID]
	return ok
}

func (ctrl *scriptPlayerWrapper) IsPartyLeader() bool {
	if ctrl.InParty() {
		return ctrl.plr.party.players[0] == ctrl.plr
	}

	return false
}

func (ctrl *scriptPlayerWrapper) IsLeader() bool {
	return ctrl.IsPartyLeader()
}

func (ctrl *scriptPlayerWrapper) PartyMembersOnMapCount() int {
	if !ctrl.InParty() {
		return 0
	}

	count := 0
	for _, v := range ctrl.plr.party.players {
		if v != nil && v.mapID == ctrl.plr.mapID {
			count++
		}
	}

	return count
}

func (ctrl *scriptPlayerWrapper) PartyMembersOnMap() []scriptPlayerWrapper {
	if !ctrl.InParty() {
		return []scriptPlayerWrapper{}
	}

	members := make([]scriptPlayerWrapper, 0, constant.MaxPartySize)

	for _, v := range ctrl.plr.party.players {
		if v != nil && v.mapID == ctrl.plr.mapID && v.inst != nil && ctrl.plr.inst != nil && v.inst.id == ctrl.plr.inst.id {
			members = append(members, scriptPlayerWrapper{plr: v, server: ctrl.server})
		}
	}

	return members
}

func (ctrl *scriptPlayerWrapper) LogEvent(msg string) {
	if ctrl == nil || ctrl.plr == nil || ctrl.plr.event == nil {
		return
	}
	log.Printf("event[%d] player[%d:%s] %s", ctrl.plr.event.id, ctrl.plr.ID, ctrl.plr.Name, msg)
}

func (ctrl *scriptPlayerWrapper) EventMembersOnMap(id int32) bool {
	if ctrl.plr.event == nil {
		return false
	}

	return ctrl.plr.event.IsParticipantsOnMap(id)
}

func (ctrl *scriptPlayerWrapper) GetEventProperty(key string) interface{} {
	if ctrl.plr == nil {
		return nil
	}
	if ctrl.plr.event != nil {
		if value, ok := ctrl.plr.event.properties[key]; ok {
			return value
		}
		return nil
	}
	if ctrl.plr.inst == nil {
		return nil
	}
	if value, ok := ctrl.plr.inst.properties[key]; ok {
		return value
	}
	return nil
}

func (ctrl *scriptPlayerWrapper) SetEventProperty(key string, value interface{}) interface{} {
	if ctrl.plr == nil {
		return nil
	}
	prev := ctrl.GetEventProperty(key)
	if ctrl.plr.event != nil {
		ctrl.plr.event.properties[key] = value
		return prev
	}
	if ctrl.plr.inst == nil {
		return nil
	}
	ctrl.plr.inst.properties[key] = value
	return prev
}

func (ctrl *scriptPlayerWrapper) FinishEvent() {
	if ctrl.plr.event != nil {
		ctrl.plr.event.Finished()
	}
}

func (ctrl *scriptPlayerWrapper) LeaveEvent() {
	if ctrl.plr.event != nil && ctrl.plr.event.playerLeaveEventCallback != nil {
		ctrl.plr.event.playerLeaveEventCallback(*ctrl)
	}
}

func (ctrl *scriptPlayerWrapper) WarpEventMembers(id int32) {
	if ctrl.plr.event != nil {
		ctrl.plr.event.WarpPlayers(id)
	}
}

func (ctrl *scriptPlayerWrapper) WarpEventMembersToPortal(id int32, portalName string) {
	if ctrl.plr.event != nil {
		ctrl.plr.event.WarpPlayersToPortal(id, portalName)
	}
}

func (ctrl *scriptPlayerWrapper) CountMonster() int {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return 0
	}
	return ctrl.plr.inst.lifePool.mobCount()
}

func (ctrl *scriptPlayerWrapper) SpawnMonster(id int32, x int16, y int16) {
	if ctrl.plr == nil || ctrl.plr.inst == nil {
		return
	}
	ctrl.plr.inst.lifePool.spawnMobFromID(id, newPos(x, y, 0), false, true, true, constant.MobSummonTypeInstant, 0)
}

func (ctrl *scriptPlayerWrapper) PartyGiveExp(val int32) {
	if !ctrl.InParty() {
		return
	}

	for _, plr := range ctrl.plr.party.players {
		if plr != nil {
			plr.giveEXP(val, false, false)
		}
	}
}

func (ctrl *scriptPlayerWrapper) DisbandGuild() {
	if ctrl.plr.guild == nil {
		return
	}

	ctrl.server.world.Send(internal.PacketGuildDisband(ctrl.plr.guild.id))
}

func (ctrl *scriptPlayerWrapper) GainGuildPoints(points int32) {
	if ctrl.plr.guild == nil {
		return
	}
	ctrl.plr.Send(packetMessageGuildPointsChange(points))
	ctrl.server.world.Send(internal.PacketGuildPointsUpdate(ctrl.plr.guild.id, ctrl.plr.guild.points+points))
}

func (ctrl *scriptPlayerWrapper) GetLevel() int {
	return int(ctrl.plr.level)
}

func (ctrl *scriptPlayerWrapper) GetQuestStatus(id int16) int {
	// 2 = completed, 1 = in progress, 0 = not started
	for _, q := range ctrl.plr.quests.completed {
		if q.id == id {
			return 2
		}
	}
	for _, q := range ctrl.plr.quests.inProgress {
		if q.id == id {
			return 1
		}
	}
	return 0
}

func (ctrl *scriptPlayerWrapper) CheckQuestStatus(id int16, status int) bool {
	return ctrl.GetQuestStatus(id) == status
}

func (ctrl *scriptPlayerWrapper) QuestStarted(id int16) bool {
	return ctrl.GetQuestStatus(id) == 1
}

func (ctrl *scriptPlayerWrapper) QuestCompleted(id int16) bool {
	return ctrl.GetQuestStatus(id) == 2
}

func (ctrl *scriptPlayerWrapper) QuestNotStarted(id int16) bool {
	return ctrl.GetQuestStatus(id) == 0
}

func (ctrl *scriptPlayerWrapper) QuestData(id int16) string {
	if q, ok := ctrl.plr.quests.inProgress[id]; ok {
		return q.name
	}
	return ""
}

func (ctrl *scriptPlayerWrapper) CheckQuestData(id int16, data string) bool {
	return ctrl.QuestData(id) == data
}

func (ctrl *scriptPlayerWrapper) SetQuestData(id int16, data string) {
	// Only allow setting data if quest is in-progress; if not, start it or ignore.
	if _, ok := ctrl.plr.quests.inProgress[id]; !ok {
		// You may choose to implicitly start; here we upsert and add in-memory if needed.
		ctrl.plr.quests.add(id, data)
	} else {
		// update in-memory
		q := ctrl.plr.quests.inProgress[id]
		q.name = data
		ctrl.plr.quests.inProgress[id] = q
	}
	// Persist + notify client
	upsertQuestRecord(ctrl.plr.ID, id, data)
	ctrl.plr.Send(packetQuestUpdate(id, data))
}

func (ctrl *scriptPlayerWrapper) StartQuest(id int16) bool {
	return ctrl.plr.tryStartQuest(id)
}

func (ctrl *scriptPlayerWrapper) CompleteQuest(id int16) bool {
	return ctrl.plr.tryCompleteQuest(id)
}

func (ctrl *scriptPlayerWrapper) ForfeitQuest(id int16) {
	if !ctrl.plr.quests.hasInProgress(id) {
		return
	}
	ctrl.plr.quests.remove(id)
	delete(ctrl.plr.quests.mobKills, id)
	deleteQuest(ctrl.plr.ID, id)
	ctrl.plr.Send(packetQuestRemove(id))
	clearQuestMobKills(ctrl.plr.ID, id)
}

func (ctrl *scriptPlayerWrapper) TakeItem(id int32, slot int16, amount int16, invID byte) bool {
	_, err := ctrl.plr.takeItem(id, slot, amount, invID)
	return err == nil
}

func (ctrl *scriptPlayerWrapper) RemoveItemsByID(id int32, count int32) bool {
	return ctrl.plr.removeItemsByID(id, count, false)
}

func (ctrl *scriptPlayerWrapper) RemoveItemsByIDSilent(id int32, count int32) bool {
	return ctrl.plr.removeItemsByID(id, count, true)
}

func (ctrl *scriptPlayerWrapper) RemoveAll(id int32) bool {
	count := ctrl.plr.countItem(id)
	if count <= 0 {
		return true
	}
	return ctrl.plr.removeItemsByID(id, count, true)
}

func (ctrl *scriptPlayerWrapper) ItemCount(id int32) int32 {
	return ctrl.plr.countItem(id)
}

func (ctrl *scriptPlayerWrapper) HaveItem(id int32, quantity int32) bool {
	if quantity <= 0 {
		quantity = 1
	}
	return ctrl.plr.countItem(id) >= quantity
}

func (ctrl *scriptPlayerWrapper) IsWearingItem(id int32) bool {
	for i := range ctrl.plr.equip {
		if ctrl.plr.equip[i].ID == id && ctrl.plr.equip[i].slotID < 0 {
			return true
		}
	}
	return false
}

func (ctrl *scriptPlayerWrapper) CanHold(id int32, amount int16) bool {
	if amount <= 0 {
		return true
	}
	item, err := CreateItemFromID(id, amount)
	if err != nil {
		return false
	}
	return ctrl.plr.CanReceiveItems([]Item{item})
}

func (ctrl *scriptPlayerWrapper) CanHoldAll(items [][]int32) bool {
	converted := make([]Item, 0, len(items))
	for _, spec := range items {
		if len(spec) == 0 {
			continue
		}
		id := spec[0]
		amount := int16(1)
		if len(spec) > 1 {
			amount = int16(spec[1])
		}
		if amount <= 0 {
			continue
		}
		item, err := CreateItemFromID(id, amount)
		if err != nil {
			return false
		}
		converted = append(converted, item)
	}
	return ctrl.plr.CanReceiveItems(converted)
}

func (ctrl *scriptPlayerWrapper) GetEquipInventoryFreeSlot() int16 {
	return int16(ctrl.plr.equipSlotSize) - int16(len(ctrl.plr.equip))
}

func (ctrl *scriptPlayerWrapper) GetUseInventoryFreeSlot() int16 {
	return int16(ctrl.plr.useSlotSize) - int16(len(ctrl.plr.use))
}

func (ctrl *scriptPlayerWrapper) GetSetupInventoryFreeSlot() int16 {
	return int16(ctrl.plr.setupSlotSize) - int16(len(ctrl.plr.setUp))
}

func (ctrl *scriptPlayerWrapper) GetEtcInventoryFreeSlot() int16 {
	return int16(ctrl.plr.etcSlotSize) - int16(len(ctrl.plr.etc))
}

func (ctrl *scriptPlayerWrapper) GetCashInventoryFreeSlot() int16 {
	return int16(ctrl.plr.cashSlotSize) - int16(len(ctrl.plr.cash))
}

func (ctrl *scriptPlayerWrapper) TakeMesos(amount int32) {
	ctrl.plr.takeMesos(amount)
}

func (ctrl *scriptPlayerWrapper) IsGM() bool {
	return ctrl.plr.admin()
}

func (ctrl *scriptPlayerWrapper) GetNX() int32 {
	return ctrl.plr.GetNX()
}

func (ctrl *scriptPlayerWrapper) SetNX(nx int32) {
	ctrl.plr.SetNX(nx)
}

func (ctrl *scriptPlayerWrapper) GetMaplePoints() int32 {
	return ctrl.plr.GetMaplePoints()
}

func (ctrl *scriptPlayerWrapper) SetMaplePoints(points int32) {
	ctrl.plr.SetMaplePoints(points)
}

func (ctrl *scriptPlayerWrapper) SetFame(value int16) {
	ctrl.plr.setFame(value)
}

func (ctrl *scriptPlayerWrapper) GiveFame(delta int16) {
	ctrl.plr.setFame(ctrl.plr.fame + delta)
}

func (ctrl *scriptPlayerWrapper) GiveAP(amount int16) {
	ctrl.plr.giveAP(amount)
}

func (ctrl *scriptPlayerWrapper) GetRemainingAP() int16 {
	return ctrl.plr.ap
}

func (ctrl *scriptPlayerWrapper) SetRemainingAP(amount int16) {
	ctrl.plr.setAP(amount)
}

func (ctrl *scriptPlayerWrapper) GetRemainingSP() int16 {
	return ctrl.plr.sp
}

func (ctrl *scriptPlayerWrapper) GiveSP(amount int16) {
	ctrl.plr.giveSP(amount)
}

func (ctrl *scriptPlayerWrapper) GiveEXP(amount int32) {
	ctrl.plr.giveEXP(amount, false, false)
}

func (ctrl *scriptPlayerWrapper) GiveHP(amount int16) {
	ctrl.plr.giveHP(amount)
}

func (ctrl *scriptPlayerWrapper) GiveMP(amount int16) {
	ctrl.plr.giveMP(amount)
}

func (ctrl *scriptPlayerWrapper) HealToFull() {
	ctrl.plr.setHP(ctrl.plr.maxHP)
	ctrl.plr.setMP(ctrl.plr.maxMP)
}

func (ctrl *scriptPlayerWrapper) SetHP(amount int16) {
	ctrl.plr.setHP(amount)
}

func (ctrl *scriptPlayerWrapper) GetStr() int16 {
	return ctrl.plr.str
}

func (ctrl *scriptPlayerWrapper) SetStr(amount int16) {
	ctrl.plr.setStr(amount)
}

func (ctrl *scriptPlayerWrapper) GetDex() int16 {
	return ctrl.plr.dex
}

func (ctrl *scriptPlayerWrapper) SetDex(amount int16) {
	ctrl.plr.setDex(amount)
}

func (ctrl *scriptPlayerWrapper) GetInt() int16 {
	return ctrl.plr.intt
}

func (ctrl *scriptPlayerWrapper) SetInt(amount int16) {
	ctrl.plr.setInt(amount)
}

func (ctrl *scriptPlayerWrapper) GetLuk() int16 {
	return ctrl.plr.luk
}

func (ctrl *scriptPlayerWrapper) SetLuk(amount int16) {
	ctrl.plr.setLuk(amount)
}

func (ctrl *scriptPlayerWrapper) Gender() byte {
	return ctrl.plr.gender
}

func (ctrl *scriptPlayerWrapper) GetGender() byte {
	return ctrl.Gender()
}

// Hair returns the current hair ID
func (ctrl *scriptPlayerWrapper) Hair() int32 {
	return ctrl.plr.hair
}

func (ctrl *scriptPlayerWrapper) GetHair() int32 {
	return ctrl.Hair()
}

// SetHair updates the player's hair and refreshes the client appearance
func (ctrl *scriptPlayerWrapper) SetHair(id int32) {
	if ctrl.plr.hair == id {
		return
	}
	err := ctrl.plr.setHair(id)
	if err != nil {
		return
	}
}

func (ctrl *scriptPlayerWrapper) Face() int32 {
	return ctrl.plr.face
}

func (ctrl *scriptPlayerWrapper) GetFace() int32 {
	return ctrl.Face()
}

func (ctrl *scriptPlayerWrapper) SetFace(id int32) {
	if ctrl.plr.face == id {
		return
	}
	err := ctrl.plr.setFace(id)
	if err != nil {
		return
	}
}

// Skin returns the current skin tone (0..n)
func (ctrl *scriptPlayerWrapper) Skin() byte {
	return ctrl.plr.skin
}

// SetSkinColor updates the player's skin tone and refreshes the client appearance
func (ctrl *scriptPlayerWrapper) SetSkinColor(skin byte) {
	if ctrl.plr.skin == skin {
		return
	}
	err := ctrl.plr.setSkin(skin)
	if err != nil {
		return
	}
}

func (ctrl *scriptPlayerWrapper) SetSkin(skin byte) {
	ctrl.SetSkinColor(skin)
}

func (ctrl *scriptPlayerWrapper) GetSkin() byte {
	return ctrl.Skin()
}

type scriptQuestView struct {
	Data   string `json:"data"`
	Status int    `json:"status"` // 0,1,2 same as GetQuestStatus
}

func (ctrl *scriptPlayerWrapper) Quest(id int16) scriptQuestView {
	status := ctrl.GetQuestStatus(id) // already 0/1/2
	return scriptQuestView{
		Data:   ctrl.QuestData(id),
		Status: status,
	}
}

func (ctrl *scriptPlayerWrapper) PreviousMap() int32 {
	return ctrl.plr.previousMap
}

func (ctrl *scriptPlayerWrapper) SaveLocation(slot string) {
	if ctrl.plr.savedMaps == nil {
		ctrl.plr.savedMaps = make(map[string]int32)
	}
	ctrl.plr.savedMaps[strings.ToUpper(slot)] = ctrl.plr.mapID
}

func (ctrl *scriptPlayerWrapper) GetSavedLocation(slot string) int32 {
	if ctrl.plr.savedMaps == nil {
		return -1
	}
	if id, ok := ctrl.plr.savedMaps[strings.ToUpper(slot)]; ok {
		return id
	}
	return -1
}

func (ctrl *scriptPlayerWrapper) ClearSavedLocation(slot string) {
	if ctrl.plr.savedMaps == nil {
		return
	}
	delete(ctrl.plr.savedMaps, strings.ToUpper(slot))
}

func (ctrl *scriptPlayerWrapper) MapID() int32 {
	return ctrl.plr.mapID
}

func (ctrl *scriptPlayerWrapper) Position() map[string]int16 {
	return map[string]int16{
		"x": ctrl.plr.pos.x,
		"y": ctrl.plr.pos.y,
	}
}

func (ctrl *scriptPlayerWrapper) Name() string {
	return ctrl.plr.Name
}

func (ctrl *scriptPlayerWrapper) InventoryExchange(itemSource int32, srcCount int32, itemExchangeFor int32, count int16) bool {
	if !ctrl.plr.removeItemsByID(itemSource, srcCount, false) {
		return false
	}

	item, err := CreateItemFromID(itemExchangeFor, count)
	if err != nil {
		return false
	}
	if _, err = ctrl.plr.GiveItem(item); err != nil {
		return false
	}
	return true
}

func (ctrl *scriptPlayerWrapper) ShowCountdown(seconds int32) {
	ctrl.plr.Send(packetShowCountdown(seconds))
}

func (ctrl *scriptPlayerWrapper) HideCountdown() {
	ctrl.plr.Send(packetHideCountdown())
}

func (ctrl *scriptPlayerWrapper) ShowNpcOk(npcID int32, msg string) {
	ctrl.plr.Send(packetNpcChatOk(npcID, msg))
}

func (ctrl *scriptPlayerWrapper) PortalEffect(path string) {
	ctrl.plr.Send(packetPortalEffectt(2, path))
}

func (ctrl *scriptPlayerWrapper) StartPartyQuest(name string, instID int) {
	if ctrl.plr.party == nil {
		return
	}

	if _, ok := ctrl.server.events[ctrl.plr.party.ID]; ok {
		return
	}

	if instID <= 0 {
		field, ok := ctrl.server.fields[ctrl.plr.mapID]
		if !ok {
			return
		}
		instID = field.createInstance(&ctrl.server.rates, ctrl.server)
	}

	program, ok := ctrl.server.eventScriptStore.scripts[name]

	if !ok {
		return
	}

	ids := []int32{}

	if ctrl.plr.party != nil {
		for i, id := range ctrl.plr.party.PlayerID {
			if ctrl.plr.mapID == ctrl.plr.party.MapID[i] && ctrl.plr.party.players[i] != nil && ctrl.plr.inst != nil && ctrl.plr.party.players[i].inst != nil {
				if ctrl.plr.inst.id == ctrl.plr.party.players[i].inst.id {
					ids = append(ids, id)
				}
			}
		}
	} else {
		ids = append(ids, ctrl.plr.ID)
	}

	if len(ids) == 0 {
		return
	}

	event, err := createEvent(ctrl.plr.party.ID, instID, ids, ctrl.server, program)

	if err != nil {
		log.Println(err)
		return
	}

	ctrl.server.events[ctrl.plr.party.ID] = event
	event.start()
}

func (ctrl *scriptPlayerWrapper) StartGuildQuest(name string, instID int) {
	if ctrl.plr.guild == nil {
		return
	}
	program, ok := ctrl.server.eventScriptStore.scripts[name]
	if !ok {
		return
	}
	ids := []int32{}
	ctrl.server.players.observe(func(other *Player) {
		if other == nil || other.guild == nil || other.guild.id != ctrl.plr.guild.id {
			return
		}
		if other.mapID != ctrl.plr.mapID || other.inst == nil || ctrl.plr.inst == nil || other.inst.id != ctrl.plr.inst.id {
			return
		}
		ids = append(ids, other.ID)
	})
	event, err := createEvent(int32(ctrl.plr.guild.id), instID, ids, ctrl.server, program)
	if err != nil {
		log.Println(err)
		return
	}
	ctrl.server.events[int32(ctrl.plr.guild.id)] = event
	event.start()
}

func (ctrl *scriptPlayerWrapper) JoinGuildQuest() bool {
	if ctrl.plr.guild == nil {
		return false
	}
	event, ok := ctrl.server.events[int32(ctrl.plr.guild.id)]
	if !ok {
		return false
	}
	event.AddPlayer(*ctrl)
	return true
}

func (ctrl *scriptPlayerWrapper) LeavePartyQuest() {
	if ctrl.plr.party == nil {
		return
	}

	if event, ok := ctrl.server.events[ctrl.plr.party.ID]; ok {
		event.playerLeaveEventCallback(scriptPlayerWrapper{plr: ctrl.plr, server: ctrl.server})
	}
}

type scriptMapWrapper struct {
	inst   *fieldInstance
	server *Server
}

func (ctrl *scriptMapWrapper) PlayerCount(mapID int32, instID int) int {
	f, ok := ctrl.server.fields[mapID]

	if !ok {
		return 0
	}

	inst, err := f.getInstance(instID)

	if err != nil {
		return 0
	}

	return len(inst.players)
}

func (ctrl *scriptMapWrapper) PlayerCountInMap(mapID int32) int {
	if ctrl == nil || ctrl.inst == nil {
		return 0
	}
	return ctrl.PlayerCount(mapID, ctrl.inst.id)
}

func (ctrl *scriptMapWrapper) PlaySound(path string) {
	ctrl.inst.send(packetPlaySound(path))
}

func (ctrl *scriptMapWrapper) ShowEffect(path string) {
	ctrl.inst.send(packetShowScreenEffect(path))
}

func (ctrl *scriptMapWrapper) PortalEffect(path string) {
	ctrl.inst.send(packetPortalEffectt(2, path))
}

func (ctrl *scriptMapWrapper) Message(msg string) {
	ctrl.inst.send(packetMessageNotice(msg))
}

func (ctrl *scriptMapWrapper) PortalEnabled(enable bool, name string) {
	ctrl.inst.setPortalEnabled(name, enable)
}

func (ctrl *scriptMapWrapper) SetPortalScript(name string, script string) {
	for i := range ctrl.inst.portals {
		if ctrl.inst.portals[i].name == name {
			ctrl.inst.portals[i].script = script
			return
		}
	}
}

func (ctrl *scriptMapWrapper) SetPortalScriptByID(id int, script string) {
	if id >= 0 && id < len(ctrl.inst.portals) {
		ctrl.inst.portals[id].script = script
	}
}

func (ctrl *scriptMapWrapper) IsPortalEnabled(name string) bool {
	return ctrl.inst.getPortalEnabled(name)
}

func (ctrl *scriptMapWrapper) Properties() map[string]interface{} {
	return ctrl.inst.properties
}

func (ctrl *scriptMapWrapper) ClearProperties() {
	for k := range ctrl.inst.properties {
		delete(ctrl.inst.properties, k)
	}
}

func (ctrl *scriptMapWrapper) PlayersInArea(id int) int {
	areas := nx.GetMaps()[ctrl.inst.fieldID].Areas
	count := 0

	for _, plr := range ctrl.inst.players {
		if areas[id].Inside(plr.pos.x, plr.pos.y) {
			count++
		}

	}

	return count
}

func (ctrl *scriptMapWrapper) MalePlayersInArea(id int) int {
	areas := nx.GetMaps()[ctrl.inst.fieldID].Areas
	if id < 0 || id >= len(areas) {
		return 0
	}
	count := 0
	for _, plr := range ctrl.inst.players {
		if plr.gender == 0 && areas[id].Inside(plr.pos.x, plr.pos.y) {
			count++
		}
	}
	return count
}

func (ctrl *scriptMapWrapper) FemalePlayersInArea(id int) int {
	areas := nx.GetMaps()[ctrl.inst.fieldID].Areas
	if id < 0 || id >= len(areas) {
		return 0
	}
	count := 0
	for _, plr := range ctrl.inst.players {
		if plr.gender == 1 && areas[id].Inside(plr.pos.x, plr.pos.y) {
			count++
		}
	}
	return count
}

func (ctrl *scriptMapWrapper) GroundItemsInArea(id int) []int32 {
	areas := nx.GetMaps()[ctrl.inst.fieldID].Areas
	if id < 0 || id >= len(areas) {
		return []int32{}
	}
	out := make([]int32, 0)
	for _, drop := range ctrl.inst.dropPool.drops {
		if drop.mesos > 0 {
			continue
		}
		if areas[id].Inside(drop.finalPos.x, drop.finalPos.y) {
			out = append(out, drop.item.ID)
		}
	}
	return out
}

func (ctrl *scriptMapWrapper) MobCount() int {
	return ctrl.inst.lifePool.mobCount()
}

func (ctrl *scriptMapWrapper) MobCountByID(id int32) int {
	return ctrl.inst.lifePool.mobCountByTemplate(id)
}

func (ctrl *scriptMapWrapper) RemoveAllMobs() {
	ctrl.inst.lifePool.eraseMobs()
}

func (ctrl *scriptMapWrapper) RemoveMobsByID(id int32) {
	ctrl.inst.lifePool.removeMobsByTemplate(id)
}

func (ctrl *scriptMapWrapper) SetMobSpawnEnabled(id int32, enable bool) {
	ctrl.inst.lifePool.setMobSpawnEnabled(id, enable)
}

func (ctrl *scriptMapWrapper) RemoveDrops() {
	ctrl.inst.dropPool.clearDrops()
}

func (ctrl *scriptMapWrapper) Reset() {
	ctrl.inst.lifePool.eraseMobs()
	ctrl.inst.lifePool.attemptMobSpawn(true)
	ctrl.inst.dropPool.eraseDrops()
	ctrl.inst.reactorPool.reset(false)
}

func (ctrl *scriptMapWrapper) HitReactorByName(name string) bool {
	for _, r := range ctrl.inst.reactorPool.reactors {
		if r.name == name {
			ctrl.inst.reactorPool.triggerHit(r.spawnID, 0, ctrl.server, nil)
			return true
		}
	}
	return false
}

func (ctrl *scriptMapWrapper) ReactorNames() []string {
	out := make([]string, 0, len(ctrl.inst.reactorPool.reactors))
	for _, r := range ctrl.inst.reactorPool.reactors {
		out = append(out, r.name)
	}
	return out
}

func (ctrl *scriptMapWrapper) ReactorNamesExcluding(exclude string) []string {
	out := make([]string, 0, len(ctrl.inst.reactorPool.reactors))
	for _, r := range ctrl.inst.reactorPool.reactors {
		if r.name != exclude {
			out = append(out, r.name)
		}
	}
	return out
}

func (ctrl *scriptMapWrapper) ReactorStateByName(name string) int {
	for _, r := range ctrl.inst.reactorPool.reactors {
		if r.name == name {
			return int(r.state)
		}
	}
	return -1
}

func (ctrl *scriptMapWrapper) SetReactorStateByName(name string, state int) bool {
	for _, r := range ctrl.inst.reactorPool.reactors {
		if r.name == name {
			r.state = byte(state)
			r.frameDelay = 0
			ctrl.inst.send(packetMapReactorChangeState(r.spawnID, r.state, r.pos.x, r.pos.y, r.frameDelay, r.faceLeft, 0))
			return true
		}
	}
	return false
}

func (ctrl *scriptMapWrapper) HitReactorByTemplate(id int32) bool {
	for _, r := range ctrl.inst.reactorPool.reactors {
		if r.templateID == id {
			ctrl.inst.reactorPool.triggerHit(r.spawnID, 0, ctrl.server, nil)
			return true
		}
	}
	return false
}

func (ctrl *scriptMapWrapper) SpawnNpc(id int32, x int16, y int16) bool {
	npcData := nx.Life{ID: id, Type: "n", X: x, Y: y, FaceLeft: false, Foothold: 0}
	spawnID, err := ctrl.inst.lifePool.nextNpcID()
	if err != nil {
		return false
	}
	val := createNpcFromData(spawnID, npcData)
	ctrl.inst.lifePool.npcs[spawnID] = &val
	ctrl.inst.send(packetNpcShow(&val))
	return true
}

func (ctrl *scriptMapWrapper) RemoveNpcByTemplate(id int32) {
	removeIDs := make([]int32, 0)
	for spawnID, n := range ctrl.inst.lifePool.npcs {
		if n != nil && n.id == id {
			if n.controller != nil {
				n.removeController()
			}
			ctrl.inst.send(packetNpcRemove(n.spawnID))
			removeIDs = append(removeIDs, spawnID)
		}
	}
	for _, spawnID := range removeIDs {
		delete(ctrl.inst.lifePool.npcs, spawnID)
	}
}

func (ctrl *scriptMapWrapper) RevealReactorsByName(names []string, startDelayMs int32, stepDelayMs int32) {
	ordered := make([]*fieldReactor, 0, len(names))
	for _, name := range names {
		for _, r := range ctrl.inst.reactorPool.reactors {
			if r.name == name {
				ordered = append(ordered, r)
				break
			}
		}
	}
	for i, r := range ordered {
		delay := time.Duration(startDelayMs+int32(i)*stepDelayMs) * time.Millisecond
		time.AfterFunc(delay, func(rr *fieldReactor) func() {
			return func() {
				ctrl.server.dispatch <- func() {
					ctrl.inst.reactorPool.triggerHit(rr.spawnID, 0, ctrl.server, nil)
				}
			}
		}(r))
	}
}

func (ctrl *scriptMapWrapper) GetMap(id int32, instID int) scriptMapWrapper {
	if field, ok := ctrl.server.fields[id]; ok {
		inst, err := field.getInstance(instID)

		if err != nil {
			instID = field.createInstance(&ctrl.server.rates, ctrl.server)
			inst, err = field.getInstance(instID)

			if err != nil {
				return scriptMapWrapper{}
			}

			return scriptMapWrapper{inst: inst, server: ctrl.server}
		}

		return scriptMapWrapper{inst: inst, server: ctrl.server}
	}

	return scriptMapWrapper{}
}

func (m *scriptMapWrapper) GetMapID() int32 {
	if m == nil || m.inst == nil {
		return 0
	}
	return m.inst.fieldID
}

type npcChatController struct {
	npcID int32
	conn  mnet.Client

	goods [][]int32

	stateTracker npcChatStateTracker

	vm      *goja.Runtime
	program *goja.Program

	selectionCalls int
	disposed       bool
}

type portalScriptController struct {
	plr     *Player
	server  *Server
	portal  portal
	warped  bool
	blocked bool
}

type reactorScriptController struct {
	server  *Server
	inst    *fieldInstance
	reactor *fieldReactor
}

func (ctrl *reactorScriptController) MapMessage(msgType int, msg string) {
	ctrl.inst.send(packetMessageNotice(msg))
}

func (ctrl *reactorScriptController) ShowEffect(path string) {
	ctrl.inst.send(packetShowScreenEffect(path))
}

func (ctrl *reactorScriptController) PlaySound(path string) {
	ctrl.inst.send(packetPlaySound(path))
}

func (ctrl *reactorScriptController) SpawnNpc(id int32, x int16, y int16) bool {
	npcData := nx.Life{ID: id, Type: "n", X: x, Y: y, FaceLeft: false, Foothold: 0}
	spawnID, err := ctrl.inst.lifePool.nextNpcID()
	if err != nil {
		return false
	}
	val := createNpcFromData(spawnID, npcData)
	ctrl.inst.lifePool.npcs[spawnID] = &val
	ctrl.inst.send(packetNpcShow(&val))
	return true
}

func (ctrl *reactorScriptController) SpawnNpcAtReactor(id int32) bool {
	return ctrl.SpawnNpc(id, ctrl.reactor.pos.x, ctrl.reactor.pos.y)
}

func (ctrl *reactorScriptController) SpawnMonster(id int32, x int16, y int16) {
	ctrl.inst.lifePool.spawnMobFromID(id, newPos(x, y, 0), false, true, true, constant.MobSummonTypeInstant, 0)
}

func (ctrl *reactorScriptController) GainGuildPoints(points int32) {
	if len(ctrl.inst.players) == 0 {
		return
	}
	for _, plr := range ctrl.inst.players {
		if plr != nil && plr.guild != nil {
			plr.Send(packetMessageGuildPointsChange(points))
			ctrl.server.world.Send(internal.PacketGuildPointsUpdate(plr.guild.id, plr.guild.points+points))
			return
		}
	}
}

func (ctrl *reactorScriptController) HitMapReactorByName(mapID int32, name string) bool {
	field, ok := ctrl.server.fields[mapID]
	if !ok {
		return false
	}
	inst, err := field.getInstance(ctrl.inst.id)
	if err != nil {
		inst, err = field.getInstance(0)
		if err != nil {
			return false
		}
	}
	m := scriptMapWrapper{inst: inst, server: ctrl.server}
	return m.HitReactorByName(name)
}

func (ctrl *reactorScriptController) DropItems(args ...int32) {
	mesos, items := buildDropRewards(ctrl.inst.lifePool.rNumber, reactorDropTable[ctrl.reactor.info.ID], ctrl.inst.dropPool.rates.drop, nil)
	if mesos > 0 || len(items) > 0 {
		ctrl.inst.dropPool.createDrop(dropSpawnNormal, dropFreeForAll, mesos, ctrl.reactor.pos, true, false, 0, 0, items...)
	}
}

func (ctrl *reactorScriptController) SetMapMobSpawnEnabled(mapID int32, mobID int32, enabled bool) bool {
	field, ok := ctrl.server.fields[mapID]
	if !ok {
		return false
	}
	inst, err := field.getInstance(ctrl.inst.id)
	if err != nil {
		inst, err = field.getInstance(0)
		if err != nil {
			return false
		}
	}
	inst.lifePool.setMobSpawnEnabled(mobID, enabled)
	return true
}

func runReactorScript(program *goja.Program, server *Server, inst *fieldInstance, reactor *fieldReactor) error {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	rm := &reactorScriptController{server: server, inst: inst, reactor: reactor}
	_ = vm.Set("rm", rm)
	_, err := vm.RunProgram(program)
	if err != nil {
		return err
	}
	if fn := vm.Get("act"); fn != nil && !goja.IsUndefined(fn) && !goja.IsNull(fn) {
		var act func()
		if err := vm.ExportTo(fn, &act); err != nil {
			return err
		}
		if act != nil {
			act()
		}
	}
	return nil
}

func choosePortalForWarp(dstInst *fieldInstance, backToMapID int32, srcName, preferName string) (portal, error) {
	if preferName != "" {
		if p, e := dstInst.getPortalFromName(preferName); e == nil {
			return p, nil
		}
	}
	for _, p := range dstInst.portals {
		if p.destFieldID == backToMapID && p.destName == srcName {
			return p, nil
		}
	}
	for _, p := range dstInst.portals {
		if p.destFieldID == backToMapID {
			return p, nil
		}
	}
	return dstInst.getRandomSpawnPortal()
}

func (ctrl *portalScriptController) Warp(mapID int32, portalName string) bool {
	dstField, ok := ctrl.server.fields[mapID]
	if !ok {
		ctrl.blocked = true
		return false
	}
	var (
		inst *fieldInstance
		err  error
	)
	if ctrl.plr.event != nil {
		inst, err = ensureFieldInstance(dstField, ctrl.plr.event.instanceID, &ctrl.server.rates, ctrl.server)
	} else {
		inst, err = dstField.getInstance(ctrl.plr.inst.id)
		if err != nil {
			inst, err = dstField.getInstance(0)
		}
	}
	if err != nil {
		ctrl.blocked = true
		return false
	}
	var portal portal
	if portalName != "" {
		portal, err = inst.getPortalFromName(portalName)
		if err != nil {
			ctrl.blocked = true
			return false
		}
	} else {
		portal, err = choosePortalForWarp(inst, ctrl.plr.mapID, "", "")
		if err != nil {
			ctrl.blocked = true
			return false
		}
	}
	_ = ctrl.server.warpPlayer(ctrl.plr, dstField, portal, true)
	ctrl.warped = true
	return true
}

func (ctrl *portalScriptController) Message(msg string) {
	ctrl.plr.Send(packetMessageRedText(msg))
}

func (ctrl *portalScriptController) Id() int {
	return int(ctrl.portal.id)
}

func (ctrl *portalScriptController) Name() string {
	return ctrl.portal.name
}

func (ctrl *portalScriptController) Block(msg string) bool {
	if msg != "" {
		ctrl.Message(msg)
	}
	ctrl.blocked = true
	return false
}

func runPortalScript(program *goja.Program, plr *Player, server *Server, src portal) (warped bool, blocked bool, err error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	plrCtrl := &scriptPlayerWrapper{plr: plr, server: server}
	mapWrapper := &scriptMapWrapper{inst: plr.inst, server: server}
	portalCtrl := &portalScriptController{plr: plr, server: server, portal: src}
	_ = vm.Set("plr", plrCtrl)
	_ = vm.Set("map", mapWrapper)
	_ = vm.Set("portal", portalCtrl)
	_, err = vm.RunProgram(program)
	return portalCtrl.warped, portalCtrl.blocked, err
}

func createNpcChatController(npcID int32, conn mnet.Client, program *goja.Program, plr *Player, server *Server) (*npcChatController, error) {
	ctrl := &npcChatController{
		npcID:   npcID,
		conn:    conn,
		vm:      goja.New(),
		program: program,
	}

	plrCtrl := &scriptPlayerWrapper{
		plr:    plr,
		server: server,
	}

	mapWrapper := &scriptMapWrapper{
		inst:   plr.inst,
		server: server,
	}

	ctrl.vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	_ = ctrl.vm.Set("npc", ctrl)
	_ = ctrl.vm.Set("plr", plrCtrl)
	_ = ctrl.vm.Set("map", mapWrapper)

	return ctrl, nil
}

func (ctrl *npcChatController) Id() int32 {
	return ctrl.npcID
}

// SendNext simple next packet to Player
func (ctrl *npcChatController) SendNext(text string) int {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatBackNext(ctrl.npcID, text, true, false))
		ctrl.vm.Interrupt("Send")
		return 0
	}
	if ctrl.stateTracker.getCurrentState() == npcNextState {
		return 1
	}
	return 0
}

// SendBackNext packet to Player
func (ctrl *npcChatController) SendBackNext(msg string) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatBackNext(ctrl.npcID, msg, true, true))
		ctrl.vm.Interrupt("SendBackNext")
	}
}

// SendBackNext packet to Player
func (ctrl *npcChatController) SendBack(msg string) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatBackNext(ctrl.npcID, msg, false, true))
		ctrl.vm.Interrupt("SendBackNext")
	}
}

// SendOK packet to Player
func (ctrl *npcChatController) SendOk(msg string) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatOk(ctrl.npcID, msg))
		ctrl.vm.Interrupt("SendOk")
	}
}

// SendYesNo packet to Player
func (ctrl *npcChatController) SendYesNo(msg string) bool {
	state := ctrl.stateTracker.getCurrentState()

	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatYesNo(ctrl.npcID, msg))
		ctrl.vm.Interrupt("SendYesNo")
		return false
	}

	if state == npcYesState {
		return true
	} else if state == npcNoState {
		return false
	}

	return false
}

func (ctrl *npcChatController) SendAcceptDecline(msg string) bool {
	return ctrl.SendYesNo(msg)
}

// SendInputText packet to Player
func (ctrl *npcChatController) SendInputText(msg, defaultInput string, minLength, maxLength int16) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatUserString(ctrl.npcID, msg, defaultInput, minLength, maxLength))
		ctrl.vm.Interrupt("SendInputText")
	}
}

// SendInputNumber packet to Player
func (ctrl *npcChatController) SendInputNumber(msg string, defaultInput, minLength, maxLength int32) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatUserNumber(ctrl.npcID, msg, defaultInput, minLength, maxLength))
		ctrl.vm.Interrupt("SendInputNumber")
	}
}

// SendSelection packet to Player
func (ctrl *npcChatController) SendSelection(msg string) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatSelection(ctrl.npcID, msg))
		ctrl.vm.Interrupt("SendSelection")
	}
}

// SendStyles packet to Player
func (ctrl *npcChatController) SendStyles(msg string, styles []int32) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.stateTracker.addState(npcSelectionState)
		ctrl.conn.Send(packetNpcChatStyleWindow(ctrl.npcID, msg, styles))
		ctrl.vm.Interrupt("SendStyles")
	}
}

func (ctrl *npcChatController) SendAvatar(text string, avatars ...int32) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.stateTracker.addState(npcSelectionState)
		ctrl.conn.Send(packetNpcChatStyleWindow(ctrl.npcID, text, avatars))
		ctrl.vm.Interrupt("SendAvatar")
	}
}

func (ctrl *npcChatController) AskAvatar(text string, avatars ...int32) int {
	return ctrl.SendStyleChoice(text, avatars)
}

func (ctrl *npcChatController) SendStyleChoice(text string, avatars []int32) int {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.stateTracker.addState(npcSelectionState)
		ctrl.conn.Send(packetNpcChatStyleWindow(ctrl.npcID, text, avatars))
		ctrl.vm.Interrupt("SendStyleChoice")
		return -1
	}
	if len(ctrl.stateTracker.selections) > ctrl.stateTracker.selection {
		val := ctrl.stateTracker.selections[ctrl.stateTracker.selection]
		ctrl.stateTracker.selection++
		return int(val)
	}
	return -1
}

// SendGuildCreation
func (ctrl *npcChatController) SendGuildCreation() {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetGuildEnterName())
		ctrl.vm.Interrupt("SendGuildCreation")
	}
}

// SendGuildEmblemEditor
func (ctrl *npcChatController) SendGuildEmblemEditor() {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetGuildEmblemEditor())
		ctrl.vm.Interrupt("SendGuildEmblemEditor")
	}
}

// SendShop packet to Player
func (ctrl *npcChatController) SendShop(goods [][]int32) {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.goods = goods
		log.Println("sending shop")
		ctrl.conn.Send(packetNpcShop(ctrl.npcID, goods))
		ctrl.vm.Interrupt("SendShop")
	}
}

func (ctrl *npcChatController) SendStorage(npcID int32) {
	if ctrl.stateTracker.performInterrupt() {
		var storageMesos int32
		var storageSlots byte
		var allItems []Item

		accountID := ctrl.conn.GetAccountID()
		if accountID != 0 {
			st := new(storage)
			if err := st.load(accountID); err == nil {
				storageMesos = st.mesos
				storageSlots = st.maxSlots
				allItems = st.getAllItems()
			}
		}

		ctrl.conn.Send(packetNpcStorageShow(npcID, storageMesos, storageSlots, allItems))
		ctrl.vm.Interrupt("SendStorage")
	}
}

func (ctrl *npcChatController) SendMenu(baseText string, selections ...string) int {
	msg := baseText
	if len(selections) > 0 {
		var b strings.Builder
		if len(msg) > 0 {
			b.WriteString(msg)
			if msg[len(msg)-1] != '\n' {
				b.WriteByte('\n')
			}
		}
		for i, s := range selections {
			fmt.Fprintf(&b, "#L%d#%s#l\n", i, s)
		}
		msg = b.String()
	}

	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatSelection(ctrl.npcID, msg))
		ctrl.vm.Interrupt("SendMenu")
		return -1
	}
	if len(ctrl.stateTracker.selections) > ctrl.stateTracker.selection {
		val := ctrl.stateTracker.selections[ctrl.stateTracker.selection]
		ctrl.stateTracker.selection++
		return int(val)
	}
	return -1
}

func (ctrl *npcChatController) AskMenu(baseText string, selections ...string) int {
	return ctrl.SendMenu(baseText, selections...)
}

func (ctrl *npcChatController) SendImage(imagePath string) {
	img := fmt.Sprintf("#f%s#", imagePath)
	ctrl.SendOk(img)
}

func (ctrl *npcChatController) SendNumber(text string, def, min, max int) int {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatUserNumber(ctrl.npcID, text, int32(def), int32(min), int32(max)))
		ctrl.vm.Interrupt("SendNumber")
		return def
	}
	if len(ctrl.stateTracker.numbers) > ctrl.stateTracker.number {
		val := ctrl.stateTracker.numbers[ctrl.stateTracker.number]
		ctrl.stateTracker.number++
		return int(val)
	}
	return def
}

func (ctrl *npcChatController) SendBoxText(askMsg, defaultAnswer string, column, line int) string {
	max := column * line
	if max <= 0 {
		max = 200
	}
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatUserString(ctrl.npcID, askMsg, defaultAnswer, 0, int16(max)))
		ctrl.vm.Interrupt("SendBoxText")
		return defaultAnswer
	}
	if len(ctrl.stateTracker.inputs) > ctrl.stateTracker.input {
		val := ctrl.stateTracker.inputs[ctrl.stateTracker.input]
		ctrl.stateTracker.input++
		return val
	}
	return defaultAnswer
}

func (ctrl *npcChatController) SendQuiz(text, problem, hint string, inputMin, inputMax, _ int) string {
	prompt := text
	if problem != "" {
		if len(prompt) > 0 {
			prompt += "\n"
		}
		prompt += problem
	}
	if hint != "" {
		prompt += "\n(" + hint + ")"
	}
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatUserString(ctrl.npcID, prompt, "", int16(inputMin), int16(inputMax)))
		ctrl.vm.Interrupt("SendQuiz")
		return ""
	}
	if len(ctrl.stateTracker.inputs) > ctrl.stateTracker.input {
		val := ctrl.stateTracker.inputs[ctrl.stateTracker.input]
		ctrl.stateTracker.input++
		return val
	}
	return ""
}

func (ctrl *npcChatController) SendSlideMenu(text string) int {
	if ctrl.stateTracker.performInterrupt() {
		ctrl.conn.Send(packetNpcChatSelection(ctrl.npcID, text))
		ctrl.vm.Interrupt("SendSlideMenu")
		return -1
	}
	if len(ctrl.stateTracker.selections) > ctrl.stateTracker.selection {
		val := ctrl.stateTracker.selections[ctrl.stateTracker.selection]
		ctrl.stateTracker.selection++
		return int(val)
	}
	return -1
}

func (ctrl *npcChatController) clearUserInput() {
	// Reset counters but preserve the data
	ctrl.stateTracker.selection = 0
	ctrl.stateTracker.input = 0
	ctrl.stateTracker.number = 0
}

func (ctrl *npcChatController) Dispose() {
	ctrl.disposed = true
}

// Selection value
func (ctrl *npcChatController) Selection() int32 {
	if len(ctrl.stateTracker.selections) == 0 {
		return -1
	}
	if ctrl.stateTracker.selection >= len(ctrl.stateTracker.selections) {
		return ctrl.stateTracker.selections[len(ctrl.stateTracker.selections)-1]
	}
	val := ctrl.stateTracker.selections[ctrl.stateTracker.selection]
	ctrl.stateTracker.selection++
	return val
}

func (ctrl *npcChatController) InputString() string {
	if len(ctrl.stateTracker.inputs) == 0 {
		return ""
	}
	if ctrl.stateTracker.input >= len(ctrl.stateTracker.inputs) {
		return ctrl.stateTracker.inputs[len(ctrl.stateTracker.inputs)-1]
	}
	val := ctrl.stateTracker.inputs[ctrl.stateTracker.input]
	ctrl.stateTracker.input++
	return val
}

func (ctrl *npcChatController) InputNumber() int32 {
	if len(ctrl.stateTracker.numbers) == 0 {
		return 0
	}
	if ctrl.stateTracker.number >= len(ctrl.stateTracker.numbers) {
		return ctrl.stateTracker.numbers[len(ctrl.stateTracker.numbers)-1]
	}
	val := ctrl.stateTracker.numbers[ctrl.stateTracker.number]
	ctrl.stateTracker.number++
	return val
}

func (ctrl *npcChatController) run() bool {
	currentConversationPos := ctrl.stateTracker.currentPos
	ctrl.selectionCalls = 0

	if currentConversationPos == 0 && ctrl.stateTracker.lastPos == 0 {
		ctrl.stateTracker.selections = ctrl.stateTracker.selections[:0]
	} else {
		ctrl.stateTracker.currentPos = 0
	}

	if ctrl.vm == nil || ctrl.program == nil {
		return true
	}

	_, err := ctrl.vm.RunProgram(ctrl.program)

	if err != nil {
		if _, isInterrupted := err.(*goja.InterruptedError); isInterrupted {
			if ctrl.disposed {
				return true
			}
			return false
		}
		return true
	}

	if ctrl.stateTracker.currentPos >= ctrl.stateTracker.lastPos {
		return true
	}
	return false
}
