package main

import (
	"google.golang.org/protobuf/runtime/protoimpl"
)

func main() {

	//	The server sends out the lock status and if a door has to be opened
	tcpPacketOut := LOCK_STATUS{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		DoorOpenInQue: false,
		LockStatus:    LOCK_STATUS_UNLOCKED,
	}

	//	The server receives a package from the embedded device, which contains the RFID data only.
	tcpPacketIn := RFID_MESSAGE{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		RFID_MESSAGE:  "",
	}

	go tcpListenerLoop(&tcpPacketIn)
	handleDatabase()
	createAccounts()
	createRDIF()
	handleHTTP(&tcpPacketOut)
}
