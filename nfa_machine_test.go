package roughfa_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/berquerant/roughfa"
	"github.com/berquerant/roughfa/internal/set"
	"github.com/stretchr/testify/assert"
)

type expandEpsilonTestcase struct {
	name        string
	transitions map[string]map[rune][]string
	want        map[string]map[rune][]string
}

func (s expandEpsilonTestcase) test(t *testing.T) {
	ts := make(map[string]map[rune]set.StringSet, len(s.transitions))
	for k, x := range s.transitions {
		ts[k] = make(map[rune]set.StringSet, len(x))
		for kx, kv := range x {
			ts[k][kx] = set.NewStringSet(kv...)
		}
	}
	got := roughfa.ExpandEpsilon(ts)
	t.Logf("want: %v", s.want)
	t.Logf("got : %v", got)
	assert.Equal(t, len(s.want), len(got), "length of transitions")
	for wk := range s.want {
		wx := s.want[wk]
		gx, ok := got[wk]
		if !assert.True(t, ok, "%s transition not exists", wk) {
			continue
		}
		assert.Equal(t, len(wx), len(gx), "%s transition length", wk)
		for wkx := range wx {
			wkv := set.NewStringSet(wx[wkx]...)
			gkv, ok := gx[wkx]
			if !assert.True(t, ok, "%q transition not exists", wkx) {
				continue
			}
			assert.True(t, wkv.Equal(gkv), "%s %q: %v != %v", wk, wkx, gkv.Unwrap(), wkv.Unwrap())
		}
	}
}

func TestExpandEpsilon(t *testing.T) {
	for _, tc := range []*expandEpsilonTestcase{
		{
			name: "[a-c]",
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
			want: map[string]map[rune][]string{
				"a_start": {
					'a': {"a_end", "end"},
				},
				"b_start": {
					'b': {"b_end", "end"},
				},
				"c_start": {
					'c': {"c_end", "end"},
				},
			},
		},
		{
			name: "first",
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
			want: map[string]map[rune][]string{
				"0": {
					'0': {"1", "2"},
				},
				"1": {
					'1': {"1", "2", "3"},
				},
				"2": {
					'0': {"4"},
				},
			},
		},
		{
			name: "a(b|c)*d",
			transitions: map[string]map[rune][]string{
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
			},
			want: map[string]map[rune][]string{
				"a-start": {
					'a': {"b-start", "c-start", "d-start", "bc-start", "a-end"},
				},
				"b-start": {
					'b': {"b-start", "c-start", "d-start", "b-end", "bc-end", "bc-start"},
				},
				"c-start": {
					'c': {"b-start", "c-start", "d-start", "c-end", "bc-end", "bc-start"},
				},
				"d-start": {
					'd': {"d-end"},
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

type (
	nfaMachineSyncTestcaseTry struct {
		name  string
		input string
	}
	nfaMachineSyncTestcase struct {
		name         string
		builders     []func(m roughfa.NFAMachine) (roughfa.NFAMachine, error)
		states       []string
		startStates  []string
		acceptStates []string
		transitions  map[string]map[rune][]string
		tries        []*nfaMachineSyncTestcaseTry
	}
)

func (s nfaMachineSyncTestcaseTry) generatetestcase(
	m roughfa.NFAMachine,
	builders []func(roughfa.NFAMachine) (roughfa.NFAMachine, error),
) func(*testing.T) {
	return func(t *testing.T) {
		var (
			ms              = make([]roughfa.NFAMachine, len(builders))
			gotErrors       = make([]bool, len(builders))
			gotAcceptedList = make([]bool, len(builders))
		)
		for i, b := range builders {
			x, err := b(m)
			if !assert.Nil(t, err, "builders[%d]", i) {
				return
			}
			x.Reset()
			ms[i] = x
		}
		for i, m := range ms {
			b, _ := m.ToShell().ToJSON()
			t.Logf("%s", b)
			for _, c := range s.input {
				cs := m.States()
				err := m.Put(c)
				t.Logf("[%d]\t%v\t=[%q]=>\t%v\t%v\t%v", i, cs, c, m.States(), m.IsAccepted(), err)
				if !gotErrors[i] {
					gotErrors[i] = err != nil
				}
			}
		}
		for i, m := range ms {
			gotAcceptedList[i] = m.IsAccepted()
		}
		t.Logf("gotAcceptedList: %v", gotAcceptedList)
		t.Logf("gotErrors: %v", gotErrors)
		for i := 1; i < len(gotErrors); i++ {
			if !assert.Equal(t, gotErrors[0], gotErrors[i]) {
				break
			}
		}
		for i := 1; i < len(ms); i++ {
			if !assert.Equal(t, gotAcceptedList[0], gotAcceptedList[i]) {
				break
			}
		}
	}
}

func (s nfaMachineSyncTestcase) test(t *testing.T) {
	if !assert.True(t, len(s.builders) > 1, "requires 2 builders at least") {
		return
	}
	m, err := roughfa.NewNFAMachineBuilder().
		States(s.states).
		StartStates(s.startStates).
		AcceptStates(s.acceptStates).
		Transitions(s.transitions).
		Build()
	assert.Nil(t, err)
	for _, tr := range s.tries {
		t.Run(tr.name, tr.generatetestcase(m, s.builders))
	}
}

func TestSyncMachine(t *testing.T) {
	for _, tc := range []*nfaMachineSyncTestcase{
		{
			name: "a(b|c)*d",
			builders: []func(roughfa.NFAMachine) (roughfa.NFAMachine, error){
				func(m roughfa.NFAMachine) (roughfa.NFAMachine, error) { return m, nil },
				func(m roughfa.NFAMachine) (roughfa.NFAMachine, error) { return m.ApplyEpsilonExpansion(), nil },
				func(m roughfa.NFAMachine) (roughfa.NFAMachine, error) {
					return m.ApplyEpsilonExpansion().ApplyPowersetConstruction()
				},
				func(m roughfa.NFAMachine) (roughfa.NFAMachine, error) { return m.Minimize() },
			},
			states: []string{
				"a-start", "a-end",
				"bc-start",
				"b-start", "b-end",
				"c-start", "c-end",
				"bc-end",
				"d-start",
				"d-end",
			},
			startStates:  []string{"a-start"},
			acceptStates: []string{"d-end"},
			transitions: map[string]map[rune][]string{
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
			},
			tries: []*nfaMachineSyncTestcaseTry{
				{
					name: "no input",
				},
				{
					name:  "a",
					input: "a",
				},
				{
					name:  "ab",
					input: "ab",
				},
				{
					name:  "abb",
					input: "abb",
				},
				{
					name:  "ac",
					input: "ac",
				},
				{
					name:  "ad",
					input: "ad",
				},
				{
					name:  "abcd",
					input: "abcd",
				},
				{
					name:  "x",
					input: "x",
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

type (
	nfaMachineTestcaseTry struct {
		name         string
		input        string
		wantStates   []string
		wantAccepted bool
		wantError    bool
	}
	nfaMachineTestcase struct {
		name         string
		states       []string
		startStates  []string
		acceptStates []string
		transitions  map[string]map[rune][]string
		tries        []*nfaMachineTestcaseTry
	}
)

func (s nfaMachineTestcaseTry) generateTestcase(m roughfa.NFAMachine) func(*testing.T) {
	return func(t *testing.T) {
		m.Reset()
		var (
			isError bool
			bs      []string
			as      []string
		)
		for _, c := range s.input {
			bs = m.States()
			err := m.Put(c)
			if !isError {
				isError = err != nil
			}
			as = m.States()
			sort.Strings(bs)
			sort.Strings(as)
			t.Logf("%v\t={%q}=>\t%v\t%v", bs, c, as, err)
		}
		assert.Equal(t, s.wantError, isError)
		sort.Strings(s.wantStates)
		assert.True(
			t,
			set.NewStringSet(s.wantStates...).Equal(set.NewStringSet(m.States()...)),
			fmt.Sprintf("%v != %v", s.wantStates, as),
		)
		assert.Equal(t, s.wantAccepted, m.IsAccepted())
	}
}

func (s nfaMachineTestcase) test(t *testing.T) {
	m, err := roughfa.NewNFAMachineBuilder().
		States(s.states).
		StartStates(s.startStates).
		AcceptStates(s.acceptStates).
		Transitions(s.transitions).
		Build()
	assert.Nil(t, err)
	for _, tr := range s.tries {
		t.Run(tr.name, tr.generateTestcase(m))
	}
}

func TestNFAMachine(t *testing.T) {
	for _, tc := range []*nfaMachineTestcase{
		{
			name: "a(b|c)*d",
			states: []string{
				"a-start", "a-end",
				"bc-start",
				"b-start", "b-end",
				"c-start", "c-end",
				"bc-end",
				"d-start",
				"d-end",
			},
			startStates:  []string{"a-start"},
			acceptStates: []string{"d-end"},
			transitions: map[string]map[rune][]string{
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
			},
			tries: []*nfaMachineTestcaseTry{
				{
					name:       "no input",
					wantStates: []string{"a-start"},
				},
				{
					name:       "a",
					input:      "a",
					wantStates: []string{"b-start", "c-start", "d-start"},
				},
				{
					name:       "ab",
					input:      "ab",
					wantStates: []string{"b-start", "c-start", "d-start"},
				},
				{
					name:       "abb",
					input:      "abb",
					wantStates: []string{"b-start", "c-start", "d-start"},
				},
				{
					name:       "ac",
					input:      "ac",
					wantStates: []string{"b-start", "c-start", "d-start"},
				},
				{
					name:         "ad",
					input:        "ad",
					wantStates:   []string{"d-end"},
					wantAccepted: true,
				},
				{
					name:         "abcd",
					input:        "abcd",
					wantStates:   []string{"d-end"},
					wantAccepted: true,
				},
				{
					name:       "x",
					input:      "x",
					wantStates: []string{},
					wantError:  true,
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}
