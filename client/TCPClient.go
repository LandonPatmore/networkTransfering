package main

import (
	"bufio"
	"fmt"
	"goNetworkTransfering/shared"
	"log"
	"net"
	"strconv"
)

func main() {
	var serverAddress string
	fmt.Println("Server address: ")
	_, _ = fmt.Scanf("%s", &serverAddress)

	conn, connError := net.Dial("tcp", serverAddress+":8274")

	shared.ErrorValidation(connError)

	for {
		// read in input from stdin
		fmt.Println("How many bytes to send: ")

		var bytesToSend int
		_, e := fmt.Scanf("%d", &bytesToSend)

		shared.ErrorValidation(e)

		// send to TCP Socket
		roundTripTime := shared.SendData(conn, bytesToSend)

		// listen for reply from Server
		_, _ = bufio.NewReader(conn).ReadBytes('\n')

		roundTripTime.GetInfo()
		log.Println("Throughput: " + strconv.FormatFloat(measureThroughput(roundTripTime, bytesToSend), 'f', 6, 64))
	}
}

func measureThroughput(rtt shared.RTT, numBytesSent int) float64 {
	bitsSent := float64(8 * numBytesSent)
	rttInSeconds := rtt.Difference()
	return (bitsSent / (rttInSeconds / 2)) * 1000 // Megabits/sec
}
