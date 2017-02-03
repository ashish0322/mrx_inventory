package main

/**
The original outline of this was taken from here:

https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go
*/

import (
	"fmt"
	"inventory"
	"net"
	"os"
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
	packetBus := make(chan inventory.RequestPacket)
	go inventory.ProcessInputs(packetBus,
		func(s inventory.State) {
			print("\033[H\033[2J")
			fmt.Println(inventory.RenderState(s))
		},
	)
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
func handleRequest(conn net.Conn, packetBus chan inventory.RequestPacket) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	readLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	r := inventory.NewRequest(string(buf[0:readLen]))
	packetBus <- r
	conn.Write([]byte(<-r.Response + "\n"))
	conn.Close()
}
