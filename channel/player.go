package channel

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	mathrand "math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/common/opcode"
	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/constant/skill"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

type buddy struct {
	id        int32
	name      string
	channelID int32
	status    byte  // 0 - online, 1 - buddy request, 2 - offline
	cashShop  int32 // > 0 means is in cash shop
}

type playerSkill struct {
	ID             int32
	Level, Mastery byte
	Cooldown       int16
	CooldownTime   int16
	TimeLastUsed   int64
}

type funcKeyMapState struct {
	Entries [89]constant.FuncKeyMapped
	Loaded  bool
}

var starterMapVisualEquipBySlot = map[byte]int32{
	1: constant.StarterMapVisualHatID,
	5: constant.StarterMapVisualOverallID,
	7: constant.StarterMapVisualShoesID,
}

func isStarterEquipOverrideMap(mapID int32) bool {
	switch mapID {
	case constant.StarterMapBeginnerTown,
		constant.StarterMapAmherst,
		constant.StarterMapSouthperry,
		constant.StarterMapTutorialExit:
		return true
	default:
		return false
	}
}

func applyStarterMapVisualEquips(mapID int32, visible, masked map[byte]int32) {
	if !isStarterEquipOverrideMap(mapID) {
		return
	}

	for slot, itemID := range starterMapVisualEquipBySlot {
		if equipped, ok := visible[slot]; ok {
			masked[slot] = equipped
		}
		visible[slot] = itemID
	}
}

func starterMapVisualEquipSlot(slotID int16) (byte, bool) {
	if slotID >= 0 {
		return 0, false
	}

	if slotID < -100 {
		slotID += 100
	}

	slot := byte(-slotID)
	_, ok := starterMapVisualEquipBySlot[slot]
	return slot, ok
}

func starterMapVisualEquipItem(slot byte) Item {
	return Item{
		invID:      constant.InventoryEquip,
		slotID:     -int16(slot),
		ID:         starterMapVisualEquipBySlot[slot],
		amount:     1,
		expireTime: neverExpire,
	}
}

func starterMapVisualEquipNegativeSlot(slot byte) int16 {
	return -int16(slot)
}

func (plr *Player) refreshStarterVisualEquipSlots() {
	if plr == nil {
		return
	}
	for slot := range starterMapVisualEquipBySlot {
		negSlot := starterMapVisualEquipNegativeSlot(slot)
		found := false
		for _, it := range plr.equip {
			if it.slotID != negSlot || it.cash {
				continue
			}
			plr.Send(packetInventoryAddItem(it, true))
			found = true
			break
		}
		if !found {
			plr.Send(packetInventoryRemoveItem(Item{invID: constant.InventoryEquip, slotID: negSlot}))
		}
	}
}

func writeSetFieldEquippedItems(p *mpacket.Packet, plr Player) {
	if !isStarterEquipOverrideMap(plr.mapID) {
		for _, it := range plr.equip {
			if it.slotID < 0 && !it.cash {
				p.WriteBytes(it.InventoryBytes())
			}
		}
		p.WriteByte(0)

		for _, it := range plr.equip {
			if it.slotID < 0 && it.cash {
				p.WriteBytes(it.InventoryBytes())
			}
		}
		p.WriteByte(0)
		return
	}

	writtenOverride := make(map[byte]bool, len(starterMapVisualEquipBySlot))

	for _, it := range plr.equip {
		if it.slotID >= 0 || it.cash {
			continue
		}

		if slot, ok := starterMapVisualEquipSlot(it.slotID); ok {
			if !writtenOverride[slot] {
				p.WriteBytes(starterMapVisualEquipItem(slot).InventoryBytes())
				writtenOverride[slot] = true
			}
			continue
		}

		p.WriteBytes(it.InventoryBytes())
	}

	for slot := range starterMapVisualEquipBySlot {
		if !writtenOverride[slot] {
			p.WriteBytes(starterMapVisualEquipItem(slot).InventoryBytes())
		}
	}
	p.WriteByte(0)

	for _, it := range plr.equip {
		if it.slotID >= 0 || !it.cash {
			continue
		}

		if _, ok := starterMapVisualEquipSlot(it.slotID); ok {
			continue
		}

		p.WriteBytes(it.InventoryBytes())
	}
	p.WriteByte(0)
}

func createPlayerSkillFromData(ID int32, level byte) (playerSkill, error) {
	skill, err := nx.GetPlayerSkill(ID)
	if err != nil {
		return playerSkill{}, fmt.Errorf("invalid skill ID %d (level %d): %w", ID, level, err)
	}
	if level == 0 || int(level) > len(skill) {
		return playerSkill{}, fmt.Errorf("invalid skill level %d for skill ID %d (max %d)", level, ID, len(skill))
	}

	return playerSkill{
		ID:           ID,
		Level:        level,
		Mastery:      byte(skill[level-1].Mastery),
		Cooldown:     0,
		CooldownTime: int16(skill[level-1].Cooltime),
		TimeLastUsed: 0,
	}, nil
}

func getSkillsFromCharID(id int32) []playerSkill {
	skills := []playerSkill{}

	const filter = "skillID, level, cooldown"
	row, err := common.DB.Query("SELECT "+filter+" FROM skills WHERE characterID=?", id)
	if err != nil {
		log.Printf("getSkillsFromCharID: query failed for character %d: %v", id, err)
		return skills
	}
	defer row.Close()

	for row.Next() {
		var ps playerSkill
		if err := row.Scan(&ps.ID, &ps.Level, &ps.Cooldown); err != nil {
			log.Printf("getSkillsFromCharID: scan failed for character %d: %v", id, err)
			continue
		}

		skillData, err := nx.GetPlayerSkill(ps.ID)
		if err != nil {
			log.Printf("getSkillsFromCharID: missing nx data for skill %d: %v", ps.ID, err)
			continue
		}
		if ps.Level == 0 || int(ps.Level) > len(skillData) {
			log.Printf("getSkillsFromCharID: invalid level %d for skill %d (max %d), skipping", ps.Level, ps.ID, len(skillData))
			continue
		}

		ps.CooldownTime = int16(skillData[ps.Level-1].Cooltime)
		skills = append(skills, ps)
	}

	if err := row.Err(); err != nil {
		log.Printf("getSkillsFromCharID: rows error for character %d: %v", id, err)
	}

	return skills
}

type updatePartyInfoFunc func(partyID, playerID, job, level, mapID int32, name string)

type Players struct {
	conn map[mnet.Client]*Player
	id   map[int32]*Player
	name map[string]*Player
}

// getItemSlotMax retrieves the actual slotMax for an item from NX data
// Falls back to constant.MaxItemStack if not found
func getItemSlotMax(itemID int32) int16 {
	if nxInfo, err := nx.GetItem(itemID); err == nil && nxInfo.SlotMax > 0 {
		return nxInfo.SlotMax
	}
	return constant.MaxItemStack
}

func NewPlayers() Players {
	return Players{
		conn: make(map[mnet.Client]*Player),
		id:   make(map[int32]*Player),
		name: make(map[string]*Player),
	}
}

func (p *Players) Add(plr *Player) {
	p.conn[plr.Conn] = plr
	p.id[plr.ID] = plr
	p.name[plr.Name] = plr
}

func (p Players) count() int {
	return len(p.conn)
}

func (p Players) observe(f func(*Player)) {
	for _, plr := range p.conn {
		f(plr)
	}
}

func (p Players) GetFromConn(conn mnet.Client) (*Player, error) {
	if v, ok := p.conn[conn]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("Player not found for connection")
}

func (p Players) GetFromID(id int32) (*Player, error) {
	if v, ok := p.id[id]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("Player not found for ID: %d", id)
}

func (p Players) GetFromName(name string) (*Player, error) {
	if v, ok := p.name[name]; ok {
		return v, nil
	}

	return nil, fmt.Errorf("Player not found for Name: %s", name)
}

func (p *Players) RemoveFromConn(conn mnet.Client) error {
	if plr, ok := p.conn[conn]; ok {
		delete(p.id, plr.ID)
		delete(p.name, plr.Name)
		delete(p.conn, conn)

		return nil
	}

	return fmt.Errorf("Player not found for removal")
}

func (p Players) broadcast(packet mpacket.Packet) {
	for conn := range p.conn {
		conn.Send(packet)
	}
}

func (p Players) Flush() {
	for _, plr := range p.conn {
		flushNow(plr)
	}
}

type Player struct {
	Conn mnet.Client
	inst *fieldInstance

	ID          int32 // Unique identifier of the character
	accountID   int32
	accountName string
	worldID     byte
	ChannelID   byte

	mapID       int32
	mapPos      byte
	previousMap int32
	savedMaps   map[string]int32
	portalCount byte

	job int16

	level byte
	str   int16
	dex   int16
	intt  int16
	luk   int16
	hp    int16
	maxHP int16
	mp    int16
	maxMP int16
	ap    int16
	sp    int16
	exp   int32
	fame  int16

	totalStr      int16
	totalDex      int16
	totalInt      int16
	totalLuk      int16
	totalWatk     int16
	totalMatk     int16
	totalAccuracy int16

	Name      string
	gender    byte
	skin      byte
	face      int32
	hair      int32
	chairID   int32
	petCashID int64
	stance    byte
	pos       pos

	equipSlotSize byte
	useSlotSize   byte
	setupSlotSize byte
	etcSlotSize   byte
	cashSlotSize  byte

	equip []Item
	use   []Item
	setUp []Item
	etc   []Item
	cash  []Item

	mesos          int32
	nx             int32
	maplepoints    int32
	partnerID      int32
	marriageID     int32
	marriageItemID int32
	divorceUntil   int64

	storageInventory *storage

	skills map[int32]playerSkill

	miniGameWins, miniGameDraw, miniGameLoss, miniGamePoints int32

	lastAttackPacketTime int64
	nextMapDamageAtMs    int64

	buddyListSize      byte
	buddyList          []buddy
	funcKeyMap         funcKeyMapState
	quickslotKeys      [2]int32
	petConsumeItemID   int32
	petConsumeMPItemID int32

	regTeleportRocks []int32 // 5 regular teleport rocks
	vipTeleportRocks []int32 // 10 VIP teleport rocks

	chalkboardText   string
	chalkboardActive bool
	pendingMerchant  merchantPermitState
	storeBankShopID  int64
	storeBankNpcID   int32
	storeBankOpen    bool

	expCouponItemID     int32
	expCouponExpiresAt  int64
	dropCouponItemID    int32
	dropCouponExpiresAt int64

	party *party
	guild *guild

	UpdatePartyInfo updatePartyInfoFunc

	rates *rates

	buffs           *CharacterBuffs
	rings           map[int32]ringRecord
	ringEffectState map[int32]bool

	quests quests

	summons *summonState
	pet     *pet

	// Mystic Door tracking
	doorMapID       int32
	doorSpawnID     int32
	doorPortalIndex int
	townDoorMapID   int32
	townDoorSpawnID int32
	townPortalIndex int

	// Per-Player RNG for deterministic randomness
	rng *mathrand.Rand

	// write-behind persistence
	dirty DirtyBits

	lastChairHeal time.Time

	deadlyAttackActive bool
	deadlyAttackTime   int64

	// Safety charm flag - prevents exp loss on death
	hasSafetyCharm bool

	event *event
}

// Helper: mark dirty and schedule debounced save.
func (d *Player) MarkDirty(bits DirtyBits, debounce time.Duration) {
	d.dirty |= bits
	scheduleSave(d, debounce)
}

// Helper: clear dirty bits after successful flush (kept for future; saver currently doesn't feed back)
func (d *Player) clearDirty(bits DirtyBits) {
	d.dirty &^= bits
}

func (d *Player) FlushNow() {
	flushNow(d)
}

// SeedRNGDeterministic seeds the per-Player RNG using stable identifiers so
// gain sequences are reproducible across restarts and processes.
func (d *Player) SeedRNGDeterministic() {
	// Compose as uint64 to avoid int64 constant overflow, then cast at runtime.
	const gamma uint64 = 0x9e3779b97f4a7c15
	seed64 := gamma ^
		(uint64(uint32(d.ID)) << 1) ^
		(uint64(uint32(d.accountID)) << 33) ^
		(uint64(d.worldID) << 52)

	seed := int64(seed64) // two's complement wrapping is fine for rand.Source
	d.rng = mathrand.New(mathrand.NewSource(seed))
}

// ensureRNG guarantees d.rng is initialized. If a deterministic seed has not
// been set yet, it will use a time-based seed (non-deterministic).
func (d *Player) ensureRNG() {
	if d.rng == nil {
		// Default to deterministic seeding for stability unless you want variability:
		d.SeedRNGDeterministic()
	}
}

func (d *Player) randIntn(n int) int {
	d.ensureRNG()
	return d.rng.Intn(n)
}

// levelUpGains returns (hpGain, mpGain) using per-Player RNG and job family.
// The random component uses a small range similar to legacy behavior.
// Tweak the constants to match your balance targets if needed.
func (d *Player) levelUpGains() (int16, int16) {
	r := int16(d.randIntn(3) + 1) // legacy-style variance 1..3

	mainClass := d.job / 100
	switch {
	case d.job == 0 || mainClass == 0: // Beginner and pre-advancement
		// Balanced but modest growth
		return r + 12, r + 10
	case mainClass == 1: // Warrior
		// High HP, low MP growth
		return r + 24, r + 4
	case mainClass == 2: // Magician
		// Low HP, high MP growth
		return r + 10, r + 22
	case mainClass == 3 || mainClass == 4: // Bowman / Thief
		// Moderate HP/MP growth
		return r + 20, r + 14
	default:
		// Fallback for any other jobs/classes
		return r + 16, r + 12
	}
}

// getPassiveHPBonus returns the HP bonus from Improved MaxHP Increase skill
func (d *Player) getPassiveHPBonus() int16 {
	if ps, ok := d.skills[int32(skill.ImprovedMaxHpIncrease)]; ok {
		skillData, err := nx.GetPlayerSkill(int32(skill.ImprovedMaxHpIncrease))
		if err == nil && ps.Level > 0 && int(ps.Level) <= len(skillData) {
			return int16(skillData[ps.Level-1].X)
		}
	}
	return 0
}

// getPassiveMPBonus returns the MP bonus from Improved MaxMP Increase skill
func (d *Player) getPassiveMPBonus() int16 {
	if ps, ok := d.skills[int32(skill.ImprovedMaxMpIncrease)]; ok {
		skillData, err := nx.GetPlayerSkill(int32(skill.ImprovedMaxMpIncrease))
		if err == nil && ps.Level > 0 && int(ps.Level) <= len(skillData) {
			return int16(skillData[ps.Level-1].X)
		}
	}
	return 0
}

// getHPGainForJob returns the HP gain when manually allocating AP to MaxHP
func (d *Player) getHPGainForJob() int16 {
	mainClass := d.job / constant.JobClassDivisor
	switch {
	case d.job == constant.BeginnerJobID || mainClass == 0: // Beginner
		return constant.BeginnerApHpGain
	case mainClass == 1: // Warrior
		return constant.WarriorApHpGain
	case mainClass == 2: // Magician
		return constant.MagicianApHpGain
	case mainClass == 3: // Bowman
		return constant.BowmanApHpGain
	case mainClass == 4: // Thief
		return constant.ThiefApHpGain
	default:
		return constant.DefaultApHpGain
	}
}

// getMPGainForJob returns the MP gain when manually allocating AP to MaxMP
func (d *Player) getMPGainForJob() int16 {
	mainClass := d.job / constant.JobClassDivisor
	baseMp := int16(constant.DefaultApMpGain)

	switch {
	case d.job == constant.BeginnerJobID || mainClass == 0: // Beginner
		baseMp = constant.BeginnerApMpGain
	case mainClass == 1: // Warrior
		baseMp = constant.WarriorApMpGain
	case mainClass == 2: // Magician
		baseMp = constant.MagicianApMpGain
	case mainClass == 3: // Bowman
		baseMp = constant.BowmanApMpGain
	case mainClass == 4: // Thief
		baseMp = constant.ThiefApMpGain
	}

	// Add INT bonus for MP (INT * multiplier / divisor)
	// Magicians get more MP from INT
	intMultiplier := int16(constant.IntMpMultiplierNormal)
	if mainClass == 2 {
		intMultiplier = constant.IntMpMultiplierMagician
	}
	baseMp += (d.intt * intMultiplier) / int16(constant.IntMpDivisor)

	return baseMp
}

