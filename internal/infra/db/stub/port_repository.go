package stub

import (
	"context"
	"errors"
	"sync"

	"github.com/fabricioandreis/ports-app/internal/contracts/repository"
	"github.com/fabricioandreis/ports-app/internal/domain/ports"
)

type PortRepository struct {
	db  sync.Map
	cfg config
}

type config struct {
	errToReturn error
}

type option func(c *config)

func WithError(err error) option {
	return func(c *config) {
		c.errToReturn = err
	}
}

func NewPortRepository(opts ...option) repository.Port {
	repo := &PortRepository{db: sync.Map{}}
	for _, opt := range opts {
		opt(&repo.cfg)
	}
	return repo
}

func (repo *PortRepository) Get(ctx context.Context, portID string) (*ports.Port, error) {
	if repo.cfg.errToReturn != nil {
		return nil, repo.cfg.errToReturn
	}

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
	if repo.cfg.errToReturn != nil {
		return repo.cfg.errToReturn
	}

	repo.db.Store(port.ID, port)
	return nil
}
