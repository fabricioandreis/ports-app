package storing

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

func (p *parser) parseStream(ctx context.Context, jsonStream io.Reader, ports chan<- domain.Port, errs chan<- error) {
	defer close(ports)
	defer close(errs)

	done := make(chan bool)
	go func() {
		select {
		case <-ctx.Done():
			p.handleError(ctx.Err(), errs)
			return
		case <-done:
		}
	}()

	dec := json.NewDecoder(jsonStream)
	_, err := dec.Token() // consume opening curly braces
	if err != nil {
		p.handleError(err, errs)
		return
	}
	for dec.More() {
		// port ID is a different token on each item, therefore we just consume it as a string token
		tokenID, err := dec.Token()
		if err != nil {
			p.handleError(err, errs)
			return
		}
		id, ok := tokenID.(string)
		if !ok {
			p.handleError(ErrCastTokenIDString, errs)
			return
		}
		port := domain.Port{ID: id}

		// consume remaining part of the item and unmarshal it into Port entity
		err = dec.Decode(&port)
		if err != nil {
			p.handleError(err, errs)
			return
		}
		ports <- port
	}
	_, err = dec.Token() // consume closing curly braces
	if err != nil {
		p.handleError(err, errs)
		return
	}
	done <- true // finishes go routine that checks for either completion or context cancellation
}

func (p *parser) handleError(err error, errs chan<- error) {
	log.Printf("[ERROR] " + err.Error())
	errs <- err
}
