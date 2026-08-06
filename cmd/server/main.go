package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()
	log.Println("Redis server is listening on 6379")
	var connection net.Conn
	for {
		connection, err = listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go connectionshandle(connection)
	}
}

func connectionshandle(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				break
			}
			return
		}
		parser := NewParser(buffer[:n])

		for parser.pos < len(parser.data) {
			cmd, err := parser.Parse()
			if err != nil {
				_, _ = conn.Write([]byte("-ERR" + err.Error() + "\r\n"))
				break
			}
			res, err := Execute(cmd)
			if err != nil {
				_, _ = conn.Write([]byte("-ERR" + err.Error() + "\r\n"))
				continue
			}
			_, err = conn.Write(res)
			if err != nil {
				log.Println("Write Error", err)
				return
			}

		}
	}
}
