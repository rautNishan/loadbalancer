package main

import "fmt"

func main() {
	backendList := [3]string{"localhost:3001", "localhost:3002", "localhost:3003"}
	for i := 0; i < len(backendList); i++ {
		fmt.Println(backendList[i])
	}
}
