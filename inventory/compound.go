package inventory

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
