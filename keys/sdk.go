package keys

import (
	"regexp"
	"time"
)

const (
	// Copy represents an action that copies a key
	Copy = iota

	// REPLACE represents an action that replace a key
	REPLACE
)

// Keys represents the keys store
type Keys interface {
	Exists(key ...string) int
	Del(key ...string) int
	Rename(key string, toKey string) error
	RenameNX(key string, toKey string) error
	Marshal(key string, ptr interface{}) error
	Expire(key string, duration time.Duration) bool
	ExpireAt(key string, timestamp time.Time) bool
	Keys(pattern *regexp.Regexp) []string
	Migrate(host string, destinationDB string, timeout time.Duration, flag int, keys ...string) error
	Move(key string, db string) bool
	Persist(key string) bool
	TimeToLive(key string) (int, error)
	Restore(key string, ttl time.Duration, serializedValue string, replace bool) bool
}
