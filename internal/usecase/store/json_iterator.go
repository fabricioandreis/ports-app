package store

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/fabricioandreis/ports-app/internal/domain"
)

var ErrCastTokenIDString = errors.New("unable to cast token ID to string")

// A jsonIterator returns the next port in the input JSON stream.
// When it finished reading, it returns nil.
// If the input JSON stream is not valid, it returns an error.
type jsonIterator struct {
	started bool
	dec     *json.Decoder
}

func newJsonIterator(jsonStream io.Reader) jsonIterator {
	return jsonIterator{
		started: false,
		dec:     json.NewDecoder(jsonStream),
	}
}

func (it *jsonIterator) next() (*domain.Port, error) {
	if !it.started {
		_, err := it.dec.Token() // read opening curly bracket
		if err != nil {
			return nil, err
		}
		it.started = true
	}
	if !it.dec.More() {
		_, err := it.dec.Token() // read closing curly bracket
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	// port ID is a different token on each item, therefore we just read it as a string token
	tokenID, err := it.dec.Token()
	if err != nil {
		return nil, err
	}
	id, ok := tokenID.(string)
	if !ok {
		return nil, ErrCastTokenIDString
	}
	port := domain.Port{ID: id}

	// read remaining part of the item and unmarshal it into Port entity
	err = it.dec.Decode(&port)
	if err != nil {
		return nil, err
	}

	return &port, nil
}
