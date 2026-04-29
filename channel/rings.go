package channel

import (
	"database/sql"
	"fmt"
	"math"
	"sort"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/mpacket"
)

const ringEffectMaxDistance = 160

type ringKind byte

const (
	ringKindNone ringKind = iota
	ringKindSoloEffect
	ringKindCouple
	ringKindLabel
	ringKindQuote
	ringKindFriendship
	ringKindWedding
)

type ringRecord struct {
	ID            int32
	ItemID        int32
	OwnerID       int32
	PartnerRingID int32
	PartnerID     int32
	PartnerName   string
	Kind          ringKind
}

func classifyRingItem(itemID int32) ringKind {
	switch {
	case itemID == 1112000:
		return ringKindSoloEffect
	case itemID == 1112001 || itemID == 1112002 || itemID == 1112003 || itemID == 1112005 || itemID == 1112006:
		return ringKindCouple
	case itemID >= 1112100 && itemID <= 1112120:
		return ringKindLabel
	case (itemID >= 1112200 && itemID <= 1112230) || itemID == 1112808:
		return ringKindQuote
	case itemID == 1112800 || itemID == 1112801 || itemID == 1112802:
		return ringKindFriendship
	case itemID == 1112803 || itemID == 1112806 || itemID == 1112807 || itemID == 1112809:
		return ringKindWedding
	default:
		return ringKindNone
	}
}

func (k ringKind) paired() bool {
	return k == ringKindCouple || k == ringKindFriendship || k == ringKindWedding
}

func loadRingRecords(items []Item) map[int32]ringRecord {
	ids := make([]int32, 0, len(items))
	seen := make(map[int32]struct{})
	for _, item := range items {
		if item.ringID <= 0 {
			continue
		}
		if _, ok := seen[item.ringID]; ok {
			continue
		}
		seen[item.ringID] = struct{}{}
		ids = append(ids, item.ringID)
	}
	if len(ids) == 0 {
		return map[int32]ringRecord{}
	}

	query := "SELECT id, itemID, ownerCharacterID, partnerRingID, partnerCharacterID, partnerName, ringType FROM rings WHERE id IN (?"
	args := make([]any, 0, len(ids))
	args = append(args, ids[0])
	for _, id := range ids[1:] {
		query += ",?"
		args = append(args, id)
	}
	query += ")"

	rows, err := common.DB.Query(query, args...)
	if err != nil {
		return map[int32]ringRecord{}
	}
	defer rows.Close()

	records := make(map[int32]ringRecord, len(ids))
	for rows.Next() {
		var rec ringRecord
		var kind int32
		if err := rows.Scan(&rec.ID, &rec.ItemID, &rec.OwnerID, &rec.PartnerRingID, &rec.PartnerID, &rec.PartnerName, &kind); err != nil {
			continue
		}
		rec.Kind = ringKind(kind)
		records[rec.ID] = rec
	}
	return records
}

func createPairedRingRecords(itemID int32, owner, partner *Player) (int32, int32, error) {
	kind := classifyRingItem(itemID)
	if !kind.paired() {
		return 0, 0, fmt.Errorf("item %d is not a paired ring", itemID)
	}

	tx, err := common.DB.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	ownerRes, err := tx.Exec("INSERT INTO rings (itemID, ownerCharacterID, partnerRingID, partnerCharacterID, partnerName, ringType) VALUES (?, ?, 0, ?, ?, ?)", itemID, owner.ID, partner.ID, partner.Name, int(kind))
	if err != nil {
		return 0, 0, err
	}
	ownerRingID64, err := ownerRes.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	partnerRes, err := tx.Exec("INSERT INTO rings (itemID, ownerCharacterID, partnerRingID, partnerCharacterID, partnerName, ringType) VALUES (?, ?, ?, ?, ?, ?)", itemID, partner.ID, ownerRingID64, owner.ID, owner.Name, int(kind))
	if err != nil {
		return 0, 0, err
	}
	partnerRingID64, err := partnerRes.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	if _, err = tx.Exec("UPDATE rings SET partnerRingID=? WHERE id=?", partnerRingID64, ownerRingID64); err != nil {
		return 0, 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}
	return int32(ownerRingID64), int32(partnerRingID64), nil
}

