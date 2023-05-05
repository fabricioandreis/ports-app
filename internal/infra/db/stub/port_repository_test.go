package stub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fabricioandreis/ports-app/internal/domain"
	"github.com/fabricioandreis/ports-app/internal/infra/db/stub"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var portRioGrande = domain.Port{
	ID:          "BRRIG",
	Code:        "35173",
	Name:        "RioGrande",
	City:        "RioGrande",
	Province:    "RioGrandedoSul",
	Country:     "Brazil",
	Timezone:    "America/Sao_Paulo",
	Alias:       []string{},
	Coordinates: []float32{-52.1075802, -32.0353776},
	Regions:     []string{},
	Unlocs:      []string{"BRRIG"},
}

func TestPortRepository(t *testing.T) {
	t.Run("Should be able to store port and retrieve it", func(t *testing.T) {
		repo := stub.NewPortRepository()

		err1 := repo.Put(context.Background(), portRioGrande)
		found, err2 := repo.Get(context.Background(), portRioGrande.ID)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		if !cmp.Equal(portRioGrande, *found) {
			assert.Fail(t, fmt.Sprintf("port not as expected: %s", cmp.Diff(portRioGrande, *found)))
		}
	})
}
