package main

import (
	"net"
	"fmt"
	"encoding/json"
//	"encoding/binary"
	"encoding/binary"
)

type requestError struct {
	Message string `json:"message,omitempty"`
}

type requestObject struct {
	City     string  `json:"city"`
	Date     int64   `json:"date"`
	Pressure float64 `json:"pressure,omitempty"`
	Humidity float64 `json:"humidity,omitempty"`
	Temp     float64 `json:"temp,omitempty"`
}

type requestInfo struct {
	Command string        `json:"command"`
	Params  requestObject `json:"params,omitempty"`
}

type requestReturn struct {
	Command string        `json:"command"`
	Error   requestError  `json:"error,omitempty"`
	Object  requestObject `json:"object,omitempty"`
}


func main() {
	param := new(requestInfo)
	param.Command = "GetWeather"
	param.Params.City = "Moscow"
	param.Params.Date = 1532255871

	conn, err := net.Dial("tcp", "localhost:7777")
	if err != nil {
		fmt.Println("error connection: ", err)
	}
	defer func() {
		conn.Close()
		fmt.Println("Disconnected")
	} ()

	requestByteArray := make([]byte, 4)
	var amount uint32
	amount = 5

	myByteArray, err := json.Marshal(param)
	binary.BigEndian.PutUint32(requestByteArray, amount)
	for _, x := range myByteArray {
		requestByteArray = append(requestByteArray, x)
	}
	x, err := conn.Write(requestByteArray)
	fmt.Println(x, "bytes sending ", string(requestByteArray))
	x, err = conn.Read(requestByteArray)
	fmt.Println(string(requestByteArray))



}
