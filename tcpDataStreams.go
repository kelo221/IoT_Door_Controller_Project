package main

import (
	"bytes"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
)

func getTpcPackage(conn net.Conn, container *Lock) {

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	result := bytes.NewBuffer(nil)
	var buf [1024]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				//	fmt.Println(err)
				break
			}
		} else {
			newMessage := Lock{}
			err = proto.Unmarshal(result.Bytes(), &newMessage)
			if err != nil {
				panic(err)
			}

			// Seems to work even with mutex warning
			*container = newMessage

		}
		result.Reset()
	}

}
