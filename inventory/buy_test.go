package inventory

import (
	"errors"
	"testing"
)

func TestBuyNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewBuyOrder("A", 1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		nil,
	)
	//This case show that buying a negative qty generates an error
	_TestOrderNextState(t,
		NewBuyOrder("A", -1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		errors.New("Negative Buy Attempted: updateBuy A -1"),
	)
	//This case shows that buying a non-existent items generates an error
	_TestOrderNextState(t,
		NewBuyOrder("B", 1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		errors.New("Buying Non-Existent Item Attempted: updateBuy B 1"),
	)
	//This case shows that the negative buy error trumps the non-existent buy error
	_TestOrderNextState(t,
		NewBuyOrder("B", -1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		errors.New("Negative Buy Attempted: updateBuy B -1"),
	)
}
