package main

import (
	"log"
	"net"
)

func forwardConnection(conn net.Conn, backendList []string) {

}

func main() {
	backendList := []string{"localhost:3001", "localhost:3002", "localhost:3003"}
	// Accept tcp connection
	listner, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go forwardConnection(conn, backendList)
	}
}
