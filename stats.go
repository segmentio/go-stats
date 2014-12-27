package stats

import "github.com/dustin/go-humanize"
import "time"
import "sync"
import "log"
import "os"

// printfer interface.
type printfer interface {
	Printf(string, ...interface{})
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
	s.TickEveryTo(d, log.New(os.Stderr, "", log.LstdFlags))
}

// TickEveryTo `d` to the given Printf-er.
func (s *Stats) TickEveryTo(d time.Duration, p printfer) {
	s.tick = time.NewTicker(d)
	for _ = range s.tick.C {
		s.Write(p)
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

// Reset statistics.
func (s *Stats) Reset() {
	s.Lock()
	defer s.Unlock()

	for k := range s.m {
		s.m[k] = 0
	}

	s.lastReset = time.Now()
}

// Write to the given printer.
func (s *Stats) Write(p printfer) {
	s.Lock()

	defer s.Reset()
	defer s.Unlock()

	secs := time.Since(s.lastReset).Seconds()

	for k, v := range s.m {
		total := humanize.Comma(s.t[k])
		p.Printf("stats: %s %.2f/s (%s)\n", k, float64(v)/secs, total)
	}
}
