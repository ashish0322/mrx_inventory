package inventory

import (
	"testing"
)

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
