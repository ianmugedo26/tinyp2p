package main

import (
	"fmt"
	"os"
	"net"
	"bufio"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8001")
	if err != nil {
		fmt.Println("Failed to connect to node A:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to node A")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter message: ")

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println("Node A responded:", response)
	}
}
