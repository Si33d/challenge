package main

import (
	"fmt"
	"io"
	"log"
	"net"
	//"strings"
	//"text/template/parse"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Redis server is listening on 6379")
	var connection net.Conn
	for {
		connection, err = listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go connectionshandle(connection)
	}
}

func connectionshandle(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			}
			return
		}
		res,err:=parseRequest(buffer[:n])
		conn.Write(res)
	}
}
func parseRequest(req []byte) ([]byte, error){

	parser := NewParser(req)
	cmd, err := parser.Parse()
	if err != nil {
		fmt.Println("Parse error:", err)
		return nil,err
	}
	fmt.Println(cmd.Name)
	fmt.Println(cmd.Args)
	return Execute(cmd)
}
