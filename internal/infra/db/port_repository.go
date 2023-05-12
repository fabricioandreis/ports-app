package db

import (
	"context"
	"log"

	"github.com/fabricioandreis/ports-app/internal/domain/ports"
	"github.com/fabricioandreis/ports-app/internal/infra/db/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type PortRepository struct {
	client Client
}

func NewPortRepository(client Client) *PortRepository {
	return &PortRepository{client}
}

func (repo *PortRepository) Get(ctx context.Context, portID string) (*ports.Port, error) {
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

func (repo *PortRepository) Put(ctx context.Context, port ports.Port) error {
	dbModel, err := repo.entityToDBModel(port)
	if err != nil {
		return err
	}

	repo.client.Set(ctx, port.ID, dbModel, 0)
	log.Println("Saved port " + port.ID)

	return nil
}

func (repo *PortRepository) entityToDBModel(port ports.Port) ([]byte, error) {
	message := &proto.Port{
		ID:       port.ID,
		Code:     port.Code,
		Name:     port.Name,
		City:     port.City,
		Province: port.Province,
		Country:  port.Country,
		Timezone: port.Timezone,
		Alias:    port.Alias,
		Coordinates: &proto.Coordinates{
			Latitude:  port.Coordinates.Lat,
			Longitude: port.Coordinates.Long,
		},
		Regions: port.Regions,
		Unlocs:  port.Unlocs,
	}

	return protobuf.Marshal(message)
}

func (repo *PortRepository) dbModelToEntity(dbModel []byte) (*ports.Port, error) {
	port := proto.Port{}

	err := protobuf.Unmarshal(dbModel, &port)
	if err != nil {
		return nil, err
	}

	return &ports.Port{
		ID:       port.ID,
		Code:     port.Code,
		Name:     port.Name,
		City:     port.City,
		Province: port.Province,
		Country:  port.Country,
		Timezone: port.Timezone,
		Alias:    port.Alias,
		Coordinates: ports.Coordinates{
			Lat:  port.Coordinates.Latitude,
			Long: port.Coordinates.Longitude,
		},
		Regions: port.Regions,
		Unlocs:  port.Unlocs,
	}, nil
}
