package mutex

import (
	redis "gopkg.in/redis.v5"
)

var (
	Client *redis.Client
)

func NewClient(address, passwd string, poolSize int) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: passwd,
		PoolSize: poolSize,
	})

	if _, err := Client.Ping().Result(); err != nil {
		return err
	}

	return nil
}
