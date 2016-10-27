package main

import "github.com/segmentio/go-stats"
import "github.com/segmentio/go-log"
import "math/rand"
import "time"

func main() {
	s := stats.New()

	go func() {
		for {
			log.Info("doing stuff")
			s.IncrBy("messages", 5)
			s.Incr("errors")
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		}
	}()

	s.TickEvery(5 * time.Second)
	time.Sleep(30 * time.Second)
	s.Stop()
}
