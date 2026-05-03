package channel

import (
	"strconv"
	"strings"
	"time"

	"github.com/Hucaru/Valhalla/nx"
)

var pacificLocation = func() *time.Location {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.FixedZone("PST8PDT", -8*60*60)
	}
	return loc
}()

func isExpCouponItem(itemID int32) bool {
	return itemID >= 5211000 && itemID <= 5211046
}

func isDropCouponItem(itemID int32) bool {
	return itemID >= 5360000 && itemID <= 5360012
}

func (p *Player) expCouponMultiplier(now time.Time) float32 {
	return p.rateCouponMultiplier(now, isExpCouponItem)
}

func (p *Player) dropCouponMultiplier(now time.Time) float32 {
	return p.rateCouponMultiplier(now, isDropCouponItem)
}

func (p *Player) rateCouponMultiplier(now time.Time, match func(int32) bool) float32 {
	best := float32(1)

	for _, item := range p.cash {
		if !match(item.ID) || !itemHasFiniteExpiry(item) || itemExpired(item, now) {
			continue
		}

		nxItem, err := nx.GetItem(item.ID)
		if err != nil || nxItem.Rate <= 0 || !couponScheduleActive(item.ID, now) {
			continue
		}

		rate := float32(nxItem.Rate)
		if rate > best {
			best = rate
		}
	}

	return best
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
