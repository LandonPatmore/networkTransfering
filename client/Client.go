package main

import (
	"fmt"
	"goNetworkTransfering/utils"
	"net"
	"os"
	"time"
)

func main() {
	t := determineClientType()
	createClient(t)
}

func determineClientType() int {
	var clientType int
	for {
		fmt.Println("What mode?\n1. TCP\n2. UDP")
		_, _ = fmt.Scanf("%d", &clientType)

		switch clientType {
		case 1:
			fallthrough
		case 2:
			return clientType
		default:
			fmt.Println("Not a client type...")
		}
	}
}

func createClient(t int) {
	outputData := utils.CreateOutputData()

	var serverAddress string
	fmt.Println("Server address: ")
	_, _ = fmt.Scanf("%s", &serverAddress)

	if t == 1 {
		outputData.Server = serverAddress
		outputData.Type = "TCP"

		conn, connError := net.Dial("tcp", serverAddress+":8721")

		utils.ErrorValidation(connError)

		ClientMode(conn, outputData, "TCP")
	} else {
		outputData.Server = serverAddress
		outputData.Type = "UDP"

		conn, connError := net.Dial("udp", serverAddress+":8274")

		utils.ErrorValidation(connError)

		ClientMode(conn, outputData, "UDP")
	}
}

// Determines the mode to put the client into
func ClientMode(conn net.Conn, outputData utils.OutputData, t string) {
	var mode int

	for {
		fmt.Println(t)
		fmt.Println("What mode?\n1. Measure RTT\n2. Measure Throughput (TCP Only)\n3. Measure Total Time\n4. Output\n5. Exit")
		_, e := fmt.Scanf("%d", &mode)
		utils.ErrorValidation(e)
		switch mode {
		case 1:
			fmt.Println("Measuring Time...")
			changeServerMode(conn, true)
			rtts := utils.MeasureRTT(conn)
			outputData.RTT = append(outputData.RTT, rtts...)
			break
		case 2:
			fmt.Println("Measuring Throughput...(TCP Only)")
			changeServerMode(conn, true)
			throughput := utils.MeasureThroughput(conn)
			outputData.Throughput = append(outputData.Throughput, throughput...)
			break
		case 3:
			fmt.Println("Measuring Total Time...")
			changeServerMode(conn, false)
			messageSizeTime := utils.MeasureTotalTime(conn)
			outputData.MessageSizeTime = append(outputData.MessageSizeTime, messageSizeTime...)
			break
		case 4:
			fmt.Println("Outputting Data...")
			utils.CreateJSONFile(t+"_"+outputData.Server, outputData)
			fmt.Println("Data Output...")
		case 5:
			fmt.Println("Exiting...")
			utils.ErrorValidation(conn.Close())
			os.Exit(0)
		default:
			fmt.Println("Error, not a selectable mode.")
		}
	}
}

func changeServerMode(conn net.Conn, echo bool) {
	if echo {
		conn.Write([] byte{1, 10})
		fmt.Println("Changed server mode to echo mode")
	} else {
		conn.Write([] byte{2, 10})
		fmt.Println("Changed server mode to acknowledge mode")
	}

	time.Sleep(2 * time.Second)
}
