package main

import (
	"fmt"
	"goNetworkTransfering/shared"
	"net"
)

func main() {
	createUDPServer()
}

func createUDPServer() {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":8274")

	shared.ErrorValidation(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	shared.ErrorValidation(err)

	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		fmt.Printf("Accepting packets\n")

		n, addr, err := ServerConn.ReadFromUDP(buf)

		shared.ErrorValidation(err)

		fmt.Printf("received: %d bytes from: %s\n", n, addr)

		_, sendError := ServerConn.WriteTo(buf[0:n], addr)

		shared.ErrorValidation(sendError)
	}
}
