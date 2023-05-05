package store

import (
	"context"
	"io"

	"github.com/fabricioandreis/ports-app/internal/contracts/repository"
	"github.com/fabricioandreis/ports-app/internal/domain"
)

type StoreUsecase struct {
	repoPort repository.Port
}

func NewStoreUsecase(repoPort repository.Port) *StoreUsecase {
	return &StoreUsecase{repoPort}
}

func (usc *StoreUsecase) Store(ctx context.Context, jsonStream io.Reader) (int, error) {
	parser := newParser()
	ports := make(chan domain.Port)
	errs := make(chan error)
	go parser.parseStream(ctx, jsonStream, ports, errs)

	count := 0
	for {
		select {
		case port, ok := <-ports:
			if !ok {
				return count, nil
			}
			usc.repoPort.Put(ctx, port)
			count++
		case err := <-errs:
			return count, err
		}
	}
}
