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

	bytes := make([]byte, 8)
	eof := false
	currentLine := ""

	for {
		bytesRead, err := file.Read(bytes)
		if errors.Is(err, io.EOF) {
			eof = true
		} else if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}

		parts := strings.Split(string(bytes[:bytesRead]), "\n")
		currentLine += parts[0]
		if len(parts) > 1 { // we hit a newline
			fmt.Printf("read: %s\n", currentLine)
			currentLine = parts[1] // set currentLine to remaining part
		}

		if eof {
			if len(currentLine) != 0 {
				fmt.Printf("read: %s\n", currentLine) // output anything that's left
			}
			break
		}
		bytes = make([]byte, 8) // reset bytes for the next read
	}
}
