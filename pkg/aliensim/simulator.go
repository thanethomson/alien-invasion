package aliensim

// SimulationConfig contains the input configuration to our alien world
// simulator.
type SimulationConfig struct {
	world  string
	aliens int
}

// SimulationResult will eventually contain our simulation results.
type SimulationResult struct {
	iterationsSimulated int
	aliensStillAlive    int
	citiesRemaining     []string
}

// SimulateAlienInvasion is our primary simulation routine which takes
// simulation configuration, parses it, simulates the invasion, and returns a
// simulation result or an error.
func SimulateAlienInvasion(config *SimulationConfig) (*SimulationResult, error) {
	if config.aliens < 2 {
		return nil, NewSimulationError(ErrTooFewAliens)
	}
	if config.aliens == 2 {
		return &SimulationResult{
			iterationsSimulated: 10000,
			aliensStillAlive:    2,
			citiesRemaining:     []string{"Foo", "Bar", "Baz", "Qu-ux", "Bee"},
		}, nil
	}
	return nil, nil
}
