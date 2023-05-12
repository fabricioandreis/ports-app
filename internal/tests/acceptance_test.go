package tests

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/fabricioandreis/ports-app/internal/infra/config"
	"github.com/fabricioandreis/ports-app/internal/infra/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type port struct {
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

func TestAcceptance(t *testing.T) {
	config := config.Load()

	dbClient := db.NewClient(config.RedisAddress, config.RedisPassword)
	repo := db.NewPortRepository(dbClient)

	expected := expectedPorts(config.InputJSONFilePath)

	for _, e := range expected {
		t.Logf("Finding port %s in database", e.ID)
		found, err := repo.Get(context.Background(), e.ID)

		assert.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, e.ID, found.ID)
		assert.Equal(t, e.Code, found.Code)
		assert.Equal(t, e.Name, found.Name)
		assert.Equal(t, e.City, found.City)
		assert.Equal(t, e.Province, found.Province)
		assert.Equal(t, e.Country, found.Country)
		assert.Equal(t, e.Timezone, found.Timezone)
		assert.ElementsMatch(t, e.Alias, found.Alias)
		assert.ElementsMatch(t, e.Regions, found.Regions)
		assert.ElementsMatch(t, e.Unlocs, found.Unlocs)
		if len(e.Coordinates) == 2 {
			assert.Equal(t, e.Coordinates[0], found.Coordinates.Lat)
			assert.Equal(t, e.Coordinates[01], found.Coordinates.Long)
		}

		t.Logf("Found port %s in database and all its fields match the expected values", e.ID)
	}
}

func expectedPorts(filepath string) []port {
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("unable to read file '%s': %s", filepath, err.Error())
	}

	expected := []port{}
	err = json.Unmarshal(content, &expected)
	if err != nil {
		log.Fatalf("unable to unmarshal file '%s': %s", filepath, err.Error())
	}

	return expected
}
