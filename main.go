package main

import (
	"fmt"
	"google.golang.org/protobuf/runtime/protoimpl"
	"net"
)

func main() {

	doorStatus := Lock{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		DoorOpenInQue: false,
		LockStatus:    Lock_UNLOCKED,
	}

	go tcpListenerLoop(&doorStatus)
	handleDatabase()
	createAccounts()
	handleHTTP(&doorStatus)
}

func tcpListenerLoop(doorStatus *Lock) {

	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			go getTpcPackage(conn, doorStatus)
		}
	}

}
