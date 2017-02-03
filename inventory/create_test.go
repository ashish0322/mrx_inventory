package inventory

import (
	"errors"
	"testing"
)

func TestCreateNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewCreate("A", 1, 1),
		State{Items: map[string]Item{}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		nil,
	)
	//This tests the case where an item is being inserted to an empty state.  Notice that the buy & sell price are NOT updated
	_TestOrderNextState(t,
		NewCreate("A", 2, 2),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		errors.New("Cannot create previously existing item create A 0.02 0.02"),
	)
	//This tests the case where a new item is being inserted to an empty state, but with a negative buy price.  It is rejected
	_TestOrderNextState(t,
		NewCreate("A", -1, 1),
		State{Items: map[string]Item{}},
		State{Items: map[string]Item{}},
		errors.New("Cannot set negative buy price create A -0.01 0.01"),
	)
	//This tests the case where a new item is being inserted to an empty state, but with a negative sell price. It is rejected
	_TestOrderNextState(t,
		NewCreate("A", 1, -1),
		State{Items: map[string]Item{}},
		State{Items: map[string]Item{}},
		errors.New("Cannot set negative sell price create A 0.01 -0.01"),
	)
}
