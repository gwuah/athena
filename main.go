package main

import (
	"github.com/electra-systems/athena/services"

	"github.com/electra-systems/athena/storage"
)

func main() {
	var storageInstance = *storage.Init(storage.InitConfig{Driver: storage.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	}, Car: storage.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	}})

	services.Init(storageInstance)

}
