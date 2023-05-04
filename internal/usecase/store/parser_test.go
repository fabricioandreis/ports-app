package storing

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/fabricioandreis/ports-app/internal/domain"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("Should be able to parse a simple JSON file containing ports in a known format", func(t *testing.T) {
		p := newParser()

		tests := []struct {
			input  io.Reader
			output []domain.Port
		}{
			{
				input: strings.NewReader(
					`{"BRPNG":{"name":"Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}, "BRITQ":{"name":"Itaqui","city":"Itaqui","province":"RioGrandedoSul","country":"Brazil","alias":[],"regions":[],"coordinates":[-56.5481122,-29.1294007],"timezone":"America/Sao_Paulo","unlocs":["BRITQ"],"code":"35135"}}`),
				output: []domain.Port{
					{
						ID:          "BRPNG",
						Code:        "35159",
						Name:        "Paranagua",
						City:        "Paranaguá",
						Country:     "Brazil",
						Province:    "Paraná",
						Timezone:    "America/Sao_Paulo",
						Alias:       []string{"br_par_01", "br_par_001"},
						Coordinates: []float32{-48.5, -25.52},
						Regions:     []string{"America", "Latin America"},
						Unlocs:      []string{"BRPNG"},
					},
					{
						ID:          "BRITQ",
						Code:        "35135",
						Name:        "Itaqui",
						City:        "Itaqui",
						Country:     "Brazil",
						Province:    "RioGrandedoSul",
						Timezone:    "America/Sao_Paulo",
						Alias:       []string{},
						Coordinates: []float32{-56.5481122, -29.1294007},
						Regions:     []string{},
						Unlocs:      []string{"BRITQ"},
					},
				},
			},
		}

		for i, data := range tests {
			t.Run(fmt.Sprintf("Test #%v", i), func(t *testing.T) {
				ports := make(chan domain.Port)
				errors := make(chan error)

				go p.parseStream(data.input, ports, errors)

				res := []domain.Port{}
				for p := range ports {
					res = append(res, p)
				}

				if !cmp.Equal(data.output, res) {
					assert.Fail(t, fmt.Sprintf("Ports are not as expected: %s", cmp.Diff(data.output, res)))
				}
				select {
				case err, ok := <-errors:
					if ok {
						assert.Fail(t, fmt.Sprintf("Should not have produced error: %s", err.Error()))
					}
				default:
				}
			})
		}
	})
}
