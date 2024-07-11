package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database connection instance.
var DB *gorm.DB

// Database connection parameters
var (
	DatabaseHost     = ""
	DatabasePort     = ""
	DatabaseUser     = ""
	DatabasePassword = ""
	DatabaseName     = ""
)

// ConnectDB establishes a connection to the database using the provided credentials.
// It also performs auto-migration for the ECU and Battery models.
func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", DatabaseUser, DatabasePassword, DatabaseHost, DatabasePort, DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db
	log.Println("Connected to database")
	db.AutoMigrate(ECU{}, Battery{})
	log.Println("AutoMigration completed")
}

// CreateECU creates a new ECU record in the database.
// It also pushes the created ECU data to the registered callback functions.
// Returns the created ECU and any error encountered.
func CreateECU(ecu ECU) (ECU, error) {
	result := DB.Create(&ecu)
	if result.Error != nil {
		return ECU{}, result.Error
	}
	PushECU(ecu)
	return ecu, nil
}

// CreateBattery creates a new Battery record in the database.
// It also pushes the created Battery data to the registered callback functions.
// Returns the created Battery and any error encountered.
func CreateBattery(battery Battery) (Battery, error) {
	result := DB.Create(&battery)
	if result.Error != nil {
		return Battery{}, result.Error
	}
	PushBattery(battery)
	return battery, nil
}
