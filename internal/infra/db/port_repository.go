package db

import (
	"context"
	"log"
	"time"

	"github.com/fabricioandreis/ports-app/internal/contracts/repository"
	"github.com/fabricioandreis/ports-app/internal/domain"
	"github.com/fabricioandreis/ports-app/internal/infra/db/proto"
	redis "github.com/redis/go-redis/v9"
	protobuf "google.golang.org/protobuf/proto"
)

type PortRepository struct {
	client redis.Client
}

func NewPortRepository(address, password string) repository.Port {
	repo := &PortRepository{
		client: *redis.NewClient(&redis.Options{
			Addr:         address,
			Password:     password,
			DB:           0,
			DialTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
		})}
	log.Println("Connected to Redis database")
	return repo
}

func (repo *PortRepository) Get(ctx context.Context, portID string) (*domain.Port, error) {
	dbModel, err := repo.client.Get(ctx, portID).Bytes()
	if err != nil {
		return nil, err
	}

	port, err := repo.dbModelToEntity(dbModel)
	if err != nil {
		return nil, err
	}

	return port, nil
}
func (repo *PortRepository) Put(ctx context.Context, port domain.Port) error {
	dbModel, err := repo.entityToDBModel(port)
	if err != nil {
		return err
	}
	repo.client.Set(ctx, port.ID, dbModel, 0)
	return nil
}

func (repo *PortRepository) entityToDBModel(port domain.Port) ([]byte, error) {
	message := &proto.Port{
		ID:          port.ID,
		Code:        port.Code,
		Name:        port.Name,
		City:        port.City,
		Province:    port.Province,
		Country:     port.Country,
		Timezone:    port.Timezone,
		Alias:       port.Alias,
		Coordinates: port.Coordinates,
		Regions:     port.Regions,
		Unlocs:      port.Unlocs,
	}
	return protobuf.Marshal(message)
}

func (repo *PortRepository) dbModelToEntity(dbModel []byte) (*domain.Port, error) {
	port := proto.Port{}
	err := protobuf.Unmarshal(dbModel, &port)
	if err != nil {
		return nil, err
	}

	return &domain.Port{
		ID:          port.ID,
		Code:        port.Code,
		Name:        port.Name,
		City:        port.City,
		Province:    port.Province,
		Country:     port.Country,
		Timezone:    port.Timezone,
		Alias:       port.Alias,
		Coordinates: port.Coordinates,
		Regions:     port.Regions,
		Unlocs:      port.Unlocs,
	}, nil
}
