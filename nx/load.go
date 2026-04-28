package nx

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hucaru/gonx"
)

var items map[int32]Item
var maps map[int32]Map
var mobs map[int32]Mob
var quests map[int16]Quest
var playerSkills map[int32][]PlayerSkill
var mobSkills map[byte][]MobSkill
var commodities map[int32]Commodity
var packages map[int32][]int32
var itemIDToSN map[int32]int32
var bestItems = make(map[FeaturedKey]int32)
var reactorInfos map[int32]ReactorInfo

// LoadFile into useable types
func LoadFile(fname string) {
	nodes, textLookup, err := loadNXSource(fname)

	if err != nil {
		panic(err)
	}

	items = extractItems(nodes, textLookup)
	maps = extractMaps(nodes, textLookup)
	applyMapNames(maps, nodes, textLookup)
	mobs = extractMobs(nodes, textLookup)
	playerSkills, mobSkills = extractSkills(nodes, textLookup)
	quests = extractQuests(nodes, textLookup)
	commodities = extractCommodities(nodes, textLookup)
	packages = extractPackages(nodes, textLookup)
	reactorInfos = extractReactors(nodes, textLookup)

	loadBestItems()
}

func loadNXSource(path string) ([]gonx.Node, []string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}

	if !info.IsDir() {
		nodes, textLookup, _, _, err := gonx.Parse(path)
		return nodes, textLookup, err
	}

	return loadNXDirectory(path)
}

func loadNXDirectory(dir string) ([]gonx.Node, []string, error) {
	manifest := []string{
		"Base.nx",
		"Character.nx",
		"Effect.nx",
		"Etc.nx",
		"Item.nx",
		"Map.nx",
		"Mob.nx",
		"Morph.nx",
		"Npc.nx",
		"Quest.nx",
		"Reactor.nx",
		"Skill.nx",
		"Sound.nx",
		"String.nx",
		"TamingMob.nx",
		"UI.nx",
	}

	builder := newNXBuilder(len(manifest))

	for i, name := range manifest {
		filePath := filepath.Join(dir, name)
		nodes, textLookup, _, _, err := gonx.Parse(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("parse %s: %w", name, err)
		}

		wrapper := strings.TrimSuffix(name, filepath.Ext(name))
		builder.addWrappedTree(i, wrapper, nodes, textLookup)
	}

	return builder.nodes, builder.text, nil
}

type nxBuilder struct {
	nodes []gonx.Node
	text  []string
	ids   map[string]uint32
}

func newNXBuilder(wrapperCount int) *nxBuilder {
	b := &nxBuilder{
		nodes: []gonx.Node{{ChildCount: uint16(wrapperCount)}},
		text:  []string{""},
		ids:   map[string]uint32{"": 0},
	}

	b.nodes[0].ChildID = uint32(len(b.nodes))
	b.nodes = append(b.nodes, make([]gonx.Node, wrapperCount)...)

	return b
}

func (b *nxBuilder) addText(s string) uint32 {
	if id, ok := b.ids[s]; ok {
		return id
	}

	id := uint32(len(b.text))
	b.text = append(b.text, s)
	b.ids[s] = id
	return id
}

func (b *nxBuilder) addWrappedTree(slot int, name string, srcNodes []gonx.Node, srcText []string) {
	root := &srcNodes[0]
	wrapperIdx := int(b.nodes[0].ChildID) + slot
	b.fillWrapper(wrapperIdx, name, root, srcNodes, srcText)
}

func (b *nxBuilder) fillWrapper(dstIdx int, name string, srcRoot *gonx.Node, srcNodes []gonx.Node, srcText []string) {
	wrapper := &b.nodes[dstIdx]
	wrapper.NameID = b.addText(name)
	wrapper.Type = 0

	if srcRoot.ChildCount == 0 {
		return
	}

	wrapper.ChildCount = srcRoot.ChildCount
	wrapper.ChildID = uint32(len(b.nodes))
	b.nodes = append(b.nodes, make([]gonx.Node, srcRoot.ChildCount)...)

	for i := uint32(0); i < uint32(srcRoot.ChildCount); i++ {
		b.fillNode(int(wrapper.ChildID+i), int(srcRoot.ChildID+i), srcNodes, srcText)
	}
}

func (b *nxBuilder) fillNode(dstIdx int, srcIdx int, srcNodes []gonx.Node, srcText []string) {
	src := srcNodes[srcIdx]
	dst := &b.nodes[dstIdx]
	dst.NameID = b.addText(srcText[src.NameID])
	dst.Type = src.Type
	dst.Data = src.Data
	if src.Type == 3 {
		ref := srcText[gonx.DataToUint32(src.Data)]
		binary.LittleEndian.PutUint32(dst.Data[:4], b.addText(ref))
		for i := 4; i < len(dst.Data); i++ {
			dst.Data[i] = 0
		}
	}

	if src.ChildCount == 0 {
		return
	}

	dst.ChildCount = src.ChildCount
	dst.ChildID = uint32(len(b.nodes))
	b.nodes = append(b.nodes, make([]gonx.Node, src.ChildCount)...)

	for i := uint32(0); i < uint32(src.ChildCount); i++ {
		b.fillNode(int(dst.ChildID+i), int(src.ChildID+i), srcNodes, srcText)
	}
}

// GetItem from loaded nx
func GetItem(id int32) (Item, error) {
	if _, ok := items[id]; !ok {
		return Item{}, fmt.Errorf("invalid item id: %v", id)
	}

	return items[id], nil
}

// GetMap from loaded nx
func GetMap(id int32) (Map, error) {
	if _, ok := maps[id]; !ok {
		return Map{}, fmt.Errorf("invalid map id: %v", id)
	}

	return maps[id], nil
}

// GetMaps from loaded nx
func GetMaps() map[int32]Map {
	return maps
}

// GetMob from loaded nx
func GetMob(id int32) (Mob, error) {
	if _, ok := mobs[id]; !ok {
		return Mob{}, fmt.Errorf("invalid mob id: %v", id)
	}

	return mobs[id], nil
}

func GetQuests() map[int16]Quest {
	return quests
}

func GetQuest(id int16) (Quest, error) {
	if _, ok := quests[id]; !ok {
		return Quest{}, fmt.Errorf("invalid quest id: %v", id)
	}
	return quests[id], nil
}

// GetPlayerSkill from loaded nx
func GetPlayerSkill(id int32) ([]PlayerSkill, error) {
	if _, ok := playerSkills[id]; !ok {
		return []PlayerSkill{}, fmt.Errorf("Invalid player skill id: %v", id)
	}

	return playerSkills[id], nil
}

// GetMobSkill from loaded nx
func GetMobSkill(id byte) ([]MobSkill, error) {
	if _, ok := mobSkills[id]; !ok {
		return []MobSkill{}, fmt.Errorf("Invalid mob skill id: %v", id)
	}

	return mobSkills[id], nil
}

// GetMobSkills from loaded nx
func GetMobSkills(id int32) map[byte]byte {
	mob, err := GetMob(id)
	if err != nil {
		log.Println(err)
		return nil
	}

	return mob.Skills
}

// GetReactorInfo from loaded nx
func GetReactorInfo(id int32) (ReactorInfo, error) {
	if _, ok := reactorInfos[id]; !ok {
		return ReactorInfo{}, fmt.Errorf("Invalid reactor id: %v", id)
	}
	return reactorInfos[id], nil
}

// GetReactorInfoList from loaded nx
func GetReactorInfoList() map[int32]ReactorInfo {
	return reactorInfos
}
