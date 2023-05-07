package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/chrisbodhi/bcsi/distributed-systems/kv-store/server"
	"github.com/chrisbodhi/bcsi/distributed-systems/kv-store/utils"
)

func main() {
	fmt.Println("Starting servers...")

	go handleRequests(":8888")

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		server.Start(utils.PORTS[0])
		defer wg.Done()
	}()

	go func() {
		server.Start(utils.PORTS[1])
		defer wg.Done()
	}()

	wg.Wait()
}

// Open a port at 8888 and listen for incoming connections.
// When a connection is received, pass it to handleConnection.
func handleRequests(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("!!", err)
	}
	fmt.Printf("LB listening on %s...\n", port)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("err with", port, ":", err)
		}
		go handleConnection(conn)
	}
}

// handleConnection sends the request to one of our servers, and then
// routes the response back to the original connection.
func handleConnection(conn net.Conn) {
	// Pick a number between 0 and the length of PORTS.
	// This will be used to select a port to connect to.
	portIndex := utils.Random(0, len(utils.PORTS))
	port := utils.PORTS[portIndex]

	// Create a connection to either 8889 or 8890.
	// 8889 and 8890 better be listening for incoming connections!
	// TODO: handle the case where 8889 and 8890 are not listening.
	serverConn, err := net.Dial("tcp", port)
	if err != nil {
		log.Fatal("!!", err)
	}
	defer serverConn.Close()

	// Read the incoming connection from the client.
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	// Write the incoming connection to the server.
	serverConn.Write(buf[:n])

	// Read the response from the server.
	buf2 := make([]byte, 1024)
	n2, err := serverConn.Read(buf2)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Received response from server at", port)

	// Write the response from server to the client.
	conn.Write(buf2[:n2])
}
