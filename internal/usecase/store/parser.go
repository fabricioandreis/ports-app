package store

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/fabricioandreis/ports-app/internal/domain"
)

var ErrCastTokenIDString = errors.New("unable to cast token ID to string")

// A parser parses an input JSON stream of a known file format for Ports.
type parser struct {
	jsonStream io.Reader
}

type result struct {
	port domain.Port
	err  error
}

func newParser(jsonStream io.Reader) *parser {
	return &parser{jsonStream}
}

// parse produces Ports from an input stream into an output channel.
// If an error occurs when trying to unmarshal the JSON stream, an error is sent to another output channel.
// The method handles context cancellation by writing an error right away into the output channel.
func (p *parser) parse(ctx context.Context, results chan<- result) {
	defer func() {
		log.Println("Closing parser channels")
		close(results)
		log.Println("Closed parser channels")
	}()

	iterator := newJsonIterator(p.jsonStream)
	for {
		port, err := iterator.next(ctx)

		select {
		case <-ctx.Done():
			p.handleError(ctx.Err(), results)
			log.Println("Context cancelled, finishing parser...")
			return
		default:
			if err != nil {
				p.handleError(err, results)
				return
			}
			if port == nil {
				log.Println("Parsed input JSON stream")
				return
			}
			log.Println("Read port " + port.ID)
			results <- result{port: *port}
		}
	}
}

func (p *parser) handleError(err error, results chan<- result) {
	err = errors.Join(errors.New("error when parsing input JSON stream"), err)
	log.Println(err.Error())
	results <- result{err: err}
}
