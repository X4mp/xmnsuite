package sets

// Sets represents the sets store
type Sets interface {
	Add(key string, members ...[]byte) int
	Card(key string) int
	Diff(key ...string) [][]byte
	DiffStore(destination string, keys ...string) int
	Inter(keys ...string) [][]byte
	InterStore(destination string, keys ...string) int
	IsMember(key string, member []byte) bool
	Members(key string) [][]byte
	Move(source string, destination string, member []byte) bool
	Pop(key string, count int) [][]byte
	RandMember(key string) []byte
	Rem(key string, members ...[]byte) int
	Union(keys ...string) [][]byte
	UnionStore(destination string, keys ...string) int
}
