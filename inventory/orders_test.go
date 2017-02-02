package inventory

import (
	"errors"
	"testing"
)

func TestIdentityRenderEntry(t *testing.T) {
	i := NewIdentityOrder()
	if i.RenderEntry() != "identity" {
		t.Errorf("The Message Renders Wrong:%v %v", i, i.RenderEntry())
	}
}

func TestBuyOrderRenderEntry(t *testing.T) {
	_TestOrderRenderEntry(t, NewBuyOrder("Test", 1), "updateBuy Test 1")
}

func TestCreateOrderRenderEntry(t *testing.T) {
	_TestOrderRenderEntry(t, NewCreate("Test", 1, 1), "create Test 0.01 0.01")
}

func _TestOrderRenderEntry(t *testing.T, order StateEntry, expected string) {
	if order.RenderEntry() != expected {
		t.Errorf("The Message Renders Wrong:%v %v", order, order.RenderEntry())
	}
}

func TestIdentityNextState(t *testing.T) {
	_TestOrderNextState(t,
		NewIdentityOrder(),
		map[string]Item{},
		map[string]Item{},
		0,
		nil,
	)
	_TestOrderNextState(t,
		NewIdentityOrder(),
		map[string]Item{"A": Item{}},
		map[string]Item{"A": Item{}},
		0,
		nil,
	)
}

func TestCreateNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewCreate("A", 1, 1),
		map[string]Item{},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		nil,
	)
	//This tests the case where an item is being inserted to an empty state.  Notice that the buy & sell price are NOT updated
	_TestOrderNextState(t,
		NewCreate("A", 2, 2),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		errors.New("Cannot create previously existing item create A 0.02 0.02"),
	)
}

func TestDeleteNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewDelete("A"),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{},
		0,
		nil,
	)
	//This tests the case where an item is being inserted to an empty state.  Notice that the buy & sell price are NOT updated
	_TestOrderNextState(t,
		NewDelete("A"),
		map[string]Item{},
		map[string]Item{},
		0,
		errors.New("Cannot delete non-existent item delete A"),
	)
}

func TestBuyNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewBuyOrder("A", 1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}},
		-1,
		nil,
	)
	_TestOrderNextState(t,
		NewBuyOrder("A", -1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		errors.New("Negative Buy Attempted: updateBuy A -1"),
	)
	//This case shows that buying a non-existent items generates an error
	_TestOrderNextState(t,
		NewBuyOrder("B", 1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		errors.New("Buying Non-Existent Item Attempted: updateBuy B 1"),
	)
	//This case shows that the negative buy error trumps the non-existent buy error
	_TestOrderNextState(t,
		NewBuyOrder("B", -1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		errors.New("Negative Buy Attempted: updateBuy B -1"),
	)
}

func TestSellNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewSellOrder("A", 1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		1,
		nil,
	)
	//This case shows that selling a negative number of items generates an error
	_TestOrderNextState(t,
		NewSellOrder("A", -1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		errors.New("Negative Sell Attempted: updateSell A -1"),
	)
	//This case shows that selling non-existent items raises an error
	_TestOrderNextState(t,
		NewSellOrder("B", 1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		errors.New("Selling Non-Existent Item Attempted: updateSell B 1"),
	)
	//This case shows that the negative buy error trumps the non-existent buy error
	_TestOrderNextState(t,
		NewSellOrder("B", -1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		errors.New("Negative Sell Attempted: updateSell B -1"),
	)
	//This case shows that selling too many of an item generates an error
	_TestOrderNextState(t,
		NewSellOrder("A", 1),
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}},
		0,
		errors.New("Selling too Many Units Attempted: updateSell A 1"),
	)
}

/*
This is a worker function that does a lot of the heavy lifting for the testing.  It takes in the input state, and the next step, and has expect output state, delta, and error
*/
func _TestOrderNextState(t *testing.T, order StateEntry, input, expectedAccum map[string]Item, expectedDelta int, expectedError error) {
	actualAccum, actualDelta, actualError := order.NextState(input)
	if actualDelta != expectedDelta {
		t.Errorf("Expected Delta:%d,\tActual Delta%d", expectedDelta, actualDelta)
	}
	for actualKey, actualValue := range actualAccum {
		expectedValue, ok := expectedAccum[actualKey]
		if ok == false {
			t.Errorf("key %v not found in expected :(", actualKey)
		}
		if actualValue != expectedValue {
			t.Errorf("Actual Value %#v\tExpected Value %#v", actualValue, expectedValue)
		}
	}
	//This error message checking is a bit of a hack, but good enough to provide some level of testing assurance.
	actualMessage, expectedMessage := "", ""
	if actualError != nil {
		actualMessage = actualError.Error()
	}
	if expectedError != nil {
		expectedMessage = expectedError.Error()
	}
	if actualMessage != expectedMessage {
		t.Errorf("Actual Error %#v\t,Expected Error%#v", actualError, expectedError)
	}
}

func TestRenderCurrency(t *testing.T) {
	_TestRenderCurrency(t, 1, "0.01")
	_TestRenderCurrency(t, 0, "0.00")
	_TestRenderCurrency(t, 100, "1.00")
}

func _TestRenderCurrency(t *testing.T, currency int, expected string) {
	actual := RenderCurrency(currency)
	if actual != expected {
		t.Errorf("Currency:\t%d,\tActual String:\t%s,\tExpected String:\t%s", currency, actual, expected)
	}
}
