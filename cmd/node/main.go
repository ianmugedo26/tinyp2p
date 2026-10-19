package main

import (
	"fmt"
	"net"
	"io"

)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to listen on port 8080:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}


func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Peer connected:", conn.RemoteAddr())
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Peer disconnected:", conn.RemoteAddr())

			} else {
				fmt.Println("Failed to read from peer:", err)
			}
			return
		}
		message := string(buffer[:n])
		fmt.Println("Received:", message)

		response := "Hello, from node 1!\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Failed to send response:", err)
			return
		}

		fmt.Println("Response Sent to:", conn.RemoteAddr())
	}
}

