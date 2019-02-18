package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	var serverAddress string
	fmt.Println("Server address: ")
	_, _ = fmt.Scanf("%s", &serverAddress)

	conn, connError := net.Dial("tcp", serverAddress+":8274")

	if connError != nil {
		log.Fatal(connError)
	}

	for {
		// read in input from stdin
		fmt.Println("How many bytes to send: ")
		var converted int
		_, e := fmt.Scanf("%d", &converted)

		if e != nil {
			log.Fatal(e)
		}

		// send to socket
		roundTripTime := sendData(conn, converted)

		// listen for reply
		m, _ := bufio.NewReader(conn).ReadBytes('\n')
		roundTripTime.finalTime = currentTimeNano()
		fmt.Println("Bytes from server: ", m)
		fmt.Printf("Round Trip Time:\n\nIntitial Time: %d nanoseconds\nFinal Time: %d nanoseconds\nTotal Time: %f milliseconds\n", roundTripTime.initialTime, roundTripTime.finalTime, roundTripTime.difference())
	}
}

func sendData(conn net.Conn, converted int) rtt {
	go conn.Write(createFilledByteArray(converted))
	return rtt{initialTime: currentTimeNano()}
}

func createFilledByteArray(size int) [] byte {
	filledArray := make([] byte, size-1)

	filledArray = append(filledArray, '\n')

	return filledArray
}

func currentTimeNano() int64 {
	return time.Now().UnixNano()
}

func (r rtt) difference() float64 {
	return float64(r.finalTime-r.initialTime) / 1000
}

type rtt struct {
	initialTime int64
	finalTime   int64
}