// getWeaponMasterySkillID returns the mastery skill ID based on equipped weapon and job
func (d *Player) getWeaponMasterySkillID() int32 {
	// Find equipped weapon
	var weaponID int32 = 0
	for _, item := range d.equip {
		if item.slotID == constant.WeaponSlot {
			weaponID = item.ID
			break
		}
	}

	if weaponID == 0 {
		return 0
	}

	weaponType := weaponID / 10000

	// Map weapon types to mastery skills based on job
	// Use exact job ID ranges to avoid ambiguity
	switch {
	case d.job >= constant.FighterJobID && d.job < constant.PageJobID: // Fighter
		if weaponType == constant.WeaponType1HSword || weaponType == constant.WeaponType2HSword {
			return int32(skill.SwordMastery)
		} else if weaponType == constant.WeaponType1HAxe || weaponType == constant.WeaponType2HAxe {
			return int32(skill.AxeMastery)
		}
	case d.job >= constant.PageJobID && d.job < constant.SpearmanJobID: // Page / White Knight
		if weaponType == constant.WeaponType1HSword || weaponType == constant.WeaponType2HSword {
			return int32(skill.PageSwordMastery)
		} else if weaponType == constant.WeaponType1HBW || weaponType == constant.WeaponType2HBW {
			return int32(skill.BwMastery)
		}
	case d.job >= constant.SpearmanJobID && d.job < constant.MagicianJobID: // Spearman / Dragon Knight
		if weaponType == constant.WeaponTypeSpear {
			return int32(skill.SpearMastery)
		} else if weaponType == constant.WeaponTypePolearm {
			return int32(skill.PolearmMastery)
		}
	case d.job >= constant.HunterJobID && d.job < constant.CrossbowmanJobID: // Hunter / Ranger
		if weaponType == constant.WeaponTypeBow {
			return int32(skill.BowMastery)
		}
	case d.job >= constant.CrossbowmanJobID && d.job < constant.ThiefJobID: // Crossbowman / Sniper
		if weaponType == constant.WeaponTypeCrossbow {
			return int32(skill.CrossbowMastery)
		}
	case d.job >= constant.AssassinJobID && d.job < constant.BanditJobID: // Assassin / Hermit
		if weaponType == constant.WeaponTypeClaw {
			return int32(skill.ClawMastery)
		}
	case d.job >= constant.BanditJobID && d.job < constant.GmJobID: // Bandit / Chief Bandit
		if weaponType == constant.WeaponTypeDagger {
			return int32(skill.DaggerMastery)
		}
	}

	return 0
}

// getMasteryDisplay returns the mastery display value for attack packets
// Formula from reference: (skillLevel + 1) / 2
func (d *Player) getMasteryDisplay() byte {
	masterySkillID := d.getWeaponMasterySkillID()
	if masterySkillID == 0 {
		return 0
	}

	if ps, ok := d.skills[masterySkillID]; ok {
		// Apply formula: (level + 1) / 2
		return (ps.Level + 1) / constant.MasteryDisplayDivisor
	}

	return 0
}

// getRechargeBonus returns the bonus amount when recharging items
func (d *Player) getRechargeBonus() int16 {
	// Assassin and Hermit get bonus from Claw Mastery (level * multiplier)
	jobFamily := d.job / constant.JobClassDivisor
	jobBranch := (d.job % constant.JobClassDivisor) / constant.JobBranchDivisor
	if jobFamily == 4 && jobBranch == 1 { // Assassin or Hermit (both are 41x)
		if ps, ok := d.skills[int32(skill.ClawMastery)]; ok {
			return int16(ps.Level * constant.ClawMasteryRechargeMultiplier)
		}
	}

	return 0
}

func (d *Player) getCriticalSkillAndRate() (int32, int) {
	var weaponID int32 = 0
	for _, item := range d.equip {
		if item.slotID == constant.WeaponSlot {
			weaponID = item.ID
			break
		}
	}

	if weaponID == 0 {
		return 0, 0
	}

	weaponType := weaponID / 10000

	if weaponType == constant.WeaponTypeBow || weaponType == constant.WeaponTypeCrossbow {
		if ps, ok := d.skills[int32(skill.CriticalShot)]; ok {
			critRate := constant.CriticalBaseRate + (int(ps.Level) * constant.CriticalRatePerLevel)
			if critRate > constant.CriticalRateMax {
				critRate = constant.CriticalRateMax
			}
			return int32(skill.CriticalShot), critRate
		}
	}

	if weaponType == constant.WeaponTypeClaw {
		if ps, ok := d.skills[int32(skill.CriticalThrow)]; ok {
			critRate := constant.CriticalBaseRate + (int(ps.Level) * constant.CriticalRatePerLevel)
			if critRate > constant.CriticalRateMax {
				critRate = constant.CriticalRateMax
			}
			return int32(skill.CriticalThrow), critRate
		}
	}

	return 0, 0
}

func (d *Player) rollCritical(attackType int) bool {
	if attackType != attackRanged {
		return false
	}

	skillID, critRate := d.getCriticalSkillAndRate()
	if skillID == 0 || critRate == 0 {
		return false
	}

	roll := d.randIntn(100)
	return roll < critRate
}

// Send the Data a packet
func (d *Player) Send(packet mpacket.Packet) {
	if d == nil || d.Conn == nil {
		return
	}
	d.Conn.Send(packet)
}

func (d *Player) setJob(id int16) {
	d.job = id
	d.Conn.Send(packetPlayerStatChange(true, constant.JobID, int32(id)))
	d.MarkDirty(DirtyJob, 300*time.Millisecond)

	if d.party != nil {
		d.UpdatePartyInfo(d.party.ID, d.ID, int32(d.job), int32(d.level), d.mapID, d.Name)
	}
}

func (d *Player) levelUp() {
	d.giveAP(5)
	if d.level < 10 {
		d.giveSP(1)
	} else {
		d.giveSP(3)
	}

	// Use per-Player RNG and job-based helper for deterministic gains.
	hpGain, mpGain := d.levelUpGains()

	// Apply passive skill bonuses for levelup
	hpGain += d.getPassiveHPBonus()
	mpGain += d.getPassiveMPBonus()

	newMaxHP := d.maxHP + hpGain
	newMaxMP := d.maxMP + mpGain
	if newMaxHP < 1 {
		newMaxHP = 1
	}
	if newMaxMP < 0 {
		newMaxMP = 0
	}

	d.setMaxHP(newMaxHP)
	d.setMaxMP(newMaxMP)

	d.setHP(newMaxHP)
	d.setMP(newMaxMP)

	d.giveLevel(1)
}

func (d *Player) setEXP(amount int32) {
	if d.level >= 200 {
		d.exp = amount
		d.Send(packetPlayerStatChange(false, constant.ExpID, int32(amount)))
		d.MarkDirty(DirtyEXP, 800*time.Millisecond)
		return
	}

	for {
		if d.level >= 200 {
			d.exp = amount
			break
		}

		expForLevel := constant.ExpTable[d.level-1]
		remainder := amount - expForLevel
		if remainder >= 0 {
			d.levelUp()
			amount = remainder
		} else {
			d.exp = amount
			break
		}
	}

	d.Send(packetPlayerStatChange(false, constant.ExpID, d.exp))
	d.MarkDirty(DirtyEXP, 800*time.Millisecond)
}

func (d *Player) giveEXP(amount int32, fromMob, fromParty bool) {
	amount = int32(d.rates.exp * d.expCouponMultiplier(time.Now()) * float32(amount))

	switch {
	case fromMob:
		d.Send(packetMessageExpGained(true, false, amount))
	case fromParty:
		d.Send(packetMessageExpGained(false, false, amount))
	default:
		d.Send(packetMessageExpGained(false, true, amount))
	}

	d.setEXP(d.exp + amount)
}

func (d *Player) GetAccountName() string {
	return d.accountName
}

func (d *Player) GetGender() byte { return d.gender }
func (d *Player) GetSkin() byte   { return d.skin }
func (d *Player) GetFace() int32  { return d.face }
func (d *Player) GetHair() int32  { return d.hair }
func (d *Player) GetLevel() byte  { return d.level }
func (d *Player) GetJob() int16   { return d.job }
func (d *Player) GetStr() int16   { return d.str }
func (d *Player) GetDex() int16   { return d.dex }
func (d *Player) GetInt() int16   { return d.intt }
func (d *Player) GetLuk() int16   { return d.luk }
func (d *Player) GetHP() int16    { return d.hp }
func (d *Player) GetMaxHP() int16 { return d.maxHP }
func (d *Player) GetMP() int16    { return d.mp }
func (d *Player) GetMaxMP() int16 { return d.maxMP }
func (d *Player) GetAP() int16    { return d.ap }
func (d *Player) GetSP() int16    { return d.sp }
func (d *Player) GetEXP() int32   { return d.exp }
func (d *Player) GetFame() int16  { return d.fame }
func (d *Player) GetMapID() int32 { return d.mapID }
func (d *Player) GetMapPos() byte { return d.mapPos }
func (d *Player) GetMesos() int32 { return d.mesos }

func (d *Player) setLevel(amount byte) {
	d.level = amount
	d.Send(packetPlayerStatChange(false, constant.LevelID, int32(amount)))
	d.inst.send(packetPlayerLevelUpAnimation(d.ID))
	d.MarkDirty(DirtyLevel, 300*time.Millisecond)

	if d.party != nil {
		d.UpdatePartyInfo(d.party.ID, d.ID, int32(d.job), int32(d.level), d.mapID, d.Name)
	}
}

func (d *Player) giveLevel(amount byte) {
	d.setLevel(d.level + amount)
}

func (d *Player) setAP(amount int16) {
	d.ap = amount
	d.Send(packetPlayerStatChange(false, constant.ApID, int32(amount)))
	d.MarkDirty(DirtyAP, 300*time.Millisecond)
}

func (d *Player) giveAP(amount int16) {
	d.setAP(d.ap + amount)
}

func (d *Player) setSP(amount int16) {
	d.sp = amount
	d.Send(packetPlayerStatChange(true, constant.SpID, int32(amount)))
	d.MarkDirty(DirtySP, 300*time.Millisecond)
}

func (d *Player) giveSP(amount int16) {
	d.setSP(d.sp + amount)
}

func (d *Player) setStr(amount int16) {
	d.str = amount
	d.recalculateTotalStats()
	d.Send(packetPlayerStatChange(true, constant.StrID, int32(amount)))
	d.MarkDirty(DirtyStr, 500*time.Millisecond)
}

func (d *Player) giveStr(amount int16) {
	d.setStr(d.str + amount)
}

func (d *Player) setDex(amount int16) {
	d.dex = amount
	d.recalculateTotalStats()
	d.Send(packetPlayerStatChange(true, constant.DexID, int32(amount)))
	d.MarkDirty(DirtyDex, 500*time.Millisecond)
}

func (d *Player) giveDex(amount int16) {
	d.setDex(d.dex + amount)
}

func (d *Player) setInt(amount int16) {
	d.intt = amount
	d.recalculateTotalStats()
	d.Send(packetPlayerStatChange(true, constant.IntID, int32(amount)))
	d.MarkDirty(DirtyInt, 500*time.Millisecond)
}

func (d *Player) giveInt(amount int16) {
	d.setInt(d.intt + amount)
}

func (d *Player) setLuk(amount int16) {
	d.luk = amount
	d.recalculateTotalStats()
	d.Send(packetPlayerStatChange(true, constant.LukID, int32(amount)))
	d.MarkDirty(DirtyLuk, 500*time.Millisecond)
}

func (d *Player) giveLuk(amount int16) {
	d.setLuk(d.luk + amount)
}

func (d *Player) giveHP(amount int16) {
	target := int(d.hp) + int(amount)
	maxAllowed := int(d.effectiveMaxHP())
	if target < 0 {
		target = 0
	} else if target > maxAllowed {
		target = maxAllowed
	}

	d.setHP(int16(target))

	if d.party != nil {
		d.party.broadcast(packetPlayerHpChange(d.ID, int32(d.hp), int32(d.maxHP)))
	}
}

func (d *Player) setHP(amount int16) {
	if amount < 0 {
		amount = 0
	}
	effMax := d.effectiveMaxHP()
	if amount > constant.MaxHpValue {
		amount = constant.MaxHpValue
	}
	if amount > effMax {
		amount = effMax
	}
	d.hp = amount
	d.Send(packetPlayerStatChange(true, constant.HpID, int32(amount)))
	d.MarkDirty(DirtyHP, 500*time.Millisecond)

	if d.party != nil {
		d.party.broadcast(packetPlayerHpChange(d.ID, int32(d.hp), int32(d.maxHP)))
	}
}

func (d *Player) setMaxHP(amount int16) {
	if amount > constant.MaxHpValue {
		amount = constant.MaxHpValue
	}
	d.maxHP = amount
	d.Send(packetPlayerStatChange(true, constant.MaxHpID, int32(amount)))
	d.MarkDirty(DirtyMaxHP, 500*time.Millisecond)
}

func (d *Player) giveMaxHP(amount int16) {
	d.setMaxHP(d.maxHP + amount)
}

func (d *Player) giveMP(amount int16) {
	target := int(d.mp) + int(amount)
	maxAllowed := int(d.effectiveMaxMP())
	if target < 0 {
		target = 0
	} else if target > maxAllowed {
		target = maxAllowed
	}
	d.setMP(int16(target))
}

func (d *Player) setMP(amount int16) {
	if amount < 0 {
		amount = 0
	}
	effMax := d.effectiveMaxMP()
	if amount > constant.MaxMpValue {
		amount = constant.MaxMpValue
	}
	if amount > effMax {
		amount = effMax
	}
	d.mp = amount
	d.Send(packetPlayerStatChange(true, constant.MpID, int32(amount)))
	d.MarkDirty(DirtyMP, 500*time.Millisecond)
}

func (d *Player) setMaxMP(amount int16) {
	if amount > constant.MaxMpValue {
		amount = constant.MaxMpValue
	}
	d.maxMP = amount
	d.Send(packetPlayerStatChange(true, constant.MaxMpID, int32(amount)))
	d.MarkDirty(DirtyMaxMP, 500*time.Millisecond)
}

func (d *Player) giveMaxMP(amount int16) {
	d.setMaxMP(d.maxMP + amount)
}

func (d *Player) effectiveMaxHP() int16 {

	var itemsHP int32 = 0
	for _, item := range d.equip {
		if item.hp == 0 || item.slotID > 0 {
			continue
		}

		itemsHP += int32(item.hp)
	}

	base32 := int32(d.maxHP)
	hpPct, _ := d.hyperBodyPercents()

	res := base32 + itemsHP

	if hpPct > 0 {
		res += (base32 * int32(hpPct)) / 100
	}
	if res > int32(constant.MaxHpValue) {
		res = int32(constant.MaxHpValue)
	}
	if res < 0 {
		res = 0
	}
	return int16(res)
}

func (d *Player) effectiveMaxMP() int16 {

	var itemsMP int32 = 0
	for _, item := range d.equip {
		if item.mp == 0 || item.slotID > 0 {
			continue
		}
		itemsMP += int32(item.mp)
	}

	base32 := int32(d.maxMP)
	_, mpPct := d.hyperBodyPercents()
	res := base32 + itemsMP
	if mpPct > 0 {
		res += (base32 * int32(mpPct)) / 100
	}
	if res > int32(constant.MaxMpValue) {
		res = int32(constant.MaxMpValue)
	}
	if res < 0 {
		res = 0
	}
	return int16(res)
}

func (d *Player) hyperBodyPercents() (int16, int16) {
	if d == nil || d.buffs == nil {
		return 0, 0
	}
	lvl, ok := d.buffs.activeSkillLevels[int32(skill.HyperBody)]
	if !ok || lvl == 0 {
		return 0, 0
	}
	data, err := nx.GetPlayerSkill(int32(skill.HyperBody))
	if err != nil || int(lvl) < 1 || int(lvl) > len(data) {
		return 0, 0
	}
	sl := data[lvl-1]
	return int16(sl.X), int16(sl.Y)
}

func (d *Player) setFame(amount int16) {
	d.fame = amount
	d.Send(packetPlayerStatChange(true, constant.FameID, int32(amount)))

	_, err := common.DB.Exec("UPDATE characters SET fame=? WHERE ID=?", d.fame, d.ID)
	if err != nil {
		log.Printf("setFame: failed to save fame for character %d: %v", d.ID, err)
	}
}

func (d *Player) setMesos(amount int32) {
	d.mesos = amount
	d.Send(packetPlayerStatChange(true, constant.MesosID, amount))
	// write-behind instead of immediate DB write
	d.MarkDirty(DirtyMesos, 200*time.Millisecond)
}

func (d *Player) giveMesos(amount int32) {
	d.setMesos(d.mesos + amount)
}

func (d *Player) takeMesos(amount int32) {
	d.setMesos(d.mesos - amount)
}

func (d *Player) saveMesos() error {
	query := "UPDATE characters SET mesos=? WHERE accountID=? and Name=?"
	_, err := common.DB.Exec(query, d.mesos, d.accountID, d.Name)
	return err
}

func (d *Player) setHair(id int32) error {
	query := "UPDATE characters SET hair=? WHERE ID=?"
	_, err := common.DB.Exec(query, id, d.ID)
	d.hair = id
	d.Send(packetPlayerStatChange(true, constant.HairID, id))
	return err
}

func (d *Player) setFace(id int32) error {
	query := "UPDATE characters SET face=? WHERE ID=?"
	_, err := common.DB.Exec(query, id, d.ID)
	d.face = id
	d.Send(packetPlayerStatChange(true, constant.FaceID, id))
	return err
}

func (d *Player) setSkin(id byte) error {
	query := "UPDATE characters SET skin=? WHERE ID=?"
	_, err := common.DB.Exec(query, id, d.ID)
	d.skin = id
	d.Send(packetPlayerStatChange(true, constant.SkinID, int32(id)))
	return err
}

func (d *Player) setPartnerID(id int32) error {
	_, err := common.DB.Exec("UPDATE characters SET partnerID=? WHERE ID=?", id, d.ID)
	if err == nil {
		d.partnerID = id
	}
	return err
}

func (d *Player) setMarriageItemID(id int32) error {
	_, err := common.DB.Exec("UPDATE characters SET marriageItemID=? WHERE ID=?", id, d.ID)
	if err == nil {
		d.marriageItemID = id
	}
	return err
}

