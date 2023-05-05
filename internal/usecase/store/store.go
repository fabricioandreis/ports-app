package store

import (
	"context"
	"io"
	"log"

	"github.com/fabricioandreis/ports-app/internal/contracts/repository"
)

type StoreUsecase struct {
	repoPort repository.Port
}

func NewStoreUsecase(repoPort repository.Port) *StoreUsecase {
	return &StoreUsecase{repoPort}
}

func (usc *StoreUsecase) Store(ctx context.Context, jsonStream io.Reader) (int, error) {
	parser := newParser(jsonStream)
	results := make(chan result)
	go parser.parseStream(ctx, results)

	count := 0
	for {
		select {
		case res, ok := <-results:
			if !ok {
				log.Println("Store use case finished")
				return count, res.err
			}
			usc.repoPort.Put(ctx, res.port)
			count++
		case <-ctx.Done():
			// Busy wait after context is Done until parser.parseStream closes its channels
			log.Println("Busy waiting for parseStream to close its channels")
			_, resultsOpen := <-results
			if !resultsOpen {
				log.Println("parseStream closed its channels, returning")
				return count, ctx.Err()
			}
		}
	}
}
