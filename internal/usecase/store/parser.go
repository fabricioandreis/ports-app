package store

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/fabricioandreis/ports-app/internal/domain/ports"
)

const channelBufferLength = 100

var ErrParsingInputJSONStream = errors.New("error when parsing input JSON stream")

// A parser parses an input JSON stream of a known file format for Ports.
type parser struct {
	jsonStream io.Reader
}

func newParser(jsonStream io.Reader) *parser {
	return &parser{jsonStream}
}

type result struct {
	port ports.Port
	err  error
}

// parse produces Ports from an input stream and returns a readonly output channel for the consumer to read the results.
// If an error occurs when trying to unmarshal the JSON stream, an error is set to the result into the output channel.
// The method handles context cancellation by writing an error right away into the output channel.
func (p *parser) parse(ctx context.Context) <-chan result {
	results := make(chan result, channelBufferLength)

	go func() {
		defer func() {
			log.Println("Closing parser channel...")
			close(results)
			log.Println("Closed parser channel")
		}()

		iterator := newJSONIterator(p.jsonStream)

		for {
			port, err := iterator.next()

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
	}()

	return results
}

func (p *parser) handleError(err error, results chan<- result) {
	err = errors.Join(ErrParsingInputJSONStream, err)
	log.Println(err.Error())
	results <- result{err: err}
}