func (d *Player) setDivorceUntil(ts int64) error {
	_, err := common.DB.Exec("UPDATE characters SET divorceUntil=? WHERE ID=?", ts, d.ID)
	if err == nil {
		d.divorceUntil = ts
	}
	return err
}

func (d *Player) married() bool {
	return d.partnerID > 0 && d.marriageItemID >= 1112800 && d.marriageItemID <= 1112809
}

func (d *Player) underMarriageCooldown() bool {
	return d.divorceUntil > time.Now().Unix()
}

// UpdateMovement - update Data from position data
func (d *Player) UpdateMovement(frag movementFrag) {
	d.pos.x = frag.x
	d.pos.y = frag.y
	d.pos.foothold = frag.foothold
	d.stance = frag.stance
	if d.inst != nil {
		d.inst.refreshPairRingEffectsFor(d)
	}
}

// SetPos of Data
func (d *Player) SetPos(pos pos) {
	d.pos = pos
}

// checks Data is within a certain range of a position
func (d Player) checkPos(pos pos, xRange, yRange int16) bool {
	var xValid, yValid bool

	if xRange == 0 {
		xValid = d.pos.x == pos.x
	} else {
		xValid = (pos.x-xRange < d.pos.x && d.pos.x < pos.x+xRange)
	}

	if yRange == 0 {
		xValid = d.pos.y == pos.y
	} else {
		yValid = (pos.y-yRange < d.pos.y && d.pos.y < pos.y+yRange)
	}

	return xValid && yValid
}

func isExcludedMap(id int32) bool {
	// Free Market range (inclusive)
	return id >= 910000000 && id <= 910000022
}

func (d *Player) setMapID(id int32) {
	// Never set previousMap to a FM ID
	if !isExcludedMap(d.mapID) {
		d.previousMap = d.mapID
	}

	d.mapID = id

	if d.party != nil {
		d.UpdatePartyInfo(d.party.ID, d.ID, int32(d.job), int32(d.level), d.mapID, d.Name)
	}

	// write-behind for mapID/pos (mapPos updated on save())
	d.MarkDirty(DirtyMap, 500*time.Millisecond)
	d.MarkDirty(DirtyPrevMap, 500*time.Millisecond)
}

func (d Player) noChange() {
	d.Send(packetInventoryNoChange())
}

func (d *Player) GetNX() int32 {
	return d.nx
}

func (d *Player) SetNX(nx int32) {
	d.nx = nx
	d.MarkDirty(DirtyNX, 300*time.Millisecond)
}

func (d *Player) GetMaplePoints() int32 {
	return d.maplepoints
}

func (d *Player) SetMaplePoints(points int32) {
	d.maplepoints = points
	d.MarkDirty(DirtyMaplePoints, 300*time.Millisecond)
}

// addStackableItemToInventory handles adding stackable items to a specific inventory type
// It merges with existing stacks when possible and creates new slots as needed
// Note: If inventory becomes full mid-operation, already-added items remain (potential partial reward issue)
func (d *Player) addStackableItemToInventory(newItem Item, items *[]Item, slotSize byte, findSlotFunc func([]Item, byte) (int16, error)) error {
	slotMax := getItemSlotMax(newItem.ID)
	remaining := newItem.amount

	for i := range *items {
		if remaining == 0 {
			break
		}
		if (*items)[i].ID != newItem.ID || (*items)[i].amount >= slotMax {
			continue
		}

		canAdd := slotMax - (*items)[i].amount
		if canAdd > remaining {
			canAdd = remaining
		}

		(*items)[i].amount += canAdd
		remaining -= canAdd
		(*items)[i].save(d.ID)
		d.Send(packetInventoryAddItem((*items)[i], true))
	}

	for remaining > 0 {
		slotID, err := findSlotFunc(*items, slotSize)
		if err != nil {
			return err
		}

		newSlotAmount := remaining
		if newSlotAmount > slotMax {
			newSlotAmount = slotMax
		}

		newSlotItem := newItem
		newSlotItem.dbID = 0
		newSlotItem.amount = newSlotAmount
		newSlotItem.slotID = slotID
		newSlotItem.save(d.ID)
		*items = append(*items, newSlotItem)
		d.Send(packetInventoryAddItem(newSlotItem, true))
		remaining -= newSlotAmount
	}

	return nil
}

// GiveItem grants the given item to a player and returns the item
func (d *Player) GiveItem(newItem Item) (Item, error) { // TODO: Refactor
	isRechargeable := func(itemID int32) bool {
		base := itemID / 10000
		return base == 207
	}

	newItem.dbID = 0
	if newItem.cash {
		newItem.EnsureCashMetadata(0, 0)
	}

	findFirstEmptySlot := func(items []Item, size byte) (int16, error) {
		slotsUsed := make([]bool, size)
		for _, v := range items {
			if v.slotID > 0 {
				slotsUsed[v.slotID-1] = true
			}
		}
		slot := 0
		for i, v := range slotsUsed {
			if !v {
				slot = i + 1
				break
			}
		}
		if slot == 0 {
			slot = len(slotsUsed) + 1
		}
		if byte(slot) > size {
			return 0, fmt.Errorf("No empty Item slot left")
		}
		return int16(slot), nil
	}

	switch newItem.invID {
	case constant.InventoryEquip: // Equip
		slotID, err := findFirstEmptySlot(d.equip, d.equipSlotSize)
		if err != nil {
			return Item{}, err
		}
		newItem.slotID = slotID
		newItem.amount = 1
		newItem.save(d.ID)
		d.equip = append(d.equip, newItem)
		d.Send(packetInventoryAddItem(newItem, true))

	case constant.InventoryUse: // Use
		if isRechargeable(newItem.ID) {
			slotID, err := findFirstEmptySlot(d.use, d.useSlotSize)
			if err != nil {
				return Item{}, err
			}
			newItem.slotID = slotID
			newItem.save(d.ID)
			d.use = append(d.use, newItem)
			d.Send(packetInventoryAddItem(newItem, true))
			return newItem, err
		}

		// Non-rechargeable stackable items
		if err := d.addStackableItemToInventory(newItem, &d.use, d.useSlotSize, findFirstEmptySlot); err != nil {
			return Item{}, err
		}

	case constant.InventorySetup: // Set-up
		slotID, err := findFirstEmptySlot(d.setUp, d.setupSlotSize)
		if err != nil {
			return Item{}, err
		}
		newItem.slotID = slotID
		newItem.save(d.ID)
		d.setUp = append(d.setUp, newItem)
		d.Send(packetInventoryAddItem(newItem, true))

	case constant.InventoryEtc: // Etc
		if err := d.addStackableItemToInventory(newItem, &d.etc, d.etcSlotSize, findFirstEmptySlot); err != nil {
			return Item{}, err
		}

	case constant.InventoryCash: // Cash
		slotID, err := findFirstEmptySlot(d.cash, d.cashSlotSize)
		if err != nil {
			return Item{}, err
		}
		newItem.slotID = slotID
		newItem.save(d.ID)
		d.cash = append(d.cash, newItem)
		d.Send(packetInventoryAddItem(newItem, true))

	default:
		return Item{}, fmt.Errorf("Unknown inventory ID: %d", newItem.invID)
	}
	if newItem.ringID > 0 {
		d.refreshRingRecords()
	}

	return newItem, nil
}

func (p *Player) CanReceiveItems(items []Item) bool {
	invCounts := map[byte]int{}
	for _, item := range items {
		invType := byte(item.ID / 1000000)
		invCounts[invType]++
	}

	for invType, needed := range invCounts {
		var cur, max byte
		switch invType {
		case constant.InventoryEquip:
			cur = byte(len(p.equip))
			max = p.equipSlotSize
		case constant.InventoryUse:
			cur = byte(len(p.use))
			max = p.useSlotSize
		case constant.InventorySetup:
			cur = byte(len(p.setUp))
			max = p.setupSlotSize
		case constant.InventoryEtc:
			cur = byte(len(p.etc))
			max = p.etcSlotSize
		case constant.InventoryCash:
			cur = byte(len(p.cash))
			max = p.cashSlotSize
		default:
			continue
		}
		if cur+byte(needed) > max {
			return false
		}
	}
	return true
}

func (d *Player) GetSlotSize(invID byte) int16 {
	switch invID {
	case constant.InventoryEquip:
		return int16(d.equipSlotSize)
	case constant.InventoryUse:
		return int16(d.useSlotSize)
	case constant.InventorySetup:
		return int16(d.setupSlotSize)
	case constant.InventoryEtc:
		return int16(d.etcSlotSize)
	case constant.InventoryCash:
		return int16(d.cashSlotSize)
	}

	return constant.InventoryBaseSlotSize
}

// TakeItem removes an item from the player's inventory
func (d *Player) takeItem(id int32, slot int16, amount int16, invID byte) (Item, error) {
	item, err := d.getItem(invID, slot)
	if err != nil {
		return item, fmt.Errorf("item not found at inv=%d slot=%d", invID, slot)
	}

	if item.ID != id {
		return item, fmt.Errorf("Item.ID(%d) does not match ID(%d) provided", item.ID, id)
	}
	if item.invID != invID {
		return item, fmt.Errorf("inventory ID mismatch: item.invID(%d) vs provided invID(%d)", item.invID, invID)
	}
	if amount <= 0 {
		return item, fmt.Errorf("invalid amount requested: %d", amount)
	}

	if amount > item.amount {
		return item, fmt.Errorf("insufficient quantity: have=%d requested=%d", item.amount, amount)
	}

	item.amount -= amount
	if item.amount == 0 {
		d.removeItem(item, false)
	} else {
		d.updateItemStack(item, false)
	}

	return item, nil

}

func (d *Player) TakeItemSilent(id int32, slot int16, amount int16, invID byte) (Item, error) {
	item, err := d.getItem(invID, slot)
	if err != nil {
		return item, fmt.Errorf("item not found at inv=%d slot=%d", invID, slot)
	}

	if item.ID != id {
		return item, fmt.Errorf("Item.ID(%d) does not match ID(%d) provided", item.ID, id)
	}
	if item.invID != invID {
		return item, fmt.Errorf("inventory ID mismatch: item.invID(%d) vs provided invID(%d)", item.invID, invID)
	}
	if amount <= 0 {
		return item, fmt.Errorf("invalid amount requested: %d", amount)
	}

	if amount > item.amount {
		return item, fmt.Errorf("insufficient quantity: have=%d requested=%d", item.amount, amount)
	}

	item.amount -= amount
	if item.amount == 0 || item.isRechargeable() {
		d.removeItem(item, true)
	} else {
		d.updateItemStack(item, true)
	}

	return item, nil

}

func (d Player) updateItemStack(item Item, silent bool) {
	if item.amount <= 0 {
		d.removeItem(item, silent)
		return
	}

	item.save(d.ID)
	d.updateItem(item)

	if !silent {
		d.Send(packetInventoryModifyItemAmount(item))
	}
}

func (d *Player) updateItem(new Item) {
	var items = d.findItemInventory(new)

	for i, v := range items {
		if v.dbID == new.dbID {
			items[i] = new
			break
		}
	}
	d.updateItemInventory(new.invID, items)
	if new.ringID > 0 {
		d.refreshRingRecords()
	}
	if new.invID == constant.InventoryEquip && new.slotID < 0 {
		d.recalculateTotalStats()
	}
}

func (d *Player) updateItemInventory(invID byte, inventory []Item) {
	switch invID {
	case constant.InventoryEquip:
		d.equip = inventory
	case constant.InventoryUse:
		d.use = inventory
	case constant.InventorySetup:
		d.setUp = inventory
	case constant.InventoryEtc:
		d.etc = inventory
	case constant.InventoryCash:
		d.cash = inventory
	}
}

func betterInventorySlotWinner(current, candidate Item) bool {
	if current.amount <= 0 && candidate.amount > 0 {
		return true
	}
	if current.amount > 0 && candidate.amount <= 0 {
		return false
	}
	if candidate.dbID != current.dbID {
		return candidate.dbID > current.dbID
	}
	if candidate.ID != current.ID {
		return candidate.ID > current.ID
	}
	if candidate.amount != current.amount {
		return candidate.amount > current.amount
	}
	return candidate.slotID > current.slotID
}

func (d *Player) findItemInventory(item Item) []Item {
	switch item.invID {
	case constant.InventoryEquip:
		return d.equip
	case constant.InventoryUse:
		return d.use
	case constant.InventorySetup:
		return d.setUp
	case constant.InventoryEtc:
		return d.etc
	case constant.InventoryCash:
		return d.cash
	}

	return nil
}

func (d Player) getItem(invID byte, slotID int16) (Item, error) {
	var items []Item

	switch invID {
	case constant.InventoryEquip:
		items = d.equip
	case constant.InventoryUse:
		items = d.use
	case constant.InventorySetup:
		items = d.setUp
	case constant.InventoryEtc:
		items = d.etc
	case constant.InventoryCash:
		items = d.cash
	}

	var found Item
	hasFound := false
	for _, v := range items {
		if v.slotID == slotID {
			if !hasFound || betterInventorySlotWinner(found, v) {
				found = v
				hasFound = true
			}
		}
	}
	if hasFound {
		return found, nil
	}

	return Item{}, fmt.Errorf("Could not find Item")
}

// GetItem retrieves an item from the player's inventory (exported for use by other packages)
func (d *Player) GetItem(invID byte, slotID int16) (Item, error) {
	return d.getItem(invID, slotID)
}

// GetItemByCashID finds an item in the specified inventory by its cash shop ID
func (d *Player) GetItemByCashID(invID byte, cashID int64) (Item, int16, error) {
	switch invID {
	case 1:
		for i := range d.equip {
			if d.equip[i].cashID == cashID {
				return d.equip[i], d.equip[i].slotID, nil
			}
		}
	case 2:
		for i := range d.use {
			if d.use[i].cashID == cashID {
				return d.use[i], d.use[i].slotID, nil
			}
		}
	case 3:
		for i := range d.setUp {
			if d.setUp[i].cashID == cashID {
				return d.setUp[i], d.setUp[i].slotID, nil
			}
		}
	case 4:
		for i := range d.etc {
			if d.etc[i].cashID == cashID {
				return d.etc[i], d.etc[i].slotID, nil
			}
		}
	case 5:
		for i := range d.cash {
			if d.cash[i].cashID == cashID {
				return d.cash[i], d.cash[i].slotID, nil
			}
		}
	}
	return Item{}, 0, fmt.Errorf("item not found with cashID: %d", cashID)
}

func (d *Player) swapItems(item1, item2 Item, start, end int16) {
	if item1.dbID != 0 && item2.dbID != 0 {
		tx, err := common.DB.Begin()
		if err != nil {
			log.Println(err)
			d.Send(packetInventoryNoChange())
			return
		}

		tmpSlot := int32(1000000) + int32(item1.dbID)
		if _, err = tx.Exec("UPDATE items SET slotNumber=? WHERE ID=?", tmpSlot, item2.dbID); err != nil {
			_ = tx.Rollback()
			log.Println(err)
			d.Send(packetInventoryNoChange())
			return
		}
		if _, err = tx.Exec("UPDATE items SET slotNumber=? WHERE ID=?", end, item1.dbID); err != nil {
			_ = tx.Rollback()
			log.Println(err)
			d.Send(packetInventoryNoChange())
			return
		}
		if _, err = tx.Exec("UPDATE items SET slotNumber=? WHERE ID=?", start, item2.dbID); err != nil {
			_ = tx.Rollback()
			log.Println(err)
			d.Send(packetInventoryNoChange())
			return
		}
		if err = tx.Commit(); err != nil {
			log.Println(err)
			d.Send(packetInventoryNoChange())
			return
		}
	} else {
		item1.slotID = end
		item1.save(d.ID)
		item2.slotID = start
		item2.save(d.ID)
	}

	item1.slotID = end
	d.updateItem(item1)
	item2.slotID = start
	d.updateItem(item2)

	d.Send(packetInventoryChangeItemSlot(item1.invID, start, end))
}

func (d *Player) removeItem(item Item, fromStorage bool) {
	removed := make([]Item, 0, 2)

	filterBySlot := func(items []Item) []Item {
		out := items[:0]
		for _, v := range items {
			if v.slotID == item.slotID {
				removed = append(removed, v)
				continue
			}
			out = append(out, v)
		}
		return out
	}

	switch item.invID {
	case constant.InventoryEquip:
		d.equip = filterBySlot(d.equip)
	case constant.InventoryUse:
		d.use = filterBySlot(d.use)
	case constant.InventorySetup:
		d.setUp = filterBySlot(d.setUp)
	case constant.InventoryEtc:
		d.etc = filterBySlot(d.etc)
	case constant.InventoryCash:
		d.cash = filterBySlot(d.cash)
	}

	for _, removedItem := range removed {
		if removedItem.dbID == 0 {
			continue
		}
		if err := removedItem.delete(); err != nil {
			log.Println(err)
		}
	}

	if !fromStorage {
		d.Send(packetInventoryRemoveItem(item))
		if item.invID == constant.InventoryEquip && item.slotID < 0 && d.inst != nil {
			d.inst.broadcastAvatarChange(d)
			d.recalculateTotalStats()
		}
	}
	if item.ringID > 0 {
		d.refreshRingRecords()
	}
}

