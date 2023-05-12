package db

import (
	"context"
	"errors"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const redisTimeout = 200 * time.Millisecond

var (
	ErrCreatingClient     = errors.New("unable to create Redis client")
	ErrConnectingDatabase = errors.New("unable to create Redis client")
)

type Client struct {
	*redis.Client
}

func NewClient(address, password string) (Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
		DB:           0,
		DialTimeout:  redisTimeout,
		ReadTimeout:  redisTimeout,
		WriteTimeout: redisTimeout,
	})
	if redisClient == nil {
		log.Println("error creating Redis client at " + address)

		return Client{}, ErrCreatingClient
	}

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Println("error connecting to Redis database at " + address)

		return Client{}, ErrConnectingDatabase
	}

	return Client{redisClient}, nil
}

func (c Client) Close() {
	log.Println("Closing Redis client")

	if c.Client != nil {
		err := c.Client.Close()
		if err != nil {
			log.Println("unable to close Redis client: " + err.Error())
		}
	}

	log.Println("Closed Redis client")
}
