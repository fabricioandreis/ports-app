package db_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fabricioandreis/ports-app/internal/contracts/repository"
	"github.com/fabricioandreis/ports-app/internal/domain/ports"
	"github.com/fabricioandreis/ports-app/internal/infra/db"
	"github.com/fabricioandreis/ports-app/internal/infra/db/stub"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var portRioGrande = ports.Port{
	ID:          "BRRIG",
	Code:        "35173",
	Name:        "RioGrande",
	City:        "RioGrande",
	Province:    "RioGrandedoSul",
	Country:     "Brazil",
	Timezone:    "America/Sao_Paulo",
	Alias:       []string{},
	Coordinates: ports.Coordinates{Lat: -52.1075802, Long: -32.0353776},
	Regions:     []string{},
	Unlocs:      []string{"BRRIG"},
}

func TestPortRepository(t *testing.T) {
	repos := []struct {
		repo    func() repository.Port
		enabled bool
	}{
		{
			repo:    func() repository.Port { return stub.NewPortRepository() },
			enabled: true,
		},
		{
			repo:    func() repository.Port { return db.NewPortRepository(db.NewClient("localhost:6379", "")) },
			enabled: false,
		},
	}

	for i, data := range repos {
		if !data.enabled {
			continue
		}

		repo := data.repo()
		t.Run(fmt.Sprintf("Repo %v: Should be able to store port and retrieve it", i+1), func(t *testing.T) {
			err1 := repo.Put(context.Background(), portRioGrande)
			found, err2 := repo.Get(context.Background(), portRioGrande.ID)

			assert.NoError(t, err1)
			assert.NoError(t, err2)
			require.NotNil(t, found)
			if !cmp.Equal(portRioGrande, *found, cmpopts.EquateEmpty()) {
				assert.Fail(t, fmt.Sprintf("port not as expected: %s", cmp.Diff(portRioGrande, *found)))
			}
		})
	}
}
