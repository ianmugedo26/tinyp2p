package main

import (
	"fmt"
	"os"
	"net"
	"io"
	"encoding/binary"
	"crypto/sha256"
	"path/filepath"

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
	err = binary.Read(conn, binary.BigEndian, &size)
	if err != nil {
		fmt.Println("Failed to read file size:", err)
		return
	}

	// Reject an unexpectedly large transfer.
	const maxFileSize = 100 * 1024 * 1024
	if size > maxFileSize {
		fmt.Println ("Package is too large:", size, "bytes")
		return
	}

	fmt.Println ("Package size:", size, "bytes")

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
		return
	}

	hash := sha256.Sum256(data)

	fmt.Println("Package received:", outputPath)
	fmt.Printf("SHA-256: %x\n", hash)

}