func (d *Player) dropMesos(amount int32) error {
	if d.mesos < amount {
		return errors.New("not enough mesos")
	}

	d.takeMesos(amount)
	d.inst.dropPool.createDrop(dropSpawnNormal, dropFreeForAll, amount, d.pos, true, true, d.ID, d.ID)

	return nil
}

func (d *Player) moveItem(start, end, amount int16, invID byte) error {
	isRechargeable := func(itemID int32) bool {
		base := itemID / 10000
		return base == 207
	}

	if end == 0 { // drop item
		item, err := d.getItem(invID, start)
		if err != nil {
			return fmt.Errorf("Item to move doesn't exist")
		}

		if isRechargeable(item.ID) {
			amount = item.amount
		}

		dropItem := item
		dropItem.amount = amount
		dropItem.dbID = 0

		// Special case: rechargeable items with 0 amount can't use takeItem (it requires amount > 0)
		// Just remove them directly
		if item.isRechargeable() && item.amount == 0 {
			d.removeItem(item, false)
		} else {
			takenItem, err := d.takeItem(item.ID, item.slotID, amount, item.invID)
			if err != nil {
				return fmt.Errorf("unable to take Item")
			}

			// For rechargeable items that reach 0 amount after takeItem, explicitly remove them
			// to prevent duplication exploit (takeItem keeps them at 0 for skill usage)
			if takenItem.isRechargeable() && takenItem.amount == 0 {
				d.removeItem(takenItem, false)
			}
		}

		d.inst.dropPool.createDrop(dropSpawnNormal, dropFreeForAll, 0, d.pos, true, true, d.ID, 0, dropItem)

		return nil
	}

	if end < 0 { // move to equip slot
		item1, err := d.getItem(invID, start)
		if err != nil {
			return fmt.Errorf("Item to move doesn't exist")
		}

		if item1.twoHanded {
			if _, err := d.getItem(invID, -10); err == nil {
				d.Send(packetInventoryNoChange())
				return nil
			}
		} else if item1.shield() {
			if weapon, err := d.getItem(invID, -11); err == nil && weapon.twoHanded {
				d.Send(packetInventoryNoChange())
				return nil
			}
		}

		item2, err := d.getItem(invID, end)
		if err == nil {
			if item1.dbID != 0 && item2.dbID != 0 {
				tx, txErr := common.DB.Begin()
				if txErr != nil {
					return txErr
				}

				tmpSlot := int32(1000000) + int32(item1.dbID)
				if _, err = tx.Exec("UPDATE items SET slotNumber=? WHERE ID=?", tmpSlot, item2.dbID); err != nil {
					_ = tx.Rollback()
					return err
				}
				if _, err = tx.Exec("UPDATE items SET slotNumber=? WHERE ID=?", end, item1.dbID); err != nil {
					_ = tx.Rollback()
					return err
				}
				if _, err = tx.Exec("UPDATE items SET slotNumber=? WHERE ID=?", start, item2.dbID); err != nil {
					_ = tx.Rollback()
					return err
				}
				err = tx.Commit()
				if err != nil {
					return err
				}
			} else {
				item2.slotID = start
				item2.save(d.ID)
				item1.slotID = end
				item1.save(d.ID)
			}

			item2.slotID = start
			d.updateItem(item2)
			item1.slotID = end
		} else {
			item1.slotID = end
			item1.save(d.ID)
		}
		d.updateItem(item1)

		d.Send(packetInventoryChangeItemSlot(invID, start, end))
		d.inst.broadcastAvatarChange(d)
		return nil
	}

	item1, err := d.getItem(invID, start)
	if err != nil {
		return fmt.Errorf("Item to move doesn't exist")
	}

	item2, err := d.getItem(invID, end)
	if err != nil { // empty slot, simple move
		item1.slotID = end
		item1.save(d.ID)
		d.updateItem(item1)
		d.Send(packetInventoryChangeItemSlot(invID, start, end))
	} else { // destination occupied
		if (item1.isStackable() && item2.isStackable()) && (item1.ID == item2.ID) {
			slotMax := getItemSlotMax(item1.ID)
			if item1.amount == slotMax || item2.amount == slotMax { // swap
				d.swapItems(item1, item2, start, end)
			} else if item2.amount < slotMax { // try full merge
				if item2.amount+item1.amount <= slotMax {
					item2.amount = item2.amount + item1.amount
					item2.save(d.ID)
					d.updateItem(item2)
					d.Send(packetInventoryAddItem(item2, false))

					d.removeItem(item1, false)
				} else {
					d.swapItems(item1, item2, start, end)
				}
			}
		} else {
			d.swapItems(item1, item2, start, end)
		}
	}

	if start < 0 || end < 0 {
		d.inst.broadcastAvatarChange(d)
		// Recalculate stats when equipment changes
		d.recalculateTotalStats()
	}

	return nil
}

func (d *Player) updateSkill(updatedSkill playerSkill) {
	if d.skills == nil {
		d.skills = make(map[int32]playerSkill)
	}
	d.skills[updatedSkill.ID] = updatedSkill
	d.Send(packetPlayerSkillBookUpdate(updatedSkill.ID, int32(updatedSkill.Level)))
	d.MarkDirty(DirtySkills, 800*time.Millisecond)
}

func (d *Player) removeSkill(skillID int32) {
	delete(d.skills, skillID)
	d.Send(packetPlayerSkillBookUpdate(skillID, 0))
	if _, err := common.DB.Exec("DELETE FROM skills WHERE characterID=? AND skillID=?", d.ID, skillID); err != nil {
		log.Printf("removeSkill delete failed: player=%s id=%d skillID=%d err=%v", d.Name, d.ID, skillID, err)
	}
	d.MarkDirty(DirtySkills, 800*time.Millisecond)
}

func (d *Player) useSkill(id int32, level byte, projectileID int32) error {
	skillInfo, _ := nx.GetPlayerSkill(id)

	skillUsed, ok := d.skills[id]
	if !ok {
		return nil
	}

	if d.buffs.hasMobDebuff(skill.Mob.Seal) && !d.admin() {
		return errors.New("character is currently sealed")
	}

	if skillUsed.Level != level {
		d.Conn.Send(packetMessageRedText("skill level mismatch"))
		return errors.New("skill level mismatch")
	}

	idx := int(skillUsed.Level) - 1
	if idx < 0 || idx >= len(skillInfo) {
		d.Conn.Send(packetMessageRedText("invalid skill data"))
		return errors.New("invalid skill data index")
	}
	si := skillInfo[idx]

	// Resource costs
	if si.MpCon > 0 {
		d.giveMP(-int16(si.MpCon))
	}
	if si.HpCon > 0 {
		d.giveHP(-int16(si.HpCon))
	}
	if si.MoneyConsume > 0 {
		d.takeMesos(int32(si.MoneyConsume))
	}

	if si.ItemCon > 0 {
		itemID := int32(si.ItemCon)
		need := int32(si.ItemConNo)
		if need <= 0 {
			need = 1
		}
		if !d.consumeItemsByID(itemID, need) {
			d.Conn.Send(packetMessageRedText("not enough items to use this skill"))
			return errors.New("not enough required items")
		}
	}

	if projectileID > 0 {
		need := int32(si.BulletConsume)
		if need <= 0 {
			need = int32(si.BulletCount)
		}
		if need > 0 {
			if !d.consumeItemsByID(projectileID, need) {
				d.Conn.Send(packetMessageRedText("not enough projectiles to use this skill"))
				return errors.New("not enough projectiles")
			}
		}
	}

	return nil
}

func (d *Player) consumeItemsByID(itemID int32, reqCount int32) bool {
	if reqCount <= 0 {
		return true
	}
	remaining := reqCount

	drain := func(invID byte, items []Item) {
		for i := range items {
			if remaining == 0 {
				return
			}
			it := items[i]
			if it.ID != itemID || it.amount <= 0 {
				continue
			}
			take := it.amount
			if int32(take) > remaining {
				take = int16(remaining)
			}
			if _, err := d.takeItem(itemID, it.slotID, take, invID); err == nil {
				remaining -= int32(take)
			}
		}
	}
	// Order: USE, SETUP, ETC, CASH
	drain(2, d.use)
	drain(3, d.setUp)
	drain(4, d.etc)
	drain(5, d.cash)

	return remaining == 0
}

func (d Player) admin() bool { return d.Conn.GetAdminLevel() > 0 }

func (d Player) avatarLookBytes() []byte {
	pkt := mpacket.NewPacket()
	pkt.WriteByte(d.gender)
	pkt.WriteByte(d.skin)
	pkt.WriteInt32(d.face)
	pkt.WriteByte(0)
	pkt.WriteInt32(d.hair)

	visible := make(map[byte]int32)
	masked := make(map[byte]int32)
	cashWeapon := int32(0)
	unknown := make([]string, 0)

	for _, b := range d.equip {
		if b.slotID < 0 && b.slotID > -20 {
			visible[byte(-b.slotID)] = b.ID
		}
	}

	for _, b := range d.equip {
		if b.slotID >= 0 {
			continue
		}

		if b.slotID > -100 {
			continue
		}

		if b.slotID == -111 {
			cashWeapon = b.ID
			continue
		}

		cashSlot := byte(-(b.slotID + 100))
		if cashSlot >= 1 && cashSlot <= 29 {
			if equipped, ok := visible[cashSlot]; ok {
				masked[cashSlot] = equipped
			}
			visible[cashSlot] = b.ID
		} else {
			unknown = append(unknown, fmt.Sprintf("slot=%d item=%d", b.slotID, b.ID))
		}
	}

	applyStarterMapVisualEquips(d.mapID, visible, masked)

	for slot := byte(1); slot <= 29; slot++ {
		if itemID, ok := visible[slot]; ok {
			pkt.WriteByte(slot)
			pkt.WriteInt32(itemID)
		}
	}

	pkt.WriteByte(0xFF)

	for slot := byte(1); slot <= 29; slot++ {
		if itemID, ok := masked[slot]; ok {
			pkt.WriteByte(slot)
			pkt.WriteInt32(itemID)
		}
	}

	pkt.WriteByte(0xFF)
	pkt.WriteInt32(cashWeapon)
	pkt.WriteInt32(0)

	return pkt
}

func (d Player) encodeDisplayBytes(pkt *mpacket.Packet) {
	pkt.WriteBytes(d.avatarLookBytes())
}

func (d Player) remoteSpawnTempStatMask() uint64 {
	if d.buffs == nil {
		return 0
	}
	return d.buffs.remoteSpawnMask()
}

func (d Player) encodeRemoteMiniRoomBalloon(pkt *mpacket.Packet) {
	if d.inst == nil {
		pkt.WriteByte(0)
		return
	}

	r, err := d.inst.roomPool.getPlayerRoom(d.ID)
	if err != nil {
		pkt.WriteByte(0)
		return
	}

	b, ok := r.(boxDisplayer)
	if !ok {
		pkt.WriteByte(0)
		return
	}

	display := b.displayBytes()
	if len(display) < 5 {
		pkt.WriteByte(0)
		return
	}

	// UserEnterField carries the room balloon without the leading owner character id.
	pkt.WriteBytes(display[4:])
}

// Logout flushes coalesced state and does a full checkpoint save.
func (d Player) Logout() {
	if d.inst != nil {
		if pos, err := d.inst.calculateNearestSpawnPortalID(d.pos); err == nil {
			d.mapPos = pos
		}
	}

	flushNow(&d)
	d.saveBuffSnapshot()

	if err := d.save(); err != nil {
		log.Printf("Player(%d) logout save failed: %v", d.ID, err)
	}

}

func (d *Player) FlushState() {
	flushNow(d)
}

func (d *Player) Kick() {
	if d == nil || d.Conn == nil {
		return
	}

	_ = d.Conn.Close()
}

// Save data - this needs to be split to occur at relevant points in time
func (d Player) save() error {
	query := `UPDATE characters set skin=?, hair=?, face=?, level=?,
	job=?, str=?, dex=?, intt=?, luk=?, hp=?, maxHP=?, mp=?, maxMP=?,
	ap=?, sp=?, exp=?, fame=?, mapID=?, mapPos=?, mesos=?, miniGameWins=?,
	miniGameDraw=?, miniGameLoss=?, miniGamePoints=?, buddyListSize=? WHERE ID=?`

	var mapPos byte
	var err error

	if d.inst != nil {
		mapPos, err = d.inst.calculateNearestSpawnPortalID(d.pos)
	}
	if err != nil {
		return err
	}

	if nxMap, nxErr := nx.GetMap(d.mapID); nxErr == nil && len(nxMap.Portals) > 0 {
		if int(mapPos) < 0 || int(mapPos) >= len(nxMap.Portals) {
			mapPos = 0
		}
	}

	d.mapPos = mapPos

	_, err = common.DB.Exec(query,
		d.skin, d.hair, d.face, d.level, d.job, d.str, d.dex, d.intt, d.luk, d.hp, d.maxHP, d.mp,
		d.maxMP, d.ap, d.sp, d.exp, d.fame, d.mapID, d.mapPos, d.mesos, d.miniGameWins,
		d.miniGameDraw, d.miniGameLoss, d.miniGamePoints, d.buddyListSize, d.ID)
	if err != nil {
		return err
	}

	query = `INSERT INTO skills(characterID,skillID,level,cooldown)
	         VALUES(?,?,?,?)
	         ON DUPLICATE KEY UPDATE level=VALUES(level), cooldown=VALUES(cooldown)`
	for skillID, skill := range d.skills {
		if _, err := common.DB.Exec(query, d.ID, skillID, skill.Level, skill.Cooldown); err != nil {
			return err
		}
	}

	return nil
}

func (d *Player) damagePlayer(damage int16) {
	if damage <= 0 {
		return
	}

	newHP := d.hp - damage
	if newHP <= 0 {
		newHP = 0
		if d.level >= 10 && !d.removeItemsByID(constant.ItemSafetyCharm, 1, false) {
			percent := int32(10 - int(d.luk)/10)
			if percent < 5 {
				percent = 5
			}

			loss := (d.exp / 100) * percent
			if loss < 1 && d.exp > 0 {
				loss = 1
			}
			newExp := d.exp - loss
			if newExp < 0 {
				newExp = 0
			}
			d.setEXP(newExp)
		}
	}

	d.setHP(newHP)
}

func (d *Player) setInventorySlotSizes(equip, use, setup, etc, cash byte) {
	changed := (d.equipSlotSize != equip) || (d.useSlotSize != use) ||
		(d.setupSlotSize != setup) || (d.etcSlotSize != etc) || (d.cashSlotSize != cash)
	if !changed {
		return
	}
	d.equipSlotSize = equip
	d.useSlotSize = use
	d.setupSlotSize = setup
	d.etcSlotSize = etc
	d.cashSlotSize = cash
	d.MarkDirty(DirtyInvSlotSizes, 2*time.Second)
}

func (d *Player) IncreaseSlotSize(invID, amount byte) error {
	switch invID {
	case constant.InventoryEquip:
		if d.equipSlotSize+amount > constant.InventoryMaxSlotSize {
			return fmt.Errorf("cannot increase equip slot size beyond %d", constant.InventoryMaxSlotSize)
		}
		d.equipSlotSize += amount
	case constant.InventoryUse:
		if d.useSlotSize+amount > constant.InventoryMaxSlotSize {
			return fmt.Errorf("cannot increase use slot size beyond %d", constant.InventoryMaxSlotSize)
		}
		d.useSlotSize += amount
	case constant.InventorySetup:
		if d.setupSlotSize+amount > constant.InventoryMaxSlotSize {
			return fmt.Errorf("cannot increase setup slot size beyond %d", constant.InventoryMaxSlotSize)
		}
		d.setupSlotSize += amount
	case constant.InventoryEtc:
		if d.etcSlotSize+amount > constant.InventoryMaxSlotSize {
			return fmt.Errorf("cannot increase etc slot size beyond %d", constant.InventoryMaxSlotSize)
		}
		d.etcSlotSize += amount
	case constant.InventoryCash:
		if d.cashSlotSize+amount > constant.InventoryMaxSlotSize {
			return fmt.Errorf("cannot increase cash slot size beyond %d", constant.InventoryMaxSlotSize)
		}
		d.cashSlotSize += amount
	}

	d.MarkDirty(DirtyInvSlotSizes, 2*time.Second)
	return nil
}

func (d *Player) setBuddyListSize(size byte) {
	if d.buddyListSize == size {
		return
	}
	d.buddyListSize = size
	d.MarkDirty(DirtyBuddySize, 1*time.Second)
}

func (d *Player) addMiniGameWin() {
	d.miniGameWins++
	d.MarkDirty(DirtyMiniGame, 1*time.Second)
}

func (d *Player) addMiniGameDraw() {
	d.miniGameDraw++
	d.MarkDirty(DirtyMiniGame, 1*time.Second)
}

func (d *Player) addMiniGameLoss() {
	d.miniGameLoss++
	d.MarkDirty(DirtyMiniGame, 1*time.Second)
}

func (d *Player) addMiniGamePoints(delta int32) {
	d.miniGamePoints += delta
	d.MarkDirty(DirtyMiniGame, 1*time.Second)
}

func (d *Player) sendBuddyList() {
	d.Send(packetBuddyListSizeUpdate(d.buddyListSize))
	d.Send(packetBuddyInfo(d.buddyList))
}

