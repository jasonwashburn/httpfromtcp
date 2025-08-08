package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		os.Exit(1)
	}

	buff := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := buff.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Println("Error sending data:", err)
		}

	}
}
