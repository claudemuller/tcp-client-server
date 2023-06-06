package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Printf("usage: %s <host>:<port>\n", args[0])

		return
	}

	connStr := args[1]
	conn, err := net.Dial("tcp", connStr)
	if err != nil {
		fmt.Printf("error when connecting: %v\n", err)

		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error when reading from stdio: %v\n", err)
		}

		fmt.Fprintf(conn, text+"\n")

		res, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("error when reading from conn: %v\n", err)

			return
		}

		fmt.Printf("-> %s", res)

		if strings.TrimSpace(string(text)) == "QUIT" {
			fmt.Println("TCP client exiting...")

			return
		}
	}
}
