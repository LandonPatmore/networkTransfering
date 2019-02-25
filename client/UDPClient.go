package main

import (
	"bufio"
	"fmt"
	"goNetworkTransfering/shared"
	"net"
)

func main() {
	var serverAddress string
	fmt.Println("Server address: ")
	_, _ = fmt.Scanf("%s", &serverAddress)

	conn, connError := net.Dial("udp", serverAddress+":8274")

	shared.ErrorValidation(connError)

	for {
		// read in input from stdin
		fmt.Println("How many bytes to send (in KB): ")

		var bytesToSend int
		_, e := fmt.Scanf("%d", &bytesToSend)
		bytesToSend = bytesToSend * shared.BytesInKB

		shared.ErrorValidation(e)

		// send to UDP Socket
		roundTripTime := shared.SendData(conn, bytesToSend)

		// listen for reply
		_, _ = bufio.NewReader(conn).ReadBytes('\n')

		roundTripTime.GetInfo()
	}
}
