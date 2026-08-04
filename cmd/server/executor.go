package main

import (
	"fmt"
	"strings"
)

func Execute(cmd *Command) ([]byte, error) {
	str := strings.ToUpper(cmd.Name)
	
	switch str{
	case "PING":
			return []byte("+PONG\r\n"),nil
	default:
		return nil,fmt.Errorf("Unknown Command",cmd.Name)
	}

}