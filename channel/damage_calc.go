package channel

import (
	"log"
	"math"
	"math/rand"

	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/constant/skill"
	"github.com/Hucaru/Valhalla/nx"
)

// DamageRange represents the min and max damage for validation
type DamageRange struct {
	Min    float64
	Max    float64
	Valid  bool
	Reason string
}

// CalcHitResult represents the result of a hit calculation
type CalcHitResult struct {
	IsCrit            bool
	IsMiss            bool
	MinDamage         float64
	MaxDamage         float64
	ExpectedDmg       float64
	ToleranceMax      float64
	ClientDamage      int32
	IsValid           bool
	ValidationSkipped bool
	ValidationReason  string
}

type DamageRngBuffer [constant.DamageRngBufferSize]uint32

func NewDamageRngBuffer(rng *rand.Rand) *DamageRngBuffer {
	b := &DamageRngBuffer{}
	if rng == nil {
		mid := uint32(^uint32(0) / 2)
		for i := 0; i < constant.DamageRngBufferSize; i++ {
			(*b)[i] = mid
		}
		return b
	}

	for i := 0; i < constant.DamageRngBufferSize; i++ {
		(*b)[i] = rng.Uint32()
	}
	return b
}

func (b *DamageRngBuffer) GetAllRaw() [constant.DamageRngBufferSize]uint32 {
	return *b
}

func (b *DamageRngBuffer) GetRaw(index int) uint32 {
	if index < 0 {
		index = -index
	}
	return (*b)[index%constant.DamageRngBufferSize]
}

func (b *DamageRngBuffer) AccuracyRaw(hitIndex int, stepsPerHit int) uint32 {
	if stepsPerHit <= 0 {
		stepsPerHit = 3
	}
	return (*b)[(hitIndex*stepsPerHit)%constant.DamageRngBufferSize]
}

func (b *DamageRngBuffer) DamageRaw(hitIndex int, stepsPerHit int) uint32 {
	if stepsPerHit <= 0 {
		stepsPerHit = 3
	}
	return (*b)[(hitIndex*stepsPerHit+1)%constant.DamageRngBufferSize]
}

func (b *DamageRngBuffer) DefenseRaw(hitIndex int, stepsPerHit int) uint32 {
	if stepsPerHit <= 0 {
		stepsPerHit = 3
	}
	return (*b)[(hitIndex*stepsPerHit+2)%constant.DamageRngBufferSize]
}

func (b *DamageRngBuffer) CriticalRaw(hitIndex int, stepsPerHit int) uint32 {
	if stepsPerHit <= 0 {
		stepsPerHit = 4
	}
	return (*b)[(hitIndex*stepsPerHit+3)%constant.DamageRngBufferSize]
}

func NormalizeDamageRng(raw uint32) float64 {
	return float64(raw%constant.DamageRngModulo) * constant.DamageRngNormalize
}

type ElementAmpData struct {
	Magic int
	Mana  int
}

const (
	elementCodeNone byte = iota
	elementCodeFire
	elementCodeIce
	elementCodeLightning
	elementCodeHoly
	elementCodePoison
)

type DamageCalculator struct {
	player       *Player
	data         *attackData
	attackType   int
	weaponID     int32
	weaponType   constant.WeaponType
	skill        *nx.PlayerSkill
	skillID      int32
	skillLevel   byte
	isRanged     bool
	masteryMod   float64
	critSkill    *nx.PlayerSkill
	critLevel    byte
	watk         int16
	projectileID int32
	attackAction constant.AttackAction
	attackOption constant.AttackOption
}

// NewDamageCalculator creates a new damage calculator
func NewDamageCalculator(plr *Player, data *attackData, attackType int) *DamageCalculator {
	calc := &DamageCalculator{
		player:       plr,
		data:         data,
		attackType:   attackType,
		isRanged:     attackType == attackRanged,
		skillID:      data.skillID,
		skillLevel:   data.skillLevel,
		projectileID: data.projectileID,
		attackAction: constant.AttackAction(data.action),
		attackOption: constant.AttackOption(data.option),
	}
	if attackType == attackSummon && data.summonType > 0 {
		calc.skillID = data.summonType
		if calc.skillLevel == 0 {
			if summ := plr.getSummon(data.summonType); summ != nil {
				calc.skillLevel = summ.Level
			}
		}
	}

	weaponID := int32(0)
	for _, item := range plr.equip {
		if item.slotID == -11 {
			weaponID = item.ID
			break
		}
	}
	calc.weaponID = weaponID
	calc.weaponType = constant.GetWeaponType(weaponID)

	if calc.skillID > 0 {
		if skillData, err := nx.GetPlayerSkill(calc.skillID); err == nil && len(skillData) > 0 {
			if calc.skillLevel > 0 && int(calc.skillLevel) <= len(skillData) {
				calc.skill = &skillData[calc.skillLevel-1]
			}
		}
	}

	calc.masteryMod = calc.GetMasteryModifier()
	calc.critLevel, calc.critSkill = calc.GetCritSkill()
	calc.watk = calc.GetTotalWatk()

	return calc
}

