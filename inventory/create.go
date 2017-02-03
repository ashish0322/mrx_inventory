package inventory

import (
	"errors"
	"fmt"
)

/**
* CREATE FAMILY
**/
type Create struct {
	ItemName  string
	BuyPrice  int
	SellPrice int
}

/**
* The create entry will create a new type of item for sale when it's encountered in the journal.  This method will return an error if it encouners a pre-existing item in the accum table.  It always returns a change or zero.
 */
func (this *Create) NextState(accum State) (State, error) {
	if _, ok := accum.Items[this.ItemName]; ok {
		return accum, errors.New("Cannot create previously existing item " + this.RenderEntry())
	}
	if this.BuyPrice < 0 {
		return accum, errors.New("Cannot set negative buy price " + this.RenderEntry())
	}
	if this.SellPrice < 0 {
		return accum, errors.New("Cannot set negative sell price " + this.RenderEntry())
	}
	accum.Items[this.ItemName] = NewItem(this.BuyPrice, this.SellPrice)
	return accum, nil
}

func (this *Create) RenderEntry() string {
	return fmt.Sprintf("create %s %s %s", this.ItemName, RenderCurrency(this.BuyPrice), RenderCurrency(this.SellPrice))
}

func NewCreate(itemName string, buyPrice, sellPrice int) StateEntry {
	return &Create{ItemName: itemName, BuyPrice: buyPrice, SellPrice: sellPrice}
}