func (d Player) buddyListFull() bool {
	count := 0
	for _, v := range d.buddyList {
		if v.status != 1 {
			count++
		}
	}

	return count >= int(d.buddyListSize)
}

func (d *Player) addOnlineBuddy(id int32, name string, channel int32) {
	if d.buddyListFull() {
		return
	}

	for i, v := range d.buddyList {
		if v.id == id {
			d.buddyList[i].status = 0
			d.buddyList[i].channelID = channel
			d.Send(packetBuddyUpdate(id, name, d.buddyList[i].status, channel, false))
			return
		}
	}

	newBuddy := buddy{id: id, name: name, status: 0, channelID: channel}

	d.buddyList = append(d.buddyList, newBuddy)
	d.Send(packetBuddyInfo(d.buddyList))
}

func (d *Player) addOfflineBuddy(id int32, name string) {
	if d.buddyListFull() {
		return
	}

	for i, v := range d.buddyList {
		if v.id == id {
			d.buddyList[i].status = 2
			d.buddyList[i].channelID = -1
			d.Send(packetBuddyUpdate(id, name, d.buddyList[i].status, -1, false))
			return
		}
	}

	newBuddy := buddy{id: id, name: name, status: 2, channelID: -1}

	d.buddyList = append(d.buddyList, newBuddy)
	d.Send(packetBuddyInfo(d.buddyList))
}

func (d Player) hasBuddy(id int32) bool {
	for _, v := range d.buddyList {
		if v.id == id {
			return true
		}
	}

	return false
}

func (d *Player) removeBuddy(id int32) {
	for i, v := range d.buddyList {
		if v.id == id {
			d.buddyList[i] = d.buddyList[len(d.buddyList)-1]
			d.buddyList = d.buddyList[:len(d.buddyList)-1]
			d.Send(packetBuddyInfo(d.buddyList))
			return
		}
	}
}

// removeEquipAtSlot removes the equip from the given slot (equipped negative or inventory positive).
func (d *Player) removeEquipAtSlot(slot int16) bool {
	if slot < 0 {
		// Equipped Item; find and clear
		for i := range d.equip {
			if d.equip[i].slotID == slot {
				// Remove equipped Item
				d.equip[i].amount = 0
				return true
			}
		}
		return false
	}

	// Inventory equip; remove from inventory
	for i := range d.equip {
		if d.equip[i].slotID == slot {
			if d.equip[i].amount != 1 {
				return false
			}
			d.equip[i].amount = 0
			return true
		}
	}
	return false
}

// findUseItemBySlot returns the use Item (scroll) at the given slot from USE inventory.
func (d *Player) findUseItemBySlot(slot int16) *Item {
	for i := range d.use {
		if d.use[i].slotID == slot {
			return &d.use[i]
		}
	}
	return nil
}

func (d *Player) findEtcItemBySlot(slot int16) *Item {
	for i := range d.etc {
		if d.etc[i].slotID == slot {
			return &d.etc[i]
		}
	}
	return nil
}

// findEquipBySlot returns the equip by slot (negative = equipped, positive = inventory slot).
func (d *Player) findEquipBySlot(slot int16) *Item {
	for i := range d.equip {
		if d.equip[i].slotID == slot {
			return &d.equip[i]
		}
	}
	return nil
}

func (plr *Player) recalculateTotalStats() {
	// Start with base stats
	plr.totalStr = plr.str
	plr.totalDex = plr.dex
	plr.totalInt = plr.intt
	plr.totalLuk = plr.luk
	plr.totalWatk = 0
	plr.totalMatk = 0
	plr.totalAccuracy = 0

	// Add bonuses from all equipped items
	for _, item := range plr.equip {
		if item.slotID < 0 {
			plr.totalStr += item.str
			plr.totalDex += item.dex
			plr.totalInt += item.intt
			plr.totalLuk += item.luk
			plr.totalWatk += item.watk
			plr.totalMatk += item.matk
			plr.totalAccuracy += item.accuracy
		}
	}

	// Add base stat contributions to attack
	plr.totalWatk += plr.str / 10
	plr.totalMatk += plr.intt / 10

	if plr.buffs != nil {
		statBonus := plr.buffs.getStatBonuses()

		plr.totalWatk += statBonus.watk
		plr.totalMatk += statBonus.matk
		plr.totalAccuracy += statBonus.accuracy
	}
}

func LoadPlayerFromID(id int32, conn mnet.Client) Player {
	c := Player{}
	filter := "ID,accountID,worldID,Name,gender,skin,hair,face,level,job,str,dex,intt," +
		"luk,hp,maxHP,mp,maxMP,ap,sp, exp,fame,mapID,mapPos,previousMapID,mesos," +
		"equipSlotSize,useSlotSize,setupSlotSize,etcSlotSize,cashSlotSize,miniGameWins," +
		"miniGameDraw,miniGameLoss,miniGamePoints,buddyListSize,regTeleportRocks,vipTeleportRocks,partnerID,marriageItemID,divorceUntil,expCouponItemID,expCouponExpiresAt,dropCouponItemID,dropCouponExpiresAt"

	var regTeleportRocksStr, vipTeleportRocksStr sql.NullString
	var partnerID, marriageItemID, expCouponItemID, dropCouponItemID sql.NullInt32
	var divorceUntil sql.NullInt64
	var expCouponExpiresAt, dropCouponExpiresAt sql.NullInt64
	err := common.DB.QueryRow("SELECT "+filter+" FROM characters where ID=?", id).Scan(&c.ID,
		&c.accountID, &c.worldID, &c.Name, &c.gender, &c.skin, &c.hair, &c.face,
		&c.level, &c.job, &c.str, &c.dex, &c.intt, &c.luk, &c.hp, &c.maxHP, &c.mp,
		&c.maxMP, &c.ap, &c.sp, &c.exp, &c.fame, &c.mapID, &c.mapPos,
		&c.previousMap, &c.mesos, &c.equipSlotSize, &c.useSlotSize, &c.setupSlotSize,
		&c.etcSlotSize, &c.cashSlotSize, &c.miniGameWins, &c.miniGameDraw, &c.miniGameLoss,
		&c.miniGamePoints, &c.buddyListSize, &regTeleportRocksStr, &vipTeleportRocksStr, &partnerID, &marriageItemID, &divorceUntil, &expCouponItemID, &expCouponExpiresAt, &dropCouponItemID, &dropCouponExpiresAt)

	if err != nil {
		log.Println(err)
		return c
	}

	c.partnerID = -1
	if partnerID.Valid {
		c.partnerID = partnerID.Int32
	}
	c.marriageID = -1

	c.marriageItemID = -1
	if marriageItemID.Valid {
		c.marriageItemID = marriageItemID.Int32
	}

	c.divorceUntil = 0
	if divorceUntil.Valid {
		c.divorceUntil = divorceUntil.Int64
	}
	if expCouponItemID.Valid {
		c.expCouponItemID = expCouponItemID.Int32
	}
	if expCouponExpiresAt.Valid {
		c.expCouponExpiresAt = expCouponExpiresAt.Int64
	}
	if dropCouponItemID.Valid {
		c.dropCouponItemID = dropCouponItemID.Int32
	}
	if dropCouponExpiresAt.Valid {
		c.dropCouponExpiresAt = dropCouponExpiresAt.Int64
	}

	c.petCashID = 0

	if err := common.DB.QueryRow("SELECT username, nx, maplepoints FROM accounts WHERE accountID=?", c.accountID).Scan(&c.accountName, &c.nx, &c.maplepoints); err != nil {
		log.Printf("loadPlayerFromID: failed to fetch accountName for accountID=%d: %v", c.accountID, err)
	}

	c.skills = make(map[int32]playerSkill)

	for _, s := range getSkillsFromCharID(c.ID) {
		c.skills[s.ID] = s
	}

	nxMap, err := nx.GetMap(c.mapID)
	if err != nil {
		log.Println(err)
		return c
	}

	if nxMap.ForcedReturn != constant.InvalidMap {
		c.mapID = nxMap.ForcedReturn
		c.MarkDirty(DirtyMap, time.Millisecond*300)
		if m2, e2 := nx.GetMap(c.mapID); e2 == nil {
			nxMap = m2
		}
	}

	if int(c.mapPos) < 0 || int(c.mapPos) >= len(nxMap.Portals) {
		c.mapPos = 0

		if _, healErr := common.DB.Exec("UPDATE characters SET mapPos=? WHERE ID=?", c.mapPos, c.ID); healErr != nil {
			log.Printf("LoadPlayerFromID: failed to heal mapPos in DB for char %d: %v", c.ID, healErr)
		}
	}

	c.pos.x = nxMap.Portals[c.mapPos].X
	c.pos.y = nxMap.Portals[c.mapPos].Y

	c.equip, c.use, c.setUp, c.etc, c.cash = loadInventoryFromDb(c.ID, c.equipSlotSize, c.useSlotSize, c.setupSlotSize, c.etcSlotSize, c.cashSlotSize)
	c.refreshRingRecords()

	// Calculate total stats including equipment bonuses
	c.recalculateTotalStats()

	c.buddyList = getBuddyList(c.ID, c.buddyListSize)

	// Initialize teleport rocks - handle NULL values from database
	regRocksStr := ""
	if regTeleportRocksStr.Valid {
		regRocksStr = regTeleportRocksStr.String
	}
	vipRocksStr := ""
	if vipTeleportRocksStr.Valid {
		vipRocksStr = vipTeleportRocksStr.String
	}
	c.regTeleportRocks = parseTeleportRocks(regRocksStr, constant.TeleportRockRegSlots)
	c.vipTeleportRocks = parseTeleportRocks(vipRocksStr, constant.TeleportRockVIPSlots)
	c.funcKeyMap = loadFuncKeyMap(c.ID)

	c.quests = loadQuestsFromDB(c.ID)
	c.quests.init()
	c.quests.mobKills = loadQuestMobKillsFromDB(c.ID)

	// Initialize the per-Player buff manager so handlers can call plr.addBuff(...)
	NewCharacterBuffs(&c)
	c.cleanupExpiredRateCoupons(time.Now())

	c.storageInventory = new(storage)

	if err := c.storageInventory.load(c.accountID); err != nil {
		log.Printf("loadPlayerFromID: failed to load storage inventory for accountID=%d: %v", c.accountID, err)
	}

	c.quickslotKeys = loadQuickslotKeys(c.ID)

	c.Conn = conn

	return c
}

func loadFuncKeyMap(characterID int32) funcKeyMapState {
	state := funcKeyMapState{Entries: constant.DefaultFuncKeyMap(), Loaded: true}

	rows, err := common.DB.Query("SELECT tkey, type, action FROM keymap WHERE characterid=?", characterID)
	if err != nil {
		log.Printf("loadFuncKeyMap: characterID=%d err=%v", characterID, err)
		return state
	}
	defer rows.Close()

	for rows.Next() {
		var key int32
		var ty int32
		var action int32
		if err := rows.Scan(&key, &ty, &action); err != nil {
			log.Printf("loadFuncKeyMap scan: characterID=%d err=%v", characterID, err)
			continue
		}
		if key < 0 || key >= int32(len(state.Entries)) {
			continue
		}
		state.Entries[key] = constant.FuncKeyMapped{Type: byte(ty), Action: action}
	}

	state.Loaded = true
	return state
}

func saveFuncKeyMapEntry(characterID int32, index int32, mapped constant.FuncKeyMapped) {
	if _, err := common.DB.Exec(`
		INSERT INTO keymap (characterid, tkey, type, action)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE type=VALUES(type), action=VALUES(action)
	`, characterID, index, mapped.Type, mapped.Action); err != nil {
		log.Printf("saveFuncKeyMapEntry: characterID=%d index=%d err=%v", characterID, index, err)
	}
}

func packetFuncKeyMappedInit(state funcKeyMapState) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelFuncKeyMappedInit)
	if !state.Loaded {
		p.WriteByte(1)
		return p
	}

	p.WriteByte(0)
	for _, key := range state.Entries {
		p.WriteByte(key.Type)
		p.WriteInt32(key.Action)
	}
	return p
}

func loadQuickslotKeys(characterID int32) [2]int32 {
	var keys [2]int32
	if err := common.DB.QueryRow("SELECT key1, key2 FROM quickslot_keymap WHERE characterID=?", characterID).Scan(&keys[0], &keys[1]); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("loadQuickslotKeys: characterID=%d err=%v", characterID, err)
		}
	}
	return keys
}

func saveQuickslotKeys(characterID int32, keys [2]int32) {
	if _, err := common.DB.Exec(`
		INSERT INTO quickslot_keymap (characterID, key1, key2)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE key1=VALUES(key1), key2=VALUES(key2)
	`, characterID, keys[0], keys[1]); err != nil {
		log.Printf("saveQuickslotKeys: characterID=%d err=%v", characterID, err)
	}
}

func getBuddyList(playerID int32, buddySize byte) []buddy {
	buddies := make([]buddy, 0, buddySize)
	filter := "friendID,accepted"
	rows, err := common.DB.Query("SELECT "+filter+" FROM buddy where characterID=?", playerID)

	if err != nil {
		log.Fatal(err)
		return buddies
	}

	defer rows.Close()

	i := 0
	for rows.Next() {
		newBuddy := buddy{}

		var accepted bool
		rows.Scan(&newBuddy.id, &accepted)

		filter := "channelID,Name,inCashShop"
		err := common.DB.QueryRow("SELECT "+filter+" FROM characters where ID=?", newBuddy.id).Scan(&newBuddy.channelID, &newBuddy.name, &newBuddy.cashShop)

		if err != nil {
			log.Fatal(err)
			return buddies
		}

		if !accepted {
			newBuddy.status = 1 // pending buddy request
		} else if newBuddy.channelID == -1 {
			newBuddy.status = 2 // offline
		} else {
			newBuddy.status = 0 // online
		}

		buddies = append(buddies, newBuddy)

		i++
	}

	return buddies
}

func parseTeleportRocks(rocksStr string, size int) []int32 {
	rocks := make([]int32, size)
	for i := range rocks {
		rocks[i] = constant.InvalidMap
	}

	if rocksStr == "" {
		return rocks
	}

	parts := strings.Split(rocksStr, ",")
	for i, part := range parts {
		if i >= size {
			break
		}
		if part == "" {
			continue
		}
		var mapID int32
		if _, err := fmt.Sscanf(part, "%d", &mapID); err == nil {
			rocks[i] = mapID
		}
	}

	return rocks
}

func serializeTeleportRocks(rocks []int32) string {
	parts := make([]string, len(rocks))
	for i, mapID := range rocks {
		parts[i] = fmt.Sprintf("%d", mapID)
	}
	return strings.Join(parts, ",")
}

func (d *Player) addBuff(skillID int32, level byte, delay int16) {
	if d.buffs == nil {
		NewCharacterBuffs(d)
	}
	// Ensure the buff manager points to this exact Player instance (avoid stale copies).
	d.buffs.plr = d
	d.buffs.AddBuff(d.ID, skillID, level, false, delay)
}

func (d *Player) addForeignBuff(charId, skillID int32, level byte, delay int16) {
	d.buffs.AddBuff(charId, skillID, level, true, delay)
}

func (d *Player) addMobDebuff(skillID, level byte, durationSec int16) {
	d.buffs.AddMobDebuff(skillID, level, durationSec)
}

func (d *Player) removeAllCooldowns() {
	if d == nil || d.skills == nil {
		return
	}
	for _, ps := range d.skills {
		ps.Cooldown = 0
		ps.TimeLastUsed = 0
		d.updateSkill(ps)
	}
}

func (d *Player) saveBuffSnapshot() {
	if d.buffs == nil {
		return
	}

	// Ensure we don't snapshot already-stale buffs
	d.buffs.AuditAndExpireStaleBuffs()

	snaps := d.buffs.Snapshot()
	if len(snaps) == 0 {
		_, _ = common.DB.Exec("DELETE FROM character_buffs WHERE characterID=?", d.ID)
		return
	}

	tx, err := common.DB.Begin()
	if err != nil {
		log.Println("saveBuffSnapshot: begin tx:", err)
		return
	}
	defer func() { _ = tx.Commit() }()

	if _, err := tx.Exec("DELETE FROM character_buffs WHERE characterID=?", d.ID); err != nil {
		log.Println("saveBuffSnapshot: clear rows:", err)
		return
	}

	stmt, err := tx.Prepare("INSERT INTO character_buffs(characterID, sourceID, level, expiresAtMs) VALUES(?,?,?,?)")
	if err != nil {
		log.Println("saveBuffSnapshot: prepare:", err)
		return
	}
	defer stmt.Close()

	for _, s := range snaps {
		if _, err := stmt.Exec(d.ID, s.SourceID, s.Level, s.ExpiresAtMs); err != nil {
			log.Println("saveBuffSnapshot: insert:", err)
			return
		}
	}
}

