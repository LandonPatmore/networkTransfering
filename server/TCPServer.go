package main

import (
	"bufio"
	"fmt"
	"goNetworkTransfering/shared"
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

		shared.ErrorValidation(err)

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

		fmt.Printf("received: %d bytes from: %s\n", len(message), conn.RemoteAddr())

		_, err := conn.Write(message)

		shared.ErrorValidation(err)
	}
}
