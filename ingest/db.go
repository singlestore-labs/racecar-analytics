package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DatabaseHost = ""
var DatabasePort = ""
var DatabaseUser = ""
var DatabasePassword = ""
var DatabaseName = ""

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

func CreateECU(ecu ECU) (ECU, error) {
	result := DB.Create(&ecu)
	if result.Error != nil {
		return ECU{}, result.Error
	}
	return ecu, nil
}

func CreateBattery(battery Battery) (Battery, error) {
	result := DB.Create(&battery)
	if result.Error != nil {
		return Battery{}, result.Error
	}
	return battery, nil
}
