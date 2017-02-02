package main

import (
	"fmt"
	"inventory"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8333"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	packetBus := make(chan RequestPacket)
	go processInputs(packetBus)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, packetBus)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, packetBus chan RequestPacket) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	readLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	r := RequestPacket{
		Body:     string(buf[0:readLen]),
		Response: make(chan string),
	}
	packetBus <- r
	conn.Write([]byte(<-r.Response))
	conn.Close()
}

type RequestPacket struct {
	Body     string
	Response chan string
}

func processInputs(packetBus chan RequestPacket) {
	state := inventory.State{Items: map[string]inventory.Item{}}
	reportState := inventory.State{Items: map[string]inventory.Item{}}
	ticker := time.Tick(5000 * time.Millisecond)
	reportBus := make(chan inventory.State, 10)
	for {
		select {
		case packet := <-packetBus:
			entry, err := inventory.ParseLine(packet.Body, reportBus)
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
		case _ = <-ticker:
			fmt.Println(reportState)
		case s := <-reportBus:
			//fmt.Println("report")
			reportState = s
		}

	}
}
