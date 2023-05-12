package store

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/fabricioandreis/ports-app/internal/domain/ports"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

var (
	portParanagua = ports.Port{
		ID:          "BRPNG",
		Code:        "35159",
		Name:        "Paranagua",
		City:        "Paranaguá",
		Province:    "Paraná",
		Country:     "Brazil",
		Timezone:    "America/Sao_Paulo",
		Alias:       []string{"br_par_01", "br_par_001"},
		Coordinates: ports.Coordinates{Lat: -48.5, Long: -25.52},
		Regions:     []string{"America", "Latin America"},
		Unlocs:      []string{"BRPNG"},
	}
	portItaqui = ports.Port{
		ID:          "BRITQ",
		Code:        "35135",
		Name:        "Itaqui",
		City:        "Itaqui",
		Province:    "RioGrandedoSul",
		Country:     "Brazil",
		Timezone:    "America/Sao_Paulo",
		Alias:       []string{},
		Coordinates: ports.Coordinates{Lat: -56.5481122, Long: -29.1294007},
		Regions:     []string{},
		Unlocs:      []string{"BRITQ"},
	}
	portNoCoord = ports.Port{
		ID:          "PORTID",
		Coordinates: ports.Coordinates{},
	}
)

func TestParser(t *testing.T) {
	t.Run("Should be able to parse a simple JSON input stream containing ports in a known format", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			input  io.Reader
			output []ports.Port
		}{
			{
				input: strings.NewReader(
					`{"BRPNG":{"name":"Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}, "BRITQ":{"name":"Itaqui","city":"Itaqui","province":"RioGrandedoSul","country":"Brazil","alias":[],"regions":[],"coordinates":[-56.5481122,-29.1294007],"timezone":"America/Sao_Paulo","unlocs":["BRITQ"],"code":"35135"}}`),
				output: []ports.Port{portParanagua, portItaqui},
			},
			{
				input:  strings.NewReader(`{"PORTID":{"coordinates":[]}}`),
				output: []ports.Port{portNoCoord},
			},
			{
				input:  strings.NewReader(`{}`),
				output: []ports.Port{},
			},
		}

		for _, data := range tests {
			data := data
			t.Run("", func(t *testing.T) {
				t.Parallel()

				res, err := parseStream(context.Background(), data.input)

				assert.NoError(t, err)
				assert.Len(t, res, len(data.output))
				if !cmp.Equal(data.output, res) {
					assert.Fail(t, fmt.Sprintf("Ports are not as expected: %s", cmp.Diff(data.output, res)))
				}
			})
		}
	})

	t.Run("Should return error when input in an invalid JSON stream but also partially process the stream up until the syntax error", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			input  io.Reader
			output []ports.Port
			errStr string
		}{
			{
				input: strings.NewReader(
					`"BRPNG":{"name":Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}}`),
				output: []ports.Port{},
				errStr: "invalid character ':' looking for beginning of value",
			},
			{
				input:  strings.NewReader(`a`),
				output: []ports.Port{},
				errStr: "invalid character 'a' looking for beginning of value",
			},
			{
				input: strings.NewReader(
					`[{"BRPNG":{"name":Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}}]`),
				output: []ports.Port{},
				errStr: ErrCastTokenIDString.Error(),
			},
			{
				input: strings.NewReader(
					`{"BRPNG":{"name":Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}}`),
				output: []ports.Port{},
				errStr: "invalid character 'P' looking for beginning of value",
			},
			{
				input: strings.NewReader(
					`{"BRPNG":{"name":"Paranagua","coordinates":[-48.5,-25.52],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}`),
				output: []ports.Port{portParanagua},
				errStr: "EOF",
			},
			{
				input: strings.NewReader(
					`{"BRPNG":{"name":"Paranagua","coordinates":[-48.5],"city":"Paranaguá","province":"Paraná","country":"Brazil","alias":["br_par_01", "br_par_001"],"regions":["America", "Latin America"],"timezone":"America/Sao_Paulo","unlocs":["BRPNG"],"code":"35159"}}`),
				output: []ports.Port{},
				errStr: ErrInvalidCoordinates.Error(),
			},
		}

		for _, data := range tests {
			data := data
			t.Run("", func(t *testing.T) {
				t.Parallel()
				res, err := parseStream(context.Background(), data.input)

				assert.Len(t, res, len(data.output))
				if !cmp.Equal(data.output, res) {
					assert.Fail(t, fmt.Sprintf("ports are not as expected: %s", cmp.Diff(data.output, res)))
				}
				assert.ErrorContains(t, err, data.errStr)
			})
		}
	})
}

func TestStopProcessing(t *testing.T) {
	t.Run("Should gracefully stop processing when context is cancelled", func(t *testing.T) {
		t.Parallel()
		input := &slowReader{}
		ctx, cancel := context.WithCancel(context.Background())

		chPorts := make(chan []ports.Port)
		chErrs := make(chan error)

		go func() {
			res, err := parseStream(ctx, input)
			chPorts <- res
			chErrs <- err
		}()
		time.Sleep(200 * time.Millisecond) // waits some time for the go routine above to start
		cancel()

		assert.Len(t, <-chPorts, 0)
		assert.ErrorIs(t, <-chErrs, context.Canceled)
	})

}

func parseStream(ctx context.Context, jsonStream io.Reader) ([]ports.Port, error) {
	p := newParser(jsonStream)

	results := p.parse(ctx)
	ports := []ports.Port{}
	for res := range results {
		if res.err != nil {
			return ports, res.err
		}
		ports = append(ports, res.port)
	}
	return ports, nil
}

type slowReader struct{}

func (r *slowReader) Read(p []byte) (n int, err error) {
	time.Sleep(time.Second)
	return 0, io.EOF
}
