package sortedlists

// ValueScore represents a value with a score
type ValueScore struct {
	Value []byte
	Score int
}

// SortedLists represents the lists data store
type SortedLists interface {
	Add(key string, values ...*ValueScore) int
	Retrieve(key string, index int, amount int) []*ValueScore
	Len(key string) int
	Union(key ...string) []*ValueScore
	UnionStore(destination string, key ...string) int
	Inter(key ...string) []*ValueScore
	InterStore(destination string, key ...string) int
	Delete(key ...string) int
	Push(key string, values ...*ValueScore) int
	PushX(key string, values ...*ValueScore) int
	Pop(key string) []byte
	Trim(key string, start int, stop int) error
}
