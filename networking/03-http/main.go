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

// 1. Write a program which accepts a TCP connection and
//	  simply responds with whatever it reads
// 		listne, accept, recv, send

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

	if err != nil {
		log.Fatal("Cannot write: ", err)
	}

	fmt.Println("They said,", string(buffer))
	conn.Close()
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
