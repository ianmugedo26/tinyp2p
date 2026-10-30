package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"encoding/binary"
	"bufio"
	"strings"
	"bytes"
	"crypto/sha256"
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

	reader := bufio.NewReader(conn)

	request, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print("Failed to read request:", err)
		return
	}

	request = strings.TrimSpace(request)

	if request != "GET demo-app-v1.0.0.tar.gz" {
		fmt.Println("Unsupported request:", request)
		fmt.Fprintln(conn, "ERROR")
		return
	}
	fmt.Fprintln(conn, "OK")


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

	// calculate the package's SHA-256 hash

	hasher := sha256.New()

	_, err := io.Copy(hasher, file)
	if err != nil {
		fmt.Println("Failed to calculate package hash:", err)
		return
	}

	expectedHash := hasher.Sum(nil)

	// return to the begining of the file before sending it

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("Failed to reset the package file:",err)
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

