package cache

import "time"

type Cache interface {
	Get(string) string
	Set(string, string, time.Duration)
}
