package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	connection, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer connection.Close()

	consoleInput := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		outputMessage, _ := consoleInput.ReadString('\n')

		_, err = connection.Write([]byte(outputMessage))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}

		buffer := make([]byte, 1024)
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}
		fmt.Print(string(buffer[:n]))

		if strings.TrimSpace(outputMessage) == "quit" {
			break
		}
	}
}
