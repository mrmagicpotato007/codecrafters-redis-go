package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	buff := make([]byte, 128)
	_, err = conn.Read(buff)
	if err != nil {
		fmt.Println("Error while reading from the connection: ", err.Error())
	}
	conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println("Error while writing to the connection: ", err.Error())
	}
}
