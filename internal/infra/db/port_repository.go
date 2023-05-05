package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fabricioandreis/ports-app/internal/contracts/repository"
	"github.com/fabricioandreis/ports-app/internal/domain"
	redis "github.com/redis/go-redis/v9"
)

type PortRepository struct {
	client redis.Client
}

func NewPortRepository(address, password string) repository.Port {
	return &PortRepository{
		client: *redis.NewClient(&redis.Options{
			Addr:         address,
			Password:     password,
			DB:           0,
			DialTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
		}),
	}
}

func (repo *PortRepository) Get(ctx context.Context, portID string) (*domain.Port, error) {
	bytes, err := repo.client.Get(ctx, portID).Bytes()
	if err != nil {
		return nil, err
	}

	var port domain.Port
	err = json.Unmarshal(bytes, &port)
	if err != nil {
		return nil, err
	}

	return &port, nil
}
func (repo *PortRepository) Put(ctx context.Context, port domain.Port) error {
	bytes, err := json.Marshal(port)
	if err != nil {
		return err
	}
	repo.client.Set(ctx, port.ID, bytes, 0)
	return nil
}
