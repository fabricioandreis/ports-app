package tests_test

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
	t.Parallel()

	config := config.Load()

	dbClient, err := db.NewClient(config.RedisAddress, config.RedisPassword)
	require.NoError(t, err)

	repo := db.NewPortRepository(dbClient)

	expected := expectedPorts(config.InputJSONFilePath)

	for _, expect := range expected {
		t.Logf("Finding port %s in database", expect.ID)
		found, err := repo.Get(context.Background(), expect.ID)

		assert.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, expect.ID, found.ID)
		assert.Equal(t, expect.Code, found.Code)
		assert.Equal(t, expect.Name, found.Name)
		assert.Equal(t, expect.City, found.City)
		assert.Equal(t, expect.Province, found.Province)
		assert.Equal(t, expect.Country, found.Country)
		assert.Equal(t, expect.Timezone, found.Timezone)
		assert.ElementsMatch(t, expect.Alias, found.Alias)
		assert.ElementsMatch(t, expect.Regions, found.Regions)
		assert.ElementsMatch(t, expect.Unlocs, found.Unlocs)

		if len(expect.Coordinates) == 2 {
			assert.Equal(t, expect.Coordinates[0], found.Coordinates.Lat)
			assert.Equal(t, expect.Coordinates[1], found.Coordinates.Long)
		}

		t.Logf("Found port %s in database and all its fields match the expected values", expect.ID)
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
