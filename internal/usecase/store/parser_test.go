package store

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/fabricioandreis/ports-app/internal/domain"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var (
	portParanagua = domain.Port{
		ID:          "BRPNG",
		Code:        "35159",
		Name:        "Paranagua",
		City:        "Paranaguá",
		Province:    "Paraná",
		Country:     "Brazil",
		Timezone:    "America/Sao_Paulo",
		Alias:       []string{"br_par_01", "br_par_001"},
		Coordinates: []float32{-48.5, -25.52},
		Regions:     []string{"America", "Latin America"},
		Unlocs:      []string{"BRPNG"},
	}
	portItaqui = domain.Port{
		ID:          "BRITQ",
		Code:        "35135",
		Name:        "Itaqui",
		City:        "Itaqui",
		Province:    "RioGrandedoSul",
		Country:     "Brazil",
		Timezone:    "America/Sao_Paulo",
		Alias:       []string{},
		Coordinates: []float32{-56.5481122, -29.1294007},
		Regions:     []string{},
		Unlocs:      []string{"BRITQ"},
	}
)

func TestParser(t *testing.T) {
	t.Run("Should be able to parse a simple JSON input stream containing ports in a known format", func(t *testing.T) {
		input := strings.NewReader(
			`{"BRPNG":{"name":"Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}, "BRITQ":{"name":"Itaqui","city":"Itaqui","province":"RioGrandedoSul","country":"Brazil","alias":[],"regions":[],"coordinates":[-56.5481122,-29.1294007],"timezone":"America/Sao_Paulo","unlocs":["BRITQ"],"code":"35135"}}`)
		output := []domain.Port{portParanagua, portItaqui}

		res, err := parseStream(context.Background(), input)

		assert.NoError(t, err)
		assert.Len(t, res, len(output))
		if !cmp.Equal(output, res) {
			assert.Fail(t, fmt.Sprintf("Ports are not as expected: %s", cmp.Diff(output, res)))
		}
	})

	t.Run("Should return error when input in an invalid JSON stream but also partially process the stream up until the syntax error", func(t *testing.T) {
		tests := []struct {
			input  io.Reader
			output []domain.Port
			errStr string
		}{
			{
				input: strings.NewReader(
					`"BRPNG":{"name":Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}}`),
				output: []domain.Port{},
				errStr: "invalid character ':' looking for beginning of value",
			},
			{
				input:  strings.NewReader(`a`),
				output: []domain.Port{},
				errStr: "invalid character 'a' looking for beginning of value",
			},
			{
				input: strings.NewReader(
					`[{"BRPNG":{"name":Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}}]`),
				output: []domain.Port{},
				errStr: ErrCastTokenIDString.Error(),
			},
			{
				input: strings.NewReader(
					`{"BRPNG":{"name":Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}}`),
				output: []domain.Port{},
				errStr: "invalid character 'P' looking for beginning of value",
			},
			{
				input: strings.NewReader(
					`{"BRPNG":{"name":"Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}`),
				output: []domain.Port{portParanagua},
				errStr: "EOF",
			},
		}

		for i, data := range tests {
			t.Run(fmt.Sprintf("Test #%v", i+1), func(t *testing.T) {
				res, err := parseStream(context.Background(), data.input)

				assert.Len(t, res, len(data.output))
				if !cmp.Equal(data.output, res) {
					assert.Fail(t, fmt.Sprintf("ports are not as expected: %s", cmp.Diff(data.output, res)))
				}
				assert.ErrorContains(t, err, data.errStr)
			})
		}
	})

	t.Run("Should gracefully stop processing when context is cancelled", func(t *testing.T) {
		input := &blockingIOReader{}
		ctx, cancel := context.WithCancel(context.Background())

		chPorts := make(chan []domain.Port)
		chErrs := make(chan error)

		go func() {
			res, err := parseStream(ctx, input)
			chPorts <- res
			chErrs <- err
		}()
		time.Sleep(time.Second) // waits some time for the go routine above to start
		cancel()

		assert.Len(t, <-chPorts, 0)
		assert.ErrorIs(t, <-chErrs, context.Canceled)
	})

}

func parseStream(ctx context.Context, jsonStream io.Reader) ([]domain.Port, error) {
	p := newParser()
	ports := make(chan domain.Port)
	errs := make(chan error)

	go p.parseStream(ctx, jsonStream, ports, errs)
	res := []domain.Port{}
	var err error
loop:
	for {
		select {
		case p, ok := <-ports:
			if ok {
				res = append(res, p)
			}
		case err = <-errs:
			break loop
		}
	}
	return res, err
}

type blockingIOReader struct{}

func (r *blockingIOReader) Read(p []byte) (n int, err error) {
	ch := make(chan bool)
	<-ch // blocks forever
	return 0, nil
}
