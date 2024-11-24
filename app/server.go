package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

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
	store map[string]*Value
}

type Value struct {
	val          string
	ttl          int
	created_time time.Time
}

func NewKvStore() *KVStore {
	return &KVStore{store: map[string]*Value{}}
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
		//log.Println("raw data", string(buff))
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
			var ttl int
			if noOfElememts == 5 && array[8] == "px" {
				ttl, _ = strconv.Atoi(array[10])

			}
			store(array[4], array[6], ttl)
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

func store(key string, val string, ttl int) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if ttl == 0 {
		kv.store[key] = &Value{val: val}
		return nil
	}
	log.Printf("got ttl of %d", ttl)
	kv.store[key] = &Value{val: val, ttl: ttl, created_time: time.Now()}
	log.Println("stored succesfully")
	return nil
}

func get(key string) (string, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if kv.store[key] != nil {
		val := kv.store[key]
		if val.created_time.IsZero() {
			return val.val, nil
		}
		createdTime := val.created_time
		ttl := val.ttl
		elapsedTime := time.Now().UnixMilli() - createdTime.UnixMilli()
		log.Printf("time elapsed %d duration %d", elapsedTime, int64(time.Duration(ttl)))
		if elapsedTime >= int64(time.Duration(ttl)) {
			return "", errors.New("key expired")
		}
		return val.val, nil
	}
	return "", fmt.Errorf("key not found in kv store")
}
