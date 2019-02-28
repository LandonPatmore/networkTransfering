package main

import (
	"fmt"
	"goNetworkTransfering/utils"
	"net"
	"os"
)

func main() {
	UDPServerMode()
}

func UDPServerMode() {
	var mode string

	for {
		fmt.Println("What mode?\n1. Echo Server\n2. Read Data Server")
		_, e := fmt.Scanf("%s", &mode)
		utils.ErrorValidation(e)
		switch mode {
		case "1":
			fmt.Println("Echo Mode...")
			createUDPServer(true)
			break
		case "2":
			fmt.Println("Read Data Mode...")
			createUDPServer(false)
			break
		case "exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Error, not a selectable mode.")
		}
	}
}

func createUDPServer(echo bool) {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":8274")

	utils.ErrorValidation(err)

	conn, err := net.ListenUDP("udp", ServerAddr)
	utils.ErrorValidation(err)

	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buf)

		utils.ErrorValidation(err)

		fmt.Printf("received: %d bytes from: %s\n", n, addr)

		if echo {
			_, err := conn.WriteTo(buf[0:n], addr)
			utils.ErrorValidation(err)
		} else {
			_, err := conn.WriteTo([] byte{10}, addr)
			utils.ErrorValidation(err)
		}
	}
}
