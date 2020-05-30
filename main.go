package main

import (
	"log"

	"github.com/electra-systems/athena/services"
	"github.com/electra-systems/athena/storage"
	"github.com/joho/godotenv"
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

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}

	services.Init(storageInstance)
}
