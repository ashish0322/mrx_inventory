package inventory

func NewRequest(body string) RequestPacket {
	return RequestPacket{Body: body, Response: make(chan string)}
}

type RequestPacket struct {
	Body     string
	Response chan string
}

func ProcessInputs(packetBus chan RequestPacket, reportHandler func(State)) {
	state := State{Items: map[string]Item{}}
	reportBus := make(chan State, 10)
	for {
		select {
		case packet := <-packetBus:
			entry, err := ParseLine(packet.Body, reportBus)
			if err != nil {
				packet.Response <- err.Error()
				continue
			}
			state, err = entry.NextState(state)
			if err != nil {
				packet.Response <- err.Error()
				continue
			}
			packet.Response <- "OK"
		case s := <-reportBus:
			reportHandler(s)
		}

	}
}
