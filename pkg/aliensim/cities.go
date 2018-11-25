package aliensim

import "fmt"

// City direction indices.
const (
	DirNorth int = 0
	DirEast  int = 1
	DirSouth int = 2
	DirWest  int = 3
)

var (
	mapDirections = map[string]int{
		"north": DirNorth,
		"east":  DirEast,
		"south": DirSouth,
		"west":  DirWest,
	}
	mapDirectionOpposites = map[int]int{
		DirNorth: DirSouth,
		DirEast:  DirWest,
		DirSouth: DirNorth,
		DirWest:  DirEast,
	}
	directionNames = map[int]string{
		DirNorth: "North",
		DirEast:  "East",
		DirSouth: "South",
		DirWest:  "West",
	}
)

// City represents a single city on the world map. Effectively implements a
// multidimensional linked list to allow for map traversal in a relatively
// memory-efficient manner.
type City struct {
	name       string  // The name of this city, as read from the input map.
	destroyed  bool    // Has the city been destroyed yet?
	neighbours []*City // An indexed list of neighbours (0=North, 1=East, 2=South, 3=West).
	x, y       int     // The calculated coordinates of this city on the map.
}

// NewCity creates a fresh new city, not yet destroyed, with no neighbours.
func NewCity(name string) *City {
	return &City{
		name:       name,
		destroyed:  false,
		neighbours: []*City{nil, nil, nil, nil},
		x:          0,
		y:          0,
	}
}

func (c *City) String() string {
	neighbours := "nil"
	if c.neighbours != nil {
		neighbours = ""
		for i, n := range c.neighbours {
			nname := "nil"
			if n != nil {
				nname = n.name
			}
			if i > 0 {
				neighbours = fmt.Sprintf("%s, ", neighbours)
			}
			neighbours = fmt.Sprintf("%s%s: %s", neighbours, directionNames[i], nname)
		}
	}
	return fmt.Sprintf(
		"City{name: %s, destroyed: %t, neighbours: {%s}, x: %d, y: %d}",
		c.name,
		c.destroyed,
		neighbours,
		c.x,
		c.y,
	)
}

// LocateRelativeTo will ensure that the given other city is located to the
// (dir) of this city. If the direction is unrecognised, returns an error.
func (c *City) LocateRelativeTo(otherCity *City, dir string) error {
	d, ok := mapDirections[dir]
	if !ok {
		return NewExtendedSimulationError(
			ErrUnknownDirection,
			fmt.Sprintf(
				"Unknown direction %s (must be one of north, east, south or west).",
				dir,
			),
			nil,
		)
	}
	dopp := mapDirectionOpposites[d]

	existingNeighbour := c.neighbours[d]
	if existingNeighbour != nil {
		// if there's already another city in that direction by a different name
		if existingNeighbour.name != otherCity.name {
			return NewExtendedSimulationError(
				ErrCityAlreadyThere,
				fmt.Sprintf(
					"There is already a city (%s) to the %s of %s.",
					existingNeighbour.name,
					directionNames[d],
					c.name,
				),
				nil,
			)
		}
	} else {
		// put the other city here
		c.neighbours[d] = otherCity
	}

	existingNeighbour = otherCity.neighbours[dopp]
	if existingNeighbour != nil {
		if existingNeighbour.name != c.name {
			return NewExtendedSimulationError(
				ErrCityAlreadyThere,
				fmt.Sprintf(
					"There is already a city (%s) to the %s of %s.",
					existingNeighbour.name,
					directionNames[dopp],
					otherCity.name,
				),
				nil,
			)
		}
	} else {
		// locate this city in the opposite direction relative to the other city
		otherCity.neighbours[dopp] = c
	}

	return nil
}

// Depending on how the cities' relative locations are specified in the input
// file, there are times when, for example, a city to the North of our current
// city isn't aware of a city to its East, but the city to the East of our
// current city is aware of that city to its North. This refreshes the
// neighbouring cities' locations.
func (c *City) recomputeNeighbours() {
	north := c.neighbours[DirNorth]
	east := c.neighbours[DirEast]
	south := c.neighbours[DirSouth]
	west := c.neighbours[DirWest]

	// the neighbours at the corners
	var nw, ne, se, sw *City = nil, nil, nil, nil

	if north != nil {
		if north.neighbours[DirWest] != nil {
			nw = north.neighbours[DirWest]
		}
		if north.neighbours[DirEast] != nil {
			ne = north.neighbours[DirEast]
		}
	}
	if east != nil {
		if east.neighbours[DirNorth] != nil {
			ne = east.neighbours[DirNorth]
		}
		if east.neighbours[DirSouth] != nil {
			se = east.neighbours[DirSouth]
		}
	}
	if south != nil {
		if south.neighbours[DirWest] != nil {
			sw = south.neighbours[DirWest]
		}
		if south.neighbours[DirEast] != nil {
			se = south.neighbours[DirEast]
		}
	}
	if west != nil {
		if west.neighbours[DirNorth] != nil {
			nw = west.neighbours[DirNorth]
		}
		if west.neighbours[DirSouth] != nil {
			sw = west.neighbours[DirSouth]
		}
	}

	if north != nil {
		north.neighbours[DirWest] = nw
		north.neighbours[DirEast] = ne
	}
	if east != nil {
		east.neighbours[DirNorth] = ne
		east.neighbours[DirSouth] = se
	}
	if south != nil {
		south.neighbours[DirWest] = sw
		south.neighbours[DirEast] = se
	}
	if west != nil {
		west.neighbours[DirNorth] = nw
		west.neighbours[DirSouth] = sw
	}
}
