package dot

const (
	// Epsilon is for the epsilon transition.
	Epsilon = 'ε'
)

type (
	NFADotBuilder interface {
		// StartStates sets the start state.
		StartStates(startState []string) NFADotBuilder
		// States sets the states.
		States(states []string) NFADotBuilder
		// AcceptStates sets the accept states.
		AcceptStates(acceptStates []string) NFADotBuilder
		// Transitions sets the transition map.
		Transitions(transitions map[string]map[rune][]string) NFADotBuilder
		// Build generates Dot.
		// Returns an error if some contradictions exist.
		Build() (Dot, error)
	}

	nfaDotBuilder struct {
		baseFaDotBuilder
	}
)

// NewNFADotBuilder creates a new NFADotBuilder.
func NewNFADotBuilder() NFADotBuilder { return &nfaDotBuilder{} }

func (s *nfaDotBuilder) StartStates(startStates []string) NFADotBuilder {
	s.startStates = startStates
	return s
}
func (s *nfaDotBuilder) States(states []string) NFADotBuilder {
	s.states = states
	return s
}
func (s *nfaDotBuilder) AcceptStates(acceptStates []string) NFADotBuilder {
	s.acceptStates = acceptStates
	return s
}
func (s *nfaDotBuilder) Transitions(transitions map[string]map[rune][]string) NFADotBuilder {
	t := make(map[string]map[rune][]string, len(transitions))
	for k, v := range transitions {
		t[k] = make(map[rune][]string, len(v))
		for kx, kv := range v {
			t[k][kx] = kv
		}
	}
	s.transitions = t
	return s
}
