package db

import (
	"context"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewClient(address, password string) Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
		DB:           0,
		DialTimeout:  200 * time.Millisecond,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
	})
	if redisClient == nil {
		log.Fatalln("error making Redis client at " + address)
	}
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("unable to connect to Redis database at " + address)
	}
	return Client{redisClient}
}

func (c Client) Close() {
	log.Println("Closing Redis client")
	err := c.Client.Close()
	if err != nil {
		log.Println("unable to close Redis client: " + err.Error())
	}
	log.Println("Closed Redis client")
}
