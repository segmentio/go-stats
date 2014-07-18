package main

import "github.com/segmentio/go-stats"
import "time"

func main() {
	s := stats.New()

	go func() {
		for {
			s.IncrBy("messages", 5)
			s.Incr("errors")
			time.Sleep(50 * time.Millisecond)
		}
	}()

	go s.TickEvery(5 * time.Second)

	time.Sleep(time.Minute)
	s.Stop()
}
