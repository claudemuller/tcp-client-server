package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Printf("usage: %s <port>\n", args[0])

		return
	}

	port := ":" + args[1]
	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("error creating listener: %v\n", err)

		return
	}
	defer listen.Close()

	fmt.Printf("Listening on %s\n", listen.Addr().String())

	var connCount int32

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("error creating connection: %v\n", err)

			return
		}

		fmt.Printf("Connection from %s\n", conn.RemoteAddr().String())

		atomic.AddInt32(&connCount, 1)
		go func() {
			if err := handleConn(conn); err != nil {
				fmt.Printf("error occurred: %v\n", err)
			}

			atomic.AddInt32(&connCount, -1)
			conn.Close()

			fmt.Printf("%d active connections\n", connCount)
		}()

		fmt.Printf("%d active connections\n", connCount)
	}
}

func handleConn(conn net.Conn) error {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading from conn: %w", err)
		}

		if strings.TrimSpace(string(data)) == "DISCONNECT" {
			fmt.Printf("Closing connection to %s...\n", conn.RemoteAddr().String())

			if _, err = conn.Write([]byte("Disconnecting you...")); err != nil {
				return fmt.Errorf("error writing to conn: %w", err)
			}

			return nil
		}

		fmt.Print("-> ", string(data))
		t := time.Now()
		tf := t.Format(time.RFC3339) + "\n"

		if _, err = conn.Write([]byte(tf)); err != nil {
			return fmt.Errorf("error writing to conn: %w", err)
		}
	}
}
