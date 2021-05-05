package roughfa

import (
	"github.com/berquerant/roughfa/internal/dot"
	"github.com/berquerant/roughfa/internal/set"
)

type (
	// DFAMachine is a runner of the deterministic finite automaton.
	DFAMachine interface {
		// State returns the current state.
		State() string
		// SetState sets the state as the current state.
		// Returns an error if the state is not in the states.
		SetState(state string) error
		// Put inputs a character.
		// Returns an error if invalid input or no transitions.
		Put(x rune) error
		// IsAccepted returns true if the current state is in the accept states.
		IsAccepted() bool
		// Reset resets the current state to the start state.
		Reset()
		// ToDot generates Dot.
		ToDot() (dot.Dot, error)
		// ToShell generates DFAMachineShell.
		ToShell() *DFAMachineShell
	}
	dfaMachine struct {
		states       set.StringSet
		chars        set.RuneSet
		startState   string
		acceptStates set.StringSet
		transitions  map[string]map[rune]string
		currentState string
	}
)

func (s dfaMachine) ToShell() *DFAMachineShell {
	return &DFAMachineShell{
		States:       s.states.Unwrap(),
		Chars:        s.chars.Unwrap(),
		StartState:   s.startState,
		AcceptStates: s.acceptStates.Unwrap(),
		Transitions:  s.transitions,
		CurrentState: s.currentState,
	}
}
func (s dfaMachine) ToDot() (dot.Dot, error) {
	return dot.NewDFADotBuilder().
		StartState(s.startState).
		States(s.states.Unwrap()).
		AcceptStates(s.acceptStates.Unwrap()).
		Transitions(s.transitions).
		Build()
}
func (s *dfaMachine) SetState(state string) error {
	if !s.states.In(state) {
		return ErrInvalidState
	}
	s.currentState = state
	return nil
}
func (s *dfaMachine) Reset()          { s.currentState = s.startState }
func (s dfaMachine) IsAccepted() bool { return s.acceptStates.In(s.currentState) }
func (s dfaMachine) State() string    { return s.currentState }
func (s *dfaMachine) Put(x rune) error {
	if s.chars.Len() > 0 && !s.chars.In(x) {
		return ErrInvalidInputChar
	}
	t, ok := s.transitions[s.currentState]
	if !ok {
		return ErrOutOfTransition
	}
	nextState, ok := t[x]
	if !ok {
		return ErrOutOfTransition
	}
	s.currentState = nextState
	return nil
}
