package inventory

import (
	"fmt"
	"sort"
)

//This controls the line format
const LINE_FORMAT_STRING string = "%-21v %10v %10v %10v %15v"

func RenderState(state State) string {
	keys := []string{}
	for key, _ := range state.Items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	output := ""
	output += fmt.Sprintf(LINE_FORMAT_STRING, "Name", "Quantity", "Bought At", "Sold At", "Bought Value") + "\n"
	output += fmt.Sprintf(LINE_FORMAT_STRING, "-----", "-----", "-----", "-----", "-----") + "\n"
	totalBoughtValue := 0
	for _, key := range keys {
		item := state.Items[key]
		output += RenderItem(key, item) + "\n"
		totalBoughtValue += item.Qty * item.BuyPrice
	}
	output += fmt.Sprintf(LINE_FORMAT_STRING, "-----", "-----", "-----", "-----", "-----") + "\n"
	output += RenderAggregate("Inventory Value", totalBoughtValue) + "\n"
	output += RenderAggregate("Revenue Since Last Report", state.Revenue) + "\n"
	output += RenderAggregate("Cost Since Last Report", state.Cost) + "\n"
	output += RenderAggregate("Profit Since Last Report", state.Revenue-state.Cost) + "\n"
	return output
}

func RenderAggregate(key string, value int) string {
	return fmt.Sprintf("%30v%-25v%15v", "", key, RenderCurrency(value))
}

func RenderItem(name string, item Item) string {
	output := fmt.Sprintf(LINE_FORMAT_STRING, name, item.Qty, RenderCurrency(item.BuyPrice), RenderCurrency(item.SellPrice), RenderCurrency(item.Qty*item.BuyPrice))
	return output
}
