package aliensim

import (
	"io"
	"strings"
	"testing"
)

const (
	// Example world from the problem statement.
	exampleWorld string = `Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
`
)

type alienCountSimulationErrorPair struct {
	alienCount int
	simError   SimulationErrorCode
}

type alienCountSimulationResult struct {
	alienCount          int
	iterationsSimulated int
	aliensStillAlive    int
	citiesRemaining     []string
}

func newTestSimulationConfig(worldReader io.Reader, aliens int) *SimulationConfig {
	return &SimulationConfig{
		worldReader:     worldReader,
		aliens:          aliens,
		rnd:             NewSequenceGenerator(),
		maxAlienMoves:   10,
		progressHandler: nil, // no need to print out progress during testing
	}
}

// Makes sure that, if fewer than 2 aliens are specified, the ErrTooFewAliens
// error is returned from the simulator.
func TestTooFewAliens(t *testing.T) {
	tests := []alienCountSimulationErrorPair{
		{-1, ErrTooFewAliens},
		{0, ErrTooFewAliens},
		{1, ErrTooFewAliens},
	}
	for _, test := range tests {
		_, err := NewSimulation(
			newTestSimulationConfig(
				strings.NewReader(exampleWorld),
				test.alienCount,
			),
		).Simulate()
		if ferr, ok := err.(*SimulationError); ok {
			if ferr.kind != test.simError {
				t.Error(
					"For N =",
					test.alienCount,
					"expected error to be",
					test.simError,
					"but got",
					ferr,
				)
			}
		} else {
			t.Error(
				"For N =",
				test.alienCount,
				"expected an error but got none",
			)
		}
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Tests a full-blown simulation, with a non-random number generator to ensure
// consistent results.
func TestControlledSimulationOfExampleMap(t *testing.T) {
	tests := []alienCountSimulationResult{
		{2, 10, 2, []string{"Foo", "Bar", "Baz", "Qu-ux", "Bee"}},
		{3, 10, 1, []string{"Bar", "Baz", "Qu-ux", "Bee"}},
		{4, 10, 2, []string{"Foo", "Bar", "Baz", "Qu-ux"}},
	}
	for _, test := range tests {
		res, err := NewSimulation(
			newTestSimulationConfig(
				strings.NewReader(exampleWorld),
				test.alienCount,
			),
		).Simulate()
		if err != nil {
			t.Error(
				"For N =", test.alienCount,
				"and example world, expected no error, but got", err,
			)
		}
		if res.iterationsSimulated != test.iterationsSimulated {
			t.Error(
				"For N =", test.alienCount,
				"and example world, expected iterationsSimulated =", test.iterationsSimulated,
				"but got", res.iterationsSimulated,
			)
		}
		if res.aliensStillAlive != test.aliensStillAlive {
			t.Error(
				"For N =", test.alienCount,
				"and example world, expected aliensStillAlive =", test.aliensStillAlive,
				"but got", res.aliensStillAlive,
			)
		}
		if !stringSlicesEqual(res.citiesRemaining, test.citiesRemaining) {
			t.Error(
				"For N =", test.alienCount,
				"and example world, expected citiesRemaining =", test.citiesRemaining,
				"but got", res.citiesRemaining,
			)
		}
	}
}
