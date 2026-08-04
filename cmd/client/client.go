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

	req := "*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nAlice\r\n"

	_, err = conn.Write([]byte(req))
	if err != nil {
		panic(err)
	}

	fmt.Println("SET request sent")
}