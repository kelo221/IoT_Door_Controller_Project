package main

import (
	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	embeddedAddress = "192.168.1.77"
	embeddedPort    = "8082"
	tcpListenPort   = "8082"
)

var tcpPacketOut = Door_Request{
	state:         protoimpl.MessageState{},
	sizeCache:     0,
	unknownFields: nil,
	DoorRequest:   LOCK_STATUS_NO_REQUEST,
	LockStatus:    LOCK_STATUS_UNLOCKED,
}

func main() {
	go handleMQTTIn()
	handleDatabase()
	createAccounts()
	createRDIF()
	handleHTTP()
}
