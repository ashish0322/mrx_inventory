package inventory

import (
	"errors"
	"testing"
)

func TestParseLine(t *testing.T) {
	reportBus := make(chan State)
	_TestParseLine(
		t,
		"",
		reportBus,
		nil,
		errors.New("No valid input detected"),
	)
	//Report Checks
	_TestParseLine(
		t,
		"report",
		reportBus,
		NewReport(reportBus),
		nil,
	)
	_TestParseLine(
		t,
		"report Bacon",
		reportBus,
		nil,
		errors.New("The report statement is malformed: report Bacon"),
	)
	//Delete Checks
	_TestParseLine(
		t,
		"delete Bacon",
		reportBus,
		NewDelete("Bacon"),
		nil,
	)
	_TestParseLine(
		t,
		"delete",
		reportBus,
		nil,
		errors.New("The delete statement is malformed: delete"),
	)
	_TestParseLine(
		t,
		"delete Bacon Ninja",
		reportBus,
		nil,
		errors.New("The delete statement is malformed: delete Bacon Ninja"),
	)
	//Create Checks
	_TestParseLine(
		t,
		"create Bacon 0.00 0.00",
		reportBus,
		NewCreate("Bacon", 0, 0),
		nil,
	)
	_TestParseLine(
		t,
		"create Bacon 1.00 0.00",
		reportBus,
		NewCreate("Bacon", 100, 0),
		nil,
	)
	_TestParseLine(
		t,
		"create Bacon 0.00 1.00",
		reportBus,
		NewCreate("Bacon", 0, 100),
		nil,
	)
	_TestParseLine(
		t,
		"create Bacon 0.00",
		reportBus,
		nil,
		errors.New("The create statement is malformed: create Bacon 0.00"),
	)
	_TestParseLine(
		t,
		"create Bacon 0.00 0.00 0.00",
		reportBus,
		nil,
		errors.New("The create statement is malformed: create Bacon 0.00 0.00 0.00"),
	)
	_TestParseLine(
		t,
		"create Bacon Ninja 0.00",
		reportBus,
		nil,
		errors.New("Error reading the currency: create Bacon Ninja 0.00"),
	)
	_TestParseLine(
		t,
		"create Bacon 0.00 Ninja",
		reportBus,
		nil,
		errors.New("Error reading the currency: create Bacon 0.00 Ninja"),
	)
	//UpdateBuy Checks
	_TestParseLine(
		t,
		"updateBuy Bacon 1",
		reportBus,
		NewBuyOrder("Bacon", 1),
		nil,
	)
	_TestParseLine(
		t,
		"updateBuy Bacon",
		reportBus,
		nil,
		errors.New("The updateBuy statement is malformed: updateBuy Bacon"),
	)
	_TestParseLine(
		t,
		"updateBuy Bacon 1 1",
		reportBus,
		nil,
		errors.New("The updateBuy statement is malformed: updateBuy Bacon 1 1"),
	)
	_TestParseLine(
		t,
		"updateBuy Bacon Ninja",
		reportBus,
		nil,
		errors.New("Error reading the quantity: updateBuy Bacon Ninja"),
	)
	//UpdateSell Checks
	_TestParseLine(
		t,
		"updateSell Bacon 1",
		reportBus,
		NewSellOrder("Bacon", 1),
		nil,
	)
	_TestParseLine(
		t,
		"updateSell Bacon",
		reportBus,
		nil,
		errors.New("The updateSell statement is malformed: updateSell Bacon"),
	)
	_TestParseLine(
		t,
		"updateSell Bacon 1 1",
		reportBus,
		nil,
		errors.New("The updateSell statement is malformed: updateSell Bacon 1 1"),
	)
	_TestParseLine(
		t,
		"updateSell Bacon Ninja",
		reportBus,
		nil,
		errors.New("Error reading the quantity: updateSell Bacon Ninja"),
	)
}

func _TestParseLine(t *testing.T, line string, reportBus chan State, expectedEntry StateEntry, expectedError error) {
	actualEntry, actualError := ParseLine(line, reportBus)
	_CompareErrorMessage(t, expectedError, actualError)
	if actualEntry == expectedEntry {
		return
	}
	//We've compared the erros, the values are equal, we're good
	if actualEntry == nil {
		return
	}

	switch castActual := actualEntry.(type) {
	case *Report:
		castExpected, _ := expectedEntry.(*Report)
		if *castActual != *castExpected {
			t.Errorf("Uh-oh, Actual\t%#v,\tExpected\t%#v", actualEntry, expectedEntry)
		}
	case *Delete:
		castExpected, _ := expectedEntry.(*Delete)
		if *castActual != *castExpected {
			t.Errorf("Uh-oh, Actual\t%#v,\tExpected\t%#v", actualEntry, expectedEntry)
		}
	case *Create:
		castExpected, _ := expectedEntry.(*Create)
		if *castActual != *castExpected {
			t.Errorf("Uh-oh, Actual\t%#v,\tExpected\t%#v", actualEntry, expectedEntry)
		}
	case *BuyOrder:
		castExpected, _ := expectedEntry.(*BuyOrder)
		if *castActual != *castExpected {
			t.Errorf("Uh-oh, Actual\t%#v,\tExpected\t%#v", actualEntry, expectedEntry)
		}
	case *SellOrder:
		castExpected, _ := expectedEntry.(*SellOrder)
		if *castActual != *castExpected {
			t.Errorf("Uh-oh, Actual\t%#v,\tExpected\t%#v", actualEntry, expectedEntry)
		}
	}

}

func TestParseCurrency(t *testing.T) {
	_TestParseCurrency(
		t,
		"0.00",
		0,
		nil,
	)
	_TestParseCurrency(
		t,
		"-0.00",
		0,
		nil,
	)
	_TestParseCurrency(
		t,
		"1.00",
		100,
		nil,
	)
	_TestParseCurrency(
		t,
		"-1.00",
		-100,
		nil,
	)
	_TestParseCurrency(
		t,
		"1.50",
		150,
		nil,
	)
	_TestParseCurrency(
		t,
		"-1.50",
		-150,
		nil,
	)
	_TestParseCurrency(
		t,
		"0.50",
		50,
		nil,
	)
	_TestParseCurrency(
		t,
		"-0.50",
		-50,
		nil,
	)
	_TestParseCurrency(
		t,
		"50",
		0,
		errors.New("Currency is Malformed: 50"),
	)
	_TestParseCurrency(
		t,
		"Bacon",
		0,
		errors.New("Currency is Malformed: Bacon"),
	)
}

func _TestParseCurrency(t *testing.T, currency string, expectedInt int, expectedError error) {
	actualInt, actualError := ParseCurrency(currency)
	if actualInt != expectedInt {
		t.Errorf("Expected currency %v, Actual Currency %v", expectedInt, actualInt)
	}
	_CompareErrorMessage(t, expectedError, actualError)
}
