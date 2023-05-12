package store_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/fabricioandreis/ports-app/internal/domain/ports"
	"github.com/fabricioandreis/ports-app/internal/infra/db/stub"
	"github.com/fabricioandreis/ports-app/internal/usecase/store"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	portRioGrande = ports.Port{
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
)

func TestStoreUsecase(t *testing.T) {
	t.Run("Should be able to store Ports from input JSON stream into the database", func(t *testing.T) {
		repoPort := stub.NewPortRepository()
		jsonStream := strings.NewReader(`{"BRRIG":{"name":"RioGrande","city":"RioGrande","province":"RioGrandedoSul","country":"Brazil","alias":[],"regions":[],"coordinates":[-52.1075802,-32.0353776],"timezone":"America/Sao_Paulo","unlocs":["BRRIG"],"code":"35173"}}`)
		ctx := context.Background()

		usecase := store.NewStoreUsecase(repoPort)
		output, err1 := usecase.Store(ctx, jsonStream)
		found, err2 := repoPort.Get(ctx, "BRRIG")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 1, output)
		require.NotNil(t, found)

		if !cmp.Equal(portRioGrande, *found) {
			assert.Fail(t, fmt.Sprintf("port not as expected: %s", cmp.Diff(portRioGrande, *found)))
		}
	})

	t.Run("Should handle error from repository", func(t *testing.T) {
		expectedErr := errors.New("unable to store port into repository")
		repoPort := stub.NewPortRepository(stub.WithError(expectedErr))
		jsonStream := strings.NewReader(`{"BRRIG":{"name":"RioGrande","city":"RioGrande","province":"RioGrandedoSul","country":"Brazil","alias":[],"regions":[],"coordinates":[-52.1075802,-32.0353776],"timezone":"America/Sao_Paulo","unlocs":["BRRIG"],"code":"35173"}}`)
		ctx := context.Background()

		usecase := store.NewStoreUsecase(repoPort)
		output, err := usecase.Store(ctx, jsonStream)

		assert.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, 0, output)
	})
}
