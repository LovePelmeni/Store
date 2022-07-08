package models 

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"os"
	"fmt"
	_ "strings"
	"log"
	"time"
)

var (
	DebugLogger *log.Logger 
	InfoLogger *log.Logger 
	ErrorLogger *log.Logger 
	WarnLogger *log.Logger 
)
var (
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DATABASE = os.Getenv("POSTGRES_DATABASE")
	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
)

var (
	Database, error = gorm.Open(postgres.New(
		postgres.Config{
			DSN: fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DATABASE, POSTGRES_USER, POSTGRES_PASSWORD),
			PreferSimpleProtocol: true,
		},
	))
)

func init() {
	LogFile, error := os.OpenFile("databaseLogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	if error != nil {panic("Failed to Initialize Database Log File.")}

	DebugLogger = log.Logger(LogFile, "DEBUG: ")
	InfoLogger = log.Logger(LogFile, "INFO: ")
	ErrorLogger = log.Logger(LogFile, "ERROR: ")
	WarnLogger = log.Logger(LogFile, "WARNING: ")
}

type Product struct {
	gorm.Model 

	ProductName string `gorm:"VARCHAR(100) NOT NULL"`
	ProductDescription string `gorm:"VARCHAR(100) NOT NULL DEFAULT 'This Product Has No Description'"`
	ProductPrice string `gorm:"NUMERIC(10, 5) NOT NULL"`
	Currency string `gorm:"VARCHAR(10) NOT NULL"`
}

type Customer struct {
	gorm.Model 

	Username string `gorm:"VARCHAR(100) NOT NULL"`
	Password string `gorm:"VARCHAR(100) NOT NULL"`
	Email string `gorm:"VARCHAR(100) NOT NULL"`
	PurchasedProducts Product  `gorm:"foreignKey:Product;references:ProductId"`
	CreatedAt time.Time `gorm:"DATE DEFAULT CURRENT DATE"`
}

type Card struct {
	gorm.Model 
	
	Owner Customer 
	Products Product 
}