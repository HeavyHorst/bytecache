package bytecache

type Cache interface {
	Set(key string, value []byte, ttl int) error
	Get(string) ([]byte, error)
	Delete(string) error
}
