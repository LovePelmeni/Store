package models

import (
	"fmt"
	"log"
	"os"
	_ "strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
)

var customer Customer
var product Product
var cart Cart

var (
	POSTGRES_USER     = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DATABASE = os.Getenv("POSTGRES_DATABASE")
	POSTGRES_HOST     = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT     = os.Getenv("POSTGRES_PORT")
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

/// Models ...

func init() {
	LogFile, error := os.OpenFile("databaseLogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if error != nil {
		panic("Failed to Initialize Database Log File.")
	}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(LogFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type BaseModel interface {
	// ORM Model Interface with base Methods that every model need to have.
	Create(ObjectData map[string]interface{}) *BaseModel
	Update(UpdatedData map[string]interface{}) bool
	Delete(ObjectId string)
}

type Product struct {
	gorm.Model

	OwnerEmail         string  `gorm:"VARCHAR(100) NOT NULL"`
	ProductName        string  `gorm:"VARCHAR(100) NOT NULL"`
	ProductDescription string  `gorm:"VARCHAR(100) NOT NULL DEFAULT 'This Product Has No Description'"`
	ProductPrice       float64 `gorm:"NUMERIC(10, 5) NOT NULL"`
	Currency           string  `gorm:"VARCHAR(10) NOT NULL"`
}

// Create Controller...

func (this *Product) CreateObject() {
}

// Updating Controller...

func (this *Product) UpdateObject(ObjId string,
	UpdatedData struct {
		ProductName        string
		ProductDescription string
	}) bool {

	Updated := Database.Table("products").Where("id = ?", ObjId).Updates(UpdatedData)
	if Updated.Error != nil {
		Updated.Rollback()
		return false
	} else {
		Updated.Commit()
		return true
	}
}

// Deleting Controller...

func (this *Product) DeleteObject(ObjId string) bool {
	Deleted := Database.Table("products").Delete("id = ?", ObjId)
	if Deleted.Error != nil {
		Deleted.Rollback()
		return false
	} else {
		Deleted.Commit()
		return true
	}
}

type Customer struct {
	gorm.Model

	Username          string `gorm:"VARCHAR(100) NOT NULL UNIQUE";json:"Username"`
	Password          string `gorm:"VARCHAR(100) NOT NULL";json:"Password"`
	Email             string `gorm:"VARCHAR(100) NOT NULL UNIQUE";json:"Email"`
	ProductId         string
	PurchasedProducts []Product `gorm:"foreignKey:Product;references:ProductId;DEFAULT NULL;constraint:ON DELETE PROTECT;";
	 							json:"PurchasedProducts;constraint:OnDelete Protect;"`
	CreatedAt time.Time `gorm:"DATE DEFAULT CURRENT DATE";json:"CreatedAt"`
}

func (this *Customer) CreateObject(ObjectData struct {
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}, PurchasedProducts ...[]Product) *Customer {
	return &Customer{}
}

func (this *Customer) UpdateObject(ObjId string, UpdatedData struct{ Password string }) bool {
	Updated := Database.Table("customers").Updates(UpdatedData)
	if Updated.Error != nil {
		ErrorLogger.Println(
			"Failed To Update Customer.")
		return false
	} else {
		return true
	}
}

func (this *Customer) DeleteObject(ObjId string) bool {
	return true
}

type Cart struct {
	gorm.Model

	CustomerId string
	ProductId  string
	Owner      Customer  `gorm:"foreignKey:Customer;references:CustomerId;constraint:OnDelete Cascade;";json:"Owner"`
	Products   []Product `gorm:"foreignKey:Customer;references:ProductId;constraints:OnDelete PROTECT;";json:"Products"`
}

// Cart Create Controller ..

func (this *Cart) CreateObject(Customer *Customer, Products []Product) *Cart {
	// Creating Cart....
	newCart := Cart{Owner: *Customer, Products: Products}
	Saved := Database.Table("carts").Save(&newCart)
	if Saved.Error != nil {
		Saved.Rollback()
		return nil
	} else {
		Saved.Commit()
		return &newCart
	}
}

// Cart Update Controller...

func (this *Cart) UpdateObject(ObjId string, UpdatedData struct{ Products []Product }) bool {
	Updated := Database.Table("carts").Where("id = ?", ObjId).Updates(UpdatedData)
	if Updated.Error != nil {
		Updated.Rollback()
		return false
	} else {
		Updated.Commit()
		return true
	}
}

// Cart Delete Controller...

func (this *Cart) DeleteObject(ObjId string) bool {
	Deleted := Database.Table("carts").Where(
		"id = ?", ObjId).Delete("id = ?", ObjId)
	if Deleted.Error != nil {
		Deleted.Rollback()
		return false
	} else {
		Deleted.Commit()
		return true
	}
}
