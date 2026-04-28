package channel

import (
	"testing"

	"github.com/Hucaru/Valhalla/constant"
	"github.com/Hucaru/Valhalla/nx"
)

func TestBeginnerBasicAttackHasNonZeroMax(t *testing.T) {
	plr := &Player{level: 1, str: 12, dex: 5, intt: 4, luk: 4}
	plr.recalculateTotalStats()

	calc := NewDamageCalculator(plr, &attackData{skillID: 0}, attackMelee)
	info := &attackInfo{damages: []int32{1}, isCritical: []bool{false}}
	mob := &monster{id: 100100, level: 1}

	result := calc.CalculateHit(mob, info, nil, &ElementAmpData{Magic: 100, Mana: 100}, calc.GetTargetAccuracy(mob), 0, 0)

	if result.ValidationSkipped {
		t.Fatalf("expected damage validation to run, got skipped: %s", result.ValidationReason)
	}
	if result.MaxDamage < 1 {
		t.Fatalf("expected non-zero max damage for beginner basic attack, got %.0f", result.MaxDamage)
	}
	if !result.IsValid {
		t.Fatalf("expected 1 damage basic attack to be valid, got cap max %.0f", result.MaxDamage)
	}
}

func TestUnknownWeaponSkipsValidationInsteadOfReturningZeroCap(t *testing.T) {
	plr := &Player{
		level: 10,
		job:   0,
		str:   25,
		dex:   10,
		equip: []Item{{slotID: -11, ID: 1362000, watk: 40}},
	}
	plr.recalculateTotalStats()

	calc := NewDamageCalculator(plr, &attackData{skillID: 0}, attackMelee)
	info := &attackInfo{damages: []int32{10}, isCritical: []bool{false}}
	mob := &monster{id: 100100, level: 10}

	result := calc.CalculateHit(mob, info, nil, &ElementAmpData{Magic: 100, Mana: 100}, calc.GetTargetAccuracy(mob), 0, 0)

	if !result.ValidationSkipped {
		t.Fatalf("expected validation skip for unsupported weapon type, got max %.0f", result.MaxDamage)
	}
	if !result.IsValid {
		t.Fatalf("expected skipped validation to leave damage untouched")
	}
}

func TestPhysicalImmunityCapsToOneDamage(t *testing.T) {
	plr := &Player{level: 30, str: 60, dex: 20}
	plr.recalculateTotalStats()

	calc := NewDamageCalculator(plr, &attackData{skillID: 0}, attackMelee)
	info := &attackInfo{damages: []int32{5}, isCritical: []bool{false}}
	mob := &monster{id: 100100, level: 30, statBuff: 0x40000}

	result := calc.CalculateHit(mob, info, nil, &ElementAmpData{Magic: 100, Mana: 100}, calc.GetTargetAccuracy(mob), 0, 0)

	if result.ValidationSkipped {
		t.Fatalf("expected immunity validation to run, got skipped: %s", result.ValidationReason)
	}
	if result.MaxDamage != 1 {
		t.Fatalf("expected immune mob max damage to be 1, got %.0f", result.MaxDamage)
	}
	if result.IsValid {
		t.Fatalf("expected 5 damage into physical immunity to be invalid")
	}
}

func TestCreateMonsterFromDataCopiesCombatStats(t *testing.T) {
	mob := createMonsterFromData(1, nx.Life{ID: 100100}, nx.Mob{
		Level:      12,
		MDDamage:   7,
		PDDamage:   9,
		MADamage:   11,
		PADamage:   13,
		Eva:        14,
		Acc:        15,
		Invincible: 1,
	}, false, false)

	if mob.level != 12 || mob.mdDamage != 7 || mob.pdDamage != 9 || mob.maDamage != 11 || mob.paDamage != 13 {
		t.Fatalf("monster combat stats were not copied correctly: %+v", mob)
	}
	if mob.eva != 14 || mob.acc != 15 || !mob.invincible {
		t.Fatalf("monster validation fields were not copied correctly: %+v", mob)
	}
}

func TestParseMobElementModifier(t *testing.T) {
	if got := parseMobElementModifier("F3I1", 'F'); got != constant.ElementModifierOneAndHalf {
		t.Fatalf("expected fire weakness modifier, got %v", got)
	}
	if got := parseMobElementModifier("F3I1", 'I'); got != constant.ElementModifierNullify {
		t.Fatalf("expected ice immunity modifier, got %v", got)
	}
	if got := parseMobElementModifier("F3I1", 'L'); got != constant.ElementModifierNormal {
		t.Fatalf("expected normal lightning modifier, got %v", got)
	}
}

func TestExpectedHitCountDoublesForShadowPartner(t *testing.T) {
	plr := &Player{}
	calc := &DamageCalculator{
		player:       plr,
		data:         &attackData{},
		skillID:      int32(4211004),
		skill:        &nx.PlayerSkill{AttackCount: 3},
		attackOption: constant.AttackOptionShadowPartner,
	}

	hits, ok := calc.GetExpectedHitCount()
	if !ok {
		t.Fatal("expected hit count metadata to be available")
	}
	if hits != 6 {
		t.Fatalf("expected doubled hit count, got %d", hits)
	}
}
