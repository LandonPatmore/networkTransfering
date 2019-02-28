package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const BytesInKB = 1024

// Keeps track of the initial time and final time
type TimeMeasurement struct {
	InitialTime int64
	FinalTime   int64
}

// Determines the total amount of time in milliseconds
func (r TimeMeasurement) GetTotalTimeInMilliseconds() float64 {
	return float64(r.Difference() / 1e6) // Milliseconds
}

// Determines the difference between the final time and initial time in nanoseconds
func (r TimeMeasurement) Difference() float64 {
	return float64(r.FinalTime - r.InitialTime) // Nanoseconds
}

// Gets the Round-Trip Time for data sent to a server
func (r TimeMeasurement) GetRTT() {
	fmt.Printf("Round Trip Time Measurement:\n\nIntitial TimeMeasurement: %d nanoseconds\nFinal TimeMeasurement: %d nanoseconds\nTotal Total TimeMeasurement: %f nanoseconds (%f milliseconds)\n", r.InitialTime, r.FinalTime, r.Difference(), r.GetTotalTimeInMilliseconds())
}

// Creates an array filled with bytes determined by a size
// entered by a user
func CreateFilledArray(size int) [] byte {
	filledArray := make([] byte, size-1)

	filledArray = append(filledArray, '\n')

	return filledArray
}

// Gets the current time in nanoseconds
func CurrentTimeNano() int64 {
	return time.Now().UnixNano()
}

// Checks if there are any errors panics if there are
func ErrorValidation(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Determines whether or not the input should be converted
// to KB or kept in bytes
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

// Measures the Return-Trip time of a stream of bytes
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

// Measures the total time it takes to send (n) amount of bytes in
// (m) amount of messages
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

// Sends multiple messages in sequential order
func SendMultipleMessages(conn net.Conn, bytes int, messageAmount int) {
	for i := 0; i < messageAmount; i++ {
		conn.Write(CreateFilledArray(bytes))

		_, _ = bufio.NewReader(conn).ReadBytes('\n')
	}
}

// Checks if a number is a multiple of a Megabyte
func isMultipleOfOneMegabyte(bytes int) bool {
	return 1048576/bytes%2 == 0
}

// Determines the amount of messages to send for a certain
// amount of bytes
func determineMessagesToSend(bytes int) int {
	return 1048576 / bytes
}
