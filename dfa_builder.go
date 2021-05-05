package roughfa

import "github.com/berquerant/roughfa/internal/set"

type (
	// DFAMachineBuilder is a builder of DFAMachine.
	DFAMachineBuilder interface {
		// States configures the states of the machine.
		// Required.
		States(states []string) DFAMachineBuilder
		// Chars configures the set of the input characters.
		// If chars is empty or not set, means universe.
		Chars(chars []rune) DFAMachineBuilder
		// StartState configures the initial state of the machine.
		// Required.
		// Must be included in States.
		StartState(state string) DFAMachineBuilder
		// AcceptStates configures the accept states of the machine.
		// Required.
		// Must be a subset of States.
		AcceptStates(acceptStates []string) DFAMachineBuilder
		// Transitions configures the transition function of the machine.
		// Required.
		Transitions(transitions map[string]map[rune]string) DFAMachineBuilder
		// Build creates a new DFAMachine.
		// Returns an error if some validations fails.
		Build() (DFAMachine, error)
	}

	dfaMachineBuilder struct {
		states       []string
		chars        []rune
		startState   string
		acceptStates []string
		transitions  map[string]map[rune]string
	}
)

// NewDFAMachineBuilder creates a new DFAMachineBuilder.
func NewDFAMachineBuilder() DFAMachineBuilder { return &dfaMachineBuilder{} }

func (s *dfaMachineBuilder) States(states []string) DFAMachineBuilder {
	s.states = states
	return s
}
func (s *dfaMachineBuilder) Chars(chars []rune) DFAMachineBuilder {
	s.chars = chars
	return s
}
func (s *dfaMachineBuilder) StartState(state string) DFAMachineBuilder {
	s.startState = state
	return s
}
func (s *dfaMachineBuilder) AcceptStates(acceptStates []string) DFAMachineBuilder {
	s.acceptStates = acceptStates
	return s
}
func (s *dfaMachineBuilder) Transitions(transitions map[string]map[rune]string) DFAMachineBuilder {
	s.transitions = transitions
	return s
}
func (s dfaMachineBuilder) Build() (DFAMachine, error) {
	states := set.NewStringSet(s.states...)
	if !states.In(s.startState) {
		return nil, ErrInvalidStartState
	}
	if !states.In(s.acceptStates...) {
		return nil, ErrInvalidAcceptStates
	}
	var (
		tStates = set.NewStringSet()
		tChars  = set.NewRuneSet()
	)
	for state, v := range s.transitions {
		if len(v) == 0 {
			continue
		}
		tStates.Add(state)
		for char, dest := range v {
			tStates.Add(dest)
			tChars.Add(char)
		}
	}
	chars := set.NewRuneSet(s.chars...)
	if !states.In(tStates.Unwrap()...) || chars.Len() > 0 && !chars.In(tChars.Unwrap()...) {
		return nil, ErrInvalidTransitions
	}
	return &dfaMachine{
		states:       states,
		chars:        chars,
		acceptStates: set.NewStringSet(s.acceptStates...),
		startState:   s.startState,
		transitions:  s.transitions,
		currentState: s.startState,
	}, nil
}
