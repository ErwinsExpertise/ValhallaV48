package channel

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type gachaponNPCConfig struct {
	NPCID        int32
	TownName     string
	TicketItemID int32
	PoolIndex    int
	MapID        int32
	HasLevelGate bool
	MinLevel     int
}

type gachaponPoolConfig struct {
	Index       int
	TotalWeight uint32
	Rewards     []gachaponRewardConfig
}

type gachaponRewardConfig struct {
	ItemID           int32
	Weight           uint32
	CumulativeWeight uint32
	Count            int16
}

type gachaponGrantOps interface {
	CanHold(id int32, amount int16) bool
	ConsumeTicket(ticketID int32) bool
	GrantReward(id int32, amount int16) bool
}

type gachaponDataFile struct {
	ScriptEntries []gachaponNPCConfig  `json:"scriptEntries"`
	Pools         []gachaponPoolConfig `json:"pools"`
}

var gachaponNPCs []gachaponNPCConfig
var gachaponPools []gachaponPoolConfig
var gachaponNPCByID map[int32]gachaponNPCConfig

func buildGachaponLookup() {
	gachaponNPCByID = make(map[int32]gachaponNPCConfig, len(gachaponNPCs))
	for _, entry := range gachaponNPCs {
		gachaponNPCByID[entry.NPCID] = entry
	}
}

func populateGachaponTable(gachaponJSON string) error {
	data, err := loadGachaponDataFile(gachaponJSON)
	if err != nil {
		return err
	}

	gachaponNPCs = data.ScriptEntries
	gachaponPools = normalizeGachaponPools(data.Pools)
	buildGachaponLookup()
	return nil
}

func loadGachaponDataFile(path string) (gachaponDataFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return gachaponDataFile{}, err
	}

	var data gachaponDataFile
	if err := json.Unmarshal(b, &data); err != nil {
		return gachaponDataFile{}, err
	}
	return data, nil
}
func normalizeGachaponPools(pools []gachaponPoolConfig) []gachaponPoolConfig {
	maxIndex := -1
	for _, pool := range pools {
		if pool.Index > maxIndex {
			maxIndex = pool.Index
		}
	}
	if maxIndex < 0 {
		return nil
	}
	normalized := make([]gachaponPoolConfig, maxIndex+1)
	for _, pool := range pools {
		normalized[pool.Index] = pool
	}
	return normalized
}

func (ctrl *npcChatController) RunGachapon() {
	log.Printf("gachapon enter npc=%d script=%s map=%d", ctrl.npcID, ctrl.script, ctrl.plr.mapID)
	config, ok := gachaponNPCByID[ctrl.npcID]
	if !ok {
		log.Printf("gachapon missing config npc=%d script=%s", ctrl.npcID, ctrl.script)
		ctrl.SendOk("Here's Gachapon.")
		return
	}
	log.Printf("gachapon config npc=%d town=%q ticket=%d pool=%d levelGate=%t minLevel=%d", ctrl.npcID, config.TownName, config.TicketItemID, config.PoolIndex, config.HasLevelGate, config.MinLevel)

	if config.HasLevelGate && int(ctrl.plr.level) < config.MinLevel {
		log.Printf("gachapon blocked by level npc=%d level=%d required=%d", ctrl.npcID, ctrl.plr.level, config.MinLevel)
		ctrl.SendOk("You need to be at least Level 15 in order to use Gachapon.")
		return
	}

	if ctrl.plr.countItem(config.TicketItemID) < 1 {
		log.Printf("gachapon missing ticket npc=%d ticket=%d", ctrl.npcID, config.TicketItemID)
		ctrl.SendOk("Here's Gachapon.")
		return
	}

	state := ctrl.stateTracker.getCurrentState()
	if state != npcYesState {
		log.Printf("gachapon prompting yes/no npc=%d state=%d", ctrl.npcID, state)
		ctrl.SendYesNo("You may use Gachapon. Would you like to use your Gachapon ticket?")
		return
	}
	log.Printf("gachapon resumed-yes npc=%d", ctrl.npcID)

	reward, ok := drawGachaponReward(config.PoolIndex, randomGachaponRoll(config.PoolIndex))
	if !ok {
		log.Printf("gachapon failed draw npc=%d pool=%d", ctrl.npcID, config.PoolIndex)
		ctrl.SendOk("Please check your item inventory and see if you have the ticket, or if the inventory is full.")
		return
	}
	log.Printf("gachapon drew reward npc=%d item=%d count=%d", ctrl.npcID, reward.ItemID, reward.Count)

	granted, ok := executeGachaponGrant(gachaponPlayerOps{plr: ctrl.plr}, config.TicketItemID, reward)
	if !ok {
		log.Printf("gachapon failed grant npc=%d item=%d count=%d", ctrl.npcID, reward.ItemID, reward.Count)
		ctrl.SendOk("Please make room on your item inventory and then try again.")
		return
	}

	log.Printf("gachapon granted npc=%d item=%d count=%d", ctrl.npcID, granted.ItemID, granted.Count)
	ctrl.SendOk("You have obtained #b#t" + strconv.Itoa(int(granted.ItemID)) + "##k.")
}

type gachaponPlayerOps struct {
	plr *Player
}

func (ops gachaponPlayerOps) CanHold(id int32, amount int16) bool {
	if amount <= 0 {
		amount = 1
	}
	item, err := CreateItemFromID(id, amount)
	if err != nil {
		return false
	}
	return ops.plr.canReceiveQuestActItems([]Item{item})
}

func (ops gachaponPlayerOps) ConsumeTicket(ticketID int32) bool {
	return ops.plr.removeItemsByID(ticketID, 1, false)
}

func (ops gachaponPlayerOps) GrantReward(id int32, amount int16) bool {
	item, err := CreateItemFromID(id, amount)
	if err != nil {
		return false
	}
	_, err = ops.plr.GiveItem(item)
	return err == nil
}

func randomGachaponRoll(poolIndex int) uint32 {
	if poolIndex < 0 || poolIndex >= len(gachaponPools) {
		return 0
	}
	total := gachaponPools[poolIndex].TotalWeight
	if total == 0 {
		return 0
	}
	return uint32(rand.Int63n(int64(total)))
}

func drawGachaponReward(poolIndex int, roll uint32) (gachaponRewardConfig, bool) {
	if poolIndex < 0 || poolIndex >= len(gachaponPools) {
		return gachaponRewardConfig{}, false
	}
	pool := gachaponPools[poolIndex]
	if pool.TotalWeight == 0 || len(pool.Rewards) == 0 {
		return gachaponRewardConfig{}, false
	}
	for _, reward := range pool.Rewards {
		if roll < reward.CumulativeWeight {
			return reward, true
		}
	}
	return gachaponRewardConfig{}, false
}

func executeGachaponGrant(ops gachaponGrantOps, ticketID int32, reward gachaponRewardConfig) (gachaponRewardConfig, bool) {
	amount := reward.Count
	if amount <= 0 {
		amount = 1
	}
	if !ops.CanHold(reward.ItemID, amount) {
		return gachaponRewardConfig{}, false
	}
	if !ops.ConsumeTicket(ticketID) {
		return gachaponRewardConfig{}, false
	}
	if !ops.GrantReward(reward.ItemID, amount) {
		return gachaponRewardConfig{}, false
	}
	reward.Count = amount
	return reward, true
}
