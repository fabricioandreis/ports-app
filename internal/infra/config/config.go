package config

import (
	"log"
	"os"
)

type Config struct {
	RedisAddress  string
	RedisPassword string
}

func Load() Config {
	return Config{
		RedisAddress:  ifEmpty(getEnv("REDIS_ADDRESS"), "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Println("Empty value for " + key)
	}
	return value
}

func ifEmpty(str1, str2 string) string {
	if len(str1) == 0 {
		return str2
	}
	return str1
}
