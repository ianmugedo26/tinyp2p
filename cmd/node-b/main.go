package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8001")
	if err != nil {
		fmt.Println("Failed to connect to node A:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to node A")

	//Read the 8-byte file-size header
	var size uint64

	var expectedHash [sha256.Size]byte
	_, err = io.ReadFull(conn, expectedHash[:])
	if err != nil {
		fmt.Println("Failed to read expected hash:", err)
		return
	}
	err = binary.Read(conn, binary.BigEndian, &size)
	if err != nil {
		fmt.Println("Failed to read file size:", err)
		return
	}

	// Reject an unexpectedly large transfer.
	const maxFileSize = 100 * 1024 * 1024
	if size > maxFileSize {
		fmt.Println("Package is too large:", size, "bytes")
		return
	}

	fmt.Println("Package size:", size, "bytes")

	// Create a destination directory for the file
	err = os.MkdirAll("received", 0755)
	if err != nil {
		fmt.Println("Failed to create received directory:", err)
		return
	}

	// Create the destination file

	outputPath := filepath.Join("received", "demo-app-v1.0.0.tar.gz")
	out, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Failed to create a destination file:", err)
		return
	}

	// Receive exactly the advertised number of bytes

	_, copyErr := io.CopyN(out, conn, int64(size))
	closeErr := out.Close()

	if copyErr != nil {
		fmt.Println("Failed to receive package:", copyErr)
		os.Remove(outputPath)
		return
	}
	if closeErr != nil {
		fmt.Println("Failed to close destination file:", closeErr)
		os.Remove(outputPath)
		return
	}

	// Calculate the received package's SHA256 hash

	data, err := os.ReadFile(outputPath)
	if err != nil {
		fmt.Println("Failed to read received package:", err)
		os.Remove(outputPath)
		return
	}

	// Calculate The received package's SHA-256 hash
	actualHash := sha256.Sum256(data)

	// Compare the received hash with the expected hash

	if !bytes.Equal(actualHash[:], expectedHash[:]) {
		fmt.Println("PACKAGE VERIFICATION FAILED")
		fmt.Println("Expected SHA-256: %x\n", expectedHash)
		fmt.Println("Actual SHA-256: %x\n", actualHash)

		os.Remove(outputPath)
		return
	}
	hash := sha256.Sum256(data)

	fmt.Println("Package received:", outputPath)
	fmt.Printf("SHA-256: %x\n", hash)

	//sending request for the code

	_, err = fmt.Fprintln(conn, "GET demo-app-v1.0.0.tar.gz")
	if err != nil {
		fmt.Println("Failed to Send the Request:", err)
		return
	}

	// Read the provider's response
	reader := bufio.NewReader(conn)

	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	response = strings.TrimSpace(response)

	if response != "OK" {
		fmt.Println("Package request rejected:", response)
		return
	}
}