func deleteRingRecordPair(ringID int32) error {
	if ringID <= 0 {
		return nil
	}
	var partnerRingID sql.NullInt32
	if err := common.DB.QueryRow("SELECT partnerRingID FROM rings WHERE id=?", ringID).Scan(&partnerRingID); err != nil && err != sql.ErrNoRows {
		return err
	}
	if partnerRingID.Valid && partnerRingID.Int32 > 0 {
		if _, err := common.DB.Exec("DELETE FROM rings WHERE id IN (?, ?)", ringID, partnerRingID.Int32); err != nil {
			return err
		}
		return nil
	}
	_, err := common.DB.Exec("DELETE FROM rings WHERE id=?", ringID)
	return err
}

func (d *Player) refreshRingRecords() {
	all := make([]Item, 0, len(d.equip)+len(d.use)+len(d.setUp)+len(d.etc)+len(d.cash))
	all = append(all, d.equip...)
	all = append(all, d.use...)
	all = append(all, d.setUp...)
	all = append(all, d.etc...)
	all = append(all, d.cash...)
	d.rings = loadRingRecords(all)
	if d.partnerID > 0 {
		d.marriageID = lookupMarriageID(d.ID)
	} else {
		d.marriageID = -1
	}
}

func (d *Player) ringRecord(item Item) *ringRecord {
	if item.ringID <= 0 || d.rings == nil {
		return nil
	}
	rec, ok := d.rings[item.ringID]
	if !ok {
		return nil
	}
	return &rec
}

func (d *Player) activeRing(kind ringKind) *ringRecord {
	items := append([]Item(nil), d.equip...)
	sort.Slice(items, func(i, j int) bool { return items[i].slotID < items[j].slotID })
	for _, item := range items {
		if item.slotID >= 0 || classifyRingItem(item.ID) != kind {
			continue
		}
		if rec := d.ringRecord(item); rec != nil {
			return rec
		}
	}
	return nil
}

func (d *Player) ringRecordsByKind(kind ringKind) []ringRecord {
	items := append([]Item(nil), d.equip...)
	sort.Slice(items, func(i, j int) bool {
		if items[i].slotID == items[j].slotID {
			return items[i].ringID < items[j].ringID
		}
		return items[i].slotID < items[j].slotID
	})
	seen := make(map[int32]struct{})
	var out []ringRecord
	for _, item := range items {
		if classifyRingItem(item.ID) != kind || item.ringID <= 0 {
			continue
		}
		if _, ok := seen[item.ringID]; ok {
			continue
		}
		if rec := d.ringRecord(item); rec != nil {
			seen[item.ringID] = struct{}{}
			out = append(out, *rec)
		}
	}
	return out
}

func (d *Player) activeMarriageRingItemID() int32 {
	items := append([]Item(nil), d.equip...)
	sort.Slice(items, func(i, j int) bool { return items[i].slotID < items[j].slotID })
	for _, item := range items {
		if item.slotID >= 0 || classifyRingItem(item.ID) != ringKindWedding {
			continue
		}
		return item.ID
	}
	return 0
}

func (d *Player) firstRingIDByKind(kind ringKind) int32 {
	for _, item := range d.equip {
		if classifyRingItem(item.ID) == kind && item.ringID > 0 {
			return item.ringID
		}
	}
	return 0
}

func (d *Player) encodeRemoteRingLooks(pkt *mpacket.Packet) {
	encodeRemotePairRingLook(pkt, d.activeRing(ringKindCouple))
	encodeRemoteFriendRingLook(pkt, d.activeRing(ringKindFriendship))
	encodeRemoteMarriageRingLook(pkt, d)
}

func (d *Player) activePairedRingPartners() map[int32]ringKind {
	out := make(map[int32]ringKind)
	for _, kind := range []ringKind{ringKindCouple, ringKindFriendship, ringKindWedding} {
		ring := d.activeRing(kind)
		if ring == nil || ring.PartnerID <= 0 {
			continue
		}
		out[ring.PartnerID] = kind
	}
	return out
}

