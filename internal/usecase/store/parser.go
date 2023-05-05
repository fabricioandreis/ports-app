package store

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/fabricioandreis/ports-app/internal/domain"
)

var ErrCastTokenIDString = errors.New("unable to cast token ID to string")

type parser struct {
}

func newParser() *parser {
	return &parser{}
}

// parseStream parses an input JSON stream of a known file format for ports.
// Ports are sent to an output channel.
// If an error occurs when trying to unmarshal the JSON stream, an error is sent to another output channel.
// The method handles context cancellation by writing an error right away into the output channel.
func (p *parser) parseStream(ctx context.Context, jsonStream io.Reader, ports chan<- domain.Port, errs chan<- error) {
	defer close(ports)
	defer close(errs)

	done := make(chan bool)
	cancelled := make(chan bool)
	go func() {
		select {
		case <-ctx.Done():
			p.handleError(ctx.Err(), errs)
			cancelled <- true
			return
		case <-done:
		}
	}()

	nextPort := newIterator(jsonStream)
	for {
		port, err := nextPort()
		if err != nil {
			p.handleError(err, errs)
			return
		}
		if port == nil {
			done <- true // finishes go routine that checks for either completion or context cancellation
			return
		}
		select {
		case <-cancelled:
			return
		default:
			log.Println("Read port " + port.ID)
			ports <- *port
		}
	}
}

func (p *parser) handleError(err error, errs chan<- error) {
	err = errors.Join(errors.New("unable to parse input JSON stream"), err)
	log.Println(err.Error())
	errs <- err
}

// newIterator creates an iterator function.
// The iterator function returns the next port in the input JSON stream.
// When it finished reading, the iterator function returns nil.
// If the input JSON stream is not valid, the iterator function returns an error.
func newIterator(jsonStream io.Reader) func() (*domain.Port, error) {
	dec := json.NewDecoder(jsonStream)
	_, err := dec.Token() // read opening curly bracket
	if err != nil {
		return func() (*domain.Port, error) {
			return nil, err
		}
	}

	return func() (*domain.Port, error) {
		if !dec.More() {
			_, err = dec.Token() // read closing curly bracket
			if err != nil {
				return nil, err
			}
			return nil, nil
		}

		// port ID is a different token on each item, therefore we just consume it as a string token
		tokenID, err := dec.Token()
		if err != nil {
			return nil, err
		}
		id, ok := tokenID.(string)
		if !ok {
			return nil, ErrCastTokenIDString
		}
		port := domain.Port{ID: id}

		// read remaining part of the item and unmarshal it into Port entity
		err = dec.Decode(&port)
		if err != nil {
			return nil, err
		}

		return &port, nil
	}
}