// ValidateAttack validates all hits in an attack and determines critical hits
func (calc *DamageCalculator) ValidateAttack() [][]CalcHitResult {
	calc.EnforceAttackShape()
	calc.LogSuspiciousAttackShape()

	results := make([][]CalcHitResult, len(calc.data.attackInfo))

	for targetIdx := range calc.data.attackInfo {
		info := &calc.data.attackInfo[targetIdx]

		if calc.player.inst == nil {
			continue
		}
		mob, err := calc.player.inst.lifePool.getMobFromID(info.spawnID)
		if err != nil {
			continue
		}

		rngBuf := NewDamageRngBuffer(calc.player.rng)

		ampData := calc.GetElementAmplification()
		targetAccuracy := calc.GetTargetAccuracy(&mob)

		results[targetIdx] = make([]CalcHitResult, len(info.damages))
		for hitIdx := range info.damages {
			results[targetIdx][hitIdx] = calc.CalculateHit(
				&mob,
				info,
				rngBuf,
				ampData,
				targetAccuracy,
				hitIdx,
				targetIdx,
			)
		}
	}

	return results
}

func (calc *DamageCalculator) EnforceAttackShape() {
	if calc == nil || calc.data == nil {
		return
	}

	if expectedTargets, ok := calc.GetExpectedTargetCount(); ok && expectedTargets > 0 && len(calc.data.attackInfo) > expectedTargets {
		calc.data.attackInfo = calc.data.attackInfo[:expectedTargets]
		calc.data.targets = byte(expectedTargets)
	}

	expectedHits, ok := calc.GetExpectedHitCount()
	if !ok || expectedHits <= 0 {
		return
	}

	for idx := range calc.data.attackInfo {
		info := &calc.data.attackInfo[idx]
		if len(info.damages) > expectedHits {
			info.damages = info.damages[:expectedHits]
		}
		if len(info.isCritical) > expectedHits {
			info.isCritical = info.isCritical[:expectedHits]
		}
		if info.hitCount > byte(expectedHits) {
			info.hitCount = byte(expectedHits)
		}
	}
	if calc.data.hits > byte(expectedHits) {
		calc.data.hits = byte(expectedHits)
	}
}

func (calc *DamageCalculator) CalculateHit(
	mob *monster,
	info *attackInfo,
	rngBuf *DamageRngBuffer,
	ampData *ElementAmpData,
	targetAccuracy float64,
	hitIdx int,
	targetIdx int,
) CalcHitResult {
	result := CalcHitResult{
		ClientDamage: info.damages[hitIdx],
		IsValid:      false,
	}
	if calc.skill != nil && calc.skill.Fixdamage > 0 {
		fixed := float64(calc.skill.Fixdamage)
		result.MinDamage = fixed
		result.MaxDamage = fixed
		result.ExpectedDmg = fixed
		applyResultTolerance(&result, 1.0)
		return result
	}
	if mob != nil && mob.fixedDamage > 0 {
		fixed := float64(mob.fixedDamage)
		result.MinDamage = fixed
		result.MaxDamage = fixed
		result.ExpectedDmg = fixed
		applyResultTolerance(&result, 1.0)
		return result
	}

	if calc.handleSpecialSkillDamage(&result, mob, info, rngBuf, hitIdx) {
		return result
	}
	if calc.attackType == attackMagic && mob.invincible {
		result.MinDamage = 1
		result.MaxDamage = 1
		result.ToleranceMax = 1
		result.IsValid = (result.ClientDamage == 1)
		return result
	}

	rangeData := calc.CalculateBaseDamageRange(mob, hitIdx)
	if !rangeData.Valid {
		result.ValidationSkipped = true
		result.ValidationReason = rangeData.Reason
		result.IsValid = true

		if result.ClientDamage > 0 {
			log.Printf(
				"Skipped damage validation for player %s (ID: %d, level: %d, job: %d): client=%d, skill=%d, attackType=%d, weaponID=%d, weaponType=%d, STR=%d, DEX=%d, INT=%d, LUK=%d, WATK=%d, MATK=%d, mobID=%d, mobPD=%d, mobMD=%d, reason=%s",
				calc.player.Name,
				calc.player.ID,
				calc.player.level,
				calc.player.job,
				result.ClientDamage,
				calc.skillID,
				calc.attackType,
				calc.weaponID,
				calc.weaponType,
				calc.GetTotalStr(),
				calc.GetTotalDex(),
				calc.GetTotalInt(),
				calc.GetTotalLuk(),
				calc.GetTotalWatk(),
				calc.GetTotalMatk(),
				mob.id,
				mob.pdDamage,
				mob.mdDamage,
				rangeData.Reason,
			)
		}

		return result
	}

	minDmg, maxDmg := calc.ApplySkillModifiers(rangeData.Min, rangeData.Max, ampData, mob)
	result.IsCrit = calc.CheckCritical(rngBuf, hitIdx)

	afterMod := calc.GetAfterModifier(targetIdx, (minDmg+maxDmg)/2.0)
	minDmg *= afterMod
	maxDmg *= afterMod
	if calc.attackOption&constant.AttackOptionShadowPartner != 0 && len(info.damages) > 1 && hitIdx >= len(info.damages)/2 {
		minDmg *= 0.5
		maxDmg *= 0.5
	}

	if maxDmg > 0 && maxDmg < 1 {
		maxDmg = 1
	}
	if minDmg > maxDmg {
		minDmg = maxDmg
	}
	if calc.skill != nil && len(info.damages) == 1 && calc.skill.BulletCount > 1 {
		// Some multi-projectile skills can arrive as one summed client damage value.
		bulletCount := float64(calc.skill.BulletCount)
		minDmg *= bulletCount
		maxDmg *= bulletCount
	}

	critMultiplier := calc.GetCriticalDamageMultiplier()
	if critMultiplier > 1.0 {
		critMin := minDmg * critMultiplier
		critMax := maxDmg * critMultiplier
		if result.IsCrit {
			minDmg = critMin
			maxDmg = critMax
		}
	}
	allowedMax := maxDmg
	if calc.HasPhysicalOrMagicImmunity(mob) {
		if allowedMax < 1 {
			allowedMax = 1
		}
		if maxDmg < 1 {
			maxDmg = 1
		}
		if minDmg > 1 {
			minDmg = 1
		}
	}

	defMin, defMax := calc.CalculateDefenseReductionBounds(mob)
	if defMax > 0 {
		minDmg -= defMax
		maxDmg -= defMin
	}
	if minDmg < 1 {
		minDmg = 1
	}
	if maxDmg < 1 {
		maxDmg = 1
	}
	allowedMax = maxDmg
	if critMultiplier > 1.0 && !result.IsCrit {
		allowedMax = maxDmg * critMultiplier
	}
	if minDmg > maxDmg {
		minDmg = maxDmg
	}

	dmgRoll := 0.5
	defRoll := 0.5
	if rngBuf != nil {
		dmgRoll = NormalizeDamageRng(rngBuf.DamageRaw(hitIdx, 3))
		defRoll = NormalizeDamageRng(rngBuf.DefenseRaw(hitIdx, 3))
	}
	expected := minDmg + (maxDmg-minDmg)*dmgRoll
	if defMax > defMin {
		expected -= defMin + (defMax-defMin)*defRoll
	} else if defMax > 0 {
		expected -= defMax
	}
	if expected < 1 && result.ClientDamage > 0 {
		expected = 1
	}

	minDmg = math.Floor(minDmg)
	maxDmg = math.Floor(maxDmg)
	allowedMax = math.Floor(allowedMax)
	expected = math.Floor(expected)

	result.MinDamage = minDmg
	result.MaxDamage = allowedMax
	result.ExpectedDmg = expected

	tolerance := constant.DamageVarianceTolerance
	toleranceMax := math.Ceil(allowedMax * (1.0 + tolerance))
	if toleranceMax < 1 && result.ClientDamage > 0 && allowedMax > 0 {
		toleranceMax = 1
	}

	clientDmgFloat := float64(result.ClientDamage)
	result.ToleranceMax = toleranceMax
	result.IsValid = (clientDmgFloat <= toleranceMax)

	if !result.IsValid {
		result.ValidationReason = "client damage exceeds tolerated cap"
		log.Printf(
			"Suspicious high damage from player %s (ID: %d, level: %d, job: %d): client=%d, max expected=%.0f (with tolerance), base min=%.0f, base max=%.0f, skill=%d, attackType=%d, weaponID=%d, weaponType=%d, STR=%d, DEX=%d, INT=%d, LUK=%d, WATK=%d, MATK=%d, mobID=%d, mobPD=%d, mobMD=%d",
			calc.player.Name,
			calc.player.ID,
			calc.player.level,
			calc.player.job,
			result.ClientDamage,
			toleranceMax,
			result.MinDamage,
			result.MaxDamage,
			calc.skillID,
			calc.attackType,
			calc.weaponID,
			calc.weaponType,
			calc.GetTotalStr(),
			calc.GetTotalDex(),
			calc.GetTotalInt(),
			calc.GetTotalLuk(),
			calc.GetTotalWatk(),
			calc.GetTotalMatk(),
			mob.id,
			mob.pdDamage,
			mob.mdDamage,
		)
	}

	return result
}

