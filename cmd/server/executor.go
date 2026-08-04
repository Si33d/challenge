package main

import (
	"fmt"
	"strings"
)

func Execute(cmd *Command) ([]byte, error) {
	str := strings.ToUpper(cmd.Name)

	switch str {
	case "PING":
		if len(cmd.Args) == 0 {
			return []byte("+PONG\r\n"), nil
		}
		if len(cmd.Args) == 1 {
			msg := cmd.Args[0]
			return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)), nil
		}
		return []byte("-ERR Wrong no.of arguments"), nil

	case "SET":
		if len(cmd.Args)!=2{
			return []byte("-ERR Wrong no. of arguments for SET"),nil
		}
		store.Set(cmd.Args[0],cmd.Args[1])
		return []byte("+OK\r\n"),nil
		

	case "GET":
		if len(cmd.Args)!=1{
			return []byte("-ERR Wrong no. of arguments for Get"),nil
		}
		value,ok:=store.Get(cmd.Args[0])
		if !ok{
			return []byte("$-1\r\n"),nil
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)), nil
	
	case "DEL":
		if len(cmd.Args)==0{
			return []byte("-ERR Wrong no. of arguments for 'DEL'\r\n"),nil
		}
		count:=store.Delete(cmd.Args...)
		return []byte(fmt.Sprintf(":%d\r\n",count)),nil
	
	case "EXISTS":
		if len(cmd.Args) != 1 {
			return []byte("-ERR wrong number of arguments for 'EXISTS'\r\n"), nil
		}
		if store.Exists(cmd.Args[0]) {
			return []byte(":1\r\n"), nil
		}
		return []byte(":0\r\n"), nil

	default:
		return nil, fmt.Errorf("-ERR Unknown Command %s", cmd.Name)
	}

}