func (d *Player) loadAndApplyBuffSnapshot() {
	rows, err := common.DB.Query("SELECT sourceID, level, expiresAtMs FROM character_buffs WHERE characterID=?", d.ID)
	if err != nil {
		log.Println("loadBuffSnapshot:", err)
		return
	}
	defer rows.Close()

	snaps := make([]BuffSnapshot, 0, 8)
	toDelete := make([]int32, 0, 8)

	now := time.Now().UnixMilli()
	for rows.Next() {
		var s BuffSnapshot
		if err := rows.Scan(&s.SourceID, &s.Level, &s.ExpiresAtMs); err != nil {
			log.Println("loadBuffSnapshot scan:", err)
			return
		}

		if s.ExpiresAtMs == 0 {
			snaps = append(snaps, s)
			continue
		}

		normalized := s.ExpiresAtMs
		if normalized > 0 && normalized < 1000000000000 {
			normalized *= 1000
		}

		if normalized <= 0 || normalized <= now {
			toDelete = append(toDelete, s.SourceID)
			continue
		}

		s.ExpiresAtMs = normalized
		snaps = append(snaps, s)
	}
	if err := rows.Err(); err != nil {
		log.Println("loadBuffSnapshot rows:", err)
		return
	}

	if len(toDelete) > 0 {
		var b strings.Builder
		b.WriteString("DELETE FROM character_buffs WHERE characterID=? AND sourceID IN (")
		for i := range toDelete {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteByte('?')
		}
		b.WriteByte(')')

		args := make([]interface{}, 0, 1+len(toDelete))
		args = append(args, d.ID)
		for _, sid := range toDelete {
			args = append(args, sid)
		}
		if _, err := common.DB.Exec(b.String(), args...); err != nil {
			log.Println("loadBuffSnapshot delete expired:", err)
		}
	}

	if len(snaps) > 0 {
		if d.buffs == nil {
			NewCharacterBuffs(d)
		}
		d.buffs.RestoreFromSnapshot(snaps)
	}
}

func (d *Player) resetDoorInfo() {
	d.doorMapID = 0
	d.doorSpawnID = 0
	d.doorPortalIndex = 0
	d.townDoorMapID = 0
	d.townDoorSpawnID = 0
	d.townPortalIndex = 0
}

// countItem returns total count across USE/SETUP/ETC for an Item ID.
func (d *Player) countItem(itemID int32) int32 {
	var total int32
	for i := range d.use {
		if d.use[i].ID == itemID {
			total += int32(d.use[i].amount)
		}
	}
	for i := range d.setUp {
		if d.setUp[i].ID == itemID {
			total += int32(d.setUp[i].amount)
		}
	}
	for i := range d.etc {
		if d.etc[i].ID == itemID {
			total += int32(d.etc[i].amount)
		}
	}
	for i := range d.cash {
		if d.cash[i].ID == itemID {
			total += int32(d.cash[i].amount)
		}
	}
	return total
}

// removeItemsByID removes up to reqCount across USE/SETUP/ETC. Returns true if fully removed.
func (d *Player) removeItemsByID(itemID int32, reqCount int32, silent bool) bool {
	if reqCount <= 0 {
		return true
	}
	remaining := reqCount

	drain := func(invID byte, items []Item) {
		for i := range items {
			if remaining == 0 {
				return
			}
			it := items[i]
			if it.ID != itemID || it.amount <= 0 {
				continue
			}
			take := int16(it.amount)
			if int32(take) > remaining {
				take = int16(remaining)
			}
			if silent {
				if _, err := d.TakeItemSilent(itemID, it.slotID, take, invID); err == nil {
					remaining -= int32(take)
				}
			} else {
				if _, err := d.takeItem(itemID, it.slotID, take, invID); err == nil {
					remaining -= int32(take)
				}
			}
		}
	}
	drain(2, d.use)
	drain(3, d.setUp)
	drain(4, d.etc)
	drain(5, d.cash)

	return remaining == 0
}

func (d *Player) meetsPrevQuestState(req nx.QuestStateReq) bool {
	switch req.State {
	case 2: // completed
		return d.quests.hasCompleted(req.ID)
	case 1: // in progress
		return d.quests.hasInProgress(req.ID)
	default:
		return true
	}
}

// meetsQuestBlock validates prereqs/Item counts.
func (d *Player) meetsQuestBlock(blk nx.CheckBlock) bool {
	if len(blk.Jobs) > 0 && !questJobAllowed(int32(d.job), blk.Job, blk.Jobs) {
		return false
	}
	if blk.LvMin > 0 && int32(d.level) < blk.LvMin {
		return false
	}
	if blk.LvMax > 0 && int32(d.level) > blk.LvMax {
		return false
	}
	if blk.Pop > 0 && int32(d.fame) < blk.Pop {
		return false
	}

	// Previous quest states
	for _, rq := range blk.PrevQuests {
		if rq.State > 0 && !d.meetsPrevQuestState(rq) {
			return false
		}
	}

	// Item possession/turn-in counts
	for _, it := range blk.Items {
		if it.Count > 0 && d.countItem(it.ID) < it.Count {
			return false
		}
	}
	return true
}

type questActPlan struct {
	adds    []Item
	removes []nx.ActItem
}

func (d *Player) buildQuestActPlan(act nx.ActBlock, selection int) (questActPlan, error) {
	plan := questActPlan{}
	selectable := make([]nx.ActItem, 0, 2)
	randomRewards := make([]nx.ActItem, 0, 4)
	totalWeight := int32(0)

	for _, ai := range act.Items {
		if !questActItemAllowed(ai, int32(d.job), d.gender) {
			continue
		}

		if ai.Count < 0 {
			plan.removes = append(plan.removes, ai)
			continue
		}
		if ai.Count == 0 {
			continue
		}

		switch {
		case ai.Prop < 0:
			selectable = append(selectable, ai)
		case ai.Prop > 0:
			randomRewards = append(randomRewards, ai)
			totalWeight += ai.Prop
		default:
			it, err := CreateItemFromID(ai.ID, int16(ai.Count))
			if err != nil {
				return questActPlan{}, err
			}
			plan.adds = append(plan.adds, it)
		}
	}

	if len(selectable) > 0 {
		if selection < 0 || selection >= len(selectable) {
			return questActPlan{}, fmt.Errorf("missing selectable quest reward choice")
		}
		selected := selectable[selection]
		it, err := CreateItemFromID(selected.ID, int16(selected.Count))
		if err != nil {
			return questActPlan{}, err
		}
		plan.adds = append(plan.adds, it)
	}

	if totalWeight > 0 {
		var seed int64
		_ = binary.Read(rand.Reader, binary.BigEndian, &seed)
		rng := mathrand.New(mathrand.NewSource(seed))
		roll := int32(rng.Int31n(totalWeight))
		cumulative := int32(0)

		for _, ai := range randomRewards {
			cumulative += ai.Prop
			if roll >= cumulative {
				continue
			}

			it, err := CreateItemFromID(ai.ID, int16(ai.Count))
			if err != nil {
				return questActPlan{}, err
			}
			plan.adds = append(plan.adds, it)
			break
		}
	}

	return plan, nil
}

func (d *Player) canReceiveQuestActItems(items []Item) bool {
	if len(items) == 0 {
		return true
	}

	type slotState struct {
		used      int
		max       int
		stackRoom map[int32]int32
	}

	states := map[byte]*slotState{
		constant.InventoryEquip: {used: len(d.equip), max: int(d.equipSlotSize), stackRoom: map[int32]int32{}},
		constant.InventoryUse:   {used: len(d.use), max: int(d.useSlotSize), stackRoom: map[int32]int32{}},
		constant.InventorySetup: {used: len(d.setUp), max: int(d.setupSlotSize), stackRoom: map[int32]int32{}},
		constant.InventoryEtc:   {used: len(d.etc), max: int(d.etcSlotSize), stackRoom: map[int32]int32{}},
		constant.InventoryCash:  {used: len(d.cash), max: int(d.cashSlotSize), stackRoom: map[int32]int32{}},
	}

	accumulate := func(inv byte, items []Item) {
		state := states[inv]
		for _, it := range items {
			if !it.isStackable() || it.amount <= 0 {
				continue
			}
			slotMax := int32(getItemSlotMax(it.ID))
			if slotMax <= 0 {
				slotMax = constant.MaxItemStack
			}
			state.stackRoom[it.ID] += slotMax - int32(it.amount)
		}
	}

	accumulate(constant.InventoryUse, d.use)
	accumulate(constant.InventoryEtc, d.etc)

	for _, it := range items {
		state := states[it.invID]
		if state == nil {
			return false
		}

		if !it.isStackable() || it.amount <= 0 {
			state.used++
			if state.used > state.max {
				return false
			}
			continue
		}

		remaining := int32(it.amount)
		if room := state.stackRoom[it.ID]; room > 0 {
			used := room
			if used > remaining {
				used = remaining
			}
			state.stackRoom[it.ID] -= used
			remaining -= used
		}

		if remaining <= 0 {
			continue
		}

		slotMax := int32(getItemSlotMax(it.ID))
		if slotMax <= 0 {
			slotMax = constant.MaxItemStack
		}
		neededSlots := int((remaining + slotMax - 1) / slotMax)
		state.used += neededSlots
		if state.used > state.max {
			return false
		}

		if leftover := neededSlots*int(slotMax) - int(remaining); leftover > 0 {
			state.stackRoom[it.ID] += int32(leftover)
		}
	}

	return true
}

func (d *Player) applyQuestAct(act nx.ActBlock, plan questActPlan, npcID int32, questID int16) error {
	for _, ai := range plan.removes {
		if !d.removeItemsByID(ai.ID, -ai.Count, false) {
			d.Send(packetQuestActionResult(constant.QuestActionFailedRetrieveEquippedItem, questID, npcID, nil))
			return fmt.Errorf("failed to remove required quest items ID=%d count=%d", ai.ID, -ai.Count)
		}
		d.Send(packetMessageDropPickUp(false, ai.ID, ai.Count))
	}

	for _, it := range plan.adds {
		given, err := d.GiveItem(it)
		if err != nil {
			d.Send(packetQuestActionResult(constant.QuestActionInventoryFull, questID, npcID, nil))
			return err
		}
		d.Send(packetMessageDropPickUp(false, given.ID, int32(given.amount)))
	}

	if act.Money != 0 {
		if act.Money > 0 {
			d.giveMesos(act.Money)
		} else {
			d.takeMesos(-act.Money)
		}
		d.Send(packetMessageMesosChangeChat(act.Money))
	}

	if act.Exp > 0 {
		d.giveEXP(act.Exp, false, false)
	}

	if act.Pop != 0 {
		d.setFame(d.fame + int16(act.Pop))
	}

	return nil
}

func questJobAllowed(playerJob int32, legacyJob int32, jobs []int32) bool {
	if len(jobs) == 0 {
		return legacyJob == 0 || playerJob == legacyJob
	}

	for _, job := range jobs {
		if job == 0 || playerJob == job {
			return true
		}
	}

	return false
}

func questActItemAllowed(ai nx.ActItem, playerJob int32, gender byte) bool {
	if !questJobAllowed(playerJob, ai.Job, ai.Jobs) {
		return false
	}

	return ai.Gender < 0 || ai.Gender == 2 || byte(ai.Gender) == gender
}

func (d *Player) canStartQuest(q nx.Quest) bool {
	if d.quests.hasInProgress(q.ID) || d.quests.hasCompleted(q.ID) {
		return false
	}
	return d.meetsQuestBlock(q.Start)
}

func (d *Player) canCompleteQuest(q nx.Quest) bool {
	if !d.quests.hasInProgress(q.ID) {
		return false
	}
	if !d.meetsQuestBlock(q.Complete) {
		return false
	}
	return d.meetsMobKills(q.ID, q.Complete.Mobs)
}

func (d *Player) questIDsForNPC(npcID int32) (available, inProgress, completable []int16) {
	ordered := make([]nx.Quest, 0, len(nx.GetQuests()))
	for _, q := range nx.GetQuests() {
		if q.Start.NPC != npcID && q.Complete.NPC != npcID {
			continue
		}
		ordered = append(ordered, q)
	}

	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Order == ordered[j].Order {
			return ordered[i].ID < ordered[j].ID
		}
		return ordered[i].Order < ordered[j].Order
	})

	for _, q := range ordered {
		switch {
		case q.Complete.NPC == npcID && d.canCompleteQuest(q):
			completable = append(completable, q.ID)
		case q.Start.NPC == npcID && d.canStartQuest(q):
			available = append(available, q.ID)
		case q.Complete.NPC == npcID && d.quests.hasInProgress(q.ID):
			inProgress = append(inProgress, q.ID)
		}
	}

	return available, inProgress, completable
}

func (d *Player) questDisplayName(questID int16) string {
	q, err := nx.GetQuest(questID)
	if err != nil {
		return strconv.Itoa(int(questID))
	}
	if strings.TrimSpace(q.Parent) != "" {
		return q.Parent
	}
	if strings.TrimSpace(q.Name) != "" {
		return q.Name
	}
	return strconv.Itoa(int(questID))
}

func (d *Player) questSayLines(questID int16, key string) []string {
	q, err := nx.GetQuest(questID)
	if err != nil || q.Say == nil {
		return nil
	}
	if len(q.Say[key]) == 0 {
		return nil
	}
	out := make([]string, len(q.Say[key]))
	copy(out, q.Say[key])
	return out
}

func (d *Player) questIncompleteLines(questID int16) []string {
	keys := []string{"complete.stop.mob", "complete.stop.item", "complete.stop.quest", "complete.stop.0", "complete.stop"}
	for _, key := range keys {
		if lines := d.questSayLines(questID, key); len(lines) > 0 {
			return lines
		}
	}
	return []string{"You have not met the requirements for this quest yet."}
}

func (d *Player) questSelectableRewards(questID int16) []string {
	q, err := nx.GetQuest(questID)
	if err != nil {
		return nil
	}
	choices := make([]string, 0, 2)
	for _, ai := range q.ActOnComplete.Items {
		if ai.Prop >= 0 || ai.Count <= 0 || !questActItemAllowed(ai, int32(d.job), d.gender) {
			continue
		}
		name := strconv.Itoa(int(ai.ID))
		if meta, err := nx.GetItem(ai.ID); err == nil && strings.TrimSpace(meta.Name) != "" {
			name = meta.Name
		}
		choices = append(choices, name)
	}
	return choices
}

// tryStartQuest validates NX Start requirements, starts quest, applies Act(0).
func (d *Player) tryStartQuest(questID int16) bool {
	return d.tryStartQuestSelection(questID, -1)
}

func (d *Player) tryStartQuestSelection(questID int16, selection int) bool {
	q, err := nx.GetQuest(questID)
	if err != nil {
		log.Printf("[Quest] start fail nx lookup: char=%s ID=%d err=%v", d.Name, questID, err)
		return false
	}

	if !d.canStartQuest(q) {
		return false
	}

	plan, err := d.buildQuestActPlan(q.ActOnStart, selection)
	if err != nil {
		return false
	}
	if !d.canReceiveQuestActItems(plan.adds) {
		d.Send(packetQuestActionResult(constant.QuestActionInventoryFull, questID, q.Start.NPC, nil))
		return false
	}

	if err := d.applyQuestAct(q.ActOnStart, plan, q.Start.NPC, questID); err != nil {
		return false
	}

	d.quests.add(questID, "")
	upsertQuestRecord(d.ID, questID, "")
	d.Send(packetQuestUpdate(questID, ""))

	var nextQuests []int16
	if q.ActOnStart.NextQuest != 0 {
		nextQuests = append(nextQuests, q.ActOnStart.NextQuest)
	}
	d.Send(packetQuestActionResult(constant.QuestActionSuccess, questID, q.Start.NPC, nextQuests))

	return true
}

func (d *Player) tryCompleteQuest(questID int16) bool {
	return d.tryCompleteQuestSelection(questID, -1)
}

func (d *Player) tryCompleteQuestSelection(questID int16, selection int) bool {
	q, err := nx.GetQuest(questID)
	if err != nil {
		log.Printf("[Quest] complete fail nx lookup: char=%s ID=%d err=%v", d.Name, questID, err)
		return false
	}

	if !d.canCompleteQuest(q) {
		return false
	}

	plan, err := d.buildQuestActPlan(q.ActOnComplete, selection)
	if err != nil {
		return false
	}
	if !d.canReceiveQuestActItems(plan.adds) {
		d.Send(packetQuestActionResult(constant.QuestActionInventoryFull, questID, q.Complete.NPC, nil))
		return false
	}

	if err := d.applyQuestAct(q.ActOnComplete, plan, q.Complete.NPC, questID); err != nil {
		return false
	}

	nowMs := time.Now().UnixMilli()
	d.quests.complete(questID, nowMs)
	setQuestCompleted(d.ID, questID, nowMs)
	clearQuestMobKills(d.ID, q.ID)

	d.Send(packetQuestUpdate(questID, ""))
	d.Send(packetQuestComplete(questID))

	var nextQuests []int16
	if q.ActOnComplete.NextQuest != 0 {
		nextQuests = append(nextQuests, q.ActOnComplete.NextQuest)
		_ = d.tryStartQuest(q.ActOnComplete.NextQuest)
	}

	d.Send(packetQuestActionResult(constant.QuestActionSuccess, questID, q.Complete.NPC, nextQuests))

	return true
}

func (d *Player) buildQuestKillString(q nx.Quest) string {
	if d.quests.mobKills == nil {
		return ""
	}
	counts := d.quests.mobKills[q.ID]
	if counts == nil {
		return ""
	}

	out := make([]byte, 0, len(q.Complete.Mobs)*3)
	for _, req := range q.Complete.Mobs {
		val := counts[req.ID]
		if val < 0 {
			val = 0
		}
		if val > 999 {
			val = 999
		}

		a := byte('0' + (val/100)%10)
		b := byte('0' + (val/10)%10)
		c := byte('0' + (val % 10))
		out = append(out, a, b, c)
	}
	return string(out)
}

