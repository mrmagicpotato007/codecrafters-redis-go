package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"log"

	"net"

	"os"
)

const (
	CLRF = "\r\n"
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

	buff := make([]byte, 2048)

	for {

		_, err := conn.Read(buff)

		if err == io.EOF {
			break
		}
		if err != nil {

			fmt.Printf("Error reading: %#v\n", err)

			return

		}
		log.Println("raw data", string(buff))
		parseIp(string(buff), conn)
	}

}

func parseIp(ip string, conn net.Conn) {

	ipLength := len(ip)

	if ipLength == 0 {
		log.Println("invalid input")
	}

	switch ip[0] {
	//we recieved array
	//*2 \r\n $4 \r\n ECHO \r\n $3 \r\n hey \r\n
	case '*':
		array := strings.Split(ip, CLRF)
		log.Println("first ele", array[0])
		noOfElememts, err := strconv.Atoi(string(array[0][1]))
		log.Println("no of elements", noOfElememts)
		if err != nil {
			fmt.Println("Error converting to integer:", err)
		}
		if array[2] == "ECHO" {
			//$3\r\nhey\r\n
			conn.Write([]byte(array[3] + CLRF + array[4] + CLRF))
		} else {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
