package aliensim

// SimulationErrorCode is an enum to represent error codes for the simulator
type SimulationErrorCode uint32

// The possible error codes that can be generated by the simulator
const (
	ErrTooFewAliens SimulationErrorCode = 0
)

// SimulationError is returned when some aspect of our simulation fails
type SimulationError struct {
	kind SimulationErrorCode
}

func NewSimulationError(kind SimulationErrorCode) error {
	return &SimulationError{kind: kind}
}

// Error translates the simulation error enum into its human-readable form
func (e *SimulationError) Error() string {
	switch e.kind {
	case ErrTooFewAliens:
		return "Please specify at least N=2 aliens for a meaningful simulation"
	}
	return "Unrecognised error code"
}
