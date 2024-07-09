package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MQTTHost = "localhost"
var MQTTPort = "1883"

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", MQTTHost, MQTTPort))
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	for {
		ecu := generateECU()
		battery := generateBattery()

		ecuBytes := ECUToBytes(ecu)
		batteryBytes := BatteryToByteArray(battery)

		if token := client.Publish("ingest/ecu", 0, false, ecuBytes); token.Wait() && token.Error() != nil {
			log.Printf("Failed to publish ECU: %v", token.Error())
		} else {
			fmt.Println("ECU published")
		}

		if token := client.Publish("ingest/battery", 0, false, batteryBytes); token.Wait() && token.Error() != nil {
			log.Printf("Failed to publish Battery: %v", token.Error())
		} else {
			fmt.Println("Battery published")
		}

		time.Sleep(1 * time.Second)
	}
}

type ECU struct {
	MotorRPM      int
	Speed         int
	Throttle      int
	BrakePressure int
}

type Battery struct {
	ChargeLevel  int `json:"charge_level"`
	CellTemp1    int `json:"cell_temp_1"`
	CellTemp2    int `json:"cell_temp_2"`
	CellTemp3    int `json:"cell_temp_3"`
	CellTemp4    int `json:"cell_temp_4"`
	CellVoltage1 int `json:"cell_voltage_1"`
	CellVoltage2 int `json:"cell_voltage_2"`
	CellVoltage3 int `json:"cell_voltage_3"`
	CellVoltage4 int `json:"cell_voltage_4"`
}

func generateECU() ECU {
	return ECU{
		MotorRPM:      generateValue(0, 10000),
		Speed:         generateValue(0, 100),
		Throttle:      generateValue(2000, 4500),
		BrakePressure: generateValue(0, 10000),
	}
}

func generateBattery() Battery {
	return Battery{
		ChargeLevel:  generateValue(0, 100),
		CellTemp1:    generateValue(0, 100),
		CellTemp2:    generateValue(0, 100),
		CellTemp3:    generateValue(0, 100),
		CellTemp4:    generateValue(0, 100),
		CellVoltage1: generateValue(0, 100),
		CellVoltage2: generateValue(0, 100),
		CellVoltage3: generateValue(0, 100),
		CellVoltage4: generateValue(0, 100),
	}
}

func generateValue(min int, max int) int {
	return rand.Intn(max-min) + min
}

func ECUToBytes(ecu ECU) []byte {
	result := make([]byte, 8)

	// Byte 0-1: Motor RPM (uint16, big-endian)
	result[0] = byte(ecu.MotorRPM >> 8)
	result[1] = byte(ecu.MotorRPM)

	// Byte 2: Speed (uint8)
	result[2] = byte(ecu.Speed)

	// Byte 3-4: Throttle (uint16, big-endian)
	result[3] = byte(ecu.Throttle >> 8)
	result[4] = byte(ecu.Throttle)

	// Byte 5-6: Brake Pressure (uint16, big-endian)
	result[5] = byte(ecu.BrakePressure >> 8)
	result[6] = byte(ecu.BrakePressure)

	// Byte 7: Blank (left as 0)

	return result
}

func BatteryToByteArray(battery Battery) []byte {
	result := make([]byte, 16)

	// Byte 0: Charge Level (1 byte)
	result[0] = byte(battery.ChargeLevel)

	// Bytes 1-7: Blank (7 bytes)
	// These are already initialized to 0, so we can skip them

	// Byte 8: Cell Temp 1
	result[8] = byte(battery.CellTemp1)

	// Byte 9: Cell Voltage 1
	result[9] = byte(battery.CellVoltage1)

	// Byte 10: Cell Temp 2
	result[10] = byte(battery.CellTemp2)

	// Byte 11: Cell Voltage 2
	result[11] = byte(battery.CellVoltage2)

	// Byte 12: Cell Temp 3
	result[12] = byte(battery.CellTemp3)

	// Byte 13: Cell Voltage 3
	result[13] = byte(battery.CellVoltage3)

	// Byte 14: Cell Temp 4
	result[14] = byte(battery.CellTemp4)

	// Byte 15: Cell Voltage 4
	result[15] = byte(battery.CellVoltage4)

	return result
}
