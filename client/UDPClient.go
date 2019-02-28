package main

import (
	"fmt"
	"goNetworkTransfering/utils"
	"net"
	"os"
)

func main() {
	var serverAddress string
	fmt.Println("Server address: ")
	_, _ = fmt.Scanf("%s", &serverAddress)

	conn, connError := net.Dial("udp", serverAddress+":8274")

	utils.ErrorValidation(connError)

	UDPDetermineClientMode(conn)
}

// Determines the mode to put the client into
func UDPDetermineClientMode(conn net.Conn) {
	var mode string

	for {
		fmt.Println("What mode?\n1. Measure RTT\n2. Measure Total Time")
		_, e := fmt.Scanf("%s", &mode)
		utils.ErrorValidation(e)
		switch mode {
		case "1":
			fmt.Println("Measuring Time...")
			utils.MeasureRTT(conn)
			break
		case "2":
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