func (calc *DamageCalculator) handleSpecialSkillDamage(result *CalcHitResult, mob *monster, info *attackInfo, rngBuf *DamageRngBuffer, hitIdx int) bool {
	// Summons are keyed off attack type, not skill id.
	if calc.attackType == attackSummon {
		if calc.skill == nil {
			result.ValidationSkipped = true
			result.ValidationReason = "missing summon skill data"
			result.IsValid = true
			return true
		}

		switch {
		case calc.skill.Mad > 0:
			baseRange := calc.CalculateSummonMagicDamageRange()
			if !baseRange.Valid {
				result.ValidationSkipped = true
				result.ValidationReason = baseRange.Reason
				result.IsValid = true
				return true
			}
			result.MinDamage = baseRange.Min
			result.MaxDamage = baseRange.Max
		case calc.skill.Pad > 0:
			baseRange := calc.CalculateSummonPhysicalDamageRange()
			if !baseRange.Valid {
				result.ValidationSkipped = true
				result.ValidationReason = baseRange.Reason
				result.IsValid = true
				return true
			}
			padRate := float64(calc.skill.Pad) / 100.0
			result.MinDamage = baseRange.Min * padRate
			result.MaxDamage = baseRange.Max * padRate
		default:
			result.ValidationSkipped = true
			result.ValidationReason = "unsupported summon damage formula"
			result.IsValid = true
			return true
		}

		result.ExpectedDmg = (result.MinDamage + result.MaxDamage) / 2.0
		applyResultTolerance(result, 1.0+constant.DamageVarianceTolerance)
		return true
	}

	switch skill.Skill(calc.skillID) {
	case skill.MesoExplosion:
		// Calculate Meso Explosion damage based on actual meso drop amounts
		if calc.skill == nil {
			result.MinDamage = 0
			result.MaxDamage = 0
			result.ExpectedDmg = 0
			applyResultTolerance(result, 1.0)
			return true
		}

		// Get total mesos from the drops
		totalMesos := int32(0)
		if calc.player.inst != nil && len(info.mesoDropIDs) > 0 {
			for _, dropID := range info.mesoDropIDs {
				if drop, ok := calc.player.inst.dropPool.drops[dropID]; ok && drop.mesos > 0 {
					// Check for integer overflow before adding
					if totalMesos > math.MaxInt32-drop.mesos {
						totalMesos = math.MaxInt32
						break
					}
					totalMesos += drop.mesos
				}
			}
		}

		// If no mesos found, damage should be 0
		if totalMesos == 0 {
			result.MinDamage = 0
			result.MaxDamage = 0
			result.ExpectedDmg = 0
			applyResultTolerance(result, 1.0)
			return true
		}

		// Calculate damage using the correct formula
		xValue := float64(calc.skill.X)
		mesos := float64(totalMesos)
		var ratio float64

		if mesos <= constant.MesoExplosionLowMesoThreshold {
			ratio = (mesos*constant.MesoExplosionLowMesoMultiplier + constant.MesoExplosionLowMesoOffset) / constant.MesoExplosionLowMesoDivisor
		} else {
			ratio = mesos / (mesos + constant.MesoExplosionHighMesoDivisorOffset)
		}

		// MIN: (50 * xValue) * 0.5 * ratio
		// MAX: (50 * xValue) * ratio
		minDamage := (50.0 * xValue) * 0.5 * ratio
		maxDamage := (50.0 * xValue) * ratio

		result.MinDamage = minDamage
		result.MaxDamage = maxDamage
		result.ExpectedDmg = (minDamage + maxDamage) / 2.0

		// Validate with tolerance
		applyResultTolerance(result, constant.MesoExplosionDamageVarianceTolerance)

		return true

	case skill.ShadowMeso:
		if calc.skill == nil {
			return true
		}

		mesoCount := float64(calc.skill.X)
		result.MinDamage = 10.0 * mesoCount
		result.MaxDamage = 10.0 * mesoCount
		result.ExpectedDmg = 10.0 * mesoCount

		if rngBuf != nil && calc.skill.Prop > 0 {
			roll := NormalizeDamageRng(rngBuf.DamageRaw(hitIdx, 3))
			prop := float64(calc.skill.Prop) / 100.0
			if calc.skill.Prop > 100 {
				prop = float64(calc.skill.Prop) / 1000.0
			}
			if roll < prop {
				result.IsCrit = true
				bonusDmg := float64(100 + calc.skill.X)
				result.MinDamage *= bonusDmg * 0.01
				result.MaxDamage *= bonusDmg * 0.01
				result.ExpectedDmg *= bonusDmg * 0.01
			}
		}

		applyResultTolerance(result, 1.0+constant.DamageVarianceTolerance)
		return true

	case skill.ShadowWeb:
		if calc.skill == nil || calc.skillLevel <= 0 {
			return true
		}

		divisor := 50.0 - float64(calc.skillLevel)
		if divisor <= 0 {
			divisor = 1.0
		}
		dmg := float64(mob.maxHP) / divisor
		result.MinDamage = dmg
		result.MaxDamage = dmg
		result.ExpectedDmg = dmg
		applyResultTolerance(result, 1.0+constant.DamageVarianceTolerance)
		return true

	case skill.PoisonMyst:
		if calc.skillLevel <= 0 {
			return true
		}

		divisor := 70.0 - float64(calc.skillLevel)
		if divisor <= 0 {
			divisor = 1.0
		}
		dmg := float64(mob.maxHP) / divisor
		result.MinDamage = dmg
		result.MaxDamage = dmg
		result.ExpectedDmg = dmg
		applyResultTolerance(result, 1.0+constant.DamageVarianceTolerance)
		return true
	}

	return false
}

