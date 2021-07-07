package cache

import "time"

// TODO: Implement Redis
type Redis struct {
}

func NewRedis() *Redis {
	return &Redis{}
}

func (c *Redis) Get(key string) string {
	panic("implement me")
}

func (c *Redis) Set(key string, data string, expiryTime time.Duration) {
	panic("implement me")
}
