package roughfa_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/berquerant/roughfa"
	"github.com/berquerant/roughfa/internal/dot"
	"github.com/stretchr/testify/assert"
)

type renderDFATestcase struct {
	name         string
	filename     string
	startState   string
	states       []string
	acceptStates []string
	transitions  map[string]map[rune]string
}

func (s renderDFATestcase) test(t *testing.T) {
	configureRenderTestcase(t)
	p := os.Getenv("PROJECT")
	filename := fmt.Sprintf("%s/tmp/%s", p, s.filename)
	t.Logf("render to %s", filename)
	d, err := dot.NewDFADotBuilder().
		StartState(s.startState).
		States(s.states).
		AcceptStates(s.acceptStates).
		Transitions(s.transitions).
		Build()
	if !assert.Nil(t, err) {
		return
	}
	t.Log(d.AsDot())
	assert.Nil(t, roughfa.NewRenderer().
		DotCommand(os.Getenv("DOTCOMMAND")).
		Source(d.AsDot()).
		Filename(filename).
		Render(),
	)
}

func TestRenderDFABuilder(t *testing.T) {
	for _, tc := range []*renderDFATestcase{
		{
			name:         "minimum",
			filename:     "minimum-dfa.png",
			states:       []string{"s1"},
			startState:   "s1",
			acceptStates: []string{"s1"},
		},
		{
			name:         "2 states",
			filename:     "2-states-dfa.png",
			states:       []string{"even", "odd"},
			startState:   "even",
			acceptStates: []string{"odd"},
			transitions: map[string]map[rune]string{
				"even": {
					'0': "even",
					'1': "odd",
				},
				"odd": {
					'0': "odd",
					'1': "even",
				},
			},
		},
		{
			name:         "3 states",
			filename:     "3-states-dfa.png",
			states:       []string{"start", "1", "11"},
			startState:   "start",
			acceptStates: []string{"11"},
			transitions: map[string]map[rune]string{
				"start": {
					'1': "1",
					'0': "start",
				},
				"1": {
					'1': "11",
					'0': "start",
				},
				"11": {
					'1': "11",
					'0': "start",
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}
