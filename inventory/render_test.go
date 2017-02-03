package inventory

import (
	"testing"
)

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

func TestRenderState(t *testing.T) {
	s := State{Items: map[string]Item{"A": NewItem(1, 1)}}
	expected := `Name                    Quantity  Bought At    Sold At    Bought Value
-----                      -----      -----      -----           -----
A                              0       0.01       0.01            0.00
-----                      -----      -----      -----           -----
                              Inventory Value                     0.00
                              Revenue Since Last Report           0.00
                              Cost Since Last Report              0.00
                              Profit Since Last Report            0.00
`
	actual := RenderState(s)
	if expected != actual {
		t.Errorf("Expected Report:\n%v\n,Actual Report:\n%v\n", expected, actual)
	}
}
