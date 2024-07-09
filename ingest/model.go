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