package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	server_input := bufio.NewReader(conn)
	console_input := bufio.NewReader(os.Stdin)

	welcome_message := make([]byte, 1024)
	server_input.Read(welcome_message)
	fmt.Print(string(welcome_message))

	for {
		fmt.Print("> ")
		out_msg, _ := console_input.ReadString('\n')

		_, err = fmt.Fprint(conn, out_msg)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}

		buffer := make([]byte, 1024)
		_, err := server_input.Read(buffer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			break
		}

		fmt.Print(string(buffer))

		if strings.Contains(out_msg, "quit") {
			break
		}
	}
}
