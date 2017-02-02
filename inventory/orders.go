package inventory

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Item struct {
	Qty       int
	BuyPrice  int
	SellPrice int
}

func NewItem(BuyPrice, SellPrice int) Item {
	return Item{Qty: 0, BuyPrice: BuyPrice, SellPrice: SellPrice}
}

/**
A utitlity function for pretty-printing the currency
*/
func RenderCurrency(currency int) string {
	isNegative := currency < 0
	if isNegative {
		currency = -currency
	}
	output := fmt.Sprintf("%d.%02d", currency/100, currency%100)
	if isNegative {
		output = "-" + output
	}
	return output
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

/**
* IDENTITY FAMILY - For Testing
**/
type IdentityOrder struct{}

func (this *IdentityOrder) NextState(accum State) (State, error) {
	return accum, nil
}

func (this *IdentityOrder) RenderEntry() string {
	return "identity"
}

func NewIdentityOrder() StateEntry {
	return new(IdentityOrder)
}

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
	accum.Revenue += this.Quantity * i.SellPrice
	return accum, nil
}

func (this *SellOrder) RenderEntry() string {
	return fmt.Sprintf("updateSell %s %d", this.ItemName, this.Quantity)
}

func NewSellOrder(itemName string, quantity int) StateEntry {
	return &SellOrder{ItemName: itemName, Quantity: quantity}
}

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

/**
* DELETE FAMILY
**/
type Delete struct {
	ItemName string
}

/**
* The delete entry will delete an existing item for sale when it's encountered in the journal.  This method will return an error if it attempts to delete an item that doesn't exist.
 */
func (this *Delete) NextState(accum State) (State, error) {
	if _, ok := accum.Items[this.ItemName]; !ok {
		return accum, errors.New("Cannot delete non-existent item " + this.RenderEntry())
	}
	delete(accum.Items, this.ItemName)
	return accum, nil
}

func (this *Delete) RenderEntry() string {
	return fmt.Sprintf("delete %s", this.ItemName)
}

func NewDelete(itemName string) StateEntry {
	return &Delete{ItemName: itemName}
}

/**
* REPORT FAMILY
**/
type Report struct {
	ReportBus chan State
}

func (this *Report) NextState(accum State) (State, error) {
	output := State{Items: map[string]Item{}}
	for key, value := range accum.Items {
		output.Items[key] = value
	}
	this.ReportBus <- accum
	return output, nil
}

func (this *Report) RenderEntry() string {
	return fmt.Sprintf("report")
}

func NewReport(reportBus chan State) StateEntry {
	return &Report{ReportBus: reportBus}
}

/**
* COMPOUND FAMILY
**/
type Compound struct {
	Steps []StateEntry
}

func (this *Compound) NextState(accum State) (State, error) {
	backupAccum := State{
		Items:   map[string]Item{},
		Revenue: accum.Revenue,
		Cost:    accum.Cost,
	}
	for key, value := range accum.Items {
		backupAccum.Items[key] = value
	}
	nextAccum := accum
	var err error = nil

	for _, entry := range this.Steps {
		nextAccum, err = entry.NextState(nextAccum)
		if err != nil {
			return backupAccum, err
		}
	}
	return nextAccum, err
}

func (this *Compound) RenderEntry() string {
	output := "compound"
	for _, entry := range this.Steps {
		output += "\n" + entry.RenderEntry()
	}
	return output
}

func NewCompound(steps ...StateEntry) StateEntry {
	return &Compound{Steps: steps}
}

/**
* Parser
 */

func ParseLine(line string, reportBus chan State) (StateEntry, error) {
	re, _ := regexp.Compile("\\s+")
	tokens := re.Split(strings.TrimSpace(line), -1)
	switch tokens[0] {
	case "report":
		if len(tokens) != 1 {
			return nil, errors.New("The " + tokens[0] + " statement is malformed: " + line)
		}
		return NewReport(reportBus), nil
	case "updateBuy":
		if len(tokens) != 3 {
			return nil, errors.New("The " + tokens[0] + " statement is malformed: " + line)
		}
		qty, err := strconv.ParseInt(tokens[2], 10, 0)
		if err != nil {
			return nil, errors.New("Error reading the quantity: " + line)
		}
		return NewBuyOrder(tokens[1], int(qty)), nil
	case "updateSell":
		if len(tokens) != 3 {
			return nil, errors.New("The " + tokens[0] + " statement is malformed: " + line)
		}
		qty, err := strconv.ParseInt(tokens[2], 10, 0)
		if err != nil {
			return nil, errors.New("Error reading the quantity: " + line)
		}
		return NewSellOrder(tokens[1], int(qty)), nil
	case "create":
		if len(tokens) != 4 {
			return nil, errors.New("The " + tokens[0] + " statement is malformed: " + line)
		}
		buyPrice, err := ParseCurrency(tokens[2])
		if err != nil {
			return nil, errors.New("Error reading the currency: " + line)
		}
		sellPrice, err := ParseCurrency(tokens[3])
		if err != nil {
			return nil, errors.New("Error reading the currency: " + line)
		}
		return NewCreate(tokens[1], buyPrice, sellPrice), nil
	case "delete":
		if len(tokens) != 2 {
			return nil, errors.New("The " + tokens[0] + " statement is malformed: " + line)
		}
		return NewDelete(tokens[1]), nil
	}

	return nil, errors.New("No valid input detected")
}

func ParseCurrency(currency string) (int, error) {
	rightFormat, _ := regexp.Compile("\\-?[0-9]+\\.[0-9]{2}")
	match := rightFormat.MatchString(currency)
	if !match {
		return 0, errors.New("Currency is Malformed: " + currency)
	}
	re, _ := regexp.Compile("\\.")
	tokens := re.Split(currency, -1)

	//This check is required to catch things like -0.50
	isNegative := tokens[0][0] == byte('-')
	//The format is protected by the earlier regex
	dollarQty, _ := strconv.ParseInt(tokens[0], 10, 0)
	centsQty, _ := strconv.ParseInt(tokens[1], 10, 0)

	if isNegative {
		return int(100*dollarQty - centsQty), nil
	} else {
		return int(100*dollarQty + centsQty), nil
	}
}
