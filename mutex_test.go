package mutex

import (
	"sync"
	"testing"
	"time"
)

func Test_Mutex(t *testing.T) {
	if err := NewClient("127.0.0.1:6379", "", 100); err != nil {
		t.Errorf("connect redis failed, error: %s\n", err)
		return
	}

	locknumber := 0
	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l := Mutex{Key: "lock", Try: 5, TTL: 30 * time.Second, Delay: 200}
			if err := l.TryLock(); err != nil {
				t.Errorf("error: %s\n", err)
				return
			}
			defer l.Unlock()
			locknumber++
		}()
	}

	wg.Wait()
	t.Logf("locknumber: %d\n", locknumber)
}
