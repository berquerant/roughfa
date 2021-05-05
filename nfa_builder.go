package roughfa

import "github.com/berquerant/roughfa/internal/set"

type (
	// NFAMachineBuilder is a builder of NFAMachine.
	NFAMachineBuilder interface {
		// States configures the states of the machine.
		States(states []string) NFAMachineBuilder
		// Chars configures the set of the input characters.
		// If chars is empty or not set, means universe.
		Chars(chars []rune) NFAMachineBuilder
		// StartStates configures the initial states of the machine.
		// Required.
		// Must be included in States.
		StartStates(state []string) NFAMachineBuilder
		// AcceptStates configures the accept states of the machine.
		// Required.
		// Must be a subset of States.
		AcceptStates(acceptStates []string) NFAMachineBuilder
		// Transitions configures the transition function of the machine.
		// Required.
		Transitions(transitions map[string]map[rune][]string) NFAMachineBuilder
		// Build creates a new NFAMachine.
		// Returns an error if some validations fails.
		Build() (NFAMachine, error)
	}

	nfaMachineBuilder struct {
		states       []string
		chars        []rune
		startStates  []string
		acceptStates []string
		transitions  map[string]map[rune][]string
	}
)

// NewNFAMachineBuilder creates a new NFAMachineBuilder.
func NewNFAMachineBuilder() NFAMachineBuilder { return &nfaMachineBuilder{} }

func (s *nfaMachineBuilder) States(states []string) NFAMachineBuilder {
	s.states = states
	return s
}
func (s *nfaMachineBuilder) Chars(chars []rune) NFAMachineBuilder {
	s.chars = chars
	return s
}
func (s *nfaMachineBuilder) StartStates(startStates []string) NFAMachineBuilder {
	s.startStates = startStates
	return s
}
func (s *nfaMachineBuilder) AcceptStates(acceptStates []string) NFAMachineBuilder {
	s.acceptStates = acceptStates
	return s
}
func (s *nfaMachineBuilder) Transitions(transitions map[string]map[rune][]string) NFAMachineBuilder {
	s.transitions = transitions
	return s
}
func (s nfaMachineBuilder) Build() (NFAMachine, error) {
	states := set.NewStringSet(s.states...)
	if !states.In(s.startStates...) {
		return nil, ErrInvalidStartStates
	}
	if !states.In(s.acceptStates...) {
		return nil, ErrInvalidAcceptStates
	}
	var (
		tStates     = set.NewStringSet()
		tChars      = set.NewRuneSet()
		transitions = make(map[string]map[rune]set.StringSet, len(s.transitions))
	)
	for state, v := range s.transitions {
		if len(v) == 0 {
			continue
		}
		tStates.Add(state)
		transitions[state] = make(map[rune]set.StringSet, len(v))
		for char, dests := range v {
			if len(dests) == 0 {
				continue
			}
			tChars.Add(char)
			tStates.Add(dests...)
			transitions[state][char] = set.NewStringSet(dests...)
		}
	}
	chars := set.NewRuneSet(s.chars...)
	if !states.In(tStates.Unwrap()...) || chars.Len() > 0 && !chars.In(tChars.Unwrap()...) {
		return nil, ErrInvalidTransitions
	}
	return &nfaMachine{
		states:        states,
		chars:         chars,
		startStates:   set.NewStringSet(s.startStates...),
		acceptStates:  set.NewStringSet(s.acceptStates...),
		transitions:   transitions,
		currentStates: set.NewStringSet(s.startStates...),
	}, nil
}
