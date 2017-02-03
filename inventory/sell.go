package inventory

import (
	"errors"
	"fmt"
)

/**
* SELL FAMILY
**/
type SellOrder struct {
	ItemName string
	Quantity int
}

func (this *SellOrder) NextState(accum State) (State, error) {
	if this.Quantity < 0 {
		return accum, errors.New("Negative Sell Attempted: " + this.RenderEntry())
	}
	if _, ok := accum.Items[this.ItemName]; ok == false {
		return accum, errors.New("Selling Non-Existent Item Attempted: " + this.RenderEntry())
	}
	i := accum.Items[this.ItemName]
	if i.Qty < this.Quantity {
		return accum, errors.New("Selling too Many Units Attempted: " + this.RenderEntry())
	}
	i.Qty -= this.Quantity
	accum.Items[this.ItemName] = i
	accum.Revenue += this.Quantity * (i.SellPrice - i.BuyPrice)
	return accum, nil
}

func (this *SellOrder) RenderEntry() string {
	return fmt.Sprintf("updateSell %s %d", this.ItemName, this.Quantity)
}

func NewSellOrder(itemName string, quantity int) StateEntry {
	return &SellOrder{ItemName: itemName, Quantity: quantity}
}
