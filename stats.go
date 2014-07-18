package stats

import "time"
import "sync"
import "log"
import "os"

// Printfer interface.
type Printfer interface {
	Printf(string, ...interface{})
}

// Stats struct.
type Stats struct {
	t         map[string]int64
	m         map[string]int64
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

func (s *Stats) TickEvery(d time.Duration) {
	s.TickEveryTo(d, log.New(os.Stderr, "stats ", log.LstdFlags))
}

func (s *Stats) TickEveryTo(d time.Duration, p Printfer) {
	for {
		time.Sleep(d)
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
func (s *Stats) Write(p Printfer) {
	s.Lock()

	defer s.Reset()
	defer s.Unlock()

	secs := time.Since(s.lastReset).Seconds()

	for k, v := range s.m {
		p.Printf("%s %.2f/s tick=%d total=%d\n", k, float64(v)/secs, v, s.t[k])
	}
}
