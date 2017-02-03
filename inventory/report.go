package inventory

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
	return "report"
}

func NewReport(reportBus chan State) StateEntry {
	return &Report{ReportBus: reportBus}
}
