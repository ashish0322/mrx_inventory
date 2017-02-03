package inventory

import (
	"errors"
	"fmt"
)

/**
* BUY FAMILY
**/
type BuyOrder struct {
	ItemName string
	Quantity int
}

func (this *BuyOrder) NextState(accum State) (State, error) {
	if this.Quantity < 0 {
		return accum, errors.New("Negative Buy Attempted: " + this.RenderEntry())
	}
	if _, ok := accum.Items[this.ItemName]; ok == false {
		return accum, errors.New("Buying Non-Existent Item Attempted: " + this.RenderEntry())
	}
	i := accum.Items[this.ItemName]
	i.Qty += this.Quantity
	accum.Items[this.ItemName] = i
	accum.Cost += this.Quantity * i.BuyPrice
	return accum, nil
}

func (this *BuyOrder) RenderEntry() string {
	return fmt.Sprintf("updateBuy %s %d", this.ItemName, this.Quantity)
}

func NewBuyOrder(ItemName string, Quantity int) StateEntry {
	return &BuyOrder{ItemName: ItemName, Quantity: Quantity}
}
