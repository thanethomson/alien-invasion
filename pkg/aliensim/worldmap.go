package aliensim

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// WorldMap contains all of the data structures relevant to mapping out our
// simulated world.
type WorldMap struct {
	cityNames   []string         // The names of the cities in the order read from the input.
	cities      map[string]*City // A map of the cities read from the input (key=city name).
	aliens      []*Alien         // A list of our aliens.
	parsedLines uint64           // How many lines of the input have we parsed so far?
}

// NewEmptyWorldMap creates an empty world map, but initialises its structures
// so data can easily be added to it.
func NewEmptyWorldMap() *WorldMap {
	return &WorldMap{
		cityNames:   []string{},
		cities:      map[string]*City{},
		aliens:      []*Alien{},
		parsedLines: 0,
	}
}

// ParseWorldMap takes the given input reader, scans it one line at a time,
// attempts to parse each line, and produces a WorldMap data structure on
// success, or an error on failure. A reader is used to allow for greater memory
// efficiency when supplying larger input files.
func ParseWorldMap(worldReader io.Reader) (*WorldMap, error) {
	scanner := bufio.NewScanner(worldReader)
	worldMap := NewEmptyWorldMap()
	for scanner.Scan() {
		err := worldMap.ParseLine(scanner.Text())
		if err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, NewExtendedSimulationError(ErrFailedToScanWorldInput, "", err)
	}
	// now run through the map to ensure all the neighbours know about each other
	for _, cityName := range worldMap.cityNames {
		worldMap.cities[cityName].recomputeNeighbours()
	}
	return worldMap, nil
}

// ParseLine parses a single line from a world map file. On success, updates the
// world map. On failure, returns an error.
func (m *WorldMap) ParseLine(line string) error {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return NewExtendedSimulationError(
			ErrFailedToParseLine,
			fmt.Sprintf(
				"Too little map data on line %d (must at least be of the form \"CityName direction=OtherCity\").",
				m.parsedLines+1,
			),
			nil,
		)
	}
	cityName := parts[0]
	city, exists := m.cities[cityName]
	if !exists {
		city = NewCity(cityName)
		m.cities[cityName] = city
		m.cityNames = append(m.cityNames, cityName)
	}

	for i := 1; i < len(parts); i++ {
		dirParts := strings.Split(parts[i], "=")
		if len(dirParts) != 2 {
			return NewExtendedSimulationError(
				ErrFailedToParseLine,
				fmt.Sprintf(
					"Invalid city location format on line %d (must be of the form \"direction=CityName\").",
					m.parsedLines+1,
				),
				nil,
			)
		}
		dir, otherCityName := strings.ToLower(dirParts[0]), dirParts[1]
		otherCity, exists := m.cities[otherCityName]
		if !exists {
			otherCity = NewCity(otherCityName)
			m.cities[otherCityName] = otherCity
			m.cityNames = append(m.cityNames, otherCityName)
		}

		// now we situate the other city relative to our current one
		err := city.LocateRelativeTo(otherCity, dir)
		if err != nil {
			return NewExtendedSimulationError(
				ErrFailedToParseLine,
				fmt.Sprintf(
					"Parsing error on line %d.",
					m.parsedLines+1,
				),
				err,
			)
		}
	}
	m.parsedLines++
	return nil
}
