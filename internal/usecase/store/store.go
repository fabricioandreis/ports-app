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

// Store reads an input JSON stream, convert its contents do Port entities and saves them into the repository
func (usc *StoreUsecase) Store(ctx context.Context, jsonStream io.Reader) (int, error) {
	log.Println("Storing data from input JSON stream into database")

	parser := newParser(jsonStream)
	results := parser.parse(ctx)

	count := 0
	for {
		select {
		case res, ok := <-results:
			if !ok || res.err != nil {
				return count, res.err
			}
			usc.repoPort.Put(ctx, res.port)
			count++
		case <-ctx.Done():
			// Busy wait after context is Done until parser.parseStream closes its channels
			log.Println("Busy waiting for store use case to process all results in channel...")
			_, resultsOpen := <-results
			if !resultsOpen {
				log.Println("Store use case processed results, finishing store use case...")
				return count, ctx.Err()
			}
		}
	}
}
