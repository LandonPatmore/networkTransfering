package main

import (
	"bufio"
	"fmt"
	"goNetworkTransfering/utils"
	"net"
)

func main() {
	createTCPServer()
}

// Creates a server
func createTCPServer() {
	listener, err := net.Listen("tcp", ":8721")

	utils.ErrorValidation(err)

	defer listener.Close()

	fmt.Println("Server started...")

	for {
		conn, err := listener.Accept()

		utils.ErrorValidation(err)

		go readConnection(conn)
	}
}

// Reads the packets coming in from the client and then either
// echos the message back or sends acknowledgment packets of
// 1 byte to the client
func readConnection(conn net.Conn) {
	var echoMode bool

	for {
		message, connError := bufio.NewReader(conn).ReadBytes('\n')

		if connError != nil {
			fmt.Println(connError)
			utils.ErrorValidation(conn.Close())
			break
		}

		fmt.Printf("received: %d bytes from: %s\n", len(message), conn.RemoteAddr())

		if len(message) == 2 { // switch server modes
			switch message[0] {
			case 1:
				fmt.Println("Changed to echo mode...")
				echoMode = true
				break
			case 2:
				fmt.Println("Changed to acknowledge mode...")
				echoMode = false
				break
			default:
				break
			}
		}

		if echoMode {
			_, err := conn.Write(message)
			utils.ErrorValidation(err)
		} else {
			_, err := conn.Write([] byte{10})
			utils.ErrorValidation(err)
		}
	}
}
