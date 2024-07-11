package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	ConnectMQTT()
	SubscribeECU()
	SubscribeBattery()
	ConnectDB()
	StartServer()
}

var Client mqtt.Client

// MQTT connection variables
var MQTTHost = "localhost"
var MQTTPort = "1883"

// ConnectMQTT establishes a connection to the MQTT broker using the specified host and port.
func ConnectMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", MQTTHost, MQTTPort))
	opts.OnConnect = connectHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	Client = client
}

// connectHandler is a callback function that is called when the MQTT client successfully connects to the broker.
// It just prints a connection confirmation message to the console.
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("[MQ] Connected to MQTT broker!")
}

// SubscribeECU subscribes to the "ingest/ecu" topic on the MQTT broker.
// It processes incoming ECU messages by converting the payload to an ECU object,
// creating a new ECU record in the database, and logging the received message.
func SubscribeECU() {
	Client.Subscribe("ingest/ecu", 0, func(client mqtt.Client, msg mqtt.Message) {
		ecu := ECUFromBytes(msg.Payload())
		ecu, err := CreateECU(ecu)
		if err != nil {
			fmt.Printf("failed to create ecu: %v", err)
		}
		fmt.Printf("[MQ] Received ecu message: %v\n", ecu)
	})
}

// SubscribeBattery subscribes to the "ingest/battery" topic on the MQTT broker.
// It processes incoming Battery messages by converting the payload to a Battery object,
// creating a new Battery record in the database, and logging the received message.
func SubscribeBattery() {
	Client.Subscribe("ingest/battery", 0, func(client mqtt.Client, msg mqtt.Message) {
		battery := BatteryFromBytes(msg.Payload())
		battery, err := CreateBattery(battery)
		if err != nil {
			fmt.Printf("failed to create battery: %v", err)
		}
		fmt.Printf("[MQ] Received battery message: %v\n", battery)
	})
}
