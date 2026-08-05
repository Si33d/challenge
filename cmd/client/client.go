package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {

		wg.Add(1)

		go func(id int) {

			defer wg.Done()

			conn, err := net.Dial("tcp", "localhost:6379")
			if err != nil {
				fmt.Printf("Client %d: %v\n", id, err)
				return
			}
			defer conn.Close()

			var req string

			switch id % 4 {

			// SET counter value-X
			case 0:
				value := fmt.Sprintf("value-%d", id)
				req = fmt.Sprintf(
					"*3\r\n$3\r\nSET\r\n$7\r\ncounter\r\n$%d\r\n%s\r\n",
					len(value),
					value,
				)

			// GET counter
			case 1:
				req = "*2\r\n$3\r\nGET\r\n$7\r\ncounter\r\n"

			// EXISTS counter
			case 2:
				req = "*2\r\n$6\r\nEXISTS\r\n$7\r\ncounter\r\n"

			// DEL counter
			case 3:
				req = "*2\r\n$3\r\nDEL\r\n$7\r\ncounter\r\n"
			}

			_, err = conn.Write([]byte(req))
			if err != nil {
				fmt.Printf("Client %d write error: %v\n", id, err)
				return
			}

			buffer := make([]byte, 1024)

			n, err := conn.Read(buffer)
			if err != nil && err != io.EOF {
				fmt.Printf("Client %d read error: %v\n", id, err)
				return
			}

			fmt.Printf("Client %2d -> %q\n", id, string(buffer[:n]))

		}(i)
	}

	wg.Wait()

	fmt.Println("\nAll clients finished.")
}