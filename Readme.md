
# go-stats

 Go stats reporting / ticker utility.

 View the [docs](http://godoc.org/pkg/github.com/segmentio/go-stats).

## Example

```go
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

  s.TickEvery(5 * time.Second)
  defer s.Stop()

  time.Sleep(time.Minute)
}
```

 Outputs the ops/s, since last tick, and since the beginning of time.

```
stats 2014/07/18 11:24:27 messages 94.00/s tick=470 total=470
stats 2014/07/18 11:24:27 errors 18.80/s tick=94 total=94
stats 2014/07/18 11:24:32 messages 96.98/s tick=485 total=955
stats 2014/07/18 11:24:32 errors 19.40/s tick=97 total=191
stats 2014/07/18 11:24:37 messages 98.99/s tick=495 total=1450
stats 2014/07/18 11:24:37 errors 19.80/s tick=99 total=290
...
```

## License

MIT
