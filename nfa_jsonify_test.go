package roughfa_test

import (
	"testing"

	"github.com/berquerant/roughfa"
	"github.com/berquerant/roughfa/internal/set"
	"github.com/stretchr/testify/assert"
)

func TestNFAMachineShell(t *testing.T) {
	toggle := func(m roughfa.NFAMachine) (roughfa.NFAMachine, error) {
		b, err := m.ToShell().ToJSON()
		if err != nil {
			return nil, err
		}
		t.Log(string(b))
		p, err := roughfa.NewNFAMachineShellFromJSON(b)
		if err != nil {
			return nil, err
		}
		return p.ToMachine()
	}
	m, err := roughfa.NewNFAMachineBuilder().
		States([]string{
			"a-start", "a-end",
			"bc-start",
			"b-start", "b-end",
			"c-start", "c-end",
			"bc-end",
			"d-start",
			"d-end",
		}).
		StartStates([]string{"a-start"}).
		AcceptStates([]string{"d-end"}).
		Transitions(map[string]map[rune][]string{
			"a-start": {
				'a': {"a-end"},
			},
			"a-end": {
				roughfa.Epsilon: {"bc-start", "d-start"},
			},
			"bc-start": {
				roughfa.Epsilon: {"b-start", "c-start"},
			},
			"b-start": {
				'b': {"b-end"},
			},
			"c-start": {
				'c': {"c-end"},
			},
			"b-end": {
				roughfa.Epsilon: {"bc-end"},
			},
			"c-end": {
				roughfa.Epsilon: {"bc-end"},
			},
			"bc-end": {
				roughfa.Epsilon: {"bc-start", "d-start"},
			},
			"d-start": {
				'd': {"d-end"},
			},
		}).
		Build()
	assert.Nil(t, err)
	m, err = toggle(m)
	if !assert.Nil(t, err) {
		return
	}
	assert.True(
		t,
		set.NewStringSet(m.States()...).Equal(set.NewStringSet("a-start")),
	)
	assert.Nil(t, m.Put('a'))
	m, err = toggle(m)
	if !assert.Nil(t, err) {
		return
	}
	assert.True(
		t,
		set.NewStringSet(m.States()...).Equal(set.NewStringSet("b-start", "c-start", "d-start")),
	)
}
