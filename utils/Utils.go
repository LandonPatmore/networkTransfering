package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const BytesInKB = 1024

type OutputData struct {
	Client          string             `json:"client"`
	Server          string             `json:"server"`
	Type            string             `json:"type"`
	RTT             [] RTT             `json:"rtt"`
	Throughput      [] Throughput      `json:"throughput"`
	MessageSizeTime [] MessageSizeTime `json:"message_size_time"`
}

type RTT struct {
	MessageSize int     `json:"message_size"`
	TotalTime   float64 `json:"total_time"`
}

type Throughput struct {
	MessageSize int     `json:"message_size"`
	Megabits    float64 `json:"megabits"`
}

type MessageSizeTime struct {
	MessageSize   int     `json:"message_size"`
	MessageAmount int     `json:"message_amount"`
	TotalTime     float64 `json:"total_time"`
}

// Keeps track of the initial time and final time
type TimeMeasurement struct {
	InitialTime int64
	FinalTime   int64
}

// Helper to create a struct for outputting data
func CreateOutputData() OutputData {
	return OutputData{
		Client:          "",
		Server:          "",
		Type:            "",
		RTT:             make([] RTT, 0),
		Throughput:      make([] Throughput, 0),
		MessageSizeTime: make([] MessageSizeTime, 0),
	}
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
	fmt.Printf("Round Trip Time Measurement:\n\nIntitial TimeMeasurement: %d nanoseconds\nFinal TimeMeasurement: %d nanoseconds\nTotal Total TimeMeasurement: %f nanoseconds (%f milliseconds)\n\n", r.InitialTime, r.FinalTime, r.Difference(), r.GetTotalTimeInMilliseconds())
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
func MeasureRTT(conn net.Conn) [] RTT {
	var rtts [] RTT

	bytes := GetInBytesOrKiloBytes(true)

	for i := 0; i < 5; i++ {
		timeMeasurement := TimeMeasurement{InitialTime: CurrentTimeNano()}
		conn.Write(CreateFilledArray(bytes))

		_, _ = bufio.NewReader(conn).ReadBytes('\n')
		timeMeasurement.FinalTime = CurrentTimeNano()

		timeMeasurement.GetRTT()

		rtts = append(rtts, RTT{MessageSize: bytes, TotalTime: timeMeasurement.GetTotalTimeInMilliseconds()})

		time.Sleep(1 * time.Second)
	}

	return rtts
}

// Measures the throughput in Megabits/sec
func MeasureThroughput(conn net.Conn) [] Throughput {
	var throughput [] Throughput

	bytes := GetInBytesOrKiloBytes(false)

	for i := 0; i < 5; i++ {

		rtt := TimeMeasurement{InitialTime: CurrentTimeNano()}
		conn.Write(CreateFilledArray(bytes))

		_, _ = bufio.NewReader(conn).ReadBytes('\n')

		rtt.FinalTime = CurrentTimeNano()

		bitsSent := float64(8 * bytes)
		rttInNanoSeconds := rtt.Difference()
		megabits := ((bitsSent / (rttInNanoSeconds / 2)) * 1000) / 8
		fmt.Printf("Throughput: %f Megabytes/sec\n", megabits) //Converts bits/nano to megabits/sec

		throughput = append(throughput, Throughput{MessageSize: bytes, Megabits: megabits})

		time.Sleep(1 * time.Second)
	}

	return throughput
}

// Measures the total time it takes to send (n) amount of bytes in
// (m) amount of messages
func MeasureTotalTime(conn net.Conn) [] MessageSizeTime {
	var messageSizeTime [] MessageSizeTime

	bytes := GetInBytesOrKiloBytes(true)

	for i := 0; i < 5; i++ {

		if isMultipleOfOneMegabyte(bytes) {
			messageAmount := determineMessagesToSend(bytes)

			fmt.Printf("Sending %d messages\n", messageAmount)

			rtt := TimeMeasurement{InitialTime: CurrentTimeNano()}

			SendMultipleMessages(conn, bytes, messageAmount)

			rtt.FinalTime = CurrentTimeNano()

			totalTime := float64(rtt.Difference()) / 1e6
			fmt.Printf("Total TimeMeasurement: %f Milliseconds\n", totalTime)

			messageSizeTime = append(messageSizeTime, MessageSizeTime{MessageSize: bytes, MessageAmount: messageAmount, TotalTime: totalTime})

			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("Not a multiple of a MegaByte.")
			break
		}
	}

	return messageSizeTime
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

func CreateJSONFile(name string, v interface{}) {
	fileWriter, _ := os.Create(name + ".json")
	ErrorValidation(json.NewEncoder(fileWriter).Encode(v))
}
