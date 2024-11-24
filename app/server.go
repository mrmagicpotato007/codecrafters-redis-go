package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"log"

	"net"

	"os"
)

const (
	CLRF = "\r\n"
)

var kv *KVStore

type KVStore struct {
	mu    sync.Mutex
	store map[string]string
}

func NewKvStore() *KVStore {
	return &KVStore{store: map[string]string{}}
}
func init() {
	kv = NewKvStore()
}

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
		data := strings.ReplaceAll(string(buff), "\n", `\n`)
		log.Println("raw data", data)
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

	case '*':
		array := strings.Split(ip, CLRF)

		noOfElememts, err := strconv.Atoi(string(array[0][1]))
		log.Println("no of elements", noOfElememts)
		if err != nil {
			fmt.Println("Error converting to integer:", err)
		}
		if array[2] == "ECHO" {
			//ip://*2 \r\n $4 \r\n ECHO \r\n $3 \r\n hey \r\n
			//op://$3\r\nhey\r\n
			conn.Write([]byte(array[3] + CLRF + array[4] + CLRF))
		} else if array[2] == "SET" {
			//ip://"*3 \r\n $3 \r\n SET \r\n $6 \r\n orange \r\n $5 \r\n apple \r\n
			//op://+OK\r\n
			store(array[4], array[6])
			log.Printf("received key %s val %s \n", array[4], array[6])
			conn.Write([]byte("+OK\r\n"))
		} else if array[2] == "GET" {
			//ip://*2 \r\n $3 \r\n GET \r\n $9 \r\n pineapple \r\n
			val, err := get(array[4])
			if err != nil {
				conn.Write([]byte("$-1\r\n"))
				return
			}
			conn.Write([]byte("$" + strconv.Itoa(len(val)) + CLRF + val + CLRF))

		} else {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}

func store(key string, val string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = val
	log.Println("stored succesfully")
	return nil
}

func get(key string) (string, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if kv.store[key] != "" {
		return kv.store[key], nil
	}
	return "", fmt.Errorf("key not found in kv store")
}
