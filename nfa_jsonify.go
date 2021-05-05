package roughfa

import "encoding/json"

const (
	// EpsilonForJSON is for the epsilon transition for jsonify.
	EpsilonForJSON = '\a'
)

type (
	// NFAMachineShell is a serializable form of NFAMachine.
	NFAMachineShell struct {
		States         []string                       `json:"states"`
		Chars          []rune                         `json:"-"`
		RawChars       []string                       `json:"chars,omitempty"`
		StartStates    []string                       `json:"start_states"`
		AcceptStates   []string                       `json:"accept_states"`
		Transitions    map[string]map[rune][]string   `json:"-"`
		RawTransitions map[string]map[string][]string `json:"transitions"`
		CurrentStates  []string                       `json:"current_states,omitempty"`
	}
)

func (s NFAMachineShell) ToJSON() ([]byte, error) {
	rcs := make([]string, len(s.Chars))
	for i, x := range s.Chars {
		rcs[i] = string(x)
	}
	rts := make(map[string]map[string][]string, len(s.Transitions))
	for k, x := range s.Transitions {
		rts[k] = make(map[string][]string, len(x))
		for kx, kv := range x {
			if kx == Epsilon {
				rts[k][string(EpsilonForJSON)] = kv
				continue
			}
			rts[k][string(kx)] = kv
		}
	}
	return json.Marshal(NFAMachineShell{
		States:         s.States,
		RawChars:       rcs,
		StartStates:    s.StartStates,
		AcceptStates:   s.AcceptStates,
		RawTransitions: rts,
		CurrentStates:  s.CurrentStates,
	})
}

func NewNFAMachineShellFromJSON(b []byte) (*NFAMachineShell, error) {
	var s NFAMachineShell
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	cs := make([]rune, len(s.RawChars))
	for i, x := range s.RawChars {
		if len(x) != 1 {
			return nil, ErrCannotUnmarshalMachine
		}
		cs[i] = []rune(x)[0]
	}
	ts := make(map[string]map[rune][]string, len(s.RawTransitions))
	for k, x := range s.RawTransitions {
		ts[k] = make(map[rune][]string, len(x))
		for kx, kv := range x {
			if len(kx) != 1 {
				return nil, ErrCannotUnmarshalMachine
			}
			c := []rune(kx)[0]
			if c == EpsilonForJSON {
				ts[k][Epsilon] = kv
				continue
			}
			ts[k][c] = kv
		}
	}
	s.Chars = cs
	s.Transitions = ts
	return &s, nil
}

func (s NFAMachineShell) ToMachine() (NFAMachine, error) {
	m, err := NewNFAMachineBuilder().
		States(s.States).
		Chars(s.Chars).
		StartStates(s.StartStates).
		AcceptStates(s.AcceptStates).
		Transitions(s.Transitions).
		Build()
	if err != nil {
		return nil, err
	}
	if s.CurrentStates != nil {
		if err := m.SetStates(s.CurrentStates); err != nil {
			return nil, err
		}
	}
	return m, nil
}
