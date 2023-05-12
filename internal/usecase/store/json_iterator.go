package store

import (
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/fabricioandreis/ports-app/internal/domain/ports"
)

const numCoordinates = 2

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
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	Timezone    string    `json:"timezone"`
	Alias       []string  `json:"alias"`
	Coordinates []float32 `json:"coordinates"`
	Regions     []string  `json:"regions"`
	Unlocs      []string  `json:"unlocs"`
}

func newJSONIterator(jsonStream io.Reader) jsonIterator {
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
	readPort := jsonPort{}

	err = it.dec.Decode(&readPort)
	if err != nil {
		return nil, err
	}

	port := ports.Port{
		ID:       portID,
		Code:     readPort.Code,
		Name:     readPort.Name,
		City:     readPort.City,
		Province: readPort.Province,
		Country:  readPort.Country,
		Timezone: readPort.Timezone,
		Alias:    readPort.Alias,
		Regions:  readPort.Regions,
		Unlocs:   readPort.Unlocs,
	}

	// converts slice of coordinates into proper domain object
	if len(readPort.Coordinates) == 0 {
		return &port, nil
	}

	if len(readPort.Coordinates) != numCoordinates {
		log.Println("port " + portID + " has invalid coordinates")

		return nil, ErrInvalidCoordinates
	}

	port.Coordinates = ports.Coordinates{Lat: readPort.Coordinates[0], Long: readPort.Coordinates[1]}

	return &port, nil
}
