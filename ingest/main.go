package main

import (
	"fmt"
	"log"
	"github.com/eclipse/paho.mqtt.golang"
)

func main() {
	ConnectMQTT()
}

var Client mqtt.Client

var MQTTHost = "localhost"
var MQTTPort = "1883"

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

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("[MQ] Connected to MQTT broker!")
}
