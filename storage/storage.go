package storage

import "github.com/go-redis/redis"

type DriverModel interface {
	Get(key string) (string, error)
	Set(key string, data interface{}) (string, error)
}

type CarModel interface {
	Get(key string) (string, error)
	Set(key string, data interface{}) (string, error)
	RemoveFromList(key string, data interface{}) (int64, error)
	InsertIntoList(key string, data interface{}) (int64, error)
}

type StorageInstance struct {
	Driver DriverModel
	Car    CarModel
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type InitConfig struct {
	Driver, Car RedisConfig
}

func Init(config InitConfig) *StorageInstance {

	var driverClient = redis.NewClient(&redis.Options{
		Addr:     config.Driver.Addr,
		Password: config.Driver.Password,
		DB:       config.Driver.DB,
	})

	var carClient = redis.NewClient(&redis.Options{
		Addr:     config.Car.Addr,
		Password: config.Car.Password,
		DB:       config.Car.DB,
	})

	driverInstance := &Driver{db: driverClient}
	carInstance := &Car{db: carClient}

	return &StorageInstance{Driver: driverInstance, Car: carInstance}
}
