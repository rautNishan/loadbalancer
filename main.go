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

func RoundRobin(backend *Backend) string {
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

func copyData(destinationConn net.Conn, sourceConn net.Conn, wg *sync.WaitGroup, path string) {
	buffer := make([]byte, 1024)
	fmt.Println("Ready to read")
	defer wg.Done()
	for {
		_ = sourceConn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := sourceConn.Read(buffer)
		// fmt.Println("Complete Read: ", string(buffer[:n]))
		if err != nil {
			fmt.Println("This is error in read", path, err)
			break
		}
		_, err = destinationConn.Write(buffer[:n])
		if err != nil {
			fmt.Println("This is error in write: ", path, err)
			break
		}
	}
	tcpConn, isTcp := destinationConn.(*net.TCPConn)
	if isTcp {
		tcpConn.CloseWrite()
	}
}

func forwardConnection(clientConn net.Conn, backend *Backend) {
	backendAddress := backend.Choose("round-robin")
	// fmt.Println("This is choosen adddress: ", backendAddress)

	serverConn, err := net.Dial("tcp", backendAddress)
	defer func() {
		serverConn.Close()
		clientConn.Close()
		fmt.Println("Connection closed")
	}()
	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go copyData(serverConn, clientConn, &wg, "client->backend")
	go copyData(clientConn, serverConn, &wg, "backend->client")
	wg.Wait()
}

func main() {
	backendLists := []string{"localhost:3001", "localhost:3002", "localhost:3003"}
	listner, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ready to accept connection")
	b := &Backend{servers: backendLists, n: 0}
	for {
		conn, err := listner.Accept()
		fmt.Println(conn)
		if err != nil {
			fmt.Println(err)
		}
		go forwardConnection(conn, b)
	}
}
