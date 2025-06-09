package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	address, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	buff := bufio.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Print("> ")
		msg, err := buff.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Println(err)
			continue
		}
	}

}