func (calc *DamageCalculator) CalculateBaseDamageRange(mob *monster, hitIdx int) DamageRange {
	str := float64(calc.GetTotalStr())
	dex := float64(calc.GetTotalDex())
	intl := float64(calc.GetTotalInt())
	luk := float64(calc.GetTotalLuk())
	watk := float64(calc.GetTotalWatk())

	if calc.HasPhysicalOrMagicImmunity(mob) {
		return DamageRange{Min: 0, Max: 1, Valid: true}
	}

	if calc.attackType == attackMagic && calc.skillID != 0 {
		return calc.CalculateMagicDamageRange()
	}

	if watk <= 0 {
		if calc.weaponID == 0 {
			return DamageRange{Min: 1, Max: 1, Valid: true}
		}
		return DamageRange{Valid: false, Reason: "effective weapon attack is zero"}
	}

	switch skill.Skill(calc.skillID) {
	case skill.DragonRoar, skill.SuperDragonRoar:
		maxDmg := (str*4.0 + dex) * watk / 100.0
		return DamageRange{Min: math.Ceil(maxDmg * math.Max(calc.masteryMod, 0.6)), Max: math.Ceil(maxDmg), Valid: true}
	}

	isSwing := calc.attackAction >= constant.AttackActionSwing1H1 && calc.attackAction <= constant.AttackActionSwing2H7
	primaryStat := str
	secondaryStat := dex
	if calc.weaponType == constant.WeaponTypeBow2 || calc.weaponType == constant.WeaponTypeCrossbow2 {
		primaryStat = dex
		secondaryStat = str
	} else if calc.weaponType == constant.WeaponTypeClaw2 || (calc.weaponType == constant.WeaponTypeDagger2 && calc.player.job/100 == 4) {
		primaryStat = luk
		secondaryStat = str + dex
	}

	if calc.weaponID == 0 || calc.weaponType == constant.WeaponTypeNone {
		if calc.weaponID != 0 {
			return DamageRange{Valid: false, Reason: "equipped weapon type is unsupported by validator"}
		}

		if calc.player.job >= 500 && calc.player.job < 600 {
			attack := math.Min(math.Floor((2.0*float64(calc.player.level)+31.0)/3.0), 31.0)
			maxDmg := math.Ceil((str*4.2 + dex) * attack / 100.0)
			return DamageRange{Min: math.Ceil(maxDmg * math.Max(calc.masteryMod, 0.6)), Max: maxDmg, Valid: true}
		}

		return DamageRange{Min: 1, Max: 1, Valid: true}
	}

	weaponMultiplier, valid := calc.GetWeaponDamageMultiplier(isSwing)
	if !valid {
		return DamageRange{Valid: false, Reason: "missing weapon multiplier for equipped weapon"}
	}

	maxDmg := math.Ceil(((weaponMultiplier * primaryStat) + secondaryStat) * watk / 100.0)
	minDmg := math.Ceil(((weaponMultiplier * primaryStat * calc.masteryMod) + secondaryStat) * watk / 100.0)

	if calc.attackType == attackMagic {
		maxDmg = math.Ceil(((weaponMultiplier * intl) + luk) * watk / 100.0)
		minDmg = math.Ceil((((weaponMultiplier * intl) * calc.masteryMod) + luk) * watk / 100.0)
	}

	if minDmg < 1 {
		minDmg = 1
	}
	if maxDmg < 1 {
		maxDmg = 1
	}

	return DamageRange{Min: minDmg, Max: maxDmg, Valid: true}
}