func (d *Player) isPairRingEffectActiveWith(partner *Player, kind ringKind) bool {
	if d == nil || partner == nil || d.inst == nil || partner.inst == nil || d.inst != partner.inst {
		return false
	}
	selfRing := d.activeRing(kind)
	partnerRing := partner.activeRing(kind)
	if selfRing == nil || partnerRing == nil {
		return false
	}
	if selfRing.PartnerID != partner.ID || partnerRing.PartnerID != d.ID {
		return false
	}
	if selfRing.PartnerRingID != partnerRing.ID || partnerRing.PartnerRingID != selfRing.ID {
		return false
	}
	dx := int(d.pos.x) - int(partner.pos.x)
	dy := int(d.pos.y) - int(partner.pos.y)
	return math.Abs(float64(dx)) <= ringEffectMaxDistance && math.Abs(float64(dy)) <= ringEffectMaxDistance
}

func encodeRemotePairRingLook(pkt *mpacket.Packet, ring *ringRecord) {
	if ring == nil {
		pkt.WriteByte(0)
		return
	}
	pkt.WriteByte(1)
	pkt.WriteInt32(ring.ID)
	pkt.WriteInt32(0)
	pkt.WriteInt32(ring.PartnerRingID)
	pkt.WriteInt32(0)
	pkt.WriteInt32(ring.ItemID)
}

func encodeRemoteFriendRingLook(pkt *mpacket.Packet, ring *ringRecord) {
	encodeRemotePairRingLook(pkt, ring)
}

func encodeRemoteMarriageRingLook(pkt *mpacket.Packet, subject *Player) {
	itemID := subject.activeMarriageRingItemID()
	if itemID <= 0 || subject.partnerID <= 0 {
		pkt.WriteByte(0)
		return
	}
	pkt.WriteByte(1)
	pkt.WriteInt32(subject.ID)
	pkt.WriteInt32(subject.partnerID)
	pkt.WriteInt32(itemID)
}

func (d *Player) encodeLocalRingRecords(pkt *mpacket.Packet) {
	crush := d.ringRecordsByKind(ringKindCouple)
	pkt.WriteInt16(int16(len(crush)))
	for _, ring := range crush {
		pkt.WriteInt32(ring.PartnerID)
		pkt.WritePaddedString(ring.PartnerName, 13)
		pkt.WriteInt32(ring.ID)
		pkt.WriteInt32(0)
		pkt.WriteInt32(ring.PartnerRingID)
		pkt.WriteInt32(0)
	}

	friend := d.ringRecordsByKind(ringKindFriendship)
	pkt.WriteInt16(int16(len(friend)))
	for _, ring := range friend {
		pkt.WriteInt32(ring.PartnerID)
		pkt.WritePaddedString(ring.PartnerName, 13)
		pkt.WriteInt32(ring.ID)
		pkt.WriteInt32(0)
		pkt.WriteInt32(ring.PartnerRingID)
		pkt.WriteInt32(0)
		pkt.WriteInt32(ring.ItemID)
	}

	if d.partnerID <= 0 || d.marriageID <= 0 {
		pkt.WriteInt16(0)
		return
	}
	pkt.WriteInt16(1)
	pkt.WriteInt32(d.marriageID)
	if d.gender == 0 {
		pkt.WriteInt32(d.ID)
		pkt.WriteInt32(d.partnerID)
	} else {
		pkt.WriteInt32(d.partnerID)
		pkt.WriteInt32(d.ID)
	}
	activeWeddingRingID := d.activeMarriageRingItemID()
	if activeWeddingRingID > 0 {
		pkt.WriteInt16(3)
		pkt.WriteInt32(activeWeddingRingID)
		pkt.WriteInt32(activeWeddingRingID)
	} else {
		pkt.WriteInt16(1)
		predicted := weddingRingOutcome(d.marriageItemID)
		pkt.WriteInt32(predicted)
		pkt.WriteInt32(predicted)
	}
	selfName := d.Name
	partnerName := CharacterNameByID(d.partnerID, "")
	if d.gender == 0 {
		pkt.WritePaddedString(selfName, 13)
		pkt.WritePaddedString(partnerName, 13)
	} else {
		pkt.WritePaddedString(partnerName, 13)
		pkt.WritePaddedString(selfName, 13)
	}
}
