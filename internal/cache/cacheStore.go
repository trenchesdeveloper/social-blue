package cache

type CacheStore interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
}