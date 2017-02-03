package inventory

import (
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
