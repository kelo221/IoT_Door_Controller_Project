package main

import (
	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	embeddedAddress = "localhost"
	embeddedPort    = "8082"
	tcpListenPort   = "8081"
)

func main() {

	//	The server sends out the lock status and if a door has to be opened
	tcpPacketOut := Door_Request{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		DoorRequest:   LOCK_STATUS_NO_REQUEST,
		LockStatus:    LOCK_STATUS_UNLOCKED,
	}

	/*	//	The server receives a package from the embedded device, which contains the RFID data only.
		tcpPacketIn := RFID_MESSAGE{
			state:         protoimpl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			RFID_CODE:     "",
		}
	*/
	go tcpListenerLoop(&tcpPacketOut)
	handleDatabase()
	createAccounts()
	createRDIF()
	handleHTTP(&tcpPacketOut)
}
