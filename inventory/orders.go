package inventory

import (
	"errors"
	"fmt"
	//"strings"
)

type Item struct {
	Qty       int
	BuyPrice  int
	SellPrice int
}

func NewItem(BuyPrice, SellPrice int) Item {
	return Item{Qty: 0, BuyPrice: BuyPrice, SellPrice: SellPrice}
}

func RenderCurrency(Currency int) string {
	return fmt.Sprintf("%d.%02d", Currency/100, Currency%100)
}

type Command int

type StateEntry interface {
	NextState(accum map[string]Item) (map[string]Item, int, error)
	RenderEntry() string
}

//The Identity order is explicitly included for unit testing
type IdentityOrder struct{}

func (this *IdentityOrder) NextState(accum map[string]Item) (map[string]Item, int, error) {
	return accum, 0, nil
}

func (this *IdentityOrder) RenderEntry() string {
	return "identity"
}

func NewIdentityOrder() StateEntry {
	return new(IdentityOrder)
}

type BuyOrder struct {
	ItemName string
	Quantity int
}

func (this *BuyOrder) NextState(accum map[string]Item) (map[string]Item, int, error) {
	if this.Quantity < 0 {
		return accum, 0, errors.New("Negative Buy Attempted: " + this.RenderEntry())
	}
	if _, ok := accum[this.ItemName]; ok == false {
		return accum, 0, errors.New("Buying Non-Existent Item Attempted: " + this.RenderEntry())
	}
	i := accum[this.ItemName]
	i.Qty += this.Quantity
	accum[this.ItemName] = i
	return accum, -this.Quantity * i.BuyPrice, nil
}

func (this *BuyOrder) RenderEntry() string {
	return fmt.Sprintf("updateBuy %s %d", this.ItemName, this.Quantity)
}

func NewBuyOrder(ItemName string, Quantity int) StateEntry {
	return &BuyOrder{ItemName: ItemName, Quantity: Quantity}
}

type SellOrder struct {
	ItemName string
	Quantity int
}

func (this *SellOrder) NextState(accum map[string]Item) (map[string]Item, int, error) {
	if this.Quantity < 0 {
		return accum, 0, errors.New("Negative Sell Attempted: " + this.RenderEntry())
	}
	if _, ok := accum[this.ItemName]; ok == false {
		return accum, 0, errors.New("Selling Non-Existent Item Attempted: " + this.RenderEntry())
	}
	i := accum[this.ItemName]
	if i.Qty < this.Quantity {
		return accum, 0, errors.New("Selling too Many Units Attempted: " + this.RenderEntry())
	}
	i.Qty -= this.Quantity
	accum[this.ItemName] = i
	return accum, this.Quantity * i.SellPrice, nil
}

func (this *SellOrder) RenderEntry() string {
	return fmt.Sprintf("updateSell %s %d", this.ItemName, this.Quantity)
}

func NewSellOrder(ItemName string, Quantity int) StateEntry {
	return &SellOrder{ItemName: ItemName, Quantity: Quantity}
}

type Create struct {
	ItemName  string
	BuyPrice  int
	SellPrice int
}

func (this *Create) NextState(accum map[string]Item) (map[string]Item, int, error) {
	if _, ok := accum[this.ItemName]; ok {
		return accum, 0, errors.New("Cannot create previously existing item " + this.RenderEntry())
	}
	accum[this.ItemName] = NewItem(this.BuyPrice, this.SellPrice)
	return accum, 0, nil
}

func (this *Create) RenderEntry() string {
	return fmt.Sprintf("create %s %s %s", this.ItemName, RenderCurrency(this.BuyPrice), RenderCurrency(this.SellPrice))
}

func NewCreate(ItemName string, BuyPrice, SellPrice int) StateEntry {
	return &Create{ItemName: ItemName, BuyPrice: BuyPrice, SellPrice: SellPrice}
}

type Delete struct {
	ItemName string
}

func (this *Delete) NextState(accum map[string]Item) (map[string]Item, int, error) {
	if _, ok := accum[this.ItemName]; !ok {
		return accum, 0, errors.New("Cannot delete non-existent item " + this.RenderEntry())
	}
	delete(accum, this.ItemName)
	return accum, 0, nil
}

func (this *Delete) RenderEntry() string {
	return fmt.Sprintf("delete %s", this.ItemName)
}

func NewDelete(ItemName string) StateEntry {
	return &Delete{ItemName: ItemName}
}

func ProcessJournal(accum map[string]Item, entries []StateEntry) {
	revenue, cost := 0, 0
	for _, entry := range entries {
		change := 0
		var err error = nil
		accum, change, err = entry.NextState(accum)
		if err != nil {
			continue
		}
		if change > 0 {
			revenue += change
		} else if change < 0 {
			cost += change
		}
	}
}

//func ProcessEntry(accum map[string]Item, elem Entry) map[string]Item, error{
//switch e := elem.(type) {
//case Order:
//case Manage:
//default:
//return accum, errors.New("Unknown type provided");
//}
//}
