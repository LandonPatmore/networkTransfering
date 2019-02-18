package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	createTCPServer()
}

func createTCPServer() {
	listener, err := net.Listen("tcp", ":8274")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go echo(conn)
	}
}

func echo(conn net.Conn) {
	for {

		message, connError := bufio.NewReader(conn).ReadBytes('\n')

		if connError != nil {
			log.Println(connError)
			conn.Close()
			break
		}

		// output message received
		fmt.Println("Bytes from client: ", message)

		_, err := conn.Write(message)

		if err != nil {
			log.Fatal(err)
		}
	}
}
