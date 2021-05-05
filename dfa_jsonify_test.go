package roughfa_test

import (
	"testing"

	"github.com/berquerant/roughfa"
	"github.com/stretchr/testify/assert"
)

func TestDFAMachineShell(t *testing.T) {
	toggle := func(m roughfa.DFAMachine) (roughfa.DFAMachine, error) {
		b, err := m.ToShell().ToJSON()
		if err != nil {
			return nil, err
		}
		t.Log(string(b))
		p, err := roughfa.NewDFAMachineShellFromJSON(b)
		if err != nil {
			return nil, err
		}
		return p.ToMachine()
	}
	m, err := roughfa.NewDFAMachineBuilder().
		States([]string{"even", "odd"}).
		StartState("even").
		AcceptStates([]string{"odd"}).
		Transitions(map[string]map[rune]string{
			"even": {
				'0': "even",
				'1': "odd",
			},
			"odd": {
				'0': "odd",
				'1': "even",
			},
		}).
		Build()
	assert.Nil(t, err)
	m, err = toggle(m)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "even", m.State())
	assert.Nil(t, m.Put('1'))
	m, err = toggle(m)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "odd", m.State())
}
