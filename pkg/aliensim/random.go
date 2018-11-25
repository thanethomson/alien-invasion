package aliensim

import (
	"math/rand"
	"time"
)

// RandomGenerator wraps a random number generator to facilitate mocking during
// testing.
type RandomGenerator interface {
	Uint32() uint32
}

// PseudorandomGenerator uses Golang's built-in pseudorandom number generation
// facility.
type PseudorandomGenerator struct {
	r *rand.Rand
}

// SequenceGenerator generates a sequence of monotonically increasing numbers
// instead of random numbers - note that this is purely for testing purposes.
type SequenceGenerator struct {
	next uint32
}

// NewPseudorandomGenerator creates a new pseudorandom number generator and
// automatically seeds it with the current system time.
func NewPseudorandomGenerator() *PseudorandomGenerator {
	return &PseudorandomGenerator{
		r: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// NewSequenceGenerator instantiates a sequence generator starting at 0.
func NewSequenceGenerator() *SequenceGenerator {
	return &SequenceGenerator{next: 0}
}

// Uint32 generates the next number in the pseudorandom number sequence.
func (g *PseudorandomGenerator) Uint32() uint32 {
	return g.r.Uint32()
}

// Uint32 returns the next number in the sequence, incrementing the internal
// sequence counter in the process.
func (s *SequenceGenerator) Uint32() uint32 {
	s.next++
	return s.next - 1
}
