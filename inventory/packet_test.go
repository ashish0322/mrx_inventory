package inventory

import (
	"testing"
)

func TestProcessInputs(t *testing.T) {
	packetBus := make(chan RequestPacket)
	reportChan := make(chan State, 1)
	handler := func(s State) {
		reportChan <- s
	}
	go ProcessInputs(packetBus, handler)
	r := NewRequest("create Bacon 10.00 50.00")
	packetBus <- r
	resp := <-r.Response
	if resp != "OK" {
		t.Errorf("Expected OK, got: %v", resp)
	}

	//Test a controlled parse failure
	r = NewRequest("fail")
	packetBus <- r
	resp = <-r.Response
	if resp != "No valid input detected" {
		t.Errorf("Expected \"No valid input detected\", got: %v", resp)
	}

	//Test a controlled logic failure
	r = NewRequest("delete Pizza")
	packetBus <- r
	resp = <-r.Response
	if resp != "Cannot delete non-existent item delete Pizza" {
		t.Errorf("Expected \"Cannot delete non-existent item delete Pizza\", got: %v", resp)
	}

	//Test a controlled logic failure
	r = NewRequest("report")
	packetBus <- r
	resp = <-r.Response
	if resp != "OK" {
		t.Errorf("Expected OK, got: %v", resp)
	}
	s := <-reportChan
	if s.Revenue != 0 {
		t.Errorf("Wrong State Recieved")
	}
}
