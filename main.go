package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	filename := "messages.txt"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filename, err)
	}

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer f.Close()
		var currentLine string
		for {
			buf := make([]byte, 8)
			n, err := f.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Fatalf("Error reading: %v", err)
				}
			}

			parts := strings.Split(string(buf[:n]), "\n")
			currentLine += parts[0]

			if len(parts) > 1 {
				ch <- currentLine
				currentLine = ""
				currentLine += parts[1]
				continue
			}

			if n < 8 || err == io.EOF {
				ch <- currentLine
				close(ch)
				return
			}
		}
	}()

	return ch
}
