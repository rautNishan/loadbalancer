package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	servers []string
	n       uint64
}

func copyData(destination net.Conn, source net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	buffer := make([]byte, 4*1024)
	for {
		//Close the connection with client
		_ = source.SetReadDeadline(time.Now().Add(20 * time.Second))
		n, err := source.Read(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		_, err = destination.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error : ", err)
			break
		}
	}
	tcpConn, isTCP := destination.(*net.TCPConn)
	if isTCP {
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

func forwardConnection(clientConn net.Conn, backend *Backend) {
	backendAddress := backend.Choose("round-robin")
	serverConn, err := net.Dial("tcp", backendAddress)

	if err != nil {
		log.Println("copyData error:", err)
	}

	defer func() {
		serverConn.Close()
		clientConn.Close()
		fmt.Println("Connection Closed")
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go copyData(serverConn, clientConn, &wg)
	go copyData(clientConn, serverConn, &wg)
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
