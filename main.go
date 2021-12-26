package main

import (
	"google.golang.org/protobuf/runtime/protoimpl"
)

var tcpPacketOut = Door_Request{
	state:         protoimpl.MessageState{},
	sizeCache:     2,
	unknownFields: nil,
	DoorRequest:   Door_Request_NO_REQUEST,
	LockStatus:    Door_Request_UNLOCKED,
}

func main() {
	tcpPacketOut.Reset()
	go handleMQTTIn()
	handleDatabase()
	createAccounts()
	createRDIF()
	handleHTTP()
}
