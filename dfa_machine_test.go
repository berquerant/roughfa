package roughfa_test

import (
	"fmt"
	"testing"

	"github.com/berquerant/roughfa"
	"github.com/stretchr/testify/assert"
)

func ExampleDFAMachineRun() {
	m, err := roughfa.NewDFAMachineBuilder().
		States([]string{"even", "odd"}).
		StartState("even").
		AcceptStates([]string{"odd"}).
		Transitions(map[string]map[rune]string{
			"even": {
				'1': "odd",
				'0': "even",
			},
			"odd": {
				'1': "even",
				'0': "odd",
			},
		}).
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(m.State())
	m.Put('1')
	fmt.Println(m.State())
	m.Put('0')
	fmt.Println(m.State())
	fmt.Println(m.IsAccepted())
	m.Put('1')
	fmt.Println(m.State())
	fmt.Println(m.IsAccepted())
	// Output:
	// even
	// odd
	// odd
	// true
	// even
	// false
}

type (
	dfaMachineTestcaseTry struct {
		name         string
		input        string
		wantState    string
		wantAccepted bool
		wantError    bool
	}
	dfaMachineTestcase struct {
		name         string
		states       []string
		startState   string
		acceptStates []string
		transitions  map[string]map[rune]string
		tries        []dfaMachineTestcaseTry
	}
)

func (s dfaMachineTestcaseTry) generateTest(m roughfa.DFAMachine) func(*testing.T) {
	return func(t *testing.T) {
		m.Reset()
		var isError bool
		for _, c := range s.input {
			bs := m.State()
			err := m.Put(c)
			if !isError {
				isError = err != nil
			}
			t.Logf("%s\t={%q}=>\t%s\t%v", bs, c, m.State(), err)
		}
		assert.Equal(t, s.wantError, isError)
		if isError {
			return
		}
		assert.Equal(t, s.wantState, m.State())
		assert.Equal(t, s.wantAccepted, m.IsAccepted())
	}
}

func (s dfaMachineTestcase) test(t *testing.T) {
	m, err := roughfa.NewDFAMachineBuilder().
		States(s.states).
		StartState(s.startState).
		AcceptStates(s.acceptStates).
		Transitions(s.transitions).
		Build()
	if !assert.Nil(t, err) {
		return
	}
	for _, tr := range s.tries {
		t.Run(tr.name, tr.generateTest(m))
	}
}

func TestDFAMachine(t *testing.T) {
	for _, tc := range []*dfaMachineTestcase{
		{
			name:         "minimum dfa",
			states:       []string{"s1"},
			startState:   "s1",
			acceptStates: []string{"s1"},
			tries: []dfaMachineTestcaseTry{
				{
					name:         "no input",
					wantState:    "s1",
					wantAccepted: true,
				},
				{
					name:      "a char",
					input:     "a",
					wantError: true,
				},
			},
		},
		{
			name:         "2 states",
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
			tries: []dfaMachineTestcaseTry{
				{
					name:      "no input",
					wantState: "even",
				},
				{
					name:         "toggle",
					input:        "1",
					wantState:    "odd",
					wantAccepted: true,
				},
				{
					name:      "long input",
					input:     "1101100",
					wantState: "even",
				},
				{
					name:      "out of transition",
					input:     "1X",
					wantError: true,
				},
			},
		},
		{
			name:         "two sequntial 1",
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
			tries: []dfaMachineTestcaseTry{
				{
					name:      "no input",
					wantState: "start",
				},
				{
					name:      "1",
					input:     "1",
					wantState: "1",
				},
				{
					name:         "11",
					input:        "11",
					wantState:    "11",
					wantAccepted: true,
				},
				{
					name:      "10",
					input:     "10",
					wantState: "start",
				},
				{
					name:      "110",
					input:     "110",
					wantState: "start",
				},
				{
					name:         "111",
					input:        "111",
					wantState:    "11",
					wantAccepted: true,
				},
				{
					name:         "001101011",
					input:        "001101011",
					wantState:    "11",
					wantAccepted: true,
				},
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}
