package store

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/fabricioandreis/ports-app/internal/domain/ports"
)

var (
	ErrCastTokenIDString  = errors.New("unable to cast token ID to string")
	ErrInvalidCoordinates = errors.New("invalid coordinates in port")
)

// A jsonIterator returns the next port in the input JSON stream.
// When it finished reading, it returns nil.
// If the input JSON stream is not valid, it returns an error.
type jsonIterator struct {
	started bool
	dec     *json.Decoder
}

type jsonPort struct {
	ID          string
	Code        string
	Name        string
	City        string
	Province    string
	Country     string
	Timezone    string
	Alias       []string
	Coordinates []float32
	Regions     []string
	Unlocs      []string
}

func newJsonIterator(jsonStream io.Reader) jsonIterator {
	return jsonIterator{
		started: false,
		dec:     json.NewDecoder(jsonStream),
	}
}

func (it *jsonIterator) next() (*ports.Port, error) {
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

	return it.decode()
}

func (it *jsonIterator) decode() (*ports.Port, error) {
	// port ID is a different token on each item, therefore we just read it as a string token
	tokenID, err := it.dec.Token()
	if err != nil {
		return nil, err
	}
	portID, ok := tokenID.(string)
	if !ok {
		return nil, ErrCastTokenIDString
	}

	// read remaining part of the item and unmarshal it into Port entity
	jp := jsonPort{}
	err = it.dec.Decode(&jp)
	if err != nil {
		return nil, err
	}
	port := ports.Port{
		ID:       portID,
		Code:     jp.Code,
		Name:     jp.Name,
		City:     jp.City,
		Province: jp.Province,
		Country:  jp.Country,
		Timezone: jp.Timezone,
		Alias:    jp.Alias,
		Regions:  jp.Regions,
		Unlocs:   jp.Unlocs,
	}

	// converts slice of coordinates into proper domain object
	if len(jp.Coordinates) == 0 {
		return &port, nil
	}
	if len(jp.Coordinates) != 2 {
		return nil, errors.Join(ErrInvalidCoordinates, errors.New("port "+portID+" has invalid coordinates"))
	}
	port.Coordinates = ports.Coordinates{Lat: jp.Coordinates[0], Long: jp.Coordinates[1]}

	return &port, nil
}
