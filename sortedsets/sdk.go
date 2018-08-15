package sortedsets

// ScoredValue represents a scored value
type ScoredValue struct {
	Value []byte
	Score int
}

// SortedSets represents the sotted sets store
type SortedSets interface {
	Add(key string, flag int, memberScores map[byte]int) int
	Card(key string) int
	Count(key string, min int, max int) int
	IncrBy(key string, increment int, member []byte) int
	Range(key string, start int, stop int) ScoredValue
	InterStore(destination string, keysWeights map[string]int, weightFlag int) int
	UnionStore(destination string, keysWeights map[string]int, weightFlag int) int
	Rem(key string, members ...string) int
}
