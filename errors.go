package roughfa

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidStartState      = errors.New("invalid start state")
	ErrInvalidAcceptStates    = errors.New("invalid accept states")
	ErrInvalidTransitions     = errors.New("invalid transitions")
	ErrInvalidInputChar       = errors.New("invalid input character")
	ErrOutOfTransition        = errors.New("out of transition")
	ErrInvalidState           = errors.New("invalid state")
	ErrCannotUnmarshalMachine = errors.New("cannot unmarshal machine")
	ErrNotDFA                 = errors.New("not dfa")
	ErrEpsilonExists          = errors.New("epsilon exists")
	ErrEmptyStates            = errors.New("empty states")
	ErrInvalidStartStates     = errors.New("invalid start states")
	ErrNoDotSource            = errors.New("no dot source")
)

type (
	// RenderError represents a rendering error.
	RenderError struct {
		Stderr string
		Err    error
	}
)

func (s RenderError) Error() string { return fmt.Sprintf("%s: %s", s.Err.Error(), s.Stderr) }
func (s RenderError) Unwrap() error { return s.Err }
