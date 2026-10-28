package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"encoding/binary"
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
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connected to:", conn.RemoteAddr())

	// open the package we want to distribute

	file, err := os.Open("dist/demo-app-v1.0.0.tar.gz")
	if err != nil {
		fmt.Println("Failed to open package:", err)
		return
	}
	defer file.Close()

	// Get the package's size
	info, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get the package size:", err)
		return
	}

	// send the file size as an 8-byte header
	size := uint64(info.Size())
	err = binary.Write(conn, binary.BigEndian, size)
	if err != nil {
		fmt.Println("Failed to send the file:", err)
		return
	}

	// Send exactly the number of byte in the package.

	_, err = io.CopyN(conn, file,info.Size())
	if err != nil {
		fmt.Println("Failed to send the package:", err)
	}

	fmt.Println("Package sent successfully:", info.Size(), "bytes")

}

