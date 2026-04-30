package nx

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Hucaru/gonx"
)

// Mob data from nx
type Mob struct {
	HP, MP             int32 // Not in nx
	MaxHP, HPRecovery  int32
	MaxMP, MPRecovery  int32
	Level              int64
	Exp                int64
	MADamage, MDDamage int64
	PADamage, PDDamage int64
	Speed, Eva, Acc    int64
	ChaseSpeed         int64
	SummonType         int8
	SummonOption       int32
	Boss, Undead       int64
	ElemAttr           string
	Link               int64
	FlySpeed           int64
	NoRegen            int64
	Invincible         int64
	SelfDestruction    int64
	ExplosiveReward    int64
	Skills             map[byte]byte
	Revives            []int32
	Fs                 float64
	Pushed             int64
	BodyAttack         int64
	NoFlip             int64
	NotAttack          int64
	FirstAttack        int64
	RemoveQuest        int64
	RemoveAfter        string
	PublicReward       int64
	DamagedByMob       int64
	DropItemPeriod     int64
	OnlyNormalAttack   int64
	FixedDamage        int64
	HPGaugeHide        int64
	HPTagBGColor       int64
	HPTagColor         int64
	Attacks            []MobAttackInfo
}

type MobAttackInfo struct {
	ConMP        int64
	AttackAfter  int64
	EffectAfter  int64
	Type         int64
	Magic        int64
	DeadlyAttack int64
	Disease      int64
	Level        int64
}

func extractMobs(nodes []gonx.Node, textLookup []string) map[int32]Mob {
	mobs := make(map[int32]Mob)

	search := "/Mob"
	valid := gonx.FindNode(search, nodes, textLookup, func(node *gonx.Node) {
		for i := uint32(0); i < uint32(node.ChildCount); i++ {
			mobNode := nodes[node.ChildID+i]
			name := textLookup[mobNode.NameID]
			subSearch := search + "/" + name + "/info"

			var mob Mob

			valid := gonx.FindNode(subSearch, nodes, textLookup, func(node *gonx.Node) {
				mob = getMob(node, nodes, textLookup)
			})

			if !valid {
				log.Println("Invalid node search:", subSearch)
			}

			name = strings.TrimSuffix(name, filepath.Ext(name))
			mobID, err := strconv.Atoi(name)

			if err != nil {
				log.Println(err)
				continue
			}

			mobs[int32(mobID)] = mob
		}
	})

	if !valid {
		log.Println("Invalid node search:", search)
	}

	return mobs
}

