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

func TestReportOrderRenderEntry(t *testing.T) {
	_TestOrderRenderEntry(t, NewReport(nil), "report")
}

func TestCompoundOrderRenderEntry(t *testing.T) {
	_TestOrderRenderEntry(t, NewCompound(), "compound")
	_TestOrderRenderEntry(t, NewCompound(
		NewBuyOrder("A", 1),
		NewBuyOrder("B", 1),
	), "compound\nupdateBuy A 1\nupdateBuy B 1")
}

func _TestOrderRenderEntry(t *testing.T, order StateEntry, expected string) {
	if order.RenderEntry() != expected {
		t.Errorf("The Message Renders Wrong:%v %v", order, order.RenderEntry())
	}
}

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

func TestDeleteNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewDelete("A"),
		State{Items: map[string]Item{"A": Item{}}},
		State{Items: map[string]Item{}},
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

func TestBuyNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewBuyOrder("A", 1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}, Cost: 1},
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

func TestSellNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	_TestOrderNextState(t,
		NewSellOrder("A", 1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}, Revenue: 1},
		nil,
	)
	//This case show that buying a negative qty generates an error
	_TestOrderNextState(t,
		NewSellOrder("A", -1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		errors.New("Negative Sell Attempted: updateSell A -1"),
	)
	//This case shows that buying a non-existent items generates an error
	_TestOrderNextState(t,
		NewSellOrder("B", 1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		errors.New("Selling Non-Existent Item Attempted: updateSell B 1"),
	)
	//This case shows that the negative buy error trumps the non-existent buy error
	_TestOrderNextState(t,
		NewSellOrder("B", -1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		errors.New("Negative Sell Attempted: updateSell B -1"),
	)
	_TestOrderNextState(t,
		NewSellOrder("A", 1),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1}}},
		errors.New("Selling too Many Units Attempted: updateSell A 1"),
	)
}

func TestReportNextState(t *testing.T) {
	//This tests the case where a new item is being inserted to an empty state
	reportBus := make(chan State)
	//Must launch the task asyncrounously for the channel to work
	go _TestOrderNextState(t,
		NewReport(reportBus),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}, Revenue: 1, Cost: 1},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}},
		nil,
	)
	localState := <-reportBus
	if localState.Revenue != 1 {
		t.Errorf("Unexpected Revenue in recieved state, actual: %v", localState.Revenue)
	}
	if localState.Cost != 1 {
		t.Errorf("Unexpected Cost in recieved state, actual: %v", localState.Cost)
	}
}

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

/**
* This test demonstrates that once the report entry send the object on the bus, other changes can be processed while the report is being handled
 */
func TestReportThreadSafety(t *testing.T) {
	reportBus := make(chan State)
	//Must launch the task asyncrounously for the channel to work
	go _TestOrderNextState(t,
		NewCompound(
			NewReport(reportBus),
			NewBuyOrder("A", 2),
		),
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 1}}, Revenue: 1, Cost: 1},
		State{Items: map[string]Item{"A": Item{BuyPrice: 1, SellPrice: 1, Qty: 3}}, Cost: 2},
		nil,
	)
	localState := <-reportBus
	if localState.Revenue != 1 {
		t.Errorf("Unexpected Revenue in recieved state, actual: %v", localState.Revenue)
	}
	if localState.Cost != 1 {
		t.Errorf("Unexpected Cost in recieved state, actual: %v", localState.Cost)
	}
	if localState.Items["A"].Qty != 1 {
		t.Errorf("Thread Safety Error, message mutated unexpectedly: %v", localState.Items)
	}
}

/*
This is a worker function that does a lot of the heavy lifting for the testing.  It takes in the input state, and the next step, and has expect output state, delta, and error
*/
func _TestOrderNextState(t *testing.T, order StateEntry, input, expectedAccum State, expectedError error) {
	actualAccum, actualError := order.NextState(input)
	if actualAccum.Cost != expectedAccum.Cost {
		t.Errorf("Expected Cost:%d,\tActual Cost%d", expectedAccum.Cost, actualAccum.Cost)
	}
	if actualAccum.Revenue != expectedAccum.Revenue {
		t.Errorf("Expected Revenue:%d,\tActual Revenue%d", expectedAccum.Revenue, actualAccum.Revenue)
	}
	for actualKey, actualValue := range actualAccum.Items {
		expectedValue, ok := expectedAccum.Items[actualKey]
		if ok == false {
			t.Errorf("key %v not found in expected :(", actualKey)
		}
		if actualValue != expectedValue {
			t.Errorf("Actual Value %#v\tExpected Value %#v", actualValue, expectedValue)
		}
	}
	_CompareErrorMessage(t, expectedError, actualError)
}

//This error message checking is a bit of a hack, but good enough to provide some level of testing assurance.
func _CompareErrorMessage(t *testing.T, expectedError, actualError error) {
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
	_TestRenderCurrency(t, -1, "-0.01")
	_TestRenderCurrency(t, -100, "-1.00")
	_TestRenderCurrency(t, -101, "-1.01")
}

func _TestRenderCurrency(t *testing.T, currency int, expected string) {
	actual := RenderCurrency(currency)
	if actual != expected {
		t.Errorf("Currency:\t%d,\tActual String:\t%s,\tExpected String:\t%s", currency, actual, expected)
	}
}
