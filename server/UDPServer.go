package main

import (
	"fmt"
	"goNetworkTransfering/utils"
	"net"
)

func main() {
	createUDPServer()
}

// Creates a server and reads the packets coming in from the client and then either
// echos the message back or sends acknowledgment packets of
// 1 byte to the client
func createUDPServer() {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":8274")

	utils.ErrorValidation(err)

	conn, err := net.ListenUDP("udp", ServerAddr)
	utils.ErrorValidation(err)

	defer conn.Close()

	message := make([]byte, 1024)
	var echoMode bool

	fmt.Println("Server started...")

	for {
		n, addr, err := conn.ReadFromUDP(message)

		utils.ErrorValidation(err)

		fmt.Printf("received: %d bytes from: %s\n", n, addr)

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

		if echoMode {
			_, err := conn.WriteTo(message[0:n], addr)
			utils.ErrorValidation(err)
		} else {
			_, err := conn.WriteTo([] byte{10}, addr)
			utils.ErrorValidation(err)
		}
	}
}
