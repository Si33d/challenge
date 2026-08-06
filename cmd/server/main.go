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

	buffer := make([]byte, 1024)
	rembuffer:=make([]byte,0)
	
	for {

		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				break
			}
			return
		}
		rembuffer=append(rembuffer, buffer[:n]...)
		parser := NewParser(rembuffer)

		for parser.pos < len(parser.data) {
			start:=parser.pos
			cmd, err := parser.Parse()
			if err != nil {
				if err.Error()=="incomplete bulk string" || err.Error() == "unexpected end of input"{
					parser.pos=start
					break
				}
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
		rembuffer=rembuffer[parser.pos:]
	}
}
