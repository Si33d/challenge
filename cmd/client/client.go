package main

import (
	"fmt"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	req := "*1\r\n$4\r\nPING\r\n"

	_, err = conn.Write([]byte(req))
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server Response:")
	fmt.Println(string(buffer[:n]))
}