package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Backends struct {
	servers []string
	n       int
	mu      sync.Mutex
}

func forwardConnection(conn net.Conn, backend *Backends) {
	defer conn.Close()
	fmt.Println("Backend: ", backend)
}

func main() {
	backendList := []string{"localhost:3001", "localhost:3002", "localhost:3003"}
	// Accept tcp connection
	listner, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listner.Close()

	//Pointer to backend
	b := &Backends{servers: backendList, n: 0}
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go forwardConnection(conn, b)
	}
}
