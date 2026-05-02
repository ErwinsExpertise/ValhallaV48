package channel

import (
	"fmt"
	"time"

	"github.com/Hucaru/Valhalla/nx"
)

const itemExpirySweepInterval = time.Minute

const windowsFileTimeEpoch = int64(116444592000000000)

func itemExpireUnixMilli(expireTime int64) int64 {
	if expireTime == 0 || expireTime == neverExpire {
		return 0
	}
	if expireTime > windowsFileTimeEpoch {
		return (expireTime - windowsFileTimeEpoch) / 10000
	}
	return expireTime
}

func itemHasFiniteExpiry(item Item) bool {
	return itemExpireUnixMilli(item.expireTime) > 0
}

func itemExpired(item Item, now time.Time) bool {
	expiresAt := itemExpireUnixMilli(item.expireTime)
	return expiresAt > 0 && now.UnixMilli() >= expiresAt
}

func expiredItemName(item Item) string {
	meta, err := nx.GetItem(item.ID)
	if err == nil && meta.Name != "" {
		return meta.Name
	}
	return fmt.Sprintf("Item %d", item.ID)
}

func (p *Player) cleanupExpiredInventoryItems(now time.Time, notify bool) []string {
	if p == nil {
		return nil
	}

	items := make([]Item, 0, len(p.equip)+len(p.use)+len(p.setUp)+len(p.etc)+len(p.cash))
	items = append(items, p.equip...)
	items = append(items, p.use...)
	items = append(items, p.setUp...)
	items = append(items, p.etc...)
	items = append(items, p.cash...)
	removedNames := make([]string, 0)
	recalcStats := false

	for _, item := range items {
		if item.pet || !itemExpired(item, now) {
			continue
		}
		if item.invID == 1 && item.slotID < 0 {
			recalcStats = true
		}
		p.removeItem(item, false)
		name := expiredItemName(item)
		removedNames = append(removedNames, name)
		if notify && p.Conn != nil {
			p.Send(packetMessageRedText(fmt.Sprintf("%s has expired and removed from your inventory", name)))
		}
	}

	if recalcStats {
		p.recalculateTotalStats()
	}

	return removedNames
}

func scheduleItemExpiry(server *Server) {
	ticker := time.NewTicker(itemExpirySweepInterval)
	defer ticker.Stop()

	for range ticker.C {
		if server == nil || server.dispatch == nil {
			continue
		}

		select {
		case server.dispatch <- func() {
			now := time.Now()
			server.players.observe(func(plr *Player) {
				if plr == nil {
					return
				}
				plr.cleanupExpiredInventoryItems(now, true)
			})
		}:
		default:
		}
	}
}
