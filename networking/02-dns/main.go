//  you should be able to specify a domain name and query type
// (such as A or NS) and see the parsed output printed to the
// command line, as you would with a tool like dig or nslookup.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

type Header struct {
	Id [4]byte
	Qr bool // 1 bit
	Opcode [2]bool // 2 bits
	Aa bool
	Tc bool
	Rd bool
	Ra bool
	Z [3]bool // 3 bits
	Rcode [1]byte // 8 bits
	QdCount [4]byte
	AnCount [4]byte
	NsCount [4]byte
	ArCount [4]byte
}

type Question struct {
	Qname [2]byte
	Qtype [2]byte
	Qclass [2]byte
}

func main() {
	var domain string
	var queryType string
	
	flag.StringVar(&domain, "domain", "newschematic.org", "the domain to check")
	flag.StringVar(&queryType, "queryType", "A", "query type for domain to check")
	// TODO validate query type entry

	flag.Parse()

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal("Shoot! ", err)
	}

	defer conn.Close()

	n, err := conn.Write([]byte(domain))
	if err != nil {
		log.Fatal("Gosh! ", err)
	}
	fmt.Println(n)
	
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal("Dang! ", err)
	}

	fmt.Println(status)

	fmt.Println("domain:", domain)
	fmt.Println("query type:", queryType)
}
