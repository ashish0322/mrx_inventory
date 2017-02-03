package inventory

import (
	"testing"
)

func TestIdentityNextState(t *testing.T) {
	_TestOrderNextState(t,
		NewIdentityOrder(),
		State{Items: map[string]Item{}},
		State{Items: map[string]Item{}},
		nil,
	)
	_TestOrderNextState(t,
		NewIdentityOrder(),
		State{Items: map[string]Item{"A": Item{}}},
		State{Items: map[string]Item{"A": Item{}}},
		nil,
	)
}
