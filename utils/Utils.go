package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const BytesInKB = 1024

type TimeMeasurement struct {
	InitialTime int64
	FinalTime   int64
}

func (r TimeMeasurement) GetTotalTimeInMilliseconds() float64 {
	return float64(r.Difference() / 1e6) // Milliseconds
}

func (r TimeMeasurement) Difference() float64 {
	return float64(r.FinalTime - r.InitialTime) // Nanoseconds
}

func (r TimeMeasurement) GetRTT() {
	fmt.Printf("Round Trip Time Measurement:\n\nIntitial TimeMeasurement: %d nanoseconds\nFinal TimeMeasurement: %d nanoseconds\nTotal Total TimeMeasurement: %f nanoseconds (%f milliseconds)\n", r.InitialTime, r.FinalTime, r.Difference(), r.GetTotalTimeInMilliseconds())
}

func CreateFilledArray(size int) [] byte {
	filledArray := make([] byte, size-1)

	filledArray = append(filledArray, '\n')

	return filledArray
}

func CurrentTimeNano() int64 {
	return time.Now().UnixNano()
}

func ErrorValidation(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetInBytesOrKiloBytes(bytes bool) int {
	var bytesToSend int

	if bytes {
		fmt.Println("How many bytes to send: ")
		_, e := fmt.Scanf("%d", &bytesToSend)
		ErrorValidation(e)

		return bytesToSend
	}

	fmt.Println("How many bytes to send (in KB): ")
	_, e := fmt.Scanf("%d", &bytesToSend)
	ErrorValidation(e)

	return bytesToSend * BytesInKB
}

func MeasureRTT(conn net.Conn) {
	for {
		bytes := GetInBytesOrKiloBytes(true)

		rtt := TimeMeasurement{InitialTime: CurrentTimeNano()}
		conn.Write(CreateFilledArray(bytes))

		_, _ = bufio.NewReader(conn).ReadBytes('\n')
		rtt.FinalTime = CurrentTimeNano()

		rtt.GetRTT()
	}
}

func MeasureTotalTime(conn net.Conn) {
	for {
		bytes := GetInBytesOrKiloBytes(true)

		if isMultipleOfOneMegabyte(bytes) {
			messageAmount := determineMessagesToSend(bytes)

			fmt.Printf("Sending %d messages\n", messageAmount)

			rtt := TimeMeasurement{InitialTime: CurrentTimeNano()}

			SendMultipleMessages(conn, bytes, messageAmount)

			rtt.FinalTime = CurrentTimeNano()

			fmt.Printf("Total TimeMeasurement: %f Milliseconds\n", float64(rtt.Difference())/1e6)
		} else {
			fmt.Println("Not a multiple of a MegaByte.")
		}
	}
}

func SendMultipleMessages(conn net.Conn, bytes int, messageAmount int) {
	for i := 0; i < messageAmount; i++ {
		conn.Write(CreateFilledArray(bytes))

		_, _ = bufio.NewReader(conn).ReadBytes('\n')
	}
}

func isMultipleOfOneMegabyte(bytes int) bool {
	return 1048576/bytes%2 == 0
}

func determineMessagesToSend(bytes int) int {
	return 1048576 / bytes
}
