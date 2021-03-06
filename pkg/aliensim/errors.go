package aliensim

import "fmt"

// SimulationErrorCode is an enum to represent error codes for the simulator
type SimulationErrorCode uint32

// The possible error codes that can be generated by the simulator
const (
	ErrTooFewAliens           SimulationErrorCode = 0
	ErrFailedToScanWorldInput SimulationErrorCode = 1
	ErrFailedToParseLine      SimulationErrorCode = 2
	ErrUnknownDirection       SimulationErrorCode = 3
	ErrCityAlreadyThere       SimulationErrorCode = 4
)

// SimulationError is returned when some aspect of our simulation fails
type SimulationError struct {
	kind     SimulationErrorCode
	detail   string // Additional detail related to the error message, if any.
	upstream error  // The upstream error that triggered this particular error, if any.
}

func NewSimulationError(kind SimulationErrorCode) error {
	return &SimulationError{kind: kind, detail: "", upstream: nil}
}

func NewExtendedSimulationError(kind SimulationErrorCode, detail string, upstream error) error {
	return &SimulationError{kind: kind, detail: detail, upstream: upstream}
}

func (e *SimulationError) buildErrorMessage(prefix string) string {
	msg := fmt.Sprintf("%s", prefix)
	if len(e.detail) > 0 {
		msg = fmt.Sprintf("%s %s", msg, e.detail)
	}
	if e.upstream != nil {
		msg = fmt.Sprintf("%s\nCaused by: %s", msg, e.upstream.Error())
	}
	return msg
}

// Error translates the simulation error enum into its human-readable form
func (e *SimulationError) Error() string {
	switch e.kind {
	case ErrTooFewAliens:
		return e.buildErrorMessage("Please specify at least N=2 aliens for a meaningful simulation.")
	case ErrFailedToScanWorldInput:
		return e.buildErrorMessage("Failed to scan one or more lines from the world map input file.")
	case ErrFailedToParseLine:
		return e.buildErrorMessage("Failed to parse world input.")
	case ErrUnknownDirection:
		return e.buildErrorMessage("Unknown direction specified.")
	case ErrCityAlreadyThere:
		return e.buildErrorMessage("Cannot locate a city on top of another city.")
	}
	return "Unrecognised error code"
}
