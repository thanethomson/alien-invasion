package aliensim

// Alien contains the location and state of a specific alien.
type Alien struct {
	city  *City // The city in which we currently find this alien.
	alive bool  // Is this alien still alive?
}
