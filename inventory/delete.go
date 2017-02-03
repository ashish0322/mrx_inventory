package inventory

import (
	"errors"
	"fmt"
)

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
