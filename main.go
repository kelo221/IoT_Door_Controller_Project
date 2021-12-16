package main

import (
	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	embeddedAddress = "192.168.1.38"
	embeddedPort    = "8081"
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

	//	The server sends out the lock status and if a door has to be opened

	/*	//	The server receives a package from the embedded device, which contains the RFID data only.
		tcpPacketIn := RFID_MESSAGE{
			state:         protoimpl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			RFID_CODE:     "",
		}
	*/
	go tcpListenerLoop()
	handleDatabase()
	createAccounts()
	createRDIF()
	handleHTTP()
}
