package storage

import "github.com/go-redis/redis"

type Driver struct {
	db *redis.Client
}

func (instance *Driver) Get(key string) (string, error) {
	result, err := instance.db.Get(key).Result()
	return result, err
}

func (instance *Driver) Set(key string, data interface{}) (string, error) {
	result, err := instance.db.Set(key, data, 0).Result()
	return result, err
}

func (instance *Driver) MGet(keys []string) ([]interface{}, error) {
	result, err := instance.db.MGet(keys...).Result()
	return result, err
}
