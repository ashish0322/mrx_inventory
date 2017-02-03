package inventory

import (
	"errors"
	"testing"
)

func TestCompoundNextState(t *testing.T) {
	//Test an empty compound statement
	_TestOrderNextState(t,
		NewCompound(),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		nil,
	)
	//Test Single Compound Statement
	_TestOrderNextState(t,
		NewCompound(
			NewBuyOrder("A", 2),
		),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 3}}, Cost: 2},
		nil,
	)
	//Test multiple compound statement
	_TestOrderNextState(t,
		NewCompound(
			NewBuyOrder("A", 2),
			NewBuyOrder("A", 2),
		),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 5}}, Cost: 4},
		nil,
	)
	//Test error ejection
	_TestOrderNextState(t,
		NewCompound(
			NewBuyOrder("A", 2),
			NewBuyOrder("A", -2),
		),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		//State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 5}}, Cost: 4},
		errors.New("Negative Buy Attempted: updateBuy A -2"),
	)
}