func getMob(node *gonx.Node, nodes []gonx.Node, textLookup []string) Mob {
	mob := Mob{}

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		option := nodes[node.ChildID+i]
		optionName := textLookup[option.NameID]

		switch optionName {
		case "maxHP":
			mob.MaxHP = gonx.DataToInt32(option.Data)
			mob.HP = mob.MaxHP
		case "hpRecovery":
			mob.HPRecovery = gonx.DataToInt32(option.Data)
		case "maxMP":
			mob.MaxMP = gonx.DataToInt32(option.Data)
			mob.MP = mob.MaxMP
		case "mpRecovery":
			mob.MPRecovery = gonx.DataToInt32(option.Data)
		case "level":
			mob.Level = gonx.DataToInt64(option.Data)
		case "exp":
			mob.Exp = gonx.DataToInt64(option.Data)
		case "MADamage":
			mob.MADamage = gonx.DataToInt64(option.Data)
		case "MDDamage":
			mob.MDDamage = gonx.DataToInt64(option.Data)
		case "PADamage":
			mob.PADamage = gonx.DataToInt64(option.Data)
		case "PDDamage":
			mob.PDDamage = gonx.DataToInt64(option.Data)
		case "speed":
			mob.Speed = gonx.DataToInt64(option.Data)
		case "chaseSpeed":
			mob.ChaseSpeed = gonx.DataToInt64(option.Data)
		case "eva":
			mob.Eva = gonx.DataToInt64(option.Data)
		case "acc":
			mob.Acc = gonx.DataToInt64(option.Data)
		case "summonType":
			mob.SummonType = int8(option.Data[0])
		case "summonOption":
			fmt.Println("Got summon option")
			mob.SummonOption = gonx.DataToInt32(option.Data)
		case "boss":
			mob.Boss = gonx.DataToInt64(option.Data)
		case "undead":
			mob.Undead = gonx.DataToInt64(option.Data)
		case "elemAttr":
			mob.ElemAttr = textLookup[gonx.DataToUint32(option.Data)]
		case "link":
			mob.Link = gonx.DataToInt64(option.Data)
		case "flySpeed", "flyspeed":
			mob.FlySpeed = gonx.DataToInt64(option.Data)
		case "noregen": // is this for both hp/mp?
			mob.NoRegen = gonx.DataToInt64(option.Data)
		case "invincible":
			mob.Invincible = gonx.DataToInt64(option.Data)
		case "selfDestruction":
			mob.SelfDestruction = gonx.DataToInt64(option.Data)
		case "explosiveReward": // A way that mob drops can drop?
			mob.ExplosiveReward = gonx.DataToInt64(option.Data)
		case "skill":
			mob.Skills = getSkills(&option, nodes, textLookup)
		case "revive":
			mob.Revives = getRevives(&option, nodes)
		case "fs":
			mob.Fs = gonx.DataToFloat64(option.Data)
		case "pushed":
			mob.Pushed = gonx.DataToInt64(option.Data)
		case "bodyAttack":
			mob.BodyAttack = gonx.DataToInt64(option.Data)
		case "noFlip":
			mob.NoFlip = gonx.DataToInt64(option.Data)
		case "notAttack":
			mob.NotAttack = gonx.DataToInt64(option.Data)
		case "firstAttack", "firstattack":
			mob.FirstAttack = gonx.DataToInt64(option.Data)
		case "removeQuest":
			mob.RemoveQuest = gonx.DataToInt64(option.Data)
		case "removeAfter":
			idLookup := gonx.DataToUint32(option.Data)
			mob.RemoveAfter = textLookup[idLookup]
		case "publicReward":
			mob.PublicReward = gonx.DataToInt64(option.Data)
		case "damagedByMob":
			mob.DamagedByMob = gonx.DataToInt64(option.Data)
		case "dropItemPeriod":
			mob.DropItemPeriod = gonx.DataToInt64(option.Data)
		case "onlyNormalAttack":
			mob.OnlyNormalAttack = gonx.DataToInt64(option.Data)
		case "fixedDamage":
			mob.FixedDamage = gonx.DataToInt64(option.Data)
		case "HPgaugeHide":
			mob.HPGaugeHide = gonx.DataToInt64(option.Data)
		case "hpTagBgcolor":
			mob.HPTagBGColor = gonx.DataToInt64(option.Data)
		case "hpTagColor":
			mob.HPTagColor = gonx.DataToInt64(option.Data)
		default:
			if strings.HasPrefix(optionName, "attack") {
				mob.Attacks = append(mob.Attacks, getMobAttack(&option, nodes, textLookup))
				continue
			}
			log.Println("Unsupported NX mob option:", optionName, "->", option.Data)
		}
	}

	return mob
}

func getMobAttack(node *gonx.Node, nodes []gonx.Node, textLookup []string) MobAttackInfo {
	attack := MobAttackInfo{}

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		child := nodes[node.ChildID+i]
		if textLookup[child.NameID] != "info" {
			continue
		}

		for j := uint32(0); j < uint32(child.ChildCount); j++ {
			option := nodes[child.ChildID+j]
			switch textLookup[option.NameID] {
			case "conMP":
				attack.ConMP = gonx.DataToInt64(option.Data)
			case "attackAfter":
				attack.AttackAfter = gonx.DataToInt64(option.Data)
			case "effectAfter":
				attack.EffectAfter = gonx.DataToInt64(option.Data)
			case "type":
				attack.Type = gonx.DataToInt64(option.Data)
			case "magic":
				attack.Magic = gonx.DataToInt64(option.Data)
			case "deadlyAttack":
				attack.DeadlyAttack = gonx.DataToInt64(option.Data)
			case "disease":
				attack.Disease = gonx.DataToInt64(option.Data)
			case "level":
				attack.Level = gonx.DataToInt64(option.Data)
			}
		}

		break
	}

	return attack
}

func getSkills(node *gonx.Node, nodes []gonx.Node, textLookup []string) map[byte]byte {
	skills := make(map[byte]byte)

	// need to subnode the children of the children to node
	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		skillDir := nodes[node.ChildID+i]

		var id byte
		var level byte

		for j := uint32(0); j < uint32(skillDir.ChildCount); j++ {
			option := nodes[skillDir.ChildID+j]
			optionName := textLookup[option.NameID]

			switch optionName {
			case "level":
				level = option.Data[0]
			case "skill":
				id = option.Data[0]
			case "action":
			case "effectAfter":
			case "count":
			default:
				log.Println("Unsupported NX mob skill option:", optionName, "->", option.Data)
			}
		}

		skills[id] = level
	}

	return skills
}

func getRevives(node *gonx.Node, nodes []gonx.Node) []int32 {
	revives := make([]int32, node.ChildCount)

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		mobID := nodes[node.ChildID+i]
		revives[i] = gonx.DataToInt32(mobID.Data)
	}

	return revives
}
