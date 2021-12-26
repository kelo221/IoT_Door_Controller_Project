package main

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
	"regexp"
	"strings"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	fmt.Println("new package")

	newRFID := strings.TrimSpace(string(msg.Payload()))

	space := regexp.MustCompile(`\s+`)
	newRFID = space.ReplaceAllString(newRFID, " ")

	fmt.Println(newRFID)

	salt := []byte("salt")

	rfidHash := hex.EncodeToString(HashPassword([]byte(newRFID), salt))
	expectedHash := aqlToString("FOR doc IN DOOR_RFID  FILTER doc.HASHED_RFID == \"" + rfidHash + "\"  RETURN doc.HASHED_RFID")

	fmt.Println(rfidHash)

	hashFound := subtle.ConstantTimeCompare([]byte(rfidHash[:]), []byte((expectedHash[:]))) == 1

	if hashFound {
		fmt.Println("match")

		aqlNoReturn("FOR doc IN DOOR_RFID" +
			" FILTER doc.HASHED_RFID == \"" + rfidHash + "\"" +
			"   INSERT {" +
			"   name: doc.RFID_OWNER," +
			"    time: DATE_NOW()" +
			"  } INTO DOOR_HISTORY OPTIONS { ignoreErrors: true }")
		fmt.Println("MATCH")
		tcpPacketOut.DoorRequest = LOCK_STATUS_APPROVED
	} else {
		fmt.Println("NO MATCH")
		tcpPacketOut.DoorRequest = LOCK_STATUS_DISAPPROVED
	}

	data, err := proto.Marshal(&tcpPacketOut)
	if err != nil {
		fmt.Println("marshall error", err)

	}
	fmt.Println(data)

	// send the package on MQTT
	handleMQTTOut()
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func publish(client mqtt.Client) {

	for {
		token := client.Publish("door/RFID", 0, false, nil)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func sub(client mqtt.Client) {
	topic := "door/RFID"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s \n", topic)
}

func handleMQTTIn() {

	var broker = "192.168.1.77"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)
	publish(client)

	client.Disconnect(250)

}
