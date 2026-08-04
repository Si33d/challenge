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
	default:
		return nil, fmt.Errorf("Unknown Command %s", cmd.Name)
	}

}
