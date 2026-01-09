package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var (
		base   = flag.String("base", "http://localhost:8080", "base url")
		path   = flag.String("path", "/flights", "path to hit")
		rps    = flag.Int("rps", 50, "requests per second")
		secs   = flag.Int("secs", 10, "duration seconds")
		concur = flag.Int("c", 20, "max concurrency")
	)
	flag.Parse()

	client := &http.Client{Timeout: 10 * time.Second}

	var ok, bad int64
	var inFlight int64

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*secs)*time.Second)
	defer cancel()

	tick := time.NewTicker(time.Second / time.Duration(*rps))
	defer tick.Stop()

	sem := make(chan struct{}, *concur)
	var wg sync.WaitGroup

	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			dur := time.Since(start)
			fmt.Printf("done in %s\nok=%d bad=%d\n", dur, ok, bad)
			return
		case <-tick.C:
			sem <- struct{}{}
			wg.Add(1)
			atomic.AddInt64(&inFlight, 1)

			go func() {
				defer wg.Done()
				defer func() { <-sem; atomic.AddInt64(&inFlight, -1) }()

				req, _ := http.NewRequestWithContext(ctx, "GET", *base+*path, nil)
				resp, err := client.Do(req)
				if err != nil {
					atomic.AddInt64(&bad, 1)
					return
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()

				if resp.StatusCode >= 200 && resp.StatusCode < 400 {
					atomic.AddInt64(&ok, 1)
				} else {
					atomic.AddInt64(&bad, 1)
				}
			}()
		}
	}
}
