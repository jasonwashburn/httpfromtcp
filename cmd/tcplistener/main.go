package main

import (
	"fmt"
	"net"
	"os"

	"github.com/jasonwashburn/httpfromtcp/internal/request"
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

		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error reading request:", err)
			break
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", r.RequestLine.Method)
		fmt.Println("- Target:", r.RequestLine.RequestTarget)
		fmt.Println("- Version:", r.RequestLine.HttpVersion)

		fmt.Println("Client disconnected:", conn.RemoteAddr())

	}
}
