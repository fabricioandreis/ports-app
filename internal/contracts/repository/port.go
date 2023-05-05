package repository

import (
	"context"

	"github.com/fabricioandreis/ports-app/internal/domain"
)

type Port interface {
	Get(ctx context.Context, portID string) (*domain.Port, error)
	Put(ctx context.Context, port domain.Port) error
}
