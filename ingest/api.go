package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Port is the port number that the server will listen on.
var Port = "9000"

// StartServer initializes and starts the API server.
// It sets up all our endpoints.
func StartServer() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ecu", GetAllECUs)
	r.GET("/battery", GetAllBatteries)
	r.GET("/ecu/averages", GetECUAverages)
	r.GET("/battery/averages", GetBatteryAverages)
	r.GET("/ecu/stream", StreamECUs)
	r.GET("/battery/stream", StreamBatteries)
	r.Run(":" + Port)
}

// GetAllECUs handles GET requests to retrieve all ECU records.
// It returns a JSON array of all ECU data.
func GetAllECUs(c *gin.Context) {
	var ecus []ECU
	DB.Find(&ecus)
	c.JSON(http.StatusOK, ecus)
}

// GetAllBatteries handles GET requests to retrieve all Battery records.
// It returns a JSON array of all Battery data.
func GetAllBatteries(c *gin.Context) {
	var batteries []Battery
	DB.Find(&batteries)
	c.JSON(http.StatusOK, batteries)
}

// GetECUAverages handles GET requests to retrieve average values for ECU data.
// It calculates and returns average values for motor RPM, speed, throttle, and brake pressure.
func GetECUAverages(c *gin.Context) {
	var result struct {
		AvgMotorRPM      float64 `json:"avg_motor_rpm"`
		AvgSpeed         float64 `json:"avg_speed"`
		AvgThrottle      float64 `json:"avg_throttle"`
		AvgBrakePressure float64 `json:"avg_brake_pressure"`
	}

	DB.Model(&ECU{}).Select("AVG(motor_rpm) as avg_motor_rpm, AVG(speed) as avg_speed, AVG(throttle) as avg_throttle, AVG(brake_pressure) as avg_brake_pressure").Scan(&result)

	c.JSON(http.StatusOK, result)
}

// GetBatteryAverages handles GET requests to retrieve average values for Battery data.
// It calculates and returns average values for charge level, cell temperatures, and cell voltages.
func GetBatteryAverages(c *gin.Context) {
	var result struct {
		AvgChargeLevel  float64 `json:"avg_charge_level"`
		AvgCellTemp1    float64 `json:"avg_cell_temp_1"`
		AvgCellTemp2    float64 `json:"avg_cell_temp_2"`
		AvgCellTemp3    float64 `json:"avg_cell_temp_3"`
		AvgCellTemp4    float64 `json:"avg_cell_temp_4"`
		AvgCellVoltage1 float64 `json:"avg_cell_voltage_1"`
		AvgCellVoltage2 float64 `json:"avg_cell_voltage_2"`
		AvgCellVoltage3 float64 `json:"avg_cell_voltage_3"`
		AvgCellVoltage4 float64 `json:"avg_cell_voltage_4"`
	}

	DB.Model(&Battery{}).Select("AVG(charge_level) as avg_charge_level, AVG(cell_temp1) as avg_cell_temp1, AVG(cell_temp2) as avg_cell_temp2, AVG(cell_temp3) as avg_cell_temp3, AVG(cell_temp4) as avg_cell_temp4, AVG(cell_voltage1) as avg_cell_voltage1, AVG(cell_voltage2) as avg_cell_voltage2, AVG(cell_voltage3) as avg_cell_voltage3, AVG(cell_voltage4) as avg_cell_voltage4").Scan(&result)

	c.JSON(http.StatusOK, result)
}

// upgrader is used to convert a regular HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StreamECUs handles WebSocket connections for streaming real-time ECU data.
// It registers a callback to send new ECU data to connected clients.
func StreamECUs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	RegisterECUCallback(func(ecu ECU) {
		conn.WriteJSON(ecu)
	})

	for {
		_, p, _ := conn.ReadMessage()
		log.Println(string(p))
	}
}

// StreamBatteries handles WebSocket connections for streaming real-time Battery data.
// It registers a callback to send new Battery data to connected clients.
func StreamBatteries(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	RegisterBatteryCallback(func(battery Battery) {
		conn.WriteJSON(battery)
	})

	for {
		_, p, _ := conn.ReadMessage()
		log.Println(string(p))
	}
}
