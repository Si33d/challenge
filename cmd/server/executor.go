package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
		if len(cmd.Args) == 2 {
			store.Set(cmd.Args[0], cmd.Args[1], nil)
			return []byte("+OK\r\n"), nil
		}
		if len(cmd.Args) == 4 {
			if strings.ToUpper(cmd.Args[2]) != "EX" {
				return []byte("-ERR Syntax error\r\n"), nil
			}
			t, err := strconv.Atoi(cmd.Args[3])
			if err != nil {
				return []byte("-ERR Invalid expire time\r\n"), nil
			}

			expiryTime := time.Now().Add(time.Duration(t) * time.Second)

			store.Set(cmd.Args[0], cmd.Args[1], &expiryTime)
			return []byte("+OK\r\n"), nil
		}
		return []byte("-ERR wrong number of arguments for 'SET'\r\n"), nil

	case "GET":
		if len(cmd.Args) != 1 {
			return []byte("-ERR Wrong no. of arguments for Get"), nil
		}
		value, ok := store.Get(cmd.Args[0])
		if !ok {
			return []byte("$-1\r\n"), nil
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)), nil

	case "DEL":
		if len(cmd.Args) == 0 {
			return []byte("-ERR Wrong no. of arguments for 'DEL'\r\n"), nil
		}
		count := store.Delete(cmd.Args...)
		return []byte(fmt.Sprintf(":%d\r\n", count)), nil

	case "EXISTS":
		if len(cmd.Args) != 1 {
			return []byte("-ERR wrong number of arguments for 'EXISTS'\r\n"), nil
		}
		if store.Exists(cmd.Args[0]) {
			return []byte(":1\r\n"), nil
		}
		return []byte(":0\r\n"), nil

	case "EXPIRE":
		if len(cmd.Args) != 2 {
			return []byte("-ERR wrong number of arguments for EXPIRE\r\n"), nil
		}
		seconds, err := strconv.Atoi(cmd.Args[1])
		if err != nil {
			return []byte("-ERR invalid expire tiem \r\n"), nil
		}
		ok := store.Expire(cmd.Args[0], seconds)
		if ok {
			return []byte(":1\r\n"), nil
		}
		return []byte(":0\r\n"), nil

	case "TTL":
		if len(cmd.Args) != 1 {
			return []byte("-ERR wrong number of arguments for 'TTL'\r\n"), nil
		}
		ttl := store.TTL(cmd.Args[0])
		return []byte(fmt.Sprintf(":%d\r\n", ttl)), nil

	default:
		return []byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd.Name)), nil
	}

}
