package lists

// Lists represents the lists store
type Lists interface {
	Insert(key string, isBefore bool, pivot []byte, value []byte) int
	Index(key string, index int) []byte
	Pop(key string) byte
	Push(key string, values ...[]byte) int
	PushX(key string, value []byte) int
	Range(key string, start int, stop int) [][]byte
	Rem(key string, count int, value []byte) int
	Set(key string, index int, value []byte) bool
	Trim(key string, start int, stop int) error
	RPop(key string) byte
	RPush(key string, values ...[]byte) int
	RPushX(key string, value []byte) int
}
