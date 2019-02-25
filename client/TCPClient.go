package main

import (
	"bufio"
	"fmt"
	"goNetworkTransfering/shared"
	"log"
	"net"
)

func main() {
	var serverAddress string
	fmt.Println("Server address: ")
	_, _ = fmt.Scanf("%s", &serverAddress)

	conn, connError := net.Dial("tcp", serverAddress+":8274")

	shared.ErrorValidation(connError)

	for {
		// read in input from stdin
		fmt.Println("How many bytes to send (in KB): ")

		var bytesToSend int
		_, e := fmt.Scanf("%d", &bytesToSend)
		bytesToSend = bytesToSend * shared.BytesInKB

		shared.ErrorValidation(e)

		// send to TCP Socket
		roundTripTime := shared.SendData(conn, bytesToSend)

		// listen for reply from Server
		_, _ = bufio.NewReader(conn).ReadBytes('\n')

		roundTripTime.GetInfo()
		measureThroughput(roundTripTime, bytesToSend)
	}
}

func measureThroughput(rtt shared.RTT, numBytesSent int) {
	bitsSent := float64(8 * numBytesSent)
	rttInNanoSeconds := rtt.Difference()
	log.Printf("Throughput: %f Megabits/sec", (bitsSent/(rttInNanoSeconds/2))*1000)
}
