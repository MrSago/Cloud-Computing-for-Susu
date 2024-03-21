package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	serverInput := bufio.NewReader(connection)
	consoleInput := bufio.NewReader(os.Stdin)

	welcomeMessage := make([]byte, 1024)
	serverInput.Read(welcomeMessage)
	fmt.Print(string(welcomeMessage))

	for {
		fmt.Print("> ")
		outputMessage, _ := consoleInput.ReadString('\n')

		_, err = fmt.Fprint(connection, outputMessage)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}

		buffer := make([]byte, 1024)
		_, err := serverInput.Read(buffer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}

		fmt.Print(string(buffer))

		if strings.TrimSpace(outputMessage) == "quit" {
			break
		}
	}
}
