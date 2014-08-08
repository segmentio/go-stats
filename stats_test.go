package stats

import "github.com/bmizerany/assert"
import "sync/atomic"
import "testing"

type Statistics struct {
	Something int64
}

func TestIncr(t *testing.T) {
	stats := New()
	stats.Incr("messages")
	stats.Incr("messages")
	assert.Equal(t, int64(2), stats.Get("messages"))
}

func TestIncrBy(t *testing.T) {
	stats := New()
	stats.IncrBy("messages", 5)
	stats.IncrBy("messages", 10)
	assert.Equal(t, int64(15), stats.Get("messages"))
}

func TestGet(t *testing.T) {
	stats := New()

	stats.IncrBy("messages", 5)
	stats.IncrBy("messages", 10)
	assert.Equal(t, int64(15), stats.Get("messages"))

	stats.Reset()
	stats.IncrBy("messages", 5)
	stats.IncrBy("messages", 10)
	assert.Equal(t, int64(15), stats.Get("messages"))
	assert.Equal(t, int64(30), stats.GetTotal("messages"))
}

func BenchmarkIncr(b *testing.B) {
	stats := New()
	for i := 0; i < b.N; i++ {
		stats.Incr("something")
	}
}

func BenchmarkAtomic(b *testing.B) {
	var s Statistics
	for i := 0; i < b.N; i++ {
		atomic.AddInt64(&s.Something, 1)
	}
}
