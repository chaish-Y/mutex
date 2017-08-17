package mutex

import (
	"errors"
	"math/rand"
	"time"
)

const (
	DefaultExpiry = 30 * time.Second
	DefaultTries  = 3
	DefaultDelay  = 3000
)

var (
	Locked = errors.New("Locked")
)

type Mutex struct {
	Key   string
	Try   int
	Delay int
	TTL   time.Duration
}

var (
	Random *rand.Rand
)

func init() {
	Random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (m *Mutex) TryLock() error {
	var (
		i   int
		err error
	)

	for i = 0; i < m.Try; i++ {
		val := time.Now().UnixNano() / 1000000
		res, _ := Client.SetNX(m.Key, val, m.TTL).Result()
		if res == false {
			delay := time.Duration(Random.Intn(m.Delay)) * time.Millisecond
			time.Sleep(delay)
			continue
		}
		return nil
	}

	ttl, err := Client.TTL(m.Key).Result()
	if err != nil {
		return err
	}
	if ttl == -1 {
		// unset expire
		if err := Client.Expire(m.Key, m.TTL).Err(); err != nil {
			return err
		}
	} else if ttl == -2 {
		// timeout
		val := time.Now().UnixNano() / 1000000
		res, _ := Client.SetNX(m.Key, val, m.TTL).Result()
		if res == true {
			return nil
		}
	}
	return Locked
}

func (m *Mutex) Lock() error {
	for {
		val := time.Now().UnixNano() / 1000000
		res, _ := Client.SetNX(m.Key, val, m.TTL).Result()
		if res == false {
			delay := time.Duration(Random.Intn(m.Delay)) * time.Millisecond
			time.Sleep(delay)
			continue
		} else {
			return nil
		}
	}

	return nil
}

func (m *Mutex) Unlock() error {
	return Client.Del(m.Key).Err()
}
