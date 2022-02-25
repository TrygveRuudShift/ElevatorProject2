package network

import (
	//"elevators/controlunit"
	"elevators/filesystem"
	"encoding/json"
	"fmt"
	"net"
	"time"
)


func InitUDPSendingSocket(port int, sendAddr string) (net.UDPAddr, *net.UDPConn) {
	sendaddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(sendAddr),
	}
	wconn, err := net.DialUDP("udp", nil, &sendaddr) // code does not block here
	if err != nil {
		panic(err)
	}
	//defer wconn.Close() //Close at the end of program

	return sendaddr, wconn
}

func BroadcastStruct(orderState filesystem.OrderState, wconn *net.UDPConn) {
	message, _ := json.Marshal(orderState)
	BroadcastMessage(message, wconn)
}

func BroadcastMessage(message []byte, wconn *net.UDPConn) {
	_, err := wconn.Write(message)
	if err != nil {
		panic(err)
	}
	// fmt.Println("You sent: msg: ", message)
}

func TestSendAndReceive() {
	UDPPort := 20014

	var state filesystem.OrderState
	state.Dir = "up"
	state.Floor = 3
	state.Name = "Elevator"

	jsState, _ := json.Marshal(state)
	json.Unmarshal(jsState, &state)
	fmt.Println(string(jsState))
	// fmt.Println(string(state))

	//Initialize sockets
	_, wconn := InitUDPSendingSocket(UDPPort, "255.255.255.255")
	_, conn := InitUDPReceivingSocket(UDPPort)

	//Close sockets when program terminates
	defer conn.Close()
	defer wconn.Close()

	//Send and receive message
	for {
		// BroadcastMessage(json.RawMessage(`{"precomputed": true}`), wconn)
		//BroadcastMessage(string(jsState), wconn)
		BroadcastStruct(state, wconn)
		time.Sleep(time.Millisecond * 1000)
		
		msg := ReceiveUDPMessage(conn)
		// fmt.Println("You received:", msg)
		
		json.Unmarshal(msg, &state)

		//fmt.Println("Your state:", state.Name)
		fmt.Printf("state: %v\n", state)

		time.Sleep(time.Millisecond * 1000)
	}
}