package aliensim

import "fmt"

// Alien contains the location and state of a specific alien.
type Alien struct {
	id    int   // The ID of this alien.
	city  *City // The city in which we currently find this alien.
	alive bool  // Is this alien still alive?
}

// NewAlien creates a living alien in the given city.
func NewAlien(id int, city *City) *Alien {
	return &Alien{
		id:    id,
		city:  city,
		alive: true,
	}
}

func (a *Alien) String() string {
	return fmt.Sprintf("Alien %d in %s (alive=%t)", a.id, a.city.name, a.alive)
}

// MoveInRandomDirection will attempt to move this alien in a random direction,
// depending on which cities around it are not yet destroyed. If it cannot move,
// this function will return false.
func (a *Alien) MoveInRandomDirection(rnd RandomGenerator) bool {
	availCities := []*City{}
	for _, n := range a.city.neighbours {
		if n != nil && !n.destroyed {
			availCities = append(availCities, n)
		}
	}
	if len(availCities) > 0 {
		cityID := rnd.Uint32() % uint32(len(availCities))
		// move the alien to this new city
		a.city = availCities[cityID]
		return true
	}
	return false
}
