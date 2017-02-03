package inventory

import (
	"errors"
	"testing"
)

func TestDeleteNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewDelete("A"),
		State{Items: map[string]Item{"A": Item{Qty: 1, BuyPrice: 2, SellPrice: 1}}},
		State{Items: map[string]Item{}, Cost: 2},
		nil,
	)
	//This tests the case where an item is being inserted to an empty state.  Notice that the buy & sell price are NOT updated
	_TestOrderNextState(t,
		NewDelete("A"),
		State{Items: map[string]Item{}},
		State{Items: map[string]Item{}},
		errors.New("Cannot delete non-existent item delete A"),
	)
}
