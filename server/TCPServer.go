package main

import (
	"bufio"
	"fmt"
	"goNetworkTransfering/utils"
	"net"
	"os"
)

func main() {
	TCPServerMode()
}

// Determines the mode to put the server into
func TCPServerMode() {
	var mode string

	for {
		fmt.Println("What mode?\n1. Echo Server\n2. Read Data Server")
		_, e := fmt.Scanf("%s", &mode)
		utils.ErrorValidation(e)
		switch mode {
		case "1":
			fmt.Println("Echo Mode...")
			createTCPServer(true)
			break
		case "2":
			fmt.Println("Read Data Mode...")
			createTCPServer(false)
			break
		case "exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Error, not a selectable mode.")
		}
	}
}

// Creates a server
func createTCPServer(e bool) {
	listener, err := net.Listen("tcp", ":8721")

	utils.ErrorValidation(err)

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		utils.ErrorValidation(err)

		go readConnection(conn, e)
	}
}

// Reads the packets coming in from the client and then either
// echos the message back or sends acknowledgment packets of
// 1 byte to the client
func readConnection(conn net.Conn, echo bool) {
	for {
		message, connError := bufio.NewReader(conn).ReadBytes('\n')

		if connError != nil {
			fmt.Println(connError)
			utils.ErrorValidation(conn.Close())
			break
		}

		fmt.Printf("received: %d bytes from: %s\n", len(message), conn.RemoteAddr())

		if echo {
			_, err := conn.Write(message)
			utils.ErrorValidation(err)
		} else {
			_, err := conn.Write([] byte{10})
			utils.ErrorValidation(err)
		}
	}
}
