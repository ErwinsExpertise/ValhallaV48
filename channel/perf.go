package channel

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

const perfSampleWindow = 256

type perfStat struct {
	count   uint64
	totalNS int64
	maxNS   int64
	samples [perfSampleWindow]int64
	filled  int
	index   int
}

func (s *perfStat) observe(d time.Duration) {
	ns := d.Nanoseconds()
	if ns < 0 {
		ns = 0
	}
	s.count++
	s.totalNS += ns
	if ns > s.maxNS {
		s.maxNS = ns
	}
	s.samples[s.index] = ns
	if s.filled < len(s.samples) {
		s.filled++
	}
	s.index = (s.index + 1) % len(s.samples)
}

type perfSnapshot struct {
	name      string
	count     uint64
	avg       time.Duration
	p95       time.Duration
	p99       time.Duration
	max       time.Duration
	scoreP99  int64
	scoreMax  int64
	scoreAvg  int64
	sampleCnt int
}

type perfProfiler struct {
	channelID byte

	mu    sync.Mutex
	stats map[string]*perfStat
	start sync.Once
}

func newPerfProfiler(channelID byte) *perfProfiler {
	return &perfProfiler{
		channelID: channelID,
		stats:     make(map[string]*perfStat),
	}
}

func (p *perfProfiler) startReporting() {
	if p == nil {
		return
	}
	p.start.Do(func() {
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				p.report()
			}
		}()
	})
}

func (p *perfProfiler) observe(name string, d time.Duration) {
	if p == nil || name == "" {
		return
	}
	p.mu.Lock()
	stat := p.stats[name]
	if stat == nil {
		stat = &perfStat{}
		p.stats[name] = stat
	}
	stat.observe(d)
	p.mu.Unlock()
}

func (p *perfProfiler) report() {
	if p == nil {
		return
	}

	p.mu.Lock()
	rows := make([]perfSnapshot, 0, len(p.stats))
	for name, stat := range p.stats {
		if stat == nil || stat.count == 0 || stat.filled == 0 {
			continue
		}
		samples := make([]int64, stat.filled)
		copy(samples, stat.samples[:stat.filled])
		sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
		rows = append(rows, perfSnapshot{
			name:      name,
			count:     stat.count,
			avg:       time.Duration(stat.totalNS / int64(stat.count)),
			p95:       time.Duration(percentile(samples, 95)),
			p99:       time.Duration(percentile(samples, 99)),
			max:       time.Duration(stat.maxNS),
			scoreP99:  percentile(samples, 99),
			scoreMax:  stat.maxNS,
			scoreAvg:  stat.totalNS / int64(stat.count),
			sampleCnt: stat.filled,
		})
	}
	p.mu.Unlock()

	if len(rows) == 0 {
		return
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].scoreP99 != rows[j].scoreP99 {
			return rows[i].scoreP99 > rows[j].scoreP99
		}
		if rows[i].scoreMax != rows[j].scoreMax {
			return rows[i].scoreMax > rows[j].scoreMax
		}
		return rows[i].scoreAvg > rows[j].scoreAvg
	})

	limit := 12
	if len(rows) < limit {
		limit = len(rows)
	}
	for i := 0; i < limit; i++ {
		row := rows[i]
		log.Printf("[perf ch=%d] %s count=%d avg=%s p95=%s p99=%s max=%s samples=%d", p.channelID, row.name, row.count, row.avg, row.p95, row.p99, row.max, row.sampleCnt)
	}
}

func percentile(sorted []int64, pct int) int64 {
	if len(sorted) == 0 {
		return 0
	}
	if pct <= 0 {
		return sorted[0]
	}
	if pct >= 100 {
		return sorted[len(sorted)-1]
	}
	idx := (len(sorted)*pct - 1) / 100
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

func (server *Server) enablePerformanceProfiling() {
	if server == nil || server.perf != nil {
		return
	}
	server.perf = newPerfProfiler(server.id)
	server.perf.startReporting()
	log.Printf("enabled gameplay performance profiling for channel %d", server.id)
}

func (server *Server) EnablePerformanceProfiling() {
	server.enablePerformanceProfiling()
}

func (server *Server) observePerf(name string, d time.Duration) {
	if server == nil || server.perf == nil {
		return
	}
	server.perf.observe(name, d)
}

func (server *Server) observeSlowFieldUpdate(inst *fieldInstance, d time.Duration) {
	if server == nil || server.perf == nil || inst == nil || d < 50*time.Millisecond {
		return
	}
	log.Printf("[perf ch=%d] slow fieldUpdate map=%d inst=%d dur=%s players=%d mobs=%d drops=%d", server.id, inst.fieldID, inst.id, d, len(inst.players), len(inst.lifePool.mobs), len(inst.dropPool.drops))
}

func (server *Server) observeBroadcast(kind string, recipients int, maxQueueLen, maxQueueCap int, d time.Duration) {
	if server == nil || server.perf == nil {
		return
	}
	server.perf.observe("broadcast/"+kind, d)
	if maxQueueCap > 0 {
		server.perf.observe(fmt.Sprintf("send_queue_peak/%s/%d_of_%d", kind, maxQueueLen, maxQueueCap), d)
	}
	if recipients > 0 && d >= 10*time.Millisecond {
		log.Printf("[perf ch=%d] slow broadcast kind=%s dur=%s recipients=%d maxSendQueue=%d/%d", server.id, kind, d, recipients, maxQueueLen, maxQueueCap)
	}
}

func (server *Server) ObserveEventLoopWait(kind string, d time.Duration) {
	server.observePerf("eventloop_wait/"+kind, d)
}

func (server *Server) ObserveEventLoopWork(kind string, d time.Duration) {
	server.observePerf("eventloop_work/"+kind, d)
}

func observeSince(server *Server, name string, start time.Time) {
	if start.IsZero() {
		return
	}
	server.observePerf(name, time.Since(start))
}
