package main

// implement a highly simplified version of an HTTP proxy
// The primary use case we’ll have in mind is caching.
// We’ll focus on the reverse proxy use case,
// used by specific web applications to sit in front
// of their own web servers, caching a narrower set
// of websites but in service of an unlimited amount of clients.

import (
	"fmt"
	"log"
	"net"
)

// [√] 1. Write a program which accepts a TCP connection and
//	      simply responds with whatever it reads
// 			listen, accept, recv, send
// [x] 2. Write a program that simply listens on a port and
//			forwards on to another server running locally (nc)
// 			my program: port 6666 (which is where I'll send the req)
//			other program: port 7777 (which is where I'll fwd the req)

// listen
func listenTCP(address string, port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
}

// recv
func recv(conn net.Conn) {
	fmt.Println("Connection received over", conn.LocalAddr().Network())

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)

	if err != nil {
		log.Fatal("Cannot read: ", err)
	}

	// send
	_, err = conn.Write(buffer)
	fwd(buffer)

	if err != nil {
		log.Fatal("Cannot write: ", err)
	}

	fmt.Println("They said,", string(buffer))
	conn.Close()
}

func fwd(msg []byte) {
	// create connection to 7777
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		log.Fatal("cannot connect to 7777: ", err)
	}
	// close connection when done
	defer conn.Close()

	// send msg
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatal("cannot write to 7777: ", err)
	}
}

func main() {
	listener, err := listenTCP("127.0.0.1", 6666)
	if err != nil {
		log.Fatal("cannot listen: ", err)
	}

	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("cannot accept: ", err)
		}

		go recv(conn)
	}
}
