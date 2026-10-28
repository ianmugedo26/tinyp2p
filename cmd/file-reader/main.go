package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("dist/demo-app-v1.0.0.tar.gz")
	if err != nil{
		fmt.Println("Failed to read package: ", err)
		return
	}

	hasher := sha256.New()

	hasher.Write(data)

	sum := hasher.Sum(nil)

	fmt.Printf("%x\n", sum)
}
