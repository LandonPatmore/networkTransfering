package shared

import (
	"log"
	"net"
	"time"
)

const BytesInKB = 1024

func createFilledBytesArray(size int) [] byte {
	filledArray := make([] byte, size-1)

	filledArray = append(filledArray, '\n')

	return filledArray
}

func CurrentTimeNano() int64 {
	return time.Now().UnixNano()
}

func (r RTT) Difference() float64 {
	return float64(r.FinalTime - r.InitialTime)
}

func ErrorValidation(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (r *RTT) GetInfo() {
	r.FinalTime = CurrentTimeNano()
	log.Printf("Round Trip Time:\n\nIntitial Time: %d nanoseconds\nFinal Time: %d nanoseconds\nTotal Time: %f nanoseconds (%f milliseconds)\n", r.InitialTime, r.FinalTime, r.Difference(), r.Difference() / 1e6)
}

func SendData(conn net.Conn, converted int) RTT {
	go conn.Write(createFilledBytesArray(converted))
	return RTT{InitialTime: CurrentTimeNano()}
}

type RTT struct {
	InitialTime int64
	FinalTime   int64
}
