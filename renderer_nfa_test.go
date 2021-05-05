package roughfa_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/berquerant/roughfa"
	"github.com/stretchr/testify/assert"
)

type renderNFATestcase struct {
	name         string
	filename     string
	startStates  []string
	states       []string
	acceptStates []string
	transitions  map[string]map[rune][]string
}

func (s renderNFATestcase) test(t *testing.T) {
	configureRenderTestcase(t)
	p := os.Getenv("PROJECT")
	filename := fmt.Sprintf("%s/tmp/%s", p, s.filename)
	t.Logf("render to %s", filename)
	d, err := roughfa.NewNFAMachineBuilder().
		StartStates(s.startStates).
		States(s.states).
		AcceptStates(s.acceptStates).
		Transitions(s.transitions).
		Build()
	if !assert.Nil(t, err) {
		return
	}
	dd, err := d.ToDot()
	if !assert.Nil(t, err) {
		return
	}
	t.Log(dd.AsDot())
	assert.Nil(t, roughfa.NewRenderer().
		DotCommand(os.Getenv("DOTCOMMAND")).
		Source(dd.AsDot()).
		Filename(filename).
		Render(),
	)
}

func TestRenderNFABuilder(t *testing.T) {
	for _, tc := range []*renderNFATestcase{
		{
			name:         "(a|b)*a",
			filename:     "nfa-regexp-3.png",
			states:       []string{"1", "2", "3"},
			startStates:  []string{"1"},
			acceptStates: []string{"2", "3"},
			transitions: map[string]map[rune][]string{
				"1": {
					'a': {"3"},
					'b': {"1"},
				},
				"2": {
					'a': {"2"},
					'b': {"1"},
				},
				"3": {
					'a': {"2"},
					'b': {"1"},
				},
			},
		},
		{
			name:     "basic",
			filename: "nfa-basic.png",
			states: []string{
				"0", "1", "2", "3", "4",
			},
			startStates:  []string{"0"},
			acceptStates: []string{"3", "4"},
			transitions: map[string]map[rune][]string{
				"0": {
					'0': {"1"},
				},
				"1": {
					roughfa.Epsilon: {"2"},
					'1':             {"1", "3"},
				},
				"2": {
					'0': {"4"},
				},
			},
		},
		{
			name:     "[a-c]",
			filename: "nfa-regexp-2.png",
			states: []string{
				"start",
				"a_start", "b_start", "c_start",
				"a_end", "b_end", "c_end",
				"end",
			},
			startStates:  []string{"start"},
			acceptStates: []string{"end"},
			transitions: map[string]map[rune][]string{
				"start": {
					roughfa.Epsilon: {"a_start", "b_start", "c_start"},
				},
				"a_start": {
					'a': {"a_end"},
				},
				"b_start": {
					'b': {"b_end"},
				},
				"c_start": {
					'c': {"c_end"},
				},
				"a_end": {
					roughfa.Epsilon: {"end"},
				},
				"b_end": {
					roughfa.Epsilon: {"end"},
				},
				"c_end": {
					roughfa.Epsilon: {"end"},
				},
			},
		},
		{
			name:     "a(b|c)*d",
			filename: "nfa-regexp-1.png",
			states: []string{
				"a_start", "a_end",
				"bc_start",
				"b_start", "b_end",
				"c_start", "c_end",
				"bc_end",
				"d_start",
				"d_end",
			},
			startStates:  []string{"a_start"},
			acceptStates: []string{"d_end"},
			transitions: map[string]map[rune][]string{
				"a_start": {
					'a': {"a_end"},
				},
				"a_end": {
					roughfa.Epsilon: {"bc_start", "d_start"},
				},
				"bc_start": {
					roughfa.Epsilon: {"b_start", "c_start"},
				},
				"b_start": {
					'b': {"b_end"},
				},
				"c_start": {
					'c': {"c_end"},
				},
				"b_end": {
					roughfa.Epsilon: {"bc_end"},
				},
				"c_end": {
					roughfa.Epsilon: {"bc_end"},
				},
				"bc_end": {
					roughfa.Epsilon: {"bc_start", "d_start"},
				},
				"d_start": {
					'd': {"d_end"},
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

type renderNFAExpandEpsilonTestcase struct {
	name         string
	filename     string
	startStates  []string
	states       []string
	acceptStates []string
	transitions  map[string]map[rune][]string
}

func (s renderNFAExpandEpsilonTestcase) test(t *testing.T) {
	configureRenderTestcase(t)
	p := os.Getenv("PROJECT")
	filename := fmt.Sprintf("%s/tmp/%s", p, s.filename)
	t.Logf("render to %s", filename)
	d, err := roughfa.NewNFAMachineBuilder().
		StartStates(s.startStates).
		States(s.states).
		AcceptStates(s.acceptStates).
		Transitions(s.transitions).
		Build()
	if !assert.Nil(t, err) {
		return
	}
	dd := d.ApplyEpsilonExpansion()
	assert.False(t, dd.HasEpsilon())
	ddd, err := dd.ToDot()
	if !assert.Nil(t, err) {
		return
	}
	t.Log(ddd.AsDot())
	assert.Nil(t, roughfa.NewRenderer().
		DotCommand(os.Getenv("DOTCOMMAND")).
		Source(ddd.AsDot()).
		Filename(filename).
		Render(),
	)
}

func TestRenderNFAExpandEpsilonBuilder(t *testing.T) {
	for _, tc := range []*renderNFAExpandEpsilonTestcase{
		{
			name:         "(a|b)*a",
			filename:     "nfa-regexp-epsilon-expanded3.png",
			states:       []string{"1", "2", "3"},
			startStates:  []string{"1"},
			acceptStates: []string{"2", "3"},
			transitions: map[string]map[rune][]string{
				"1": {
					'a': {"3"},
					'b': {"1"},
				},
				"2": {
					'a': {"2"},
					'b': {"1"},
				},
				"3": {
					'a': {"2"},
					'b': {"1"},
				},
			},
		},
		{
			name:     "basic",
			filename: "nfa-basic-epsilon-expanded.png",
			states: []string{
				"0", "1", "2", "3", "4",
			},
			startStates:  []string{"0"},
			acceptStates: []string{"3", "4"},
			transitions: map[string]map[rune][]string{
				"0": {
					'0': {"1"},
				},
				"1": {
					roughfa.Epsilon: {"2"},
					'1':             {"1", "3"},
				},
				"2": {
					'0': {"4"},
				},
			},
		},
		{
			name:     "[a-c]",
			filename: "nfa-regexp-epsilon-expanded-2.png",
			states: []string{
				"start",
				"a_start", "b_start", "c_start",
				"a_end", "b_end", "c_end",
				"end",
			},
			startStates:  []string{"start"},
			acceptStates: []string{"end"},
			transitions: map[string]map[rune][]string{
				"start": {
					roughfa.Epsilon: {"a_start", "b_start", "c_start"},
				},
				"a_start": {
					'a': {"a_end"},
				},
				"b_start": {
					'b': {"b_end"},
				},
				"c_start": {
					'c': {"c_end"},
				},
				"a_end": {
					roughfa.Epsilon: {"end"},
				},
				"b_end": {
					roughfa.Epsilon: {"end"},
				},
				"c_end": {
					roughfa.Epsilon: {"end"},
				},
			},
		},
		{
			name:     "a(b|c)*d",
			filename: "nfa-regexp-epsilon-expanded-1.png",
			states: []string{
				"a_start", "a_end",
				"bc_start",
				"b_start", "b_end",
				"c_start", "c_end",
				"bc_end",
				"d_start",
				"d_end",
			},
			startStates:  []string{"a_start"},
			acceptStates: []string{"d_end"},
			transitions: map[string]map[rune][]string{
				"a_start": {
					'a': {"a_end"},
				},
				"a_end": {
					roughfa.Epsilon: {"bc_start", "d_start"},
				},
				"bc_start": {
					roughfa.Epsilon: {"b_start", "c_start"},
				},
				"b_start": {
					'b': {"b_end"},
				},
				"c_start": {
					'c': {"c_end"},
				},
				"b_end": {
					roughfa.Epsilon: {"bc_end"},
				},
				"c_end": {
					roughfa.Epsilon: {"bc_end"},
				},
				"bc_end": {
					roughfa.Epsilon: {"bc_start", "d_start"},
				},
				"d_start": {
					'd': {"d_end"},
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

type renderNFAPowersetConstructionTestcase struct {
	name         string
	filename     string
	startStates  []string
	states       []string
	acceptStates []string
	transitions  map[string]map[rune][]string
}

func (s renderNFAPowersetConstructionTestcase) test(t *testing.T) {
	configureRenderTestcase(t)
	p := os.Getenv("PROJECT")
	filename := fmt.Sprintf("%s/tmp/%s", p, s.filename)
	t.Logf("render to %s", filename)
	d, err := roughfa.NewNFAMachineBuilder().
		StartStates(s.startStates).
		States(s.states).
		AcceptStates(s.acceptStates).
		Transitions(s.transitions).
		Build()
	if !assert.Nil(t, err) {
		return
	}
	dd, err := d.ApplyEpsilonExpansion().ApplyPowersetConstruction()
	if !assert.Nil(t, err) {
		return
	}
	assert.True(t, dd.IsDFA())
	ddd, err := dd.ToDot()
	if !assert.Nil(t, err) {
		return
	}
	t.Log(ddd.AsDot())
	assert.Nil(t, roughfa.NewRenderer().
		DotCommand(os.Getenv("DOTCOMMAND")).
		Source(ddd.AsDot()).
		Filename(filename).
		Render(),
	)
}

func TestRenderNFAPowersetConstructionBuilder(t *testing.T) {
	for _, tc := range []*renderNFAPowersetConstructionTestcase{
		{
			name:         "(a|b)*a",
			filename:     "nfa-regexp-powerset-construction-3.png",
			states:       []string{"1", "2", "3"},
			startStates:  []string{"1"},
			acceptStates: []string{"2", "3"},
			transitions: map[string]map[rune][]string{
				"1": {
					'a': {"3"},
					'b': {"1"},
				},
				"2": {
					'a': {"2"},
					'b': {"1"},
				},
				"3": {
					'a': {"2"},
					'b': {"1"},
				},
			},
		},
		{
			name:     "basic",
			filename: "nfa-basic-powerset-construction.png",
			states: []string{
				"0", "1", "2", "3", "4",
			},
			startStates:  []string{"0"},
			acceptStates: []string{"3", "4"},
			transitions: map[string]map[rune][]string{
				"0": {
					'0': {"1"},
				},
				"1": {
					roughfa.Epsilon: {"2"},
					'1':             {"1", "3"},
				},
				"2": {
					'0': {"4"},
				},
			},
		},
		{
			name:     "[a-c]",
			filename: "nfa-regexp-powerset-construction-2.png",
			states: []string{
				"start",
				"a_start", "b_start", "c_start",
				"a_end", "b_end", "c_end",
				"end",
			},
			startStates:  []string{"start"},
			acceptStates: []string{"end"},
			transitions: map[string]map[rune][]string{
				"start": {
					roughfa.Epsilon: {"a_start", "b_start", "c_start"},
				},
				"a_start": {
					'a': {"a_end"},
				},
				"b_start": {
					'b': {"b_end"},
				},
				"c_start": {
					'c': {"c_end"},
				},
				"a_end": {
					roughfa.Epsilon: {"end"},
				},
				"b_end": {
					roughfa.Epsilon: {"end"},
				},
				"c_end": {
					roughfa.Epsilon: {"end"},
				},
			},
		},
		{
			name:     "a(b|c)*d",
			filename: "nfa-regexp-powerset-construction-1.png",
			states: []string{
				"a_start", "a_end",
				"bc_start",
				"b_start", "b_end",
				"c_start", "c_end",
				"bc_end",
				"d_start",
				"d_end",
			},
			startStates:  []string{"a_start"},
			acceptStates: []string{"d_end"},
			transitions: map[string]map[rune][]string{
				"a_start": {
					'a': {"a_end"},
				},
				"a_end": {
					roughfa.Epsilon: {"bc_start", "d_start"},
				},
				"bc_start": {
					roughfa.Epsilon: {"b_start", "c_start"},
				},
				"b_start": {
					'b': {"b_end"},
				},
				"c_start": {
					'c': {"c_end"},
				},
				"b_end": {
					roughfa.Epsilon: {"bc_end"},
				},
				"c_end": {
					roughfa.Epsilon: {"bc_end"},
				},
				"bc_end": {
					roughfa.Epsilon: {"bc_start", "d_start"},
				},
				"d_start": {
					'd': {"d_end"},
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

type renderNFAMinimizeTestcase struct {
	name         string
	filename     string
	startStates  []string
	states       []string
	acceptStates []string
	transitions  map[string]map[rune][]string
}

func (s renderNFAMinimizeTestcase) test(t *testing.T) {
	configureRenderTestcase(t)
	p := os.Getenv("PROJECT")
	filename := fmt.Sprintf("%s/tmp/%s", p, s.filename)
	t.Logf("render to %s", filename)
	d, err := roughfa.NewNFAMachineBuilder().
		StartStates(s.startStates).
		States(s.states).
		AcceptStates(s.acceptStates).
		Transitions(s.transitions).
		Build()
	if !assert.Nil(t, err) {
		return
	}
	dd, err := d.Minimize()
	if !assert.Nil(t, err) {
		return
	}
	assert.True(t, dd.IsDFA())
	ddd, err := dd.ToDot()
	if !assert.Nil(t, err) {
		return
	}
	t.Log(ddd.AsDot())
	assert.Nil(t, roughfa.NewRenderer().
		DotCommand(os.Getenv("DOTCOMMAND")).
		Source(ddd.AsDot()).
		Filename(filename).
		Render(),
	)
}

func TestRenderNFAMinimizeBuilder(t *testing.T) {
	for _, tc := range []*renderNFAMinimizeTestcase{
		{
			name:         "(a|b)*a",
			filename:     "nfa-regexp-minimize-3.png",
			states:       []string{"1", "2", "3"},
			startStates:  []string{"1"},
			acceptStates: []string{"2", "3"},
			transitions: map[string]map[rune][]string{
				"1": {
					'a': {"3"},
					'b': {"1"},
				},
				"2": {
					'a': {"2"},
					'b': {"1"},
				},
				"3": {
					'a': {"2"},
					'b': {"1"},
				},
			},
		},
		{
			name:     "basic",
			filename: "nfa-basic-minimize.png",
			states: []string{
				"0", "1", "2", "3", "4",
			},
			startStates:  []string{"0"},
			acceptStates: []string{"3", "4"},
			transitions: map[string]map[rune][]string{
				"0": {
					'0': {"1"},
				},
				"1": {
					roughfa.Epsilon: {"2"},
					'1':             {"1", "3"},
				},
				"2": {
					'0': {"4"},
				},
			},
		},
		{
			name:     "[a-c]",
			filename: "nfa-regexp-minimize-2.png",
			states: []string{
				"start",
				"a_start", "b_start", "c_start",
				"a_end", "b_end", "c_end",
				"end",
			},
			startStates:  []string{"start"},
			acceptStates: []string{"end"},
			transitions: map[string]map[rune][]string{
				"start": {
					roughfa.Epsilon: {"a_start", "b_start", "c_start"},
				},
				"a_start": {
					'a': {"a_end"},
				},
				"b_start": {
					'b': {"b_end"},
				},
				"c_start": {
					'c': {"c_end"},
				},
				"a_end": {
					roughfa.Epsilon: {"end"},
				},
				"b_end": {
					roughfa.Epsilon: {"end"},
				},
				"c_end": {
					roughfa.Epsilon: {"end"},
				},
			},
		},
		{
			name:     "a(b|c)*d",
			filename: "nfa-regexp-minimize-1.png",
			states: []string{
				"a_start", "a_end",
				"bc_start",
				"b_start", "b_end",
				"c_start", "c_end",
				"bc_end",
				"d_start",
				"d_end",
			},
			startStates:  []string{"a_start"},
			acceptStates: []string{"d_end"},
			transitions: map[string]map[rune][]string{
				"a_start": {
					'a': {"a_end"},
				},
				"a_end": {
					roughfa.Epsilon: {"bc_start", "d_start"},
				},
				"bc_start": {
					roughfa.Epsilon: {"b_start", "c_start"},
				},
				"b_start": {
					'b': {"b_end"},
				},
				"c_start": {
					'c': {"c_end"},
				},
				"b_end": {
					roughfa.Epsilon: {"bc_end"},
				},
				"c_end": {
					roughfa.Epsilon: {"bc_end"},
				},
				"bc_end": {
					roughfa.Epsilon: {"bc_start", "d_start"},
				},
				"d_start": {
					'd': {"d_end"},
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}
