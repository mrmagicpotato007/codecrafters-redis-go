package main

import (
	"fmt"

	"log"

	"net"

	"os"
)

func main() {

	l, err := net.Listen("tcp", "127.0.0.1:6379")

	if err != nil {

		fmt.Println("Failed to bind to port 6379")

		os.Exit(1)

	}

	defer l.Close()

	for {

		conn, err := l.Accept()

		if err != nil {

			fmt.Println("Error accepting connection: ", err.Error())

			os.Exit(1)

		}

		go handleRequest(conn)

	}

}

func handleRequest(conn net.Conn) {

	buff := make([]byte, 1024)

	for {
		_, err := conn.Read(buff)

		if err != nil {

			fmt.Printf("Error reading: %#v\n", err)

			return

		}

		log.Println("raw data", string(buff))
		conn.Write([]byte("+PONG\r\n"))
	}

}
