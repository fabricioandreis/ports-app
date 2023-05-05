package db

import (
	"context"
	"log"

	"github.com/fabricioandreis/ports-app/internal/contracts/repository"
	"github.com/fabricioandreis/ports-app/internal/domain"
	"github.com/fabricioandreis/ports-app/internal/infra/db/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type PortRepository struct {
	client Client
}

func NewPortRepository(client Client) repository.Port {
	return &PortRepository{client}
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
	log.Println("Saved port " + port.ID)
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
