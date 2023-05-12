package db

import (
	"context"
	"errors"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

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
		DialTimeout:  200 * time.Millisecond,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
	})
	if redisClient == nil {
		return Client{}, errors.Join(ErrCreatingClient, errors.New("error creating Redis client at "+address))
	}
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return Client{}, errors.Join(ErrConnectingDatabase, errors.New("error connecting to Redis database at "+address))
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
