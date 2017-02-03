package inventory

type Item struct {
	Qty       int
	BuyPrice  int
	SellPrice int
}

func NewItem(BuyPrice, SellPrice int) Item {
	return Item{Qty: 0, BuyPrice: BuyPrice, SellPrice: SellPrice}
}

type StateEntry interface {
	NextState(accum State) (State, error)
	RenderEntry() string
}

type State struct {
	Items   map[string]Item
	Revenue int
	Cost    int
}