func (calc *DamageCalculator) CalculateMagicDamageRange() DamageRange {
	totalMAD := float64(calc.GetTotalMatk())
	intl := float64(calc.GetTotalInt())
	luk := float64(calc.GetTotalLuk())

	if totalMAD <= 0 {
		return DamageRange{Valid: false, Reason: "effective magic attack is zero"}
	}
	if calc.skill == nil {
		return DamageRange{Valid: false, Reason: "missing skill data for magic validation"}
	}

	if skill.Skill(calc.skillID) == skill.Heal {
		baseMax := math.Round((intl*4.8 + luk*4.0) * totalMAD / 1000.0)
		if calc.skill.Hp > 0 {
			baseMax *= float64(calc.skill.Hp) / 100.0
		}
		return DamageRange{Min: math.Ceil(baseMax * 0.5), Max: math.Ceil(baseMax), Valid: true}
	}

	if calc.skill.Mad <= 0 {
		return DamageRange{Valid: false, Reason: "missing skill magic multiplier"}
	}

	// Pre-BB magic uses the spell's NX mad value as the skill multiplier.
	spellBase := math.Ceil((totalMAD*math.Ceil(totalMAD/1000.0)+totalMAD)/30.0) + math.Ceil(intl/200.0)
	baseMax := spellBase * float64(calc.skill.Mad)
	baseMin := math.Ceil(baseMax * calc.masteryMod)

	return DamageRange{Min: baseMin, Max: baseMax, Valid: true}
}

func (calc *DamageCalculator) CalculateSummonMagicDamageRange() DamageRange {
	totalMAD := float64(calc.GetTotalMatk())
	intl := float64(calc.GetTotalInt())

	if totalMAD <= 0 {
		return DamageRange{Valid: false, Reason: "effective magic attack is zero"}
	}
	if calc.skill == nil || calc.skill.Mad <= 0 {
		return DamageRange{Valid: false, Reason: "missing summon magic multiplier"}
	}

	maxDmg := math.Ceil((intl * totalMAD * float64(calc.skill.Mad)) / 10000.0)
	minDmg := math.Ceil(maxDmg * calc.masteryMod)
	if minDmg < 1 {
		minDmg = 1
	}
	if maxDmg < 1 {
		maxDmg = 1
	}
	return DamageRange{Min: minDmg, Max: maxDmg, Valid: true}
}

func (calc *DamageCalculator) CalculateSummonPhysicalDamageRange() DamageRange {
	baseRange := calc.CalculateBaseDamageRange(nil, 0)
	if !baseRange.Valid {
		return baseRange
	}
	return baseRange
}

func (calc *DamageCalculator) ApplySkillModifiers(minDmg, maxDmg float64, ampData *ElementAmpData, mob *monster) (float64, float64) {
	if calc.skill == nil {
		return minDmg, maxDmg
	}

	if calc.attackType == attackMagic {
		elemMod := float64(ampData.Magic) / 100.0
		minDmg *= elemMod
		maxDmg *= elemMod
	}

	skillElemMod := calc.GetElementalDamageModifier(mob)
	if skillElemMod == constant.ElementModifierNullify {
		return 1, 1
	}
	if skillElemMod == constant.ElementModifierOneAndHalf {
		minDmg *= 1.5
		maxDmg *= 1.5
	}

	if calc.UsesSkillDamageMultiplier() {
		skillMultiplier := float64(calc.skill.Damage) / 100.0
		if skillMultiplier > 0 {
			minDmg *= skillMultiplier
			maxDmg *= skillMultiplier
		}
	}

	return minDmg, maxDmg
}

