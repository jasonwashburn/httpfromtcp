package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	ch := getLinesChannel(f)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	currentLineContents := ""

	go func() {
		defer f.Close()
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					ch <- currentLineContents
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
				ch <- fmt.Sprintf("%s%s", currentLineContents, part)
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
		close(ch)
	}()

	return ch
}
