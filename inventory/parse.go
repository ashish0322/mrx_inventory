package inventory

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

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
