package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type Backend struct {
	servers []string
	n       uint64
}

func copyData(destination net.Conn, source net.Conn, wg *sync.WaitGroup) {
	fmt.Println("Starting copy from source to destination")
	defer fmt.Println("Finished copy from source to destination")
	defer wg.Done()
	_, err := io.Copy(destination, source)
	if err != nil {
		log.Println("copyData error:", err)
	}
	if tcpConn, ok := destination.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
}

func RoundRobin(backend *Backend) string {
	fmt.Println(backend.n)
	i := atomic.AddUint64(&backend.n, 1)
	return backend.servers[i%uint64(len(backend.servers))]
}

func (backend *Backend) Choose(strategy string) string {
	switch strategy {
	case "round-robin":
		return RoundRobin(backend)
	default:
		return ""
	}
}

func forwardConnection(conn net.Conn, backend *Backend) {
	backendAddress := backend.Choose("round-robin")
	serverConn, err := net.Dial("tcp", backendAddress)

	if err != nil {
		log.Println("copyData error:", err)
	}

	defer serverConn.Close()
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go copyData(conn, serverConn, &wg)
	go copyData(serverConn, conn, &wg)
	wg.Wait()
	fmt.Println("Wait Completes")
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
	b := &Backend{servers: backendList, n: 0}
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go forwardConnection(conn, b)
	}
}
