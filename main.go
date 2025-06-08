package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting: %v", err)
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}
		log.Printf("Closed connection from %s", conn.RemoteAddr())
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
