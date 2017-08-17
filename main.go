package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cwiggers/mutex/mutex"
)

func main() {
	if err := mutex.NewClient("127.0.0.1:6379", "", 100); err != nil {
		fmt.Printf("connect redis failed, error: %s\n", err)
		return
	}

	locknumber := 0
	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now().UnixNano()
			l := mutex.Mutex{Key: "lock", Try: 5, TTL: 30 * time.Second, Delay: 200}
			if err := l.TryLock(); err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			defer l.Unlock()
			locknumber++
			fmt.Printf("locknumber: %d\n", locknumber)
			fmt.Println((time.Now().UnixNano() - start) / 1000000)
		}()
	}

	wg.Wait()
	fmt.Printf("locknumber: %d\n", locknumber)
}
