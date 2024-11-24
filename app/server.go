package main

import (
	"fmt"

	"log"

	"net"

	"os"

	"strconv"

	"strings"
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

	length, err := conn.Read(buff)

	if err != nil {

		fmt.Printf("Error reading: %#v\n", err)

		return

	}

	rawData := string(buff[:length])
	log.Println("raw data", rawData)
	lines := strings.Split(rawData, "\n")

	// if the received data is a array

	if len(lines) > 0 && strings.HasPrefix(lines[0], "*") {

		elements := []string{}

		for i := 1; i < len(lines); i++ {

			if strings.HasPrefix(lines[i], "$") {

				elementLength, err := strconv.Atoi(strings.Trim(lines[i][1:], "\r"))

				if err != nil {

					log.Println("Error parsing element length:", err)

					return

				}

				if i+1 < len(lines) && len(strings.Trim(lines[i+1], "\r")) == elementLength {

					elements = append(elements, strings.Trim(lines[i+1], "\r"))

					i++ // Skip the next line as it is part of the current element

				}

			}

		}

		if len(elements) == 1 && elements[0] == "PING" {

			conn.Write([]byte("+PONG\r\n"))

			return

		}

	}

}
