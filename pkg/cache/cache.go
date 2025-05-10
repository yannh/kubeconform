package cache

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, schema interface{}) error
}