func (d *Player) onMobKilled(mobID int32) {
	if d == nil {
		return
	}
	for qid := range d.quests.inProgress {
		q, err := nx.GetQuest(qid)
		if err != nil {
			continue
		}

		var needed int32
		for _, rm := range q.Complete.Mobs {
			if rm.ID == mobID {
				needed = rm.Count
				break
			}
		}
		if needed == 0 {
			continue
		}

		// Init maps
		if d.quests.mobKills == nil {
			d.quests.mobKills = make(map[int16]map[int32]int32)
		}
		if d.quests.mobKills[qid] == nil {
			d.quests.mobKills[qid] = make(map[int32]int32)
		}

		cur := d.quests.mobKills[qid][mobID]
		if cur < needed {
			d.quests.mobKills[qid][mobID] = cur + 1
			upsertQuestMobKill(d.ID, qid, mobID, 1)
		}

		d.Send(packetQuestUpdateMobKills(qid, d.buildQuestKillString(q)))
	}
}

func (d *Player) meetsMobKills(questID int16, reqs []nx.ReqMob) bool {
	if len(reqs) == 0 {
		return true
	}
	m := d.quests.mobKills[questID]
	if m == nil {
		return false
	}
	for _, r := range reqs {
		if m[r.ID] < r.Count {
			return false
		}
	}
	return true
}

func (d *Player) allowsQuestDrop(qid int32) bool {
	if qid == 0 {
		return true
	}
	if d == nil {
		return false
	}
	return d.quests.hasInProgress(int16(qid))
}

func (p *Player) ensureSummonState() {
	if p.summons == nil {
		p.summons = &summonState{}
	}
}

func (p *Player) addSummon(su *summon) {
	p.ensureSummonState()
	if su == nil {
		return
	}

	switch skill.Skill(su.SkillID) {
	case skill.Puppet, skill.SniperPuppet:
		if p.summons.puppet != nil {
			p.broadcastRemoveSummon(p.summons.puppet.SkillID, constant.SummonRemoveReasonReplaced)
			p.summons.puppet = nil
		}
		p.summons.puppet = su

	case skill.SummonDragon, skill.SilverHawk, skill.GoldenEagle:
		if p.summons.summon != nil {
			p.broadcastRemoveSummon(p.summons.summon.SkillID, constant.SummonRemoveReasonReplaced)
			p.summons.summon = nil
		}
		p.summons.summon = su

	default:
		if p.summons.summon != nil {
			p.broadcastRemoveSummon(p.summons.summon.SkillID, constant.SummonRemoveReasonReplaced)
			p.summons.summon = nil
		}
		p.summons.summon = su
	}

	p.broadcastShowSummon(su)
}

func (p *Player) getSummon(skillID int32) *summon {
	p.ensureSummonState()
	if p.summons.summon != nil && p.summons.summon.SkillID == skillID {
		return p.summons.summon
	}
	if p.summons.puppet != nil && p.summons.puppet.SkillID == skillID {
		return p.summons.puppet
	}
	return nil
}

func (p *Player) expireSummons() {
	if p != nil && p.buffs != nil {
		for sid := range p.buffs.activeSkillLevels {
			switch skill.Skill(sid) {
			case skill.SilverHawk, skill.GoldenEagle, skill.SummonDragon, skill.Puppet, skill.SniperPuppet:
				p.buffs.expireBuffNow(sid)
			}
		}
	}
}

func (p *Player) removeActiveSummonForSkill(skillID int32, reason byte) {
	if p == nil {
		return
	}
	p.ensureSummonState()

	p.broadcastRemoveSummon(skillID, reason)

	switch skill.Skill(skillID) {
	case skill.Puppet, skill.SniperPuppet:
		p.summons.puppet = nil

	case skill.SummonDragon, skill.SilverHawk, skill.GoldenEagle:
		p.summons.summon = nil
	}
}

func (p *Player) broadcastRemoveSummon(summonSkillID int32, reason byte) {
	if p == nil {
		return
	}

	p.Send(packetRemoveSummon(p.ID, summonSkillID, reason))

	if p.inst != nil {
		p.inst.send(packetRemoveSummon(p.ID, summonSkillID, reason))
	}
}

func (p *Player) broadcastShowSummon(su *summon) {
	if p == nil || p.inst == nil {
		return
	}

	p.Send(packetShowSummon(p.ID, su, false))
	p.inst.sendExcept(packetShowSummon(p.ID, su, true), p.Conn)
}

func (p *Player) shouldKeepSummonOnTransfer(su *summon) bool {
	return su != nil && p.hasActiveBuff(su.SkillID)
}

func (p *Player) hasActiveBuff(skillID int32) bool {
	if p == nil || p.buffs == nil {
		return false
	}
	lvl, ok := p.buffs.activeSkillLevels[skillID]
	return ok && lvl > 0
}

func (p *Player) updatePet() {
	p.MarkDirty(DirtyPet, time.Millisecond*300)
	p.inst.send(packetPlayerPetUpdate(p.pet.lockerSN))
}

func (p *Player) petCanTakeDrop(drop fieldDrop) bool {
	if p.pet == nil {
		return false
	}

	if drop.mesos > 0 {
		if p.hasEquipped(constant.ItemMesoMagnet) {
			return true
		}
		return false
	} else {
		if p.hasEquipped(constant.ItemItemPouch) {
			return true
		}
		return false
	}
}

func (p *Player) hasEquipped(itemID int32) bool {
	if p == nil || itemID <= 0 {
		return false
	}
	for i := range p.equip {
		it := p.equip[i]
		if it.slotID < 0 && it.amount > 0 && it.ID == itemID {
			return true
		}
	}
	return false
}

func (p *Player) hasBuff(mask int) bool {
	if p == nil || p.buffs == nil {
		return false
	}
	return p.buffs.hasBuff(mask)
}

func packetPlayerReceivedDmg(charID int32, attack int8, initalAmmount, reducedAmmount, spawnID, mobID, healSkillID int32,
	stance, reflectAction byte, reflected byte, reflectX, reflectY int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerTakeDmg)
	p.WriteInt32(charID)
	p.WriteInt8(attack)
	p.WriteInt32(initalAmmount)

	// v48 remote hit packets encode the mob template here; reflected hit data,
	// including the target spawn ID, is carried later in the optional reflect block.
	p.WriteInt32(mobID)
	p.WriteByte(stance)
	p.WriteByte(reflected)

	if reflected > 0 {
		p.WriteByte(0)
		p.WriteInt32(spawnID)
		p.WriteByte(reflectAction)
		p.WriteInt16(reflectX)
		p.WriteInt16(reflectY)
	}

	p.WriteByte(0)
	p.WriteInt32(reducedAmmount)

	// Check if used
	if reducedAmmount < 0 {
		p.WriteInt32(healSkillID)
	}

	return p
}

func packetPlayerLevelUpAnimation(charID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerAnimation)
	p.WriteInt32(charID)
	p.WriteByte(constant.PlayerEffectLevelUp)

	return p
}

func packetPlayerEffectSkill(onOther bool, skillID int32, level byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerEffect)
	if onOther {
		p.WriteByte(constant.PlayerEffectSkillOnOther)
	} else {
		p.WriteByte(constant.PlayerEffectSkillOnSelf)
	}
	p.WriteInt32(skillID)
	p.WriteByte(level)
	return p
}

func packetPlayerRemoteSkillAnimationWithType(charID int32, animationType byte, skillID int32, level byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerAnimation)
	p.WriteInt32(charID)
	p.WriteByte(animationType)
	p.WriteInt32(skillID)
	p.WriteByte(level)
	return p
}

func packetPlayerRemotePrimarySkillAnimation(charID int32, skillID int32, level byte) mpacket.Packet {
	return packetPlayerRemoteSkillAnimationWithType(charID, constant.PlayerEffectSkillOnSelf, skillID, level)
}

func packetPlayerLocalSkillAnimationWithType(animationType byte, skillID int32, level byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerRemoteAnimation)
	p.WriteByte(animationType)
	p.WriteInt32(skillID)
	p.WriteByte(level)
	return p
}

func packetPlayerLocalPrimarySkillAnimation(skillID int32, level byte) mpacket.Packet {
	return packetPlayerLocalSkillAnimationWithType(constant.PlayerEffectSkillOnSelf, skillID, level)
}

func sendPrimarySkillAnimation(plr *Player, skillID int32, level byte) {
	if plr == nil {
		return
	}

	plr.Send(packetPlayerLocalPrimarySkillAnimation(skillID, level))
	if plr.inst != nil && plr.Conn != nil {
		plr.inst.sendExcept(packetPlayerRemotePrimarySkillAnimation(plr.ID, skillID, level), plr.Conn)
	}
}

func sendSecondaryRemoteSkillAnimation(plr *Player, skillID int32, level byte) {
	if plr == nil || plr.inst == nil || plr.Conn == nil {
		return
	}

	plr.inst.sendExcept(packetPlayerRemoteSkillAnimationWithType(plr.ID, constant.PlayerEffectSkillOnOther, skillID, level), plr.Conn)
}

func sendSecondarySkillAnimation(plr *Player, skillID int32, level byte) {
	if plr == nil {
		return
	}

	plr.Send(packetPlayerLocalSkillAnimationWithType(constant.PlayerEffectSkillOnOther, skillID, level))
	sendSecondaryRemoteSkillAnimation(plr, skillID, level)
}

func packetPlayerGiveBuff(mask []byte, values []byte, delay int16, extra byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelTempStatChange)

	// Masks are already encoded in the same reversed-DWORD order used by CTS flags.
	if len(mask) < 8 {
		tmp := make([]byte, 8)
		copy(tmp, mask)
		mask = tmp
	} else if len(mask) > 8 {
		mask = mask[:8]
	}
	writeExtra := buffMaskNeedsExtraByte(mask)
	p.WriteBytes(mask)

	// Local v48 forced-stat entries: short value, int32 source, short duration500.
	p.WriteBytes(values)

	// Self path: 2-byte delay
	p.WriteInt16(delay)

	if writeExtra {
		p.WriteByte(extra)
	}

	return p
}

func buffMaskNeedsExtraByte(mask []byte) bool {
	if len(mask) < 8 {
		return false
	}
	value := uint64(mask[0]) | uint64(mask[1])<<8 | uint64(mask[2])<<16 | uint64(mask[3])<<24 |
		uint64(mask[4])<<32 | uint64(mask[5])<<40 | uint64(mask[6])<<48 | uint64(mask[7])<<56
	const forcedStatExtraMask uint64 = 0x0000408B40020180
	return (value & forcedStatExtraMask) != 0
}

func packetPlayerShowCombo(count byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelShowCombo)
	p.WriteInt32(int32(count))
	return p
}

// Self-cancel using 8-byte mask
func packetPlayerCancelBuff(mask []byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelRemoveTempStat)

	// Masks are already encoded in the same reversed-DWORD order used by CTS flags.
	if len(mask) < 8 {
		tmp := make([]byte, 8)
		copy(tmp, mask)
		mask = tmp
	} else if len(mask) > 8 {
		mask = mask[:8]
	}
	p.WriteBytes(mask)
	if buffMaskNeedsExtraByte(mask) {
		p.WriteByte(0)
	}
	return p
}

func packetPlayerCancelForeignBuff(charID int32, mask []byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerResetForeignBuff)
	p.WriteInt32(charID)
	p.WriteBytes(mask)
	return p
}

func packetPlayerMove(charID int32, bytes []byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerMovement)
	p.WriteInt32(charID)
	p.WriteBytes(bytes)

	return p
}

func packetPlayerEmoticon(charID int32, emotion int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerEmoticon)
	p.WriteInt32(charID)
	p.WriteInt32(emotion)

	return p
}

func packetPlayerSkillBookUpdate(skillID int32, level int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSkillRecordUpdate)
	p.WriteByte(1)
	p.WriteInt16(1)
	p.WriteInt32(skillID)
	p.WriteInt32(level)
	p.WriteInt32(0)
	p.WriteByte(1)

	return p
}

func packetPlayerSkillCooldown(skillID int32, time int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelSkillCooldown)
	p.WriteInt32(skillID)
	p.WriteInt16(time)

	return p
}

func packetPlayerStatChange(flag bool, stat int32, value int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStatChange)
	p.WriteBool(flag)
	p.WriteInt32(stat)
	p.WriteInt32(value)

	return p
}

func packetPlayerNoChange() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelInventoryOperation)
	p.WriteByte(0x01)
	p.WriteByte(0x00)
	p.WriteByte(0x00)

	return p
}

func packetChangeChannel(ip []byte, port int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelChange)
	p.WriteBool(true)
	p.WriteBytes(ip)
	p.WriteInt16(port)

	return p
}

func packetCannotEnterCashShop() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelChangeServer)
	p.WriteByte(2)

	return p
}

const (
	setFieldCharSectionStats         int16 = 0x0001
	setFieldCharSectionMeta          int16 = 0x0002
	setFieldCharSectionSlotSizes     int16 = 0x0080
	setFieldCharSectionEquip         int16 = 0x0004
	setFieldCharSectionUse           int16 = 0x0008
	setFieldCharSectionSetup         int16 = 0x0010
	setFieldCharSectionEtc           int16 = 0x0020
	setFieldCharSectionCash          int16 = 0x0040
	setFieldCharSectionSkills        int16 = 0x0100
	setFieldCharSectionActiveQuests  int16 = 0x0200
	setFieldCharSectionMiniGames     int16 = 0x0400
	setFieldCharSectionRings         int16 = 0x0800
	setFieldCharSectionTeleportRocks int16 = 0x1000
	setFieldCharSectionCompletedQ    int16 = 0x4000
	setFieldCharSectionCooldowns     int16 = -0x8000
	setFieldCharSectionMinimal             = setFieldCharSectionStats | setFieldCharSectionTeleportRocks
	setFieldCharSectionInventory           = setFieldCharSectionEquip | setFieldCharSectionUse | setFieldCharSectionSetup | setFieldCharSectionEtc | setFieldCharSectionCash
	setFieldCharSectionSafeZero            = setFieldCharSectionSkills | setFieldCharSectionCooldowns | setFieldCharSectionActiveQuests | setFieldCharSectionCompletedQ | setFieldCharSectionMiniGames
	cashShopCharacterDataMask              = setFieldCharSectionStats | setFieldCharSectionMeta | setFieldCharSectionSlotSizes | setFieldCharSectionEquip | setFieldCharSectionUse | setFieldCharSectionSetup | setFieldCharSectionEtc | setFieldCharSectionCash
)

// AppendCashShopCharacterData mirrors the verified CharacterData::Decode read order
// used by CStage::OnSetCashShop. It intentionally emits only the sections we have
// confirmed and need for Cash Shop entry: stats, the post-stat byte, mesos, slot
// sizes, equipped items, and the inventory tabs the client still expects here.
func AppendCashShopCharacterData(p *mpacket.Packet, plr *Player) {
	p.WriteInt16(cashShopCharacterDataMask)
	writeSetFieldCharacterStats(p, *plr)
	writeSetFieldPostStatByte(p, *plr)
	writeSetFieldMesos(p, plr.mesos)
	writeSetFieldSlotSizes(p, plr)
	writeCashShopInventory(p, *plr)
}

func writeCashShopInventory(p *mpacket.Packet, plr Player) {
	for _, it := range plr.equip {
		if it.slotID < 0 && !it.cash {
			p.WriteBytes(it.InventoryBytes())
		}
	}
	p.WriteByte(0)

	for _, it := range plr.equip {
		if it.slotID < 0 && it.cash {
			p.WriteBytes(it.InventoryBytes())
		}
	}
	p.WriteByte(0)

	writeSetFieldInventoryTab(p, plr.equip)
	writeSetFieldInventoryTab(p, plr.use)
	writeSetFieldInventoryTab(p, plr.setUp)
	writeSetFieldInventoryTab(p, plr.etc)
	writeSetFieldInventoryTab(p, plr.cash)
}

func packetPlayerEnterGame(plr Player, channelID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelWarpToMap)
	p.WriteInt32(channelID)
	p.WriteByte(0)
	p.WriteByte(1)

	randomBytes := make([]byte, 12)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err.Error())
	}
	p.WriteBytes(randomBytes)

	sectionMask := setFieldCharSectionMinimal | setFieldCharSectionMeta | setFieldCharSectionSlotSizes | setFieldCharSectionInventory | setFieldCharSectionSafeZero
	if len(plr.rings) > 0 || plr.partnerID > 0 {
		sectionMask |= setFieldCharSectionRings
	}
	p.WriteInt16(sectionMask)
	writeSetFieldCharacterStats(&p, plr)
	writeSetFieldPostStatByte(&p, plr)
	writeSetFieldMesos(&p, plr.mesos)
	writeSetFieldSlotSizes(&p, &plr)
	writeSetFieldInventory(&p, plr)
	writeSetFieldSkills(&p, plr)
	writeSetFieldCooldowns(&p, plr)
	writeSetFieldActiveQuests(&p, plr)
	writeSetFieldCompletedQuests(&p, plr)
	writeSetFieldEmptyMiniGames(&p)
	if sectionMask&setFieldCharSectionRings != 0 {
		writeSetFieldRings(&p, &plr)
	}
	writeSetFieldTeleportRocks(&p, &plr)

	return p
}

func writeSetFieldCharacterStats(p *mpacket.Packet, plr Player) {
	p.WriteInt32(plr.ID)
	p.WritePaddedString(plr.Name, 13)
	p.WriteByte(plr.gender)
	p.WriteByte(plr.skin)
	p.WriteInt32(plr.face)
	p.WriteInt32(plr.hair)
	p.WriteInt64(plr.petCashID)
	p.WriteByte(plr.level)
	p.WriteInt16(plr.job)
	p.WriteInt16(plr.str)
	p.WriteInt16(plr.dex)
	p.WriteInt16(plr.intt)
	p.WriteInt16(plr.luk)
	p.WriteInt16(plr.hp)
	p.WriteInt16(plr.maxHP)
	p.WriteInt16(plr.mp)
	p.WriteInt16(plr.maxMP)
	p.WriteInt16(plr.ap)
	p.WriteInt16(plr.sp)
	p.WriteInt32(plr.exp)
	p.WriteInt16(plr.fame)
	p.WriteInt32(plr.mapID)
	p.WriteByte(plr.mapPos)

	// Client consumes one byte after the stat block before optional sections.
	// The exact meaning in this v48 SetField branch is still being verified.
}

func writeSetFieldPostStatByte(p *mpacket.Packet, plr Player) {
	// Working packet dumps show inventory begins immediately after this byte.
	// Keep the current value stable until the controlling decode branch is known.
	p.WriteByte(plr.buddyListSize)
}

func writeSetFieldTeleportRocks(p *mpacket.Packet, plr *Player) {
	for i := 0; i < constant.TeleportRockRegSlots; i++ {
		if i < len(plr.regTeleportRocks) {
			p.WriteInt32(plr.regTeleportRocks[i])
		} else {
			p.WriteInt32(constant.InvalidMap)
		}
	}
	for i := 0; i < constant.TeleportRockVIPSlots; i++ {
		if i < len(plr.vipTeleportRocks) {
			p.WriteInt32(plr.vipTeleportRocks[i])
		} else {
			p.WriteInt32(constant.InvalidMap)
		}
	}
}

func writeSetFieldRings(p *mpacket.Packet, plr *Player) {
	plr.encodeLocalRingRecords(p)
}

func writeSetFieldMesos(p *mpacket.Packet, mesos int32) {
	p.WriteInt32(mesos)
}

func writeSetFieldSlotSizes(p *mpacket.Packet, plr *Player) {

	if plr.equipSlotSize == 0 {
		plr.equipSlotSize = 24
	}
	if plr.useSlotSize == 0 {
		plr.useSlotSize = 24
	}
	if plr.setupSlotSize == 0 {
		plr.setupSlotSize = 24
	}
	if plr.etcSlotSize == 0 {
		plr.etcSlotSize = 24
	}
	if plr.cashSlotSize == 0 {
		plr.cashSlotSize = 24
	}

	p.WriteByte(plr.equipSlotSize)
	p.WriteByte(plr.useSlotSize)
	p.WriteByte(plr.setupSlotSize)
	p.WriteByte(plr.etcSlotSize)
	p.WriteByte(plr.cashSlotSize)
}

func writeSetFieldInventory(p *mpacket.Packet, plr Player) {
	writeSetFieldEquippedItems(p, plr)

	writeSetFieldInventoryTab(p, plr.equip)
	writeSetFieldInventoryTab(p, plr.use)
	writeSetFieldInventoryTab(p, plr.setUp)
	writeSetFieldInventoryTab(p, plr.etc)
	writeSetFieldInventoryTab(p, plr.cash)
}

func writeSetFieldInventoryTab(p *mpacket.Packet, items []Item) {
	cp := make([]Item, 0, len(items))
	for _, it := range items {
		if it.slotID > 0 {
			cp = append(cp, it)
		}
	}

	sort.Slice(cp, func(i, j int) bool {
		return cp[i].slotID < cp[j].slotID
	})

	for _, it := range cp {
		p.WriteBytes(it.InventoryBytes())
	}
	p.WriteByte(0)
}

func writeSetFieldSkills(p *mpacket.Packet, plr Player) {
	if len(plr.skills) == 0 {
		p.WriteInt16(0)
		return
	}

	skills := make([]playerSkill, 0, len(plr.skills))
	for _, ps := range plr.skills {
		skills = append(skills, ps)
	}

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].ID < skills[j].ID
	})

	p.WriteInt16(int16(len(skills)))
	for _, ps := range skills {
		p.WriteInt32(ps.ID)
		p.WriteInt32(int32(ps.Level))
		if isFourthJobSkill(ps.ID) {
			p.WriteInt32(int32(ps.Mastery))
		}
	}
}

func writeSetFieldCooldowns(p *mpacket.Packet, plr Player) {
	active := make([]playerSkill, 0, len(plr.skills))
	for _, ps := range plr.skills {
		if ps.Cooldown > 0 {
			active = append(active, ps)
		}
	}

	sort.Slice(active, func(i, j int) bool {
		return active[i].ID < active[j].ID
	})

	p.WriteInt16(int16(len(active)))
	for _, ps := range active {
		p.WriteInt32(ps.ID)
		p.WriteInt16(ps.Cooldown)
	}
}

func writeSetFieldActiveQuests(p *mpacket.Packet, plr Player) {
	writeActiveQuests(p, plr.quests.inProgressList())
}

func writeSetFieldCompletedQuests(p *mpacket.Packet, plr Player) {
	writeCompletedQuests(p, plr.quests.completedList())
}

func writeSetFieldEmptyMiniGames(p *mpacket.Packet) {
	p.WriteInt16(0)
}

func isFourthJobSkill(id int32) bool {
	jobID := id / 10000
	return jobID >= 100 && jobID%10 == 2
}

func packetInventoryAddItem(item Item, newItem bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelInventoryOperation)
	p.WriteByte(0x01)
	p.WriteByte(0x01)

	if newItem {
		p.WriteByte(0x00)
		p.WriteByte(item.invID)
		p.WriteInt16(item.slotID)
		p.WriteBytes(item.StorageBytes())
	} else {
		p.WriteByte(0x01)
		p.WriteByte(item.invID)
		p.WriteInt16(item.slotID)
		p.WriteInt16(item.amount)
	}

	return p
}

func packetInventoryModifyItemAmount(item Item) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelInventoryOperation)
	p.WriteByte(0x01)
	p.WriteByte(0x01)
	p.WriteByte(0x01)
	p.WriteByte(item.invID)
	p.WriteInt16(item.slotID)
	p.WriteInt16(item.amount)

	return p
}

func packetInventoryAddItems(items []Item, newItem []bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelInventoryOperation)

	p.WriteByte(0x01)
	if len(items) != len(newItem) {
		p.WriteByte(0)
		return p
	}

	p.WriteByte(byte(len(items)))

	for i, v := range items {
		if newItem[i] {
			p.WriteByte(0x00)
			p.WriteByte(v.invID)
			p.WriteInt16(v.slotID)
			p.WriteBytes(v.StorageBytes())
		} else {
			p.WriteByte(0x01)
			p.WriteByte(v.invID)
			p.WriteInt16(v.slotID)
			p.WriteInt16(v.amount)
		}
	}

	return p
}

func packetInventoryChangeItemSlot(invTabID byte, origPos, newPos int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelInventoryOperation)
	p.WriteByte(0x01)
	p.WriteByte(0x01)
	p.WriteByte(0x02)
	p.WriteByte(invTabID)
	p.WriteInt16(origPos)
	p.WriteInt16(newPos)
	if invTabID == 1 && (origPos < 0 || newPos < 0) {
		p.WriteByte(0)
	}

	return p
}

func packetInventoryRemoveItem(item Item) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelInventoryOperation)
	p.WriteByte(0x01)
	p.WriteByte(0x01)
	p.WriteByte(0x03)
	p.WriteByte(item.invID)
	p.WriteInt16(item.slotID)
	if item.invID == constant.InventoryEquip && item.slotID < 0 {
		p.WriteByte(0)
	}

	return p
}

func packetInventoryChangeEquip(chr Player) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerChangeAvatar)
	p.WriteInt32(chr.ID)
	p.WriteByte(1)
	chr.encodeDisplayBytes(&p)
	chr.encodeRemoteRingLooks(&p)
	p.WriteInt32(0)

	return p
}

func packetInventoryNoChange() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelInventoryOperation)
	p.WriteByte(0x01)
	p.WriteByte(0x00)

	return p
}

func packetInventoryDontTake() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelInventoryOperation)
	p.WriteInt16(1)

	return p
}

func packetBuddyInfo(buddyList []buddy) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelBuddyInfo)
	p.WriteByte(0x12)
	p.WriteByte(byte(len(buddyList)))

	for _, v := range buddyList {
		p.WriteInt32(v.id)
		p.WritePaddedString(v.name, 13)
		p.WriteByte(v.status)
		p.WriteInt32(v.channelID)
	}

	for _, v := range buddyList {
		p.WriteInt32(v.cashShop) // wizet mistake and this should be a bool?
	}

	return p
}

// It is possible to change ID's using this packet, however if the ID is a request it will crash the users
// client when selecting an option in notification, therefore the ID has not been allowed to change
func packetBuddyUpdate(id int32, name string, status byte, channelID int32, cashShop bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelBuddyInfo)
	p.WriteByte(0x08)
	p.WriteInt32(id) // original ID
	p.WriteInt32(id)
	p.WritePaddedString(name, 13)
	p.WriteByte(status)
	p.WriteInt32(channelID)
	p.WriteBool(cashShop)

	return p
}

func packetBuddyListSizeUpdate(size byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelBuddyInfo)
	p.WriteByte(0x15)
	p.WriteByte(size)

	return p
}

func packetPlayerAvatarSummaryWindow(charID int32, plr Player) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelAvatarInfoWindow)
	p.WriteInt32(charID)
	p.WriteByte(plr.level)
	p.WriteInt16(plr.job)
	p.WriteInt16(plr.fame)

	if plr.guild != nil {
		p.WriteString(plr.guild.name)
	} else {
		p.WriteString("")
	}

	p.WriteBool(false)
	p.WriteBool(false)
	p.WriteByte(0)

	return p
}

func packetShowCountdown(time int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelTimer)
	p.WriteByte(2)
	p.WriteInt32(time)

	return p
}

func packetHideCountdown() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelTimer)
	p.WriteByte(0)
	p.WriteInt32(0)

	return p
}

func packetBuddyUnkownError() mpacket.Packet {
	return packetBuddyRequestResult(0x16)
}

func packetBuddyPlayerFullList() mpacket.Packet {
	return packetBuddyRequestResult(0x0b)
}

func packetBuddyOtherFullList() mpacket.Packet {
	return packetBuddyRequestResult(0x0c)
}

func packetBuddyAlreadyAdded() mpacket.Packet {
	return packetBuddyRequestResult(0x0d)
}

func packetBuddyIsGM() mpacket.Packet {
	return packetBuddyRequestResult(0x0e)
}

func packetBuddyNameNotRegistered() mpacket.Packet {
	return packetBuddyRequestResult(0x0f)
}

func packetBuddyRequestResult(code byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelBuddyInfo)
	p.WriteByte(code)

	return p
}

func packetBuddyReceiveRequest(fromID int32, fromName string, fromChannelID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelBuddyInfo)
	p.WriteByte(0x9)
	p.WriteInt32(fromID)
	p.WriteString(fromName)
	p.WriteInt32(fromID)
	p.WritePaddedString(fromName, 13)
	p.WriteByte(1)
	p.WriteInt32(fromChannelID)
	p.WriteBool(false) // sender in cash shop

	return p
}

func packetBuddyOnlineStatus(id int32, channelID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelBuddyInfo)
	p.WriteByte(0x14)
	p.WriteInt32(id)
	p.WriteInt8(0)
	p.WriteInt32(channelID)

	return p
}

func packetBuddyChangeChannel(id int32, channelID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelBuddyInfo)
	p.WriteByte(0x14)
	p.WriteInt32(id)
	p.WriteInt8(1)
	p.WriteInt32(channelID)

	return p
}

func packetMapChange(mapID int32, channelID int32, mapPos byte, hp int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelWarpToMap)
	p.WriteInt32(channelID)
	p.WriteInt16(2)
	p.WriteInt32(mapID)
	p.WriteByte(mapPos)
	p.WriteInt16(hp)
	p.WriteByte(0) // flag for more reading
	p.WriteInt64(0x01FFFFFFFFFFFFFF)

	return p
}

func packetPlayerPetUpdate(sn int64) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStatChange)
	p.WriteBool(false)
	p.WriteInt32(constant.PetID)
	p.WriteUint64(uint64(sn))
	p.WriteByte(0)
	p.WriteByte(0)

	return p
}

func packetPlayerGiveForeignBuff(charID int32, mask []byte, values []byte, delay int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerGiveForeignBuff)
	p.WriteInt32(charID)

	// Masks are already encoded in the same reversed-DWORD order used by CTS flags.
	if len(mask) < 8 {
		tmp := make([]byte, 8)
		copy(tmp, mask)
		mask = tmp
	} else if len(mask) > 8 {
		mask = mask[:8]
	}
	p.WriteBytes(mask)

	// Subset payload in reference order
	p.WriteBytes(values)

	// Delay (usually 0)
	p.WriteInt16(delay)
	return p
}

func packetPlayerShowChair(plrID, chairID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerSit)
	p.WriteInt32(plrID)
	p.WriteInt32(chairID)
	return p
}

func packetPlayerChairResult(chairID int16) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerSitResult)
	p.WriteBool(chairID != -1)
	if chairID != -1 {
		p.WriteInt16(chairID)
	}
	return p
}

func packetPlayerChairUpdate() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelStatChange)
	p.WriteInt16(1)
	p.WriteInt32(0)
	return p
}

// writeAttackDamages writes damage values to packet
func writeAttackDamages(p *mpacket.Packet, info attackInfo) {
	for idx, dmg := range info.damages {
		if idx < len(info.isCritical) && info.isCritical[idx] {
			crit := int32(uint32(dmg) | 0x80000000)
			p.WriteInt32(crit)
		} else {
			p.WriteInt32(dmg)
		}
	}
}

func packetSkillMelee(char Player, ad attackData) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerUseMeleeSkill)
	p.WriteInt32(char.ID)
	p.WriteByte(ad.targets*0x10 + ad.hits)
	p.WriteByte(ad.skillLevel)

	if ad.skillLevel != 0 {
		p.WriteInt32(ad.skillID)
	}

	if ad.facesLeft {
		p.WriteByte(ad.action | (1 << 7))
	} else {
		p.WriteByte(ad.action)
	}

	p.WriteByte(ad.attackType)

	// Write mastery display value
	mastery := char.getMasteryDisplay()
	p.WriteByte(mastery)
	p.WriteInt32(ad.projectileID)

	for _, info := range ad.attackInfo {
		p.WriteInt32(info.spawnID)
		p.WriteByte(info.hitAction)

		if ad.isMesoExplosion {
			p.WriteByte(byte(len(info.damages)))
		}

		writeAttackDamages(&p, info)
	}

	return p
}

func packetSkillRanged(char Player, ad attackData) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerUseRangedSkill)
	p.WriteInt32(char.ID)
	p.WriteByte(ad.targets*0x10 + ad.hits)
	p.WriteByte(ad.skillLevel)

	if ad.skillLevel != 0 {
		p.WriteInt32(ad.skillID)
	}

	if ad.facesLeft {
		p.WriteByte(ad.action | (1 << 7))
	} else {
		p.WriteByte(ad.action | 0)
	}

	p.WriteByte(ad.attackType)

	// Write mastery display value
	mastery := char.getMasteryDisplay()
	p.WriteByte(mastery)
	p.WriteInt32(ad.projectileID)

	for _, info := range ad.attackInfo {
		p.WriteInt32(info.spawnID)
		p.WriteByte(info.hitAction)

		if ad.isMesoExplosion {
			p.WriteByte(byte(len(info.damages)))
		}

		writeAttackDamages(&p, info)
	}

	return p
}

func packetSkillMagic(char Player, ad attackData) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerUseMagicSkill)
	p.WriteInt32(char.ID)
	p.WriteByte(ad.targets*0x10 + ad.hits)
	p.WriteByte(ad.skillLevel)

	if ad.skillLevel != 0 {
		p.WriteInt32(ad.skillID)
	}

	if ad.facesLeft {
		p.WriteByte(ad.action | (1 << 7))
	} else {
		p.WriteByte(ad.action | 0)
	}

	p.WriteByte(ad.attackType)

	// Write mastery display value
	mastery := char.getMasteryDisplay()
	p.WriteByte(mastery)
	p.WriteInt32(ad.projectileID)

	for _, info := range ad.attackInfo {
		p.WriteInt32(info.spawnID)
		p.WriteByte(info.hitAction)

		if ad.isMesoExplosion {
			p.WriteByte(byte(len(info.damages)))
		}

		writeAttackDamages(&p, info)
	}

	return p
}

func packetSkillStop(plrID, skillID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerStopSkill)
	p.WriteInt32(plrID)
	p.WriteInt32(skillID)
	return p
}

func packetPlayerHpChange(plrID, hp, maxHp int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcode.SendChannelPlayerPartyHP)
	p.WriteInt32(plrID)
	p.WriteInt32(hp)
	p.WriteInt32(maxHp)

	return p
}
