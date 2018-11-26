package aliensim

import (
	"fmt"
	"strings"
	"testing"
)

const (
	contradictoryTestMap string = `Foo north=Bar east=Baz south=Qu-ux
Baz west=Bar
`
)

func TestSuccessfullyParsingExampleWorldMap(t *testing.T) {
	m, err := ParseWorldMap(strings.NewReader(ExampleWorld))
	if err != nil {
		t.Error("Parsing failed with error:", err)
	} else {
		cityLinks := map[string][]string{
			"Foo":   []string{"Bar", "", "Qu-ux", "Baz"},
			"Bar":   []string{"", "", "Foo", "Bee"},
			"Baz":   []string{"Bee", "Foo", "", ""},
			"Qu-ux": []string{"Foo", "", "", ""},
			"Bee":   []string{"", "Bar", "Baz", ""},
		}

		for cityName, expectedNeighbours := range cityLinks {
			// first make sure the city exists in the map
			city, exists := m.cities[cityName]
			if !exists {
				t.Error(
					"Expected city",
					cityName,
					"to be present, but was not",
				)
			}

			// now check that the city has the given neighbours
			for dir, expectedNeighbour := range expectedNeighbours {
				neighbour := city.neighbours[dir]
				if expectedNeighbour == "" {
					if neighbour != nil {
						t.Error(
							fmt.Sprintf(
								"City %s is supposed to have no neighbour to the %s, but has neighbour %s",
								cityName,
								directionNames[dir],
								neighbour.name,
							),
						)
					}
				} else {
					if neighbour == nil {
						t.Error(
							fmt.Sprintf(
								"City %s is supposed to have neighbour %s to the %s, but has none",
								cityName,
								expectedNeighbour,
								directionNames[dir],
							),
						)
					} else if neighbour.name != expectedNeighbour {
						t.Error(
							fmt.Sprintf(
								"City %s is supposed to have neighbour %s to the %s, but has %s",
								cityName,
								expectedNeighbour,
								directionNames[dir],
								neighbour.name,
							),
						)
					}
				}
			}
		}
	}
}

func TestParsingContradictoryMap(t *testing.T) {
	_, err := ParseWorldMap(strings.NewReader(contradictoryTestMap))
	// we should get an error
	if err == nil {
		t.Error("Expected contradictory test to result in an error")
	}
	serr, _ := err.(*SimulationError)
	if serr.kind != ErrFailedToParseLine {
		t.Error(
			"Expected error kind to be", ErrFailedToParseLine,
			"but got", serr.kind,
		)
	}
	if serr.upstream == nil {
		t.Error("Expected upstream error to be non-nil, but was nil")
	}
	uerr, _ := serr.upstream.(*SimulationError)
	if uerr.kind != ErrCityAlreadyThere {
		t.Error(
			"Expected upstream error kind to be", ErrCityAlreadyThere,
			"but was", uerr.kind,
		)
	}
}
