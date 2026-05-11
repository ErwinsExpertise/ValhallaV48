package channel

import (
	"fmt"
	"strings"

	"github.com/Hucaru/Valhalla/constant"
)

func searchLabelForItemType(searchType string) string {
	switch searchType {
	case "item":
		return "Items"
	case "equip":
		return "Equips"
	case "cash":
		return "Cash Items"
	case "use":
		return "Use Items"
	case "setup":
		return "Setup Items"
	case "etc":
		return "Etc Items"
	case "pet":
		return "Pet Items"
	case "accessory":
		return "Accessories"
	case "cap":
		return "Caps"
	case "cape":
		return "Capes"
	case "coat":
		return "Coats"
	case "face":
		return "Faces"
	case "glove":
		return "Gloves"
	case "hair":
		return "Hair"
	case "longcoat":
		return "Longcoats"
	case "pants":
		return "Pants"
	case "petequip":
		return "Pet Equips"
	case "ring":
		return "Rings"
	case "shield":
		return "Shields"
	case "shoes":
		return "Shoes"
	case "weapon":
		return "Weapons"
	default:
		return "Items"
	}
}

var warpDestinationByName = map[string]int32{
	"amherst":        1010000,
	"southperry":     60000,
	"sp":             60000,
	"lith":           104000000,
	"lithharbor":     104000000,
	"lithharbour":    104000000,
	"henesys":        100000000,
	"hene":           100000000,
	"ellinia":        101000000,
	"elli":           101000000,
	"perion":         102000000,
	"kerning":        103000000,
	"kerningcity":    103000000,
	"kc":             103000000,
	"sleepy":         105040300,
	"sleepywood":     105040300,
	"orbis":          200000000,
	"elnath":         211000000,
	"elnathtown":     211000000,
	"ludi":           220000000,
	"ludibrium":      220000000,
	"omega":          221000000,
	"omegasector":    221000000,
	"kft":            222000000,
	"koreanfolk":     222000000,
	"koreanfolktown": 222000000,
	"aqua":           230000000,
	"aquarium":       230000000,
	"aquaroad":       230000000,
	"mulung":         250000000,
	"herbtown":       251000000,
	"nlc":            600000000,
	"newleafcity":    600000000,
	"zipangu":        800000000,
	"mushroomshrine": 800000000,
	"amoria":         680000000,
	"gm":             180000000,
	"balrog":         105090900,
	"guild":          200000301,
	"pap":            constant.MapBossPapulatus,
	"pianus":         constant.MapBossPianus,
	"zakum":          constant.MapBossZakum,
	"kerningpq":      constant.MapKerningPQ,
	"ludipq":         constant.MapLudiPQ,
}

func convertMapNameToID(name string) (int32, bool) {
	id, ok := warpDestinationByName[normalizeWarpDestination(name)]
	return id, ok
}

func normalizeWarpDestination(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, "-", "")
	return name
}

func convertJobNameToID(name string) int16 {
	switch name {
	case "Beginner":
		return 0
	case "Warrior":
		return 100
	case "Fighter":
		return 110
	case "Crusader":
		return 111
	case "Page":
		return 120
	case "WhiteKnight":
		return 121
	case "Spearman":
		return 130
	case "DragonKnight":
		return 131
	case "Magician":
		return 200
	case "FirePoisonWizard":
		return 210
	case "FirePoisonMage":
		return 211
	case "IceLightWizard":
		return 220
	case "IceLightMage":
		return 221
	case "Cleric":
		return 230
	case "Priest":
		return 231
	case "Bowman":
		return 300
	case "Hunter":
		return 310
	case "Ranger":
		return 311
	case "Crossbowman":
		return 320
	case "Sniper":
		return 321
	case "Thief":
		return 400
	case "Assassin":
		return 410
	case "Hermit":
		return 411
	case "Bandit":
		return 420
	case "ChiefBandit":
		return 421
	case "Gm":
		return 500
	case "SuperGm":
		return 510
	default:
		return 0
	}
}

func covnertMobNameToID(name string) ([]int32, error) {
	switch name {
	case "balrog":
		return []int32{constant.MobBalrog}, nil
	case "cbalrog":
		return []int32{constant.MobCrimsonBalrog}, nil
	case "zakum":
		return []int32{
			constant.MobZakumArm1,
			constant.MobZakumArm2,
			constant.MobZakumArm3,
			constant.MobZakumArm4,
			constant.MobZakumArm5,
			constant.MobZakumArm6,
			constant.MobZakumArm7,
			constant.MobZakumArm8,
			constant.MobZakum1Body,
		}, nil
	case "pap":
		return []int32{constant.MobPapalatus}, nil
	case "pianus":
		return []int32{constant.MobPianus}, nil
	case "mushmom":
		return []int32{constant.MobMushmom}, nil
	case "zmushmom":
		return []int32{constant.MobZombieMushmom}, nil
	}

	return nil, fmt.Errorf("unknown mob Name")
}
