package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			break
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}
		fmt.Println("Client disconnected:", conn.RemoteAddr())

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
