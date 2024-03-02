package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())

	input := bufio.NewReader(conn)

	fmt.Fprintf(conn,
		"Welcome to the server!\r\n"+
			"Your IP address is %s\r\n"+
			"Type 'quit' to exit.\r\n",
		conn.RemoteAddr().String())

	for {
		message, err := input.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}

		fmt.Printf("[%s]: %s", conn.RemoteAddr().String(), message)

		if strings.Contains(message, "quit") {
			fmt.Printf("Closing connection from %s\n", conn.RemoteAddr().String())
			fmt.Fprint(conn, "Bye!")
			break
		}

		message = "[Server]: " + message
		_, err = fmt.Fprint(conn, message)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}
		go handleConnection(conn)
	}
}