func (calc *DamageCalculator) CalculateDefenseReductionBounds(mob *monster) (float64, float64) {
	if skill.Skill(calc.skillID) == skill.Sacrifice ||
		skill.Skill(calc.skillID) == skill.Assaulter {
		return 0, 0
	}

	var mobDef float64
	if calc.attackType == attackMagic {
		mobDef = float64(mob.mdDamage)
	} else {
		mobDef = float64(mob.pdDamage)
	}

	redMin := mobDef * 0.5
	redMax := mobDef * 0.6
	return redMin, redMax
}

func (calc *DamageCalculator) CheckCritical(rngBuf *DamageRngBuffer, hitIdx int) bool {
	if !calc.isRanged || calc.critSkill == nil || rngBuf == nil {
		return false
	}

	if skill.Skill(calc.skillID) == skill.Blizzard {
		return false
	}

	roll := NormalizeDamageRng(rngBuf.CriticalRaw(hitIdx, 4))

	prop := float64(calc.critSkill.Prop) / 100.0
	if calc.critSkill.Prop > 100 {
		prop = float64(calc.critSkill.Prop) / 1000.0
	}

	return roll < prop
}

func (calc *DamageCalculator) GetCriticalDamageMultiplier() float64 {
	if calc.critSkill == nil || calc.critSkill.Damage <= 0 {
		return 1.0
	}

	return float64(calc.critSkill.Damage) / 100.0
}

func (calc *DamageCalculator) GetAfterModifier(targetIdx int, baseDmg float64) float64 {
	if calc.skill == nil {
		return 1.0
	}

	if calc.attackOption == constant.AttackOptionSlashBlastFA {
		return constant.SlashBlastFAModifiers[targetIdx]
	}

	if calc.skillID == int32(skill.ArrowBomb) {
		if targetIdx > 0 {
			return float64(calc.skill.X) * 0.01
		}
		if baseDmg > 0 {
			return 0.5
		}
		return 0
	}

	if calc.skillID == int32(skill.IronArrow) {
		return constant.IronArrowModifiers[targetIdx]
	}

	return 1.0
}

func (calc *DamageCalculator) GetIsMiss(rngBuf *DamageRngBuffer, targetAccuracy float64, mob *monster, hitIdx int) bool {
	roll := NormalizeDamageRng(rngBuf.AccuracyRaw(hitIdx, 3))

	var minModifier, maxModifier float64
	if calc.attackType == attackMagic {
		minModifier = 0.5
		maxModifier = 1.2
	} else {
		minModifier = 0.7
		maxModifier = 1.3
	}

	minTACC := targetAccuracy * minModifier
	randTACC := minTACC + (targetAccuracy*maxModifier-minTACC)*roll
	mobAvoid := float64(mob.eva)

	return randTACC < mobAvoid
}

func (calc *DamageCalculator) GetElementAmplification() *ElementAmpData {
	jobID := calc.player.job
	ampSkillID := int32(0)

	if jobID/10 == 21 {
		ampSkillID = int32(skill.ElementAmplification)
	} else if jobID/10 == 22 {
		ampSkillID = int32(skill.ILElementAmplification)
	}

	ampData := &ElementAmpData{Magic: 100, Mana: 100}
	if ampSkillID > 0 {
		if ampSkillInfo, ok := calc.player.skills[ampSkillID]; ok {
			skillData, err := nx.GetPlayerSkill(ampSkillID)
			if err == nil && len(skillData) > 0 && ampSkillInfo.Level > 0 {
				idx := int(ampSkillInfo.Level) - 1
				if idx < len(skillData) {
					ampData.Mana = int(skillData[idx].X)
					ampData.Magic = int(skillData[idx].Y)
				}
			}
		}
	}
	return ampData
}

func (calc *DamageCalculator) GetTargetAccuracy(mob *monster) float64 {
	levelDiff := int(mob.level) - int(calc.player.level)
	if levelDiff < 0 {
		levelDiff = 0
	}

	var accuracy int
	if calc.attackType == attackMagic {
		accuracy = int(5 * (calc.player.intt/10 + calc.player.luk/10))
	} else {
		accuracy = int(calc.player.dex)
	}

	return float64(accuracy*100) / (float64(levelDiff*10) + 255.0)
}

func (calc *DamageCalculator) GetMasteryModifier() float64 {
	var mastery int
	if calc.attackType == attackMagic {
		if calc.skill != nil {
			mastery = int(calc.skill.Mastery)
		}
	} else {
		mastery = calc.GetWeaponMastery()
	}
	return (float64(mastery)*5.0 + 10.0) * 0.009000000000000001
}

func (calc *DamageCalculator) GetWeaponMastery() int {
	switch calc.weaponType {
	case constant.WeaponTypeBow2, constant.WeaponTypeCrossbow2, constant.WeaponTypeClaw2:
		if !calc.isRanged {
			return 0
		}
	default:
		if calc.isRanged {
			return 0
		}
	}

	var skillID int32
	switch calc.weaponType {
	case constant.WeaponTypeSword1H, constant.WeaponTypeSword2H:
		if calc.player.job/10 == 11 {
			skillID = int32(skill.SwordMastery)
		} else {
			skillID = int32(skill.PageSwordMastery)
		}
	case constant.WeaponTypeAxe1H, constant.WeaponTypeAxe2H:
		skillID = int32(skill.AxeMastery)
	case constant.WeaponTypeBW1H, constant.WeaponTypeBW2H:
		skillID = int32(skill.BwMastery)
	case constant.WeaponTypeDagger2:
		skillID = int32(skill.DaggerMastery)
	case constant.WeaponTypeSpear2:
		skillID = int32(skill.SpearMastery)
	case constant.WeaponTypePolearm2:
		skillID = int32(skill.PolearmMastery)
	case constant.WeaponTypeBow2:
		skillID = int32(skill.BowMastery)
	case constant.WeaponTypeCrossbow2:
		skillID = int32(skill.CrossbowMastery)
	case constant.WeaponTypeClaw2:
		skillID = int32(skill.ClawMastery)
	default:
		return 0
	}

	if skillID != 0 {
		if skillInfo, ok := calc.player.skills[skillID]; ok {
			if skillData, err := nx.GetPlayerSkill(skillID); err == nil && len(skillData) > 0 {
				if skillInfo.Level > 0 && int(skillInfo.Level) <= len(skillData) {
					return int(skillData[skillInfo.Level-1].Mastery)
				}
			}
		}
	}
	return 0
}

func (calc *DamageCalculator) GetCritSkill() (byte, *nx.PlayerSkill) {
	if !calc.isRanged {
		return 0, nil
	}

	var skillID int32
	switch calc.weaponType {
	case constant.WeaponTypeBow2, constant.WeaponTypeCrossbow2:
		skillID = int32(skill.CriticalShot)
	case constant.WeaponTypeClaw2:
		skillID = int32(skill.CriticalThrow)
	default:
		return 0, nil
	}

	if skillInfo, ok := calc.player.skills[skillID]; ok {
		if skillData, err := nx.GetPlayerSkill(skillID); err == nil && len(skillData) > 0 {
			if skillInfo.Level > 0 && int(skillInfo.Level) <= len(skillData) {
				return skillInfo.Level, &skillData[skillInfo.Level-1]
			}
		}
	}
	return 0, nil
}

func (calc *DamageCalculator) GetTotalWatk() int16 {
	watk := calc.player.totalWatk - calc.player.str/10
	if watk < 0 {
		watk = 0
	}
	if calc.projectileID != 0 {
		watk += calc.GetProjectileWatk()
	}
	return watk
}

func (calc *DamageCalculator) GetTotalMatk() int16 {
	matk := calc.player.totalMatk - calc.player.intt/10
	if matk < 0 {
		return 0
	}
	return matk
}

func (calc *DamageCalculator) GetTotalAccuracy() int16 {
	return calc.player.totalAccuracy
}

func (calc *DamageCalculator) GetProjectileWatk() int16 {
	if calc.projectileID == 0 {
		return 0
	}

	for _, item := range calc.player.use {
		if item.ID == calc.projectileID {
			return item.watk
		}
	}

	return 0
}

func (calc *DamageCalculator) GetTotalStr() int16 {
	return calc.player.totalStr
}

func (calc *DamageCalculator) GetTotalDex() int16 {
	return calc.player.totalDex
}

func (calc *DamageCalculator) GetTotalInt() int16 {
	return calc.player.totalInt
}

func (calc *DamageCalculator) GetTotalLuk() int16 {
	return calc.player.totalLuk
}

func (calc *DamageCalculator) GetWeaponDamageMultiplier(isSwing bool) (float64, bool) {
	switch calc.weaponType {
	case constant.WeaponTypeBow2:
		return 3.4, true
	case constant.WeaponTypeCrossbow2:
		return 3.6, true
	case constant.WeaponTypeClaw2:
		return 3.6, true
	case constant.WeaponTypeSword1H:
		return 4.0, true
	case constant.WeaponTypeSword2H:
		return 4.6, true
	case constant.WeaponTypeDagger2:
		if calc.player.job/100 == 4 {
			return 3.6, true
		}
		return 4.0, true
	case constant.WeaponTypeWand2, constant.WeaponTypeStaff2:
		return 3.6, true
	case constant.WeaponTypeAxe1H, constant.WeaponTypeBW1H:
		if isSwing {
			return 4.4, true
		}
		return 3.2, true
	case constant.WeaponTypeAxe2H, constant.WeaponTypeBW2H:
		if isSwing {
			return 4.8, true
		}
		return 3.4, true
	case constant.WeaponTypeSpear2:
		if isSwing {
			return 3.0, true
		}
		return 5.0, true
	case constant.WeaponTypePolearm2:
		if isSwing {
			return 5.0, true
		}
		return 3.0, true
	default:
		return 0, false
	}
}

func (calc *DamageCalculator) UsesSkillDamageMultiplier() bool {
	if calc.skill == nil || calc.skillID == 0 || calc.skill.Damage <= 0 {
		return false
	}

	switch skill.Skill(calc.skillID) {
	case skill.ShadowMeso, skill.ShadowWeb, skill.Heal:
		return false
	default:
		return true
	}
}

func (calc *DamageCalculator) HasPhysicalOrMagicImmunity(mob *monster) bool {
	if mob == nil {
		return false
	}

	if calc.attackType == attackMagic {
		return (mob.statBuff & skill.MobStat.MagicImmune) > 0
	}

	return (mob.statBuff & skill.MobStat.PhysicalImmune) > 0
}

func (calc *DamageCalculator) LogSuspiciousAttackShape() {
	expectedTargets, hasTargetLimit := calc.GetExpectedTargetCount()
	if hasTargetLimit && expectedTargets > 0 && int(calc.data.targets) > expectedTargets {
		log.Printf(
			"Suspicious attack target count from player %s (ID: %d): actualTargets=%d, expectedTargets=%d, skill=%d, attackType=%d",
			calc.player.Name,
			calc.player.ID,
			calc.data.targets,
			expectedTargets,
			calc.skillID,
			calc.attackType,
		)
	}

	expectedHits, hasHitLimit := calc.GetExpectedHitCount()
	if !hasHitLimit || expectedHits <= 0 {
		return
	}

	for _, info := range calc.data.attackInfo {
		if len(info.damages) > expectedHits {
			log.Printf(
				"Suspicious attack hit count from player %s (ID: %d): actualHits=%d, expectedHits=%d, skill=%d, attackType=%d, target=%d",
				calc.player.Name,
				calc.player.ID,
				len(info.damages),
				expectedHits,
				calc.skillID,
				calc.attackType,
				info.spawnID,
			)
		}
	}
}

func (calc *DamageCalculator) GetExpectedHitCount() (int, bool) {
	if calc.attackType == attackSummon {
		return 1, true
	}

	if calc.skillID == 0 {
		return 1, true
	}
	if calc.skill == nil {
		return 0, false
	}

	expected := int(calc.skill.AttackCount)
	if int(calc.skill.BulletCount) > expected {
		expected = int(calc.skill.BulletCount)
	}
	if expected <= 0 {
		return 0, false
	}
	if calc.attackOption&constant.AttackOptionShadowPartner != 0 {
		expected *= 2
	}

	return expected, true
}

func (calc *DamageCalculator) GetExpectedTargetCount() (int, bool) {
	if calc.attackType == attackSummon {
		return 1, true
	}
	if calc.skillID == 0 {
		return 1, true
	}
	if calc.skill == nil || calc.skill.MobCount <= 0 {
		return 0, false
	}

	return int(calc.skill.MobCount), true
}

func (calc *DamageCalculator) GetElementalDamageModifier(mob *monster) constant.ElementModifier {
	elemCode := calc.GetAttackElement()
	if elemCode == elementCodeNone {
		elemCode = calc.GetActiveChargeElement()
	}
	if elemCode == elementCodeNone {
		return constant.ElementModifierNormal
	}

	modifier := calc.GetMobElementModifier(mob, elemCode)
	if modifier == constant.ElementModifierOneAndHalf {
		return modifier
	}
	if modifier == constant.ElementModifierNullify {
		return modifier
	}
	return constant.ElementModifierNormal
}

func (calc *DamageCalculator) GetAttackElement() byte {
	switch skill.Skill(calc.skillID) {
	case skill.FireArrow, skill.Explosion, skill.Inferno, skill.ElementComposition:
		return elementCodeFire
	case skill.ColdBeam, skill.IceStrike, skill.Blizzard, skill.ILElementComposition:
		return elementCodeIce
	case skill.ThunderBolt, skill.Lightning:
		return elementCodeLightning
	case skill.Heal, skill.HolyArrow:
		return elementCodeHoly
	case skill.PoisonBreath, skill.PoisonMyst:
		return elementCodePoison
	default:
		return elementCodeNone
	}
}

func (calc *DamageCalculator) GetActiveChargeElement() byte {
	if calc.player.buffs == nil {
		return elementCodeNone
	}

	for skillID := range calc.player.buffs.activeSkillLevels {
		switch skill.Skill(skillID) {
		case skill.BwFireCharge, skill.SwordFireCharge:
			return elementCodeFire
		case skill.BwIceCharge, skill.SwordIceCharge:
			return elementCodeIce
		case skill.BwLitCharge, skill.SwordLitCharge:
			return elementCodeLightning
		}
	}

	return elementCodeNone
}

func (calc *DamageCalculator) GetMobElementModifier(mob *monster, elemCode byte) constant.ElementModifier {
	if mob == nil {
		return constant.ElementModifierNormal
	}

	mobData, err := nx.GetMob(mob.id)
	if err != nil || mobData.ElemAttr == "" {
		return constant.ElementModifierNormal
	}

	return parseMobElementModifier(mobData.ElemAttr, calc.GetElementRune(elemCode))
}

func parseMobElementModifier(elemAttr string, elemRune byte) constant.ElementModifier {
	if elemRune == 0 || elemAttr == "" {
		return constant.ElementModifierNormal
	}

	for i := 0; i+1 < len(elemAttr); i += 2 {
		if elemAttr[i] != elemRune {
			continue
		}

		switch elemAttr[i+1] {
		case '1':
			return constant.ElementModifierNullify
		case '2':
			return constant.ElementModifierHalf
		case '3':
			return constant.ElementModifierOneAndHalf
		default:
			return constant.ElementModifierNormal
		}
	}

	return constant.ElementModifierNormal
}

func (calc *DamageCalculator) GetElementRune(elemCode byte) byte {
	switch elemCode {
	case elementCodeFire:
		return 'F'
	case elementCodeIce:
		return 'I'
	case elementCodeLightning:
		return 'L'
	case elementCodeHoly:
		return 'H'
	case elementCodePoison:
		return 'P'
	default:
		return 0
	}
}

func applyResultTolerance(result *CalcHitResult, multiplier float64) {
	result.ToleranceMax = math.Ceil(result.MaxDamage * multiplier)
	if result.ToleranceMax < 1 && result.ClientDamage > 0 && result.MaxDamage > 0 {
		result.ToleranceMax = 1
	}
	result.IsValid = float64(result.ClientDamage) <= result.ToleranceMax
	if !result.IsValid {
		result.ValidationReason = "client damage exceeds tolerated cap"
	}
}
