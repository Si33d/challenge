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

	// Send raw RESP bytes for: PING
	data := "*1\r\n$4\r\nPING\r\n"

	fmt.Println("Sending:")
	fmt.Println(data)

	_, err = conn.Write([]byte(data))
	if err != nil {
		panic(err)
	}

	fmt.Println("Request sent")
}