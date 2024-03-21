package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(connection net.Conn) {
	defer connection.Close()

	fmt.Printf("New connection from %s\n", connection.RemoteAddr().String())

	reader := bufio.NewReader(connection)

	fmt.Fprintf(connection, "Welcome to the server!\r\n"+
		"Your IP address is %s\r\n"+
		"Type 'quit' to exit.\r\n",
		connection.RemoteAddr().String())

	for {
		inputMessage, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}

		fmt.Printf("[%s]: %s", connection.RemoteAddr().String(), inputMessage)

		if strings.TrimSpace(inputMessage) == "quit" {
			fmt.Printf("Closing connection from %s\n", connection.RemoteAddr().String())
			fmt.Fprint(connection, "Bye!")
			break
		}

		responseMessage := fmt.Sprintf("[Server]: %s", inputMessage)
		_, err = fmt.Fprint(connection, responseMessage)
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
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}
		go handleConnection(connection)
	}
}
