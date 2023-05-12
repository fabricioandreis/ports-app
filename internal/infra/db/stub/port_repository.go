package stub

import (
	"context"
	"errors"
	"sync"

	"github.com/fabricioandreis/ports-app/internal/domain/ports"
)

var ErrCastResult = errors.New("unable to cast result as port")

type PortRepository struct {
	db  sync.Map
	cfg config
}

type config struct {
	errToReturn error
}

type Option func(repo *PortRepository)

func WithError(err error) Option {
	return func(repo *PortRepository) {
		repo.cfg.errToReturn = err
	}
}

func NewPortRepository(opts ...Option) *PortRepository {
	repo := &PortRepository{db: sync.Map{}}
	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

func (repo *PortRepository) Get(_ context.Context, portID string) (*ports.Port, error) {
	if repo.cfg.errToReturn != nil {
		return nil, repo.cfg.errToReturn
	}

	res, found := repo.db.Load(portID)
	if !found {
		return nil, nil
	}

	port, found := res.(ports.Port)
	if !found {
		return nil, ErrCastResult
	}

	return &port, nil
}

func (repo *PortRepository) Put(_ context.Context, port ports.Port) error {
	if repo.cfg.errToReturn != nil {
		return repo.cfg.errToReturn
	}

	repo.db.Store(port.ID, port)

	return nil
}
