package state_test

import (
	"testing"

	"github.com/moynur/gateway/internal/state"
)

var testCases = []struct {
	description  string
	currentState state.State
	targetState  state.State
	expected     bool
}{
	{
		"Auth cant follow Auth",
		state.Auth,
		state.Auth,
		false,
	},
	{
		"Capture can follow Auth",
		state.Auth,
		state.Capture,
		true,
	},
	{
		"Refund cant follow Auth",
		state.Auth,
		state.Refund,
		false,
	},
	{
		"Void can follow Auth",
		state.Auth,
		state.Void,
		true,
	},
	{
		"Void cant follow Void",
		state.Void,
		state.Void,
		false,
	},
	{
		"Capture cant follow Void",
		state.Void,
		state.Capture,
		false,
	},
	{
		"Refund cant follow Void",
		state.Auth,
		state.Auth,
		false,
	},
	{
		"Auth cant follow Void",
		state.Void,
		state.Auth,
		false,
	},
	{
		"Capture can follow Capture",
		state.Capture,
		state.Capture,
		true,
	},
	{
		"Refund can follow Capture",
		state.Capture,
		state.Refund,
		true,
	},
	{
		"Auth cant follow Capture",
		state.Capture,
		state.Auth,
		false,
	},
	{
		"Void cant follow Capture",
		state.Capture,
		state.Void,
		false,
	},
	{
		"Refund can follow Refund",
		state.Refund,
		state.Refund,
		true,
	},
	{
		"Capture cant follow Refund",
		state.Refund,
		state.Capture,
		false,
	},
	{
		"Auth cant follow Refund",
		state.Refund,
		state.Auth,
		false,
	},
	{
		"Void cant follow Refund",
		state.Refund,
		state.Void,
		false,
	},
}

func TestState(t *testing.T) {
	t.Run("It should return the right value for all test cases", func(t *testing.T) {
		for _, test := range testCases {
			if actual := state.TransitionAllowed(test.currentState, test.targetState); actual != test.expected {
				t.Fatalf("\n Name: %s \n current state: %s \n Target State: %s \n Expected Outcome: %t \n Actual: %t ", test.description, test.currentState, test.targetState, test.expected, actual)
			}
		}
	})
}
