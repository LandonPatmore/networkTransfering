package main

import (
	"bufio"
	"fmt"
	"goNetworkTransfering/utils"
	"net"
	"os"
)

func main() {
	var serverAddress string
	fmt.Println("Server address: ")
	_, _ = fmt.Scanf("%s", &serverAddress)

	conn, connError := net.Dial("tcp", serverAddress+":8721")

	utils.ErrorValidation(connError)

	TCPDetermineClientMode(conn)
}

func TCPDetermineClientMode(conn net.Conn) {
	var mode string

	for {
		fmt.Println("What mode?\n1. Measure RTT\n2. Measure Throughput\n3. Measure Total Time")
		_, e := fmt.Scanf("%s", &mode)
		utils.ErrorValidation(e)
		switch mode {
		case "1":
			fmt.Println("Measuring Time...")
			utils.MeasureRTT(conn)
			break
		case "2":
			fmt.Println("Measuring Throughput...")
			measureThroughput(conn)
			break
		case "3":
			fmt.Println("Measuring Total Time...")
			utils.MeasureTotalTime(conn)
			break
		case "exit":
			fmt.Println("Exiting...")
			utils.ErrorValidation(conn.Close())
			os.Exit(0)
		default:
			fmt.Println("Error, not a selectable mode.")
		}
	}
}

func measureThroughput(conn net.Conn) {
	for {
		bytes := utils.GetInBytesOrKiloBytes(false)

		rtt := utils.TimeMeasurement{InitialTime: utils.CurrentTimeNano()}
		conn.Write(utils.CreateFilledArray(bytes))

		_, _ = bufio.NewReader(conn).ReadBytes('\n')

		rtt.FinalTime = utils.CurrentTimeNano()

		bitsSent := float64(8 * bytes)
		rttInNanoSeconds := rtt.Difference()
		fmt.Printf("Throughput: %f Megabits/sec\n", (bitsSent/(rttInNanoSeconds/2))*1000) //Converts bits/nano to megabits/sec
	}
}
