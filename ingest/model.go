package main

type ECU struct {
	ID       int `json:"id"`
	MotorRPM int `json:"motor_rpm"`
	Speed    int `json:"speed"`
	Throttle     int `json:"throttle"`
	BrakePressure int `json:"brake_pressure"`
	CreatedAt     time.Time `json:"created_at"`
}

type Battery struct {
	ID int `json:"id"`
	ChargeLevel int `json:"charge_level"`
	CellTemp1 int `json:"cell_temp_1"`
	CellTemp2 int `json:"cell_temp_2"`
	CellTemp3 int `json:"cell_temp_3"`
	CellTemp4   int `json:"cell_temp_4"`
	CellVoltage1 int `json:"cell_voltage_1"`
	CellVoltage2 int `json:"cell_voltage_2"`
	CellVoltage3 int `json:"cell_voltage_3"`
	CellVoltage4 int `json:"cell_voltage_4"`
	CreatedAt     time.Time `json:"created_at"`
}

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

	// Scale throttle to 0-100, assuming 20000 is 0 and 450000 is max
	ecu.Throttle = (ecu.Throttle - 20000) * 100 / (450000 - 20000)
	return ecu
}

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