package storing

import (
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/fabricioandreis/ports-app/internal/domain"
)

type parser struct {
}

func newParser() *parser {
	return &parser{}
}

func (p *parser) parseStream(jsonStream io.Reader, ports chan<- domain.Port, errs chan<- error) {
	defer close(ports)
	defer close(errs)

	dec := json.NewDecoder(jsonStream)
	_, err := dec.Token() // Consume opening curly braces
	if err != nil {
		p.handleError(err, errs)
		return
	}
	for dec.More() {
		// Port ID is a different token on each item, therefore we just consume it as a string token
		tokenID, err := dec.Token()
		if err != nil {
			p.handleError(err, errs)
			return
		}
		id, ok := tokenID.(string)
		if !ok {
			p.handleError(errors.New("unable to cast token ID to string"), errs)
			return
		}
		port := domain.Port{ID: id}

		// Consume remaining part of the item and unmarshal it into Port entity
		err = dec.Decode(&port)
		if err != nil {
			p.handleError(err, errs)
			return
		}
		ports <- port
	}
	_, err = dec.Token() // Consume closing curly braces
	if err != nil {
		p.handleError(err, errs)
		return
	}
}

func (p *parser) handleError(err error, errs chan<- error) {
	log.Fatal(err)
	errs <- err
}
