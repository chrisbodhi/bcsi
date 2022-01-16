package main

// implement a highly simplified version of an HTTP proxy
// The primary use case we’ll have in mind is caching.
// We’ll focus on the reverse proxy use case,
// used by specific web applications to sit in front
// of their own web servers, caching a narrower set
// of websites but in service of an unlimited amount of clients.

// client --> me (9000) --> server (7777)
// client <-- me (9000) <-- server (7777)

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// [√] 1. Write a program which accepts a TCP connection and
//	      simply responds with whatever it reads
// 			listen, accept, recv, send
// [x] 2. Write a program that simply listens on a port and
//			forwards on to another server running locally (nc)
// 			my program: port 9000 (which is where I'll send the req)
//			other program: port 7777 (which is where I'll fwd the req)
// [ ] 3. Handle a complete request and response

func main() {
	listener, err := listenTCP("127.0.0.1", 9000)
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

// listen
func listenTCP(address string, port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
}

// recv
func recv(conn net.Conn) {
	fmt.Println("Connection received over", conn.LocalAddr().Network())

	buffer := make([]byte, 1<<10)
	_, err := conn.Read(buffer)

	if err != nil {
		log.Fatal("Cannot read: ", err)
	}
	fmt.Println("Forwarding\n", string(buffer))
	// send
	res, err := fwd(buffer)
	if err != nil {
		log.Fatal("Did not get response: ", err)
	}
	fmt.Println("Writing response to client\n", string(res))
	// TODO get length of res?
	trimTo := calcHeaderLen(res)
	_, err = conn.Write(res[:trimTo])

	if err != nil {
		log.Fatal("Cannot write: ", err)
	}

	conn.Close()
}

func fwd(msg []byte) ([]byte, error) {
	// create connection to 7777: me --> server
	serverConn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		log.Fatal("cannot connect to 7777: ", err)
	}
	// close connection when done
	defer serverConn.Close()

	// send msg received from client over to server
	_, err = serverConn.Write(msg)

	// TODO new method of creating the correctly-sized buffer
	// 		otherwise, sending too much causes curl to show an error message
	//		"* Excess found in a read: excess = 827, size = 86, maxdownload = 86, bytecount = 0"
	//		(and maybe the browsers to silently fail)

	// TODO consult the HTTP spec
	// get len of res from start to the first empty line: this will be the end of the headers and the start of the body (of which there will be none in the GET request)

	res := make([]byte, 1<<10)
	_, err = serverConn.Read(res)

	if err != nil {
		log.Fatal("cannot write to 7777: ", err)
	}

	return res, nil
}

// TODO we don't want the length of the headers, we want
//		the size of the body -- which contains the headers
//		in JSON form
func calcHeaderLen(res []byte) int {
	str := string(res)
	stSp := strings.Split(str, "\n")
	totalLen := len("HTTP/1.0 200 OK\n")

	for i, s := range stSp {
		fmt.Printf("line %d | len %d | contents %s\n", i, len(s), s)
		// Skip the first line, which isn't a header
		if i == 0 {
			continue
		}
		// Empty new line, which seperates headers from body.
		// We just want the headers, so stop counting here.
		if len(s) == 1 {
			break
		}

		if len(s) > 0 {
			totalLen += len(s)
		}
	}

	return totalLen
}
