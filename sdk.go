package redismint

import "github.com/XMNBlockchain/redismint/hashtree"

// Keystore represents the keystore
type Keystore interface {
	GetHead() hashtree.HashTree
	GetHashes() Hashes
}

// Hashes represents the hashes keystore
type Hashes interface {
	// Exists:
	Exists(key string, field string) bool

	// Get:
	Get(key string, field string) []byte
	GetAll(key string) []string
	MultiGet(key string, fields ...string) [][]byte

	// Set:
	Set(key string, field string, value []byte) bool
	SetNX(key string, field string, value []byte) bool
	MultiSet(key string, keyValues ...map[string][]byte)

	// Increment:
	IncrBy(key string, field string, increment int64) int64
	IncrByFloat(key string, field string, increment float64) float64

	// Length:
	Len(key string) int64
	StrLen(key string, field string) int

	// Misc:
	Del(key string, fields ...string) int
	Keys(key string) []string
	Vals(key string) []byte
}
