package stub

import (
	"context"
	"errors"
	"sync"

	"github.com/fabricioandreis/ports-app/internal/contracts/repository"
	"github.com/fabricioandreis/ports-app/internal/domain/ports"
)

type PortRepository struct {
	db sync.Map
}

func NewPortRepository() repository.Port {
	return &PortRepository{db: sync.Map{}}
}

func (repo *PortRepository) Get(ctx context.Context, portID string) (*ports.Port, error) {
	res, ok := repo.db.Load(portID)
	if !ok {
		return nil, nil
	}
	port, ok := res.(ports.Port)
	if !ok {
		return nil, errors.New("unable to cast result as port")
	}
	return &port, nil

}
func (repo *PortRepository) Put(ctx context.Context, port ports.Port) error {
	repo.db.Store(port.ID, port)
	return nil
}
