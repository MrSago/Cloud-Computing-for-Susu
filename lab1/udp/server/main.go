package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error resolving UDP address: %s\n", err)
		return
	}

	connection, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Error listening on UDP: %s\n", err)
		return
	}
	defer connection.Close()

	fmt.Println("Listening on UDP port 8080...")

	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %s\n", err)
			continue
		}

		inputMessage := string(buffer[:n])
		fmt.Printf("[%s]: %s", remoteAddr.String(), inputMessage)

		if strings.TrimSpace(inputMessage) == "quit" {
			fmt.Printf("Closing connection from %s\n", remoteAddr.String())
			_, err = connection.WriteToUDP([]byte("Bye!"), remoteAddr)
			if err != nil {
				fmt.Printf("Error sending response: %s\n", err)
			}
			continue
		}

		responseMessage := "[Server]: " + inputMessage
		_, err = connection.WriteToUDP([]byte(responseMessage), remoteAddr)
		if err != nil {
			fmt.Printf("Error sending response: %s\n", err)
		}
	}
}
