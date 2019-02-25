package shared

import (
	"fmt"
	"log"
	"net"
	"time"
)

func createFilledBytesArray(size int) [] byte {
	filledArray := make([] byte, size-1)

	filledArray = append(filledArray, '\n')

	return filledArray
}

func CurrentTimeNano() int64 {
	return time.Now().UnixNano()
}

func (r RTT) Difference() float64 {
	return float64(r.FinalTime-r.InitialTime) / 1e6 // Milliseconds
}

func ErrorValidation(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (r *RTT) GetInfo() {
	r.FinalTime = CurrentTimeNano()
	fmt.Printf("Round Trip Time:\n\nIntitial Time: %d nanoseconds\nFinal Time: %d nanoseconds\nTotal Time: %f Milliseconds\n", r.InitialTime, r.FinalTime, r.Difference())
}

func SendData(conn net.Conn, converted int) RTT {
	go conn.Write(createFilledBytesArray(converted))
	return RTT{InitialTime: CurrentTimeNano()}
}

type RTT struct {
	InitialTime int64
	FinalTime   int64
}