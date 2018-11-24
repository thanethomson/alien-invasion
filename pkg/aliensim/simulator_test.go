package aliensim

import "testing"

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

// TestTooFewAliens makes sure that, if fewer than 2 aliens are specified, the
// ErrTooFewAliens error is returned from the simulator.
func TestTooFewAliens(t *testing.T) {
	tests := []alienCountSimulationErrorPair{
		{-1, ErrTooFewAliens},
		{0, ErrTooFewAliens},
		{1, ErrTooFewAliens},
	}
	for _, test := range tests {
		_, err := SimulateAlienInvasion(&SimulationConfig{
			world:  exampleWorld,
			aliens: test.alienCount,
		})
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
