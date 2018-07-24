package main

import (
	"net"
	"fmt"
	"encoding/json"
//	"encoding/binary"
	"encoding/binary"
	"time"
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
		return
	}
	defer func() {
		conn.Close()
		fmt.Println("Disconnecting")
	} ()
	for {


		requestByteArray := make([]byte, 4)

		var amount uint32

		myByteArray, err := json.Marshal(param)
		if err != nil {
			fmt.Println("error marshaling:", err)
			break
		}
		amount = uint32(len(myByteArray))
		binary.BigEndian.PutUint32(requestByteArray, amount)
		for _, bytePart := range myByteArray {
			requestByteArray = append(requestByteArray, bytePart)
		}
		x, err := conn.Write(requestByteArray)
		fmt.Println(x, "bytes sending ", string(requestByteArray))
		buff := make([]byte, 1024)
		x, err = conn.Read(buff)
		if err !=nil{
			fmt.Println("error receive:", err)
			break
		}
		var forecastNow requestReturn
		amount = binary.BigEndian.Uint32([]byte(buff[0:4]))
		fmt.Println("amount =", amount)
		fmt.Println("Unmarshalling", string(buff[4:amount+4]))
		err = json.Unmarshal(buff[4:amount+4], &forecastNow)
		if err != nil {
			fmt.Println("Unmarshal error: ", err)
			break
		}

		fmt.Println("Closed forecast for", time.Unix(forecastNow.Object.Date, 0), "City", forecastNow.Object.City, "Temperature is", forecastNow.Object.Temp)
		fmt.Println("received ", x, "bytes")
		fmt.Println("received string", string(buff))
		param.Command = "closeConnection"
		myByteArray, err = json.Marshal(param)
		if err != nil {
			fmt.Println("error marshaling:", err)
			break
		}
		fmt.Println("Marshaled string:", string(myByteArray))
		amount = uint32(len(myByteArray))
		fmt.Println()
		fmt.Println("amount bytes in marshalled string", amount)
		requestByteArray = make([]byte, 4)
		binary.BigEndian.PutUint32(requestByteArray, amount)
		for _, bytePart := range myByteArray {
			requestByteArray = append(requestByteArray, bytePart)
		}
		x, err = conn.Write(requestByteArray)
		fmt.Println(x, "bytes sending ", string(requestByteArray))



		break
	}



//	x, err = conn.Write([]byte(StopCharacter))
//	fmt.Println(x, "bytes sent")
//	x, err = conn.Read(returnByteArray)
//	fmt.Println(x, "bytes received")
//	fmt.Println(string(returnByteArray))



}
