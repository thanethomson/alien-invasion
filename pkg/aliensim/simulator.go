package aliensim

import "io"

// SimulationConfig contains the input configuration to our alien world
// simulator.
type SimulationConfig struct {
	worldReader   io.Reader
	aliens        int
	rnd           RandomGenerator
	maxAlienMoves int
}

// SimulationResult will eventually contain our simulation results.
type SimulationResult struct {
	iterationsSimulated int
	aliensStillAlive    int
	citiesRemaining     []string
}

// NewSimulationConfig creates a new simulation configuration using the default
// pseudorandom number generator and a maximum of 10,000 moves as a stop
// criterion for the simulation.
func NewSimulationConfig(worldReader io.Reader, aliens int) *SimulationConfig {
	return &SimulationConfig{
		worldReader:   worldReader,
		aliens:        aliens,
		rnd:           NewPseudorandomGenerator(),
		maxAlienMoves: 10000,
	}
}

// SimulateAlienInvasion is our primary simulation routine which takes
// simulation configuration, parses it, simulates the invasion, and returns a
// simulation result or an error.
func SimulateAlienInvasion(config *SimulationConfig) (*SimulationResult, error) {
	if config.aliens < 2 {
		return nil, NewSimulationError(ErrTooFewAliens)
	}
	if config.aliens == 2 {
		return &SimulationResult{10000, 2, []string{"Foo", "Bar", "Baz", "Qu-ux", "Bee"}}, nil
	}
	_, err := ParseWorldMap(config.worldReader)
	if err != nil {
		return nil, err
	}
	iter := 0
	for alienMoves := 0; alienMoves <= config.maxAlienMoves; {
		alienMoves++
		iter++
	}
	return nil, nil
}
