package channel

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Hucaru/Valhalla/common"
	"github.com/Hucaru/Valhalla/nx"
)

const rateCouponDuration = 24 * time.Hour

var pacificLocation = func() *time.Location {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.FixedZone("PST8PDT", -8*60*60)
	}
	return loc
}()

func ensureCharacterRateCouponColumns() {
	if common.DB == nil {
		return
	}
	ensureCharacterColumn("expCouponItemID", "INT NOT NULL DEFAULT 0")
	ensureCharacterColumn("expCouponExpiresAt", "BIGINT NOT NULL DEFAULT 0")
	ensureCharacterColumn("dropCouponItemID", "INT NOT NULL DEFAULT 0")
	ensureCharacterColumn("dropCouponExpiresAt", "BIGINT NOT NULL DEFAULT 0")
}

func ensureCharacterColumn(name, ddl string) {
	var count int
	err := common.DB.QueryRow(`SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'characters' AND COLUMN_NAME = ?`, name).Scan(&count)
	if err != nil {
		log.Printf("ensureCharacterColumn probe failed: %s err=%v", name, err)
		return
	}
	if count > 0 {
		return
	}
	query := fmt.Sprintf("ALTER TABLE characters ADD COLUMN %s %s", name, ddl)
	if _, err := common.DB.Exec(query); err != nil {
		log.Printf("ensureCharacterColumn alter failed: %s err=%v", name, err)
	}
}

func isExpCouponItem(itemID int32) bool {
	return itemID >= 5211004 && itemID <= 5211046
}

func isDropCouponItem(itemID int32) bool {
	return itemID >= 5360000 && itemID <= 5360013
}

func (p *Player) activateExpCoupon(itemID int32, now time.Time) {
	p.expCouponItemID = itemID
	p.expCouponExpiresAt = now.Add(rateCouponDuration).UnixMilli()
	p.MarkDirty(DirtyRateCoupons, 300*time.Millisecond)
}

func (p *Player) activateDropCoupon(itemID int32, now time.Time) {
	p.dropCouponItemID = itemID
	p.dropCouponExpiresAt = now.Add(rateCouponDuration).UnixMilli()
	p.MarkDirty(DirtyRateCoupons, 300*time.Millisecond)
}

func (p *Player) cleanupExpiredRateCoupons(now time.Time) {
	changed := false
	if p.expCouponItemID != 0 && (p.expCouponExpiresAt == 0 || now.UnixMilli() >= p.expCouponExpiresAt) {
		p.expCouponItemID = 0
		p.expCouponExpiresAt = 0
		changed = true
	}
	if p.dropCouponItemID != 0 && (p.dropCouponExpiresAt == 0 || now.UnixMilli() >= p.dropCouponExpiresAt) {
		p.dropCouponItemID = 0
		p.dropCouponExpiresAt = 0
		changed = true
	}
	if changed {
		p.MarkDirty(DirtyRateCoupons, 300*time.Millisecond)
	}
}

func (p *Player) expCouponMultiplier(now time.Time) float32 {
	p.cleanupExpiredRateCoupons(now)
	if p.expCouponItemID == 0 || !couponScheduleActive(p.expCouponItemID, now) {
		return 1
	}
	item, err := nx.GetItem(p.expCouponItemID)
	if err != nil || item.Rate <= 0 {
		return 1
	}
	return float32(item.Rate)
}

func (p *Player) dropCouponMultiplier(now time.Time) float32 {
	p.cleanupExpiredRateCoupons(now)
	if p.dropCouponItemID == 0 || !couponScheduleActive(p.dropCouponItemID, now) {
		return 1
	}
	item, err := nx.GetItem(p.dropCouponItemID)
	if err != nil || item.Rate <= 0 {
		return 1
	}
	return float32(item.Rate)
}

func couponScheduleActive(itemID int32, now time.Time) bool {
	item, err := nx.GetItem(itemID)
	if err != nil || len(item.TimeEntries) == 0 {
		return true
	}
	now = now.In(pacificLocation)
	day := now.Weekday()
	dayName := []string{"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"}[day]
	hour := now.Hour()
	for _, entry := range item.TimeEntries {
		parts := strings.Split(entry, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] != dayName && parts[0] != "HOL" {
			continue
		}
		rangeParts := strings.Split(parts[1], "-")
		if len(rangeParts) != 2 {
			continue
		}
		start, err1 := strconv.Atoi(rangeParts[0])
		end, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil {
			continue
		}
		if end == 24 {
			if hour >= start {
				return true
			}
			continue
		}
		if hour >= start && hour < end {
			return true
		}
	}
	return false
}
