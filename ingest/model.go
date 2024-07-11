package main

import (
	"encoding/binary"
	"time"
)

/*
CAN Frame Format:

+--------+-----+-----------+---------+---------+----------+---------+----------+----------+----------+
| Node   | ID  | Byte 0    | Byte 1  | Byte 2  | Byte 3   | Byte 4  | Byte 5   | Byte 6   | Byte 7   |
+--------+-----+-----------+---------+---------+----------+---------+----------+----------+----------+
| ECU    | 100 |       Motor RPM     |  Speed  |     Throttle       |    Brake Pressure   |          |
+--------+-----+-----------+---------+---------+----------+---------+----------+----------+----------+
| Battery| 200 | Charge    |         |         |          |         |          |          |          |
|        |     | Level     |         |         |          |         |          |          |          |
+--------+-----+-----------+---------+---------+----------+---------+----------+----------+----------+
| Battery| 201 | Cell 1    | Cell 1  | Cell 2  | Cell 2   | Cell 3  | Cell 3   | Cell 4   | Cell 4   |
|        |     | Temp      | Voltage | Temp    | Voltage  | Temp    | Voltage  | Temp     | Voltage  |
+--------+-----+-----------+---------+---------+----------+---------+----------+----------+----------+

*/

type ECU struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	MotorRPM      int       `json:"motor_rpm"`
	Speed         int       `json:"speed"`
	Throttle      int       `json:"throttle"`
	BrakePressure int       `json:"brake_pressure"`
	CreatedAt     time.Time `json:"created_at" gorm:"precision:6"`
}

type Battery struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	ChargeLevel  int       `json:"charge_level"`
	CellTemp1    int       `json:"cell_temp_1"`
	CellTemp2    int       `json:"cell_temp_2"`
	CellTemp3    int       `json:"cell_temp_3"`
	CellTemp4    int       `json:"cell_temp_4"`
	CellVoltage1 int       `json:"cell_voltage_1"`
	CellVoltage2 int       `json:"cell_voltage_2"`
	CellVoltage3 int       `json:"cell_voltage_3"`
	CellVoltage4 int       `json:"cell_voltage_4"`
	CreatedAt    time.Time `json:"created_at" gorm:"precision:6"`
}

// ECUFromBytes converts a byte slice to an ECU struct.
// It interprets the byte data according to the CAN frame format and scales the throttle value.
func ECUFromBytes(data []byte) ECU {
	var ecu ECU
	// MotorRPM is bytes 0-1
	ecu.MotorRPM = int(binary.BigEndian.Uint16(data[0:2]))
	// Speed is byte 2
	ecu.Speed = int(data[2])
	// Throttle is bytes 3-4
	ecu.Throttle = int(binary.BigEndian.Uint16(data[3:5]))
	// BrakePressure is bytes 5-6
	ecu.BrakePressure = int(binary.BigEndian.Uint16(data[5:7]))

	ecu.CreatedAt = time.Now()

	// Scale throttle to 0-100, assuming 2000 is 0 and 4500 is max
	ecu.Throttle = (ecu.Throttle - 2000) * 100 / (4500 - 2000)
	return ecu
}

// BatteryFromBytes converts a byte slice to a Battery struct.
// It interprets the byte data according to the CAN frame format for battery information.
func BatteryFromBytes(data []byte) Battery {
	var battery Battery
	// ChargeLevel is byte 0
	battery.ChargeLevel = int(data[0])
	// Rest of Frame 200 is empty
	// Cell 1 Temp is byte 8
	battery.CellTemp1 = int(data[8])
	// Cell 1 Voltage is byte 9
	battery.CellVoltage1 = int(data[9])
	// Cell 2 Temp is byte 10
	battery.CellTemp2 = int(data[10])
	// Cell 2 Voltage is byte 11
	battery.CellVoltage2 = int(data[11])
	// Cell 3 Temp is byte 12
	battery.CellTemp3 = int(data[12])
	// Cell 3 Voltage is byte 13
	battery.CellVoltage3 = int(data[13])
	// Cell 4 Temp is byte 14
	battery.CellTemp4 = int(data[14])
	// Cell 4 Voltage is byte 15
	battery.CellVoltage4 = int(data[15])

	battery.CreatedAt = time.Now()
	return battery
}

// ecuCallbacks is a slice of functions to be called when new ECU data is available.
var ecuCallbacks []func(ecu ECU)

// batteryCallbacks is a slice of functions to be called when new Battery data is available.
var batteryCallbacks []func(battery Battery)

// RegisterECUCallback adds a new callback function to be executed when ECU data is pushed.
func RegisterECUCallback(callback func(ecu ECU)) {
	ecuCallbacks = append(ecuCallbacks, callback)
}

// RegisterBatteryCallback adds a new callback function to be executed when Battery data is pushed.
func RegisterBatteryCallback(callback func(battery Battery)) {
	batteryCallbacks = append(batteryCallbacks, callback)
}

// PushECU executes all registered ECU callbacks with the provided ECU data.
func PushECU(ecu ECU) {
	for _, callback := range ecuCallbacks {
		callback(ecu)
	}
}

// PushBattery executes all registered Battery callbacks with the provided Battery data.
func PushBattery(battery Battery) {
	for _, callback := range batteryCallbacks {
		callback(battery)
	}
}
