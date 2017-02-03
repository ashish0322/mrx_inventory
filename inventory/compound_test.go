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
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 3}}},
		nil,
	)
	//Test multiple compound statement
	_TestOrderNextState(t,
		NewCompound(
			NewBuyOrder("A", 2),
			NewBuyOrder("A", 2),
		),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 5}}},
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
		errors.New("Negative Buy Attempted: updateBuy A -2"),
	)
}

func TestBarclayExample(t *testing.T) {
	//Test the first portion of the Barclay's README example
	_TestOrderNextState(t,
		NewCompound(
			NewCreate("Book01", 1050, 1379),
			NewCreate("Food01", 147, 398),
			NewCreate("Med01", 3063, 3429),
			NewCreate("Tab01", 5700, 8498),
			NewBuyOrder("Tab01", 100),
			NewSellOrder("Tab01", 2),
			NewBuyOrder("Food01", 500),
			NewBuyOrder("Book01", 100),
			NewBuyOrder("Med01", 100),
			NewSellOrder("Food01", 1),
			NewSellOrder("Food01", 1),
			NewSellOrder("Tab01", 2),
		),
		State{Items: map[string]Item{}},
		State{
			Items: map[string]Item{
				"Book01": Item{Qty: 100, BuyPrice: 1050, SellPrice: 1379},
				"Food01": Item{Qty: 498, BuyPrice: 147, SellPrice: 398},
				"Med01":  Item{Qty: 100, BuyPrice: 3063, SellPrice: 3429},
				"Tab01":  Item{Qty: 96, BuyPrice: 5700, SellPrice: 8498},
			},
			Revenue: 11694,
		},
		nil,
	)
	//Test the second portion of the Barclay's README example, post report
	//I really don't feel like repeating the conccurency testing setup...
	_TestOrderNextState(t,
		NewCompound(
			NewDelete("Book01"),
			NewSellOrder("Tab01", 5),
			NewCreate("Mobile01", 1051, 4456),
			NewBuyOrder("Mobile01", 250),
			NewSellOrder("Food01", 5),
			NewSellOrder("Mobile01", 4),
			NewSellOrder("Med01", 10),
		),
		State{
			//Revenue reset to zero, to simulate the report being run
			Items: map[string]Item{
				"Book01": Item{Qty: 100, BuyPrice: 1050, SellPrice: 1379},
				"Food01": Item{Qty: 498, BuyPrice: 147, SellPrice: 398},
				"Med01":  Item{Qty: 100, BuyPrice: 3063, SellPrice: 3429},
				"Tab01":  Item{Qty: 96, BuyPrice: 5700, SellPrice: 8498},
			},
		},
		State{
			Items: map[string]Item{
				"Food01":   Item{Qty: 493, BuyPrice: 147, SellPrice: 398},
				"Med01":    Item{Qty: 90, BuyPrice: 3063, SellPrice: 3429},
				"Tab01":    Item{Qty: 91, BuyPrice: 5700, SellPrice: 8498},
				"Mobile01": Item{Qty: 246, BuyPrice: 1051, SellPrice: 4456},
			},
			Cost:    105000,
			Revenue: 32525,
		},
		nil,
	)
}
