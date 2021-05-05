package roughfa

import "encoding/json"

type (
	// DFAMachineShell is a serializable form of DFAMachine.
	DFAMachineShell struct {
		States         []string                     `json:"states"`
		Chars          []rune                       `json:"-"`
		RawChars       []string                     `json:"chars,omitempty"`
		StartState     string                       `json:"start_state"`
		AcceptStates   []string                     `json:"accept_states"`
		Transitions    map[string]map[rune]string   `json:"-"`
		RawTransitions map[string]map[string]string `json:"transitions"`
		CurrentState   string                       `json:"current_state,omitempty"`
	}
)

func (s DFAMachineShell) ToJSON() ([]byte, error) {
	rcs := make([]string, len(s.Chars))
	for i, x := range s.Chars {
		rcs[i] = string(x)
	}
	rts := make(map[string]map[string]string, len(s.Transitions))
	for k, x := range s.Transitions {
		rts[k] = make(map[string]string, len(x))
		for kx, kv := range x {
			rts[k][string(kx)] = kv
		}
	}
	return json.Marshal(DFAMachineShell{
		States:         s.States,
		RawChars:       rcs,
		StartState:     s.StartState,
		AcceptStates:   s.AcceptStates,
		RawTransitions: rts,
		CurrentState:   s.CurrentState,
	})
}

func NewDFAMachineShellFromJSON(b []byte) (*DFAMachineShell, error) {
	var s DFAMachineShell
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
	ts := make(map[string]map[rune]string, len(s.RawTransitions))
	for k, x := range s.RawTransitions {
		ts[k] = make(map[rune]string, len(x))
		for kx, kv := range x {
			if len(kx) != 1 {
				return nil, ErrCannotUnmarshalMachine
			}
			ts[k][[]rune(kx)[0]] = kv
		}
	}
	s.Chars = cs
	s.Transitions = ts
	return &s, nil
}

func (s DFAMachineShell) ToMachine() (DFAMachine, error) {
	m, err := NewDFAMachineBuilder().
		States(s.States).
		Chars(s.Chars).
		StartState(s.StartState).
		AcceptStates(s.AcceptStates).
		Transitions(s.Transitions).
		Build()
	if err != nil {
		return nil, err
	}
	if s.CurrentState != "" {
		if err := m.SetState(s.CurrentState); err != nil {
			return nil, err
		}
	}
	return m, nil
}
