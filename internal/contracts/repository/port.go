package repository

import (
	"context"

	"github.com/fabricioandreis/ports-app/internal/domain/ports"
)

type Port interface {
	Get(ctx context.Context, portID string) (*ports.Port, error)
	Put(ctx context.Context, port ports.Port) error
}
