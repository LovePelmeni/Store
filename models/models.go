package models

import (
	"fmt"
	"log"
	"os"
	"regexp"
	_ "strings"

	"strconv"

	"time"

	"github.com/LovePelmeni/Store/exceptions"
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
	Database, error = gorm.Open(
		postgres.Open(
			fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
				POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DATABASE),
		),
		&gorm.Config{},
	)
)

/// Models ...

func InitializeLoggers() bool {
	LogFile, error := os.OpenFile("databaseLogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if error != nil {
		return false
	}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(LogFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	return true
}

func init() {
	Initialized := InitializeLoggers()
	if Initialized != true {
		panic("Failed to Initialize Loggers.")
	}
}

type ValidationError string

// Model Abstractions...

//go:generate -destination=StoreService/mocks/models.go --build_flags=--mod=mod . BaseModel
type BaseModel interface {
	// ORM Model Interface with base Methods that every model need to have.
	Create(ObjectData map[string]interface{}) *BaseModel
	Update(UpdatedData map[string]interface{}) bool
	Delete(ObjectId string)
}

//go:generate -destination=StoreService/mocks/models.go --build_flags=--mod=mod . BaesModelValidator
type BaseModelValidator interface {
	// Base Interface for the ORM Model, that allows to
	GetPatterns() map[string]string
	Validate(Credentials map[string]string) (map[string]string, []ValidationError)
}

// Because of microservice architecture, there is still some models in the bounded contexts,
// that requires to be initialized/created during the local transaction...
// That's why the interfaces down below helps to achieve this sort of functionality, to provide quick distributed service commication ...

type ProductModelValidator struct {
	Patterns map[string]string // Map key: Product Model Field Name, Value: Regex for validating this field.
}

func NewProductModelValidator() ProductModelValidator {
	Patterns := map[string]string{
		"OwnerEmail":         "^[a-z0-9](.?[a-z0-9]){5,}@g(oogle)?mail.com$", // default email regex.
		"ProductName":        "^[a-zA-z]{1,100}$",
		"ProductDescription": "^.*{1,300}$",
		"ProductPrice":       "^[0-9].[0-9]$",     // checks that the product price has appropriate format, like: 0.00
		"Currency":           "^[A-z]/[A-z]{10}$", // Checks that the letters is all upper's, max length 3 letters,
		// Example "USD", "EUR", "RUB" ...
	}
	return ProductModelValidator{Patterns: Patterns}
}

func (this ProductModelValidator) Validate(Data map[string]string) (map[string]string, []ValidationError) {

	ValidationErrors := []ValidationError{}
	for Property, Value := range Data {
		if Valid, Error := regexp.MatchString(this.Patterns[Property], Value); Valid == false || Error != nil {
			ValidationErrors =
				append(ValidationErrors, ValidationError(fmt.Sprintf(
					"Field `%s` is incorrect, Invalid Format.", Property)))
		}
	}
	if ValidationErrors != nil {
		return nil, ValidationErrors
	} else {
		return Data, nil
	}
}

func (this ProductModelValidator) GetPatterns() map[string]string {
	return this.Patterns
}

type Product struct {
	gorm.Model
	Id                 int     `gorm:"->;unique; BIGSERIAL NOT NULL PRIMARY KEY UNIQUE" json:"Id"`
	OwnerEmail         string  `gorm:"<-:create;unique; VARCHAR(100) NOT NULL" json:"OwnerEmail"`
	ProductName        string  `gorm:"<-;unique; VARCHAR(100) NOT NULL UNIQUE" json:"ProductName"`
	ProductDescription string  `gorm:"<-; VARCHAR(300) NOT NULL DEFAULT 'This Product Has No Description'" json:"ProductDescription"`
	ProductPrice       float64 `gorm:"<-; NUMERIC(10, 5) NOT NULL" json:"ProductPrice"`
	Currency           string  `gorm:"<-; VARCHAR(10) NOT NULL" json:"Currency"`
	CreatedAt          string  `gorm:"->;unique; DATE DEFAULT CURRENT DATE;" json:"CreatedAt"`
}

// Create Controller...

func (this *Product) CreateObject(

	ObjectData map[string]string,
	Validator BaseModelValidator,

) (bool, []ValidationError) {

	ValidatedData, Errors := Validator.Validate(ObjectData)
	if ValidatedData == nil || len(Errors) != 0 {
		return false, Errors
	}

	// Setting up Valid Values for the Object...
	this.ProductName = ValidatedData["ProductName"]
	this.OwnerEmail = ValidatedData["OwnerEmail"]
	this.ProductDescription = ValidatedData["ProductDescription"]
	this.Currency = ValidatedData["Currency"]
	this.ProductPrice, error = strconv.ParseFloat(ValidatedData["ProductPrice"], 5)

	LocalProductTransaction := Database.Table("products").Save(&this)
	if LocalProductTransaction.Error != nil {
		DebugLogger.Println("Failed to Save Product.")
		return false,
			[]ValidationError{ValidationError("Failed To Save Product.")}
	} else {
		return true, nil
	}
}

// Updating Controller...

func (this *Product) UpdateObject(ObjId string,
	UpdatedData struct {
		ProductName        string
		ProductDescription string
		ProductPrice       float64
	}, Validator BaseModelValidator) (bool, []ValidationError) {

	ValidatedData, Errors := Validator.Validate(map[string]string{
		"ProductName":        UpdatedData.ProductName,
		"ProductDescription": UpdatedData.ProductDescription})

	if len(Errors) != 0 || ValidatedData == nil {
		return false, Errors
	}

	Updated := Database.Table("products").Where("id = ?", ObjId).Updates(UpdatedData)
	if Updated.Error != nil {
		Updated.Rollback()
		return false, []ValidationError{ValidationError("Updated Failure")}
	} else {
		Updated.Commit()
		return true, nil
	}
}

// Deleting Controller...

func (this *Product) DeleteObject(ObjId string) (bool, []ValidationError) {
	Deleted := Database.Table("products").Delete("id = ?", ObjId)
	if Deleted.Error != nil {
		Deleted.Rollback()
		return false, []ValidationError{ValidationError("Deletion Failure.")}
	} else {
		Deleted.Commit()
		return true, nil
	}
}

// Customer ORM Model Validator...

type CustomerModelValidator struct {
	Patterns map[string]string
}

func NewCustomerModelValidator() CustomerModelValidator {
	Patterns := map[string]string{
		"Username": "^[a-zA-z]{1,100}$",                            // default string regex + max length
		"Password": "^[a-zA-z]{1,100}$",                            // default string regex + max length
		"Email":    "^[a-z0-9](.?[a-z0-9]){5,}@g(oogle)?mail.com$", // default email regex.
	}
	return CustomerModelValidator{Patterns: Patterns}
}

func (this CustomerModelValidator) Validate(ObjectData map[string]string) (map[string]string, []ValidationError) {
	ValidationErrors := []ValidationError{}
	for Property, Value := range ObjectData {
		if Valid, Error := regexp.MatchString(this.Patterns[Property], Value); Valid != true || Error != nil {
			ValidationErrors = append(ValidationErrors, ValidationError(
				fmt.Sprintf("Invalid Value for Field `%s`", Property)))
		}
	}
	if len(ValidationErrors) != 0 {
		return nil, ValidationErrors
	} else {
		return ObjectData, nil
	}
}

func (this CustomerModelValidator) GetPatterns() map[string]string {
	return this.Patterns
}

type Customer struct {
	gorm.Model
	Id                  int     `gorm:"->; uniqueIndex; BIGSERIAL NOT NULL UNIQUE PRIMARY KEY;"`
	Username            string  `gorm:"<-:create; unique; VARCHAR(100) NOT NULL;" json:"Username"`
	Password            string  `gorm:"<-; VARCHAR(100) NOT NULL;" json:"Password"`
	Email               string  `gorm:"<-:create; unique; VARCHAR(100) NOT NULL;" json:"Email"`
	PurchasedProductsId string  `gorm:"<-;primaryKey;DEFAULT NULL;" json:"PurchasedProductsId"`
	PurchasedProducts   Product `gorm:"-;foreignKey:Id;references:PurchasedProductsId;default:null;"`
	CreatedAt           string  `gorm:"<-:create; index; DATE DEFAULT CURRENT DATE;;" json:"CreatedAt"`
}

func (this *Customer) CreateObject(ObjectData struct {
	Username string
	Password string
	Email    string
}, Validator BaseModelValidator, PurchasedProducts ...[]Product) (*Customer, []ValidationError) {
	ValidatedData, Errors := Validator.Validate(map[string]string{
		"Username": ObjectData.Username,
		"Email":    ObjectData.Email,
		"Password": ObjectData.Password,
	})
	newCustomer := Customer{
		Username:            ValidatedData["Username"],
		Email:               ValidatedData["Email"],
		Password:            ValidatedData["Password"],
		CreatedAt:           time.Now().String(),
		PurchasedProductsId: "",
	}
	if Errors != nil {
		return nil, Errors
	}
	NewCustomerCreated := Database.Create(&newCustomer)
	if NewCustomerCreated.Error != nil {
		return nil, []ValidationError{ValidationError(
			exceptions.DatabaseOperationFailure().Error())}
	}
	return &newCustomer, nil
}

func (this *Customer) UpdateObject(

	ObjId string,
	UpdatedData struct{ Password string },
	Validator BaseModelValidator,
	// Client That is responsible for making remote transactions.
	// In this case, it is going to be used for Updated Remote Customer ORM Model Object from the `Payment Service.`

) (bool, []ValidationError) {

	ValidatedData, Errors := Validator.Validate(map[string]string{"Password": UpdatedData.Password})
	if Errors != nil {
		return false, Errors
	}
	Updated := Database.Table("customers").Updates(ValidatedData)
	if Updated.Error != nil {
		return false, []ValidationError{
			ValidationError(exceptions.DatabaseOperationFailure().Error())}
	}
	return true, nil
}

func (this *Customer) DeleteObject(ObjId string) (bool, []ValidationError) {
	Customer := Database.Table("customers").Where("id = ?", ObjId)
	if Customer.Error != nil {
		return false, []ValidationError{ValidationError(Customer.Error.Error())}
	}
	DeletionTransaction := Database.Table("customers").Delete(&Customer)

	if DeletionTransaction.Error != nil {
		return false, []ValidationError{ValidationError(
			DeletionTransaction.Error.Error())}
	} else {
		DeletionTransaction.Commit()
		return true, nil
	}
}

type CartModelValidator struct {
	// Validator Struct, that Represents Validator For `Cart` Model
	Patterns map[string]string
}

func NewCartModelValidator() *CartModelValidator {
	return &CartModelValidator{Patterns: nil}
}

func (this *CartModelValidator) Validate(ObjectData map[string]string) (map[string]string, []ValidationError) {

	ValidationErrors := []ValidationError{}
	for Property, Value := range ObjectData {

		if Matches, Error := regexp.MatchString(this.Patterns[Property], Value); Matches != true || Error != nil {
			ValidationErrors = append(
				ValidationErrors, ValidationError(fmt.Sprintf("Invalid Value for Field `%s`", Property)))
		}
	}
	if len(ValidationErrors) != 0 {
		return nil, ValidationErrors
	} else {
		return ObjectData, nil
	}
}

func (this *CartModelValidator) GetPatterns() map[string]string {
	return this.Patterns
}

type Cart struct {
	gorm.Model
	Id        int      `gorm:"->;BIGSERIAL NOT NULL PRIMARY KEY UNIQUE;"`
	OwnerId   string   `gorm:"<-:create;unique;primaryKey;" json:"OwnerId"`
	ProductId string   `gorm:"<-create;unique;primaryKey;" json:"ProductId"`
	Owner     Customer `gorm:"foreignKey:Id;references:OwnerId;constraint:OnDelete:Cascade;"`
	Products  Product  `gorm:"foreignKey:Id;references:ProductId;constraint:OnDelete:SET NULL;"`
	CreatedAt string   `gorm:"->;DATE DEFAULT CURRENT DATE;"`
}

// Cart Create Controller ..

func (this *Cart) CreateObject(Customer *Customer, Products []Product, Validator BaseModelValidator) (*Cart, []ValidationError) {
	// Creating Cart....
	newCart := Cart{Owner: *Customer, Products: Products[0]}
	Saved := Database.Table("carts").Save(&newCart)
	if Saved.Error != nil {
		Saved.Rollback()
		return nil, []ValidationError{ValidationError(Saved.Error.Error())}
	} else {
		Saved.Commit()
		return &newCart, nil
	}
}

// Cart Update Controller...

func (this *Cart) UpdateObject(ObjId string, UpdatedData struct{ Products []Product }, Validator BaseModelValidator) (bool, []ValidationError) {
	Updated := Database.Table("carts").Where("id = ?", ObjId).Updates(UpdatedData)
	if Updated.Error != nil {
		Updated.Rollback()
		return false, []ValidationError{ValidationError(Updated.Error.Error())}
	} else {
		Updated.Commit()
		return true, nil
	}
}

// Cart Delete Controller...

func (this *Cart) DeleteObject(ObjId string) (bool, []ValidationError) {
	Deleted := Database.Table("carts").Where(
		"id = ?", ObjId).Delete("id = ?", ObjId)
	if Deleted.Error != nil {
		Deleted.Rollback()
		return false, []ValidationError{ValidationError(Deleted.Error.Error())}
	} else {
		Deleted.Commit()
		return true, nil
	}
}
