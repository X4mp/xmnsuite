package hashes

// Hashes represents the hashes keystore
type Hashes interface {
	// Exists:
	Exists(key string, field string) bool

	// Get:
	Get(key string, field string) []byte
	GetAll(key string) map[string][]byte
	MultiGet(key string, fields ...string) map[string][]byte

	// Set:
	Set(key string, field string, value []byte) bool
	SetNX(key string, field string, value []byte) bool
	MultiSet(key string, keyValues ...map[string][]byte)

	// Increment:
	IncrBy(key string, field string, increment int64) int64
	IncrByFloat(key string, field string, increment float64) (float64, error)

	// Length:
	Len(key string) int64
	StrLen(key string, field string) int

	// Misc:
	Del(key string, fields ...string) int
	Keys(key string) []string
	Vals(key string) []byte
}

// SDKFunc represents the hash SDK func
var SDKFunc = struct {
	Create func() Hashes
}{
	Create: func() Hashes {
		out := createHashes()
		return out
	},
}
