package stats

import "github.com/dustin/go-humanize"
import "time"
import "sync"
import "sort"
import "log"
import "os"

// LogFunc interface.
type LogFunc func(fmt string, v ...interface{})

// stat struct.
type stat struct {
	name  string
	value int64
}

// Stats struct.
type Stats struct {
	t         map[string]int64
	m         map[string]int64
	tick      *time.Ticker
	lastReset time.Time
	sync.Mutex
}

// New stats reporter.
func New() *Stats {
	return &Stats{
		t:         make(map[string]int64),
		m:         make(map[string]int64),
		lastReset: time.Now(),
	}
}

// Stop ticker.
func (s *Stats) Stop() {
	s.tick.Stop()
}

// TickEvery `d` to stderr via the std log package.
func (s *Stats) TickEvery(d time.Duration) {
	logger := log.New(os.Stderr, "stats ", log.LstdFlags)
	s.TickEveryTo(d, logger.Printf)
}

// TickEveryTo `d` to the given logger.
func (s *Stats) TickEveryTo(d time.Duration, log LogFunc) {
	s.tick = time.NewTicker(d)
	for _ = range s.tick.C {
		s.Write(log)
	}
}

// Incr increments the stat `name`.
func (s *Stats) Incr(name string) {
	s.IncrBy(name, 1)
}

// IncrBy increments the stat `name` by `n`.
func (s *Stats) IncrBy(name string, n int64) {
	s.Lock()
	defer s.Unlock()
	s.t[name] += n
	s.m[name] += n
}

// Get the value of `name` or 0.
func (s *Stats) Get(name string) int64 {
	s.Lock()
	defer s.Unlock()
	return s.m[name]
}

// GetTotal the value of `name` or 0.
func (s *Stats) GetTotal(name string) int64 {
	s.Lock()
	defer s.Unlock()
	return s.t[name]
}

// Slice of stats.
func (s *Stats) slice() (ret []*stat) {
	for k, v := range s.m {
		ret = append(ret, &stat{k, v})
	}
	return
}

// Reset statistics.
func (s *Stats) Reset() {
	s.Lock()
	defer s.Unlock()

	for k := range s.m {
		s.m[k] = 0
	}

	s.lastReset = time.Now()
}

// Write to the given logger.
func (s *Stats) Write(log LogFunc) {
	s.Lock()

	defer s.Reset()
	defer s.Unlock()

	if len(s.m) == 0 {
		return
	}

	stats := s.slice()
	sort.Sort(byName(stats))

	secs := time.Since(s.lastReset).Seconds()

	log("")

	for _, stat := range stats {
		total := humanize.Comma(s.t[stat.name])
		log("%s %.2f/s (%s)", stat.name, float64(stat.value)/secs, total)
	}

	log("")
}
