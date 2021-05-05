package roughfa

import (
	"fmt"
	"sort"
	"strings"

	"github.com/berquerant/roughfa/internal/dot"
	"github.com/berquerant/roughfa/internal/set"
)

const (
	// Epsilon is for the epsilon transition.
	Epsilon = 'ε'
)

type (
	// NFAMachine is a runner of the non deterministic finite automaton.
	NFAMachine interface {
		// States returns the current states.
		States() []string
		// SetStates sets the states as the current states.
		// Returns an error if a invalid state in the states exists.
		SetStates(states []string) error
		// Put inputs a character.
		// Returns an error if invalid input or no transitions.
		Put(x rune) error
		// IsAccepted returns true if the current states are acceptable.
		IsAccepted() bool
		// Reset resets the current states to the start state.
		Reset()
		// ToShell generates NFAMachineShell.
		ToShell() *NFAMachineShell
		// ApplyEpsilonExpansion creates a new Machine that applied the epsilon expansion from this NFAMachine.
		// The states that are not an accept state and have no outbound transitions.
		ApplyEpsilonExpansion() NFAMachine
		// ApplyPowersetConstruction creates a new NFAMachine that applied the powerset construction.
		ApplyPowersetConstruction() (NFAMachine, error)
		// ToDot generates Dot.
		ToDot() (dot.Dot, error)
		// HasEpsilon returns true if this has an epsilon transition.
		HasEpsilon() bool
		// IsDFA returns true if this is a dfa.
		IsDFA() bool
		// ToDFA creates a dfa from this.
		// Returns an error if this is not a dfa.
		ToDFA() (DFAMachine, error)
		// Not creates a Machine whose accept states are the diff of the states and the original accept states.
		Not() NFAMachine
		// Reverse creates a NFAMachine whose start states and accept states are reversed,
		// and transitions are reversed also.
		Reverse() NFAMachine
		// Minimize minimizes NFAMachine.
		Minimize() (NFAMachine, error)
	}

	nfaMachine struct {
		states        set.StringSet
		chars         set.RuneSet
		startStates   set.StringSet
		acceptStates  set.StringSet
		transitions   map[string]map[rune]set.StringSet
		currentStates set.StringSet
	}
)

// FromDFA creates a nfa from a dfa.
func FromDFA(m DFAMachine) NFAMachine {
	s := m.ToShell()
	ts := make(map[string]map[rune]set.StringSet, len(s.Transitions))
	for fromState, x := range s.Transitions {
		ts[fromState] = make(map[rune]set.StringSet, len(x))
		for c, toState := range x {
			ts[fromState][c] = set.NewStringSet(toState)
		}
	}
	return &nfaMachine{
		states:        set.NewStringSet(s.States...),
		startStates:   set.NewStringSet(s.StartState),
		chars:         set.NewRuneSet(s.Chars...),
		acceptStates:  set.NewStringSet(s.AcceptStates...),
		currentStates: set.NewStringSet(s.StartState),
		transitions:   ts,
	}
}

// ExpandEpsilon expands epsilon transitions.
// Generated transitions has no epsilon transitions.
func ExpandEpsilon(transitions map[string]map[rune]set.StringSet) map[string]map[rune]set.StringSet {
	ts := make(map[string]map[rune]set.StringSet)
	// extract not epsilon transitions
	for fromState, v := range transitions {
		routes := make(map[rune]set.StringSet)
		for c, toStates := range v {
			if c == Epsilon || toStates.Len() == 0 {
				continue
			}
			routes[c] = set.NewStringSet(toStates.Unwrap()...)
		}
		if len(routes) == 0 {
			continue
		}
		ts[fromState] = routes
	}
	// apply epsilon
	for _, v := range ts {
		for _, toStates := range v {
			modified := true
			for modified {
				modified = false
				old := set.NewStringSet(toStates.Unwrap()...)
				for _, toState := range toStates.Unwrap() {
					t, ok := transitions[toState]
					if !ok {
						continue
					}
					nextStates, ok := t[Epsilon]
					if !ok {
						continue
					}
					toStates.Add(nextStates.Unwrap()...)
				}
				if !modified {
					modified = !old.Equal(toStates)
				}
			}
		}
	}
	return ts
}

func (s nfaMachine) Minimize() (NFAMachine, error) {
	// apply brzozowski minimization
	m, err := s.Reverse().ApplyEpsilonExpansion().ApplyPowersetConstruction()
	if err != nil {
		return nil, err
	}
	return m.Reverse().ApplyEpsilonExpansion().ApplyPowersetConstruction()
}

func (s nfaMachine) Reverse() NFAMachine {
	transitions := map[string]map[rune]set.StringSet{}
	for fromState, x := range s.transitions {
		for c, toStates := range x {
			for _, toState := range toStates.Unwrap() {
				if _, ok := transitions[toState]; !ok {
					transitions[toState] = map[rune]set.StringSet{}
				}
				if _, ok := transitions[toState][c]; !ok {
					transitions[toState][c] = set.NewStringSet()
				}
				transitions[toState][c].Add(fromState)
			}
		}
	}
	return &nfaMachine{
		states:        s.states.Clone(),
		chars:         s.chars.Clone(),
		transitions:   transitions,
		startStates:   s.acceptStates.Clone(),
		acceptStates:  s.startStates.Clone(),
		currentStates: s.acceptStates.Clone(),
	}
}

func (s nfaMachine) Not() NFAMachine {
	acceptStates := s.states.Clone()
	acceptStates.Del(s.acceptStates.Unwrap()...)
	return &nfaMachine{
		states:        s.states.Clone(),
		chars:         s.chars.Clone(),
		startStates:   s.startStates.Clone(),
		acceptStates:  acceptStates,
		transitions:   s.transitions,
		currentStates: s.startStates.Clone(),
	}
}

func (s nfaMachine) ToDFA() (DFAMachine, error) {
	if !s.IsDFA() {
		return nil, ErrNotDFA
	}
	ts := make(map[string]map[rune]string, len(s.transitions))
	for fromState, x := range s.transitions {
		ts[fromState] = make(map[rune]string, len(x))
		for c, toStates := range x {
			ts[fromState][c] = toStates.Unwrap()[0]
		}
	}
	return NewDFAMachineBuilder().
		States(s.states.Unwrap()).
		AcceptStates(s.acceptStates.Unwrap()).
		Chars(s.chars.Unwrap()).
		StartState(s.startStates.Unwrap()[0]).
		Transitions(ts).
		Build()
}

func (s nfaMachine) IsDFA() bool {
	if s.HasEpsilon() || s.startStates.Len() != 1 {
		return false
	}
	for _, x := range s.transitions {
		for _, toStates := range x {
			if toStates.Len() != 1 {
				return false
			}
		}
	}
	return true
}

func (s nfaMachine) HasEpsilon() bool {
	for _, x := range s.transitions {
		if _, ok := x[Epsilon]; ok {
			return true
		}
	}
	return false
}

func (s nfaMachine) ToDot() (dot.Dot, error) {
	t := make(map[string]map[rune][]string)
	for fromState, x := range s.transitions {
		t[fromState] = make(map[rune][]string, len(x))
		for c, toStates := range x {
			t[fromState][c] = toStates.Unwrap()
		}
	}
	return dot.NewNFADotBuilder().
		StartStates(s.startStates.Unwrap()).
		States(s.states.Unwrap()).
		AcceptStates(s.acceptStates.Unwrap()).
		Transitions(t).
		Build()
}

func (s *nfaMachine) applyEpsilon() {
	modified := true
	for modified {
		modified = false
		for _, state := range s.currentStates.Unwrap() {
			t, ok := s.transitions[state]
			if !ok {
				continue
			}
			if u, ok := t[Epsilon]; ok {
				states := set.NewStringSet(s.currentStates.Unwrap()...)
				s.currentStates.Del(state)
				s.currentStates.Add(u.Unwrap()...)
				if !modified {
					modified = !states.Equal(s.currentStates)
				}
			}
		}
	}
}

func (s nfaMachine) applyEpsilonToStartStates() set.StringSet {
	var (
		startStates = s.startStates.Clone()
		modified    = true
	)
	for modified {
		modified = false
		for _, state := range startStates.Unwrap() {
			t, ok := s.transitions[state]
			if !ok {
				continue
			}
			if toStates, ok := t[Epsilon]; ok {
				states := startStates.Clone()
				startStates.Add(toStates.Unwrap()...)
				if !modified {
					modified = states.Equal(startStates)
				}
			}
		}
	}
	return startStates
}

func (s nfaMachine) ApplyEpsilonExpansion() NFAMachine {
	var (
		transitions   = ExpandEpsilon(s.transitions)
		uselessStates = set.NewStringSet()
	)
	// extract not an accept state and no outbound transitions
	for _, state := range s.states.Unwrap() {
		if s.acceptStates.In(state) {
			continue
		}
		if _, ok := transitions[state]; ok {
			continue
		}
		uselessStates.Add(state)
	}
	var (
		states        = s.states.Clone()
		startStates   = s.applyEpsilonToStartStates()
		currentStates = s.currentStates.Clone()
	)
	// exclude useless states
	states.Del(uselessStates.Unwrap()...)
	startStates.Del(uselessStates.Unwrap()...)
	for _, x := range transitions {
		for _, toStates := range x {
			toStates.Del(uselessStates.Unwrap()...)
		}
	}
	return &nfaMachine{
		states:        states,
		chars:         s.chars.Clone(),
		startStates:   startStates,
		acceptStates:  s.acceptStates.Clone(),
		currentStates: currentStates,
		transitions:   transitions,
	}
}

func (s nfaMachine) ToShell() *NFAMachineShell {
	t := make(map[string]map[rune][]string, len(s.transitions))
	for k, x := range s.transitions {
		t[k] = make(map[rune][]string, len(x))
		for kx, kv := range x {
			t[k][kx] = kv.Unwrap()
		}
	}
	return &NFAMachineShell{
		States:        s.states.Unwrap(),
		Chars:         s.chars.Unwrap(),
		StartStates:   s.startStates.Unwrap(),
		AcceptStates:  s.acceptStates.Unwrap(),
		Transitions:   t,
		CurrentStates: s.currentStates.Unwrap(),
	}
}

func (s nfaMachine) ApplyPowersetConstruction() (NFAMachine, error) {
	// requires no epsilon transitions
	if s.HasEpsilon() {
		return nil, ErrEpsilonExists
	}

	const dfaStartState = "0"
	var (
		states       = set.NewStringSet()
		transitions  = map[string]map[rune]string{}
		acceptStates = set.NewStringSet()
		chars        = func() set.RuneSet {
			cs := set.NewRuneSet()
			if s.chars.Len() > 0 {
				cs.Add(s.chars.Unwrap()...)
				return cs
			}
			for _, x := range s.transitions {
				for c := range x {
					cs.Add(c)
				}
			}
			return cs
		}()
		stringSetToString = func(x set.StringSet) string {
			y := x.Unwrap()
			sort.Strings(y)
			// TODO: use hash
			return strings.Join(y, "_")
		}
		dfaStatesMap = map[string]string{
			stringSetToString(s.startStates): dfaStartState,
		}
		q = []set.StringSet{set.NewStringSet(s.startStates.Unwrap()...)}
	)

	for len(q) > 0 {
		dState := q[0]
		q = q[1:]
		states.Add(dfaStatesMap[stringSetToString(dState)])
		if dState.And(s.acceptStates).Len() > 0 {
			acceptStates.Add(dfaStatesMap[stringSetToString(dState)])
		}
		for _, c := range chars.Unwrap() {
			dNext := set.NewStringSet()
			for _, x := range dState.Unwrap() {
				t, ok := s.transitions[x]
				if !ok {
					continue
				}
				if y, ok := t[c]; ok {
					dNext.Add(y.Unwrap()...)
				}
			}
			if dNext.Len() == 0 {
				continue
			}
			if _, ok := dfaStatesMap[stringSetToString(dNext)]; !ok {
				q = append(q, dNext)
				newState := fmt.Sprint(len(dfaStatesMap))
				dfaStatesMap[stringSetToString(dNext)] = newState
			}
			k := dfaStatesMap[stringSetToString(dState)]
			if _, ok := transitions[k]; !ok {
				transitions[k] = map[rune]string{}
			}
			transitions[k][c] = dfaStatesMap[stringSetToString(dNext)]
		}
	}

	ts := make(map[string]map[rune]set.StringSet, len(transitions))
	for fromState, x := range transitions {
		ts[fromState] = make(map[rune]set.StringSet, len(x))
		for c, toState := range x {
			ts[fromState][c] = set.NewStringSet(toState)
		}
	}
	return &nfaMachine{
		states:        states,
		chars:         chars,
		startStates:   set.NewStringSet(dfaStartState),
		acceptStates:  acceptStates,
		transitions:   ts,
		currentStates: set.NewStringSet(dfaStartState),
	}, nil

}

func (s *nfaMachine) Put(x rune) error {
	if s.chars.Len() > 0 && !s.chars.In(x) {
		return ErrInvalidInputChar
	}
	if s.currentStates.Len() == 0 {
		return ErrEmptyStates
	}
	if s.HasEpsilon() {
		s.applyEpsilon()
	}
	css := make([]set.StringSet, s.currentStates.Len())
	for i, state := range s.currentStates.Unwrap() {
		css[i] = set.NewStringSet()
		t, ok := s.transitions[state]
		if !ok {
			continue
		}
		if u, ok := t[x]; ok {
			css[i].Add(u.Unwrap()...)
		}
	}
	nextStates := set.NewStringSet()
	for _, cs := range css {
		nextStates.Add(cs.Unwrap()...)
	}
	s.currentStates = nextStates
	if s.HasEpsilon() {
		s.applyEpsilon()
	}
	if s.currentStates.Len() == 0 {
		return ErrEmptyStates
	}
	return nil
}
func (s nfaMachine) States() []string { return s.currentStates.Unwrap() }
func (s nfaMachine) IsAccepted() bool { return s.currentStates.And(s.acceptStates).Len() > 0 }
func (s *nfaMachine) Reset()          { s.currentStates = set.NewStringSet(s.startStates.Unwrap()...) }
func (s *nfaMachine) SetStates(states []string) error {
	x := set.NewStringSet(states...)
	if !s.states.In(x.Unwrap()...) {
		return ErrInvalidState
	}
	s.currentStates = x
	return nil
}
