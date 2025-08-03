package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	bytes := make([]byte, 8)
	eof := false

	for {
		_, err := file.Read(bytes)
		if errors.Is(err, io.EOF) {
			eof = true
		} else if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}

		fmt.Printf("read: %s\n", bytes)
		if eof {
			break
		}
	}
}
