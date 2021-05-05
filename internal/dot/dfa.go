package dot

type (
	DFADotBuilder interface {
		// StartState sets the start state.
		StartState(startState string) DFADotBuilder
		// States sets the states.
		States(states []string) DFADotBuilder
		// AcceptStates sets the accept states.
		AcceptStates(acceptStates []string) DFADotBuilder
		// Transitions sets the transition map.
		Transitions(transitions map[string]map[rune]string) DFADotBuilder
		// Build generates Dot.
		// Returns an error if some contradictions exist.
		Build() (Dot, error)
	}

	dfaDotBuilder struct {
		baseFaDotBuilder
	}
)

// NewDFADotBuilder creates a new DFADotBuilder.
func NewDFADotBuilder() DFADotBuilder { return &dfaDotBuilder{} }

func (s *dfaDotBuilder) StartState(startState string) DFADotBuilder {
	s.startStates = []string{startState}
	return s
}
func (s *dfaDotBuilder) States(states []string) DFADotBuilder {
	s.states = states
	return s
}
func (s *dfaDotBuilder) AcceptStates(acceptStates []string) DFADotBuilder {
	s.acceptStates = acceptStates
	return s
}
func (s *dfaDotBuilder) Transitions(transitions map[string]map[rune]string) DFADotBuilder {
	t := make(map[string]map[rune][]string, len(transitions))
	for k, v := range transitions {
		t[k] = make(map[rune][]string, len(v))
		for kx, kv := range v {
			t[k][kx] = []string{kv}
		}
	}
	s.transitions = t
	return s
}
