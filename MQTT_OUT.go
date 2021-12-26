package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

func handleMQTTOut() {

	fmt.Println("sending new package...")

	data, err := proto.Marshal(&tcpPacketOut)
	if err != nil {
		fmt.Println("marshall error", err)

	}

	fmt.Println(tcpPacketOut.GetLockStatus())
	fmt.Println(tcpPacketOut.GetDoorRequest())
	fmt.Println(data)

	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.1.77:1883")

	opts.SetClientID("CLIENT_ID" + strconv.FormatInt(time.Now().Unix(), 10))

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := c.Publish("door/request", 0, false, data)
	token.Wait()

}
