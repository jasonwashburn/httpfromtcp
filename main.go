package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	currentLineContents := ""

	for {
		buffer := make([]byte, 8)
		n, err := file.Read(buffer)
		if err != nil {
			if currentLineContents != "" {
				fmt.Printf("read: %s\n", currentLineContents)
				currentLineContents = ""
			}

			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading file:", err)
			break
		}

		str := string(buffer[:n])
		parts := strings.Split(str, "\n")
		for _, part := range parts[:len(parts)-1] {
			fmt.Printf("read: %s%s\n", currentLineContents, part)
			currentLineContents = ""
		}
		currentLineContents += parts[len(parts)-1]
	}
}
