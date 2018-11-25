package aliensim

import (
	"fmt"
	"io"
	"strings"

	"github.com/deckarep/golang-set"
)

// SimulationIterationsHardLimit is a hard limit on the number of iterations we
// can run through in any given simulation. This is not the same as the limit on
// the maximum number of alien moves.
const SimulationIterationsHardLimit = 20000

// SimulationConfig contains the input configuration to our alien world
// simulator.
type SimulationConfig struct {
	worldReader     io.Reader
	aliens          int
	rnd             RandomGenerator
	maxAlienMoves   int
	progressHandler SimulationProgressHandler
}

// SimulationResult will eventually contain our simulation results.
type SimulationResult struct {
	iterationsSimulated int
	aliensStillAlive    int
	citiesRemaining     []string
}

// SimulationProgressHandler is a simple interface to handle the various
// different events as they are emitted by the simulation process.
type SimulationProgressHandler interface {
	CityDestroyed(cityName string, alienIDs []int)
	AllAliensTrapped()
}

type NoopSimulationProgressHandler struct{}
type StdoutSimulationProgressHandler struct{}

// Simulation is our primary structure through which we execute our simulation.
type Simulation struct {
	config          *SimulationConfig
	worldMap        *WorldMap
	aliens          []*Alien
	citiesDestroyed mapset.Set
}

// NewSimulationConfig creates a new simulation configuration using the default
// pseudorandom number generator and a maximum of 10,000 moves as a stop
// criterion for the simulation.
func NewSimulationConfig(worldReader io.Reader, aliens int) *SimulationConfig {
	return &SimulationConfig{
		worldReader:     worldReader,
		aliens:          aliens,
		rnd:             NewPseudorandomGenerator(),
		maxAlienMoves:   10000,
		progressHandler: &StdoutSimulationProgressHandler{},
	}
}

// NewSimulation constructs a new simulation from the given configuration.
func NewSimulation(config *SimulationConfig) *Simulation {
	return &Simulation{
		config:          config,
		worldMap:        nil,
		aliens:          []*Alien{},
		citiesDestroyed: mapset.NewSet(),
	}
}

// Simulate is our primary simulation routine which takes simulation
// configuration, parses it, simulates the invasion, and returns a simulation
// result or an error.
func (s *Simulation) Simulate() (*SimulationResult, error) {
	if s.config.aliens < 2 {
		return nil, NewSimulationError(ErrTooFewAliens)
	}
	worldMap, err := ParseWorldMap(s.config.worldReader)
	if err != nil {
		return nil, err
	}
	// keep track of it
	s.worldMap = worldMap
	// randomly place our aliens on the map
	s.scatterAliens()

	alienMoves := 0
	iter := 0
	for ; alienMoves < s.config.maxAlienMoves && iter < SimulationIterationsHardLimit; iter++ {
		m, destroyed := s.RunSimulationIteration()
		for cityName := range destroyed {
			s.citiesDestroyed.Add(cityName)
		}
		// if all aliens are trapped
		if m == 0 {
			if s.config.progressHandler != nil {
				s.config.progressHandler.AllAliensTrapped()
			}
			break
		}
		alienMoves += m
	}

	// count how many aliens are still alive
	aliensStillAlive := 0
	for _, alien := range s.aliens {
		if alien.alive {
			aliensStillAlive++
		}
	}

	// compute which cities are still left standing
	citiesRemaining := []string{}
	for _, cityName := range s.worldMap.cityNames {
		if !s.citiesDestroyed.Contains(cityName) {
			citiesRemaining = append(citiesRemaining, cityName)
		}
	}

	return &SimulationResult{
		iterationsSimulated: iter,
		aliensStillAlive:    aliensStillAlive,
		citiesRemaining:     citiesRemaining,
	}, nil
}

func (s *Simulation) scatterAliens() {
	cityCount := len(s.worldMap.cityNames)
	for n := 0; n < s.config.aliens; n++ {
		cityID := s.config.rnd.Uint32() % uint32(cityCount)
		s.aliens = append(
			s.aliens,
			NewAlien(s.worldMap.cities[s.worldMap.cityNames[cityID]]),
		)
	}
}

// RunSimulationIteration runs a single iteration of our simulation, returning
// the maximum number of moves that any single alien made during this iteration.
// The second return parameter is a mapping of city names to a list of aliens
// responsible for each city's destruction.
func (s *Simulation) RunSimulationIteration() (int, map[string][]int) {
	destroyed := map[string][]int{}
	aliensInCities := map[string][]int{}
	alienMoves := 0
	for alienID, alien := range s.aliens {
		if alien.alive {
			cityName := alien.city.name
			aliensInCities[cityName] = append(aliensInCities[cityName], alienID)

			// now move this alien
			couldMove := alien.MoveInRandomDirection(s.config.rnd)
			if couldMove {
				alienMoves = 1
			}
		}
	}

	for cityName, alienIDs := range aliensInCities {
		// Have two or more aliens ended up in a particular city? If so, they'll
		// destroy each other and the city.
		if len(alienIDs) > 1 {
			destroyed[cityName] = alienIDs
			for _, id := range alienIDs {
				s.aliens[id].alive = false
			}
			s.worldMap.cities[cityName].destroyed = true
			if s.config.progressHandler != nil {
				s.config.progressHandler.CityDestroyed(cityName, alienIDs)
			}
		}
	}

	return alienMoves, destroyed
}

func (h *NoopSimulationProgressHandler) CityDestroyed(cityName string, alienIDs []int) {}
func (h *NoopSimulationProgressHandler) AllAliensTrapped()                             {}

// CityDestroyed prints out the fact that a city has been destroyed to Stdout.
func (h *StdoutSimulationProgressHandler) CityDestroyed(cityName string, alienIDs []int) {
	idStrings := []string{}
	for _, id := range alienIDs {
		idStrings = append(idStrings, fmt.Sprintf("alien %d", id))
	}
	fmt.Println(
		fmt.Sprintf(
			"%s has been destroyed by %s!",
			cityName,
			strings.Join(idStrings, " and "),
		),
	)
}

func (h *StdoutSimulationProgressHandler) AllAliensTrapped() {
	fmt.Println("All aliens have been trapped! Simulation will be ended here.")
}
