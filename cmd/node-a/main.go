package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Println("Failed to start port 8001:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port 8001...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connected to:", conn.RemoteAddr())

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Node B disconnected")
				return
			} else {
				fmt.Println("Read Error:", err)
				return
			}	}

			fmt.Println("Node A received: %s\n", string(buffer[:n]))

			_, err = conn.Write([]byte("Hello from NodeA!\n"))
			if err != nil {
				fmt.Println("Write Error:", err)
				return
			}
		}
	}
