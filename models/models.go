package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	_ "strings"
	"time"

	"errors"
	"strconv"
	"sync"

	RemoteCustomerControllers "github.com/LovePelmeni/OnlineStore/StoreService/models/payment_service_customers"
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
	Validate(map[string]string) (map[string]string, []string)
}

// Because of microservice architecture, there is still some models in the bounded contexts,
// that requires to be initialized/created during the local transaction...
// That's why the interfaces down below helps to achieve this sort of functionality, to provide quick distributed service commication ...

type ProductModelValidator struct {
	Patterns map[string]string // Map key: Product Model Field Name, Value: Regex for validating this field.
}

func NewProductModelValidator() *ProductModelValidator {
	Patterns := map[string]string{
		"OwnerEmail":         "", // default email regex.
		"ProductName":        "",
		"ProductDescription": "",
		"ProductPrice":       "[0-9].[0-9]", // checks that the product price has appropriate format, like: 0.00
		"Currency":           "",            // Checks that the letters is all upper's, max length 3 letters,
		// Example "USD", "EUR", "RUB" ...
	}
	return &ProductModelValidator{Patterns: Patterns}
}

func (this *ProductModelValidator) Validate(Data map[string]string) (map[string]string, []string) {

	ValidationErrors := []string{}
	for Property, Value := range Data {
		if Valid, Error := regexp.MatchString(this.Patterns[Property], Value); Valid == false || Error != nil {
			ValidationErrors =
				append(ValidationErrors, fmt.Sprintf(
					"Field `%s` is incorrect, Invalid Format.", Property))
		}
	}
	if len(ValidationErrors) != 0 {
		return nil, ValidationErrors
	} else {
		return Data, []string{}
	}
}

func (this *ProductModelValidator) GetPatterns() map[string]string {
	return this.Patterns
}

type Product struct {
	gorm.Model
	Id                 int     `gorm:"BIGSERIAL NOT NULL PRIMARY KEY UNIQUE" json:"Id"`
	OwnerEmail         string  `gorm:"VARCHAR(100) NOT NULL UNIQUE" json:"OwnerEmail"`
	ProductName        string  `gorm:"VARCHAR(100) NOT NULL UNIQUE" json:"ProductName"`
	ProductDescription string  `gorm:"VARCHAR(100) NOT NULL DEFAULT 'This Product Has No Description'" json:"ProductDescription"`
	ProductPrice       float64 `gorm:"NUMERIC(10, 5) NOT NULL" json:"ProductPrice"`
	Currency           string  `gorm:"VARCHAR(10) NOT NULL" json:"Currency"`
}

// Create Controller...

func (this *Product) CreateObject(

	ObjectData map[string]string,
	Validator BaseModelValidator,

) (bool, []string) {

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
			[]string{"Failed To Save Product."}
	} else {
		return true, []string{}
	}
}

// Updating Controller...

func (this *Product) UpdateObject(ObjId string,
	UpdatedData struct {
		ProductName        string
		ProductDescription string
	}, Validator BaseModelValidator) (bool, []string) {

	ValidatedData, Errors := Validator.Validate(map[string]string{
		"ProductName":        UpdatedData.ProductName,
		"ProductDescription": UpdatedData.ProductDescription})

	if len(Errors) != 0 || ValidatedData == nil {
		return false, Errors
	}

	Updated := Database.Table("products").Where("id = ?", ObjId).Updates(UpdatedData)
	if Updated.Error != nil {
		Updated.Rollback()
		return false, []string{Updated.Error.Error()}
	} else {
		Updated.Commit()
		return true, []string{}
	}
}

// Deleting Controller...

func (this *Product) DeleteObject(ObjId string) (bool, []string) {
	Deleted := Database.Table("products").Delete("id = ?", ObjId)
	if Deleted.Error != nil {
		Deleted.Rollback()
		return false, []string{Deleted.Error.Error()}
	} else {
		Deleted.Commit()
		return true, []string{}
	}
}

// Customer ORM Model Validator...

type CustomerModelValidator struct {
	Patterns map[string]string
}

func NewCustomerModelValidator() *CustomerModelValidator {
	Patterns := map[string]string{
		"Username":  "", // default string regex + max length
		"Password":  "", // default string regex + max length
		"Email":     "", // default email regex.
		"ProductId": "", // valid integer
	}
	return &CustomerModelValidator{Patterns: Patterns}
}

func (this *CustomerModelValidator) Validate(ObjectData map[string]string) (map[string]string, []string) {
	return ObjectData, []string{}
}

func (this *CustomerModelValidator) GetPatterns() map[string]string {
	return this.Patterns
}

type Customer struct {
	gorm.Model
	Id                int    `gorm:"BIGSERIAL NOT NULL UNIQUE PRIMARY KEY"`
	Username          string `gorm:"VARCHAR(100) NOT NULL UNIQUE" json:"Username"`
	Password          string `gorm:"VARCHAR(100) NOT NULL" json:"Password"`
	Email             string `gorm:"VARCHAR(100) NOT NULL UNIQUE" json:"Email"`
	ProductId         string
	PurchasedProducts Product   `gorm:"foreignKey:ProductId;references:Id;default:null;" json:"PurchasedProducts"`
	CreatedAt         time.Time `gorm:"DATE DEFAULT CURRENT DATE" json:"CreatedAt"`
}

func (this *Customer) CreateObject(ObjectData struct {
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}, Validator BaseModelValidator, PurchasedProducts ...[]Product) (*Customer, []string) {
	return &Customer{}, []string{}
}

func (this *Customer) UpdateObject(

	ObjId string,
	UpdatedData struct{ Password string },
	Validator BaseModelValidator,
	CustomerGrpcClient RemoteCustomerControllers.PaymentServiceCustomerControllerInterface,
	// Client That is responsible for making remote transactions.
	// In this case, it is going to be used for Updated Remote Customer ORM Model Object from the `Payment Service.`

) (bool, []string) {

	ValidatedData, Errors := Validator.Validate(map[string]string{"Password": UpdatedData.Password})
	if Errors != nil {
		return false, Errors
	}
	group := sync.WaitGroup{}
	RequestContext, Cancel := context.WithCancel(context.Background())
	Updated := Database.Table("customers").Updates(ValidatedData)

	go func(RequestContext context.Context) {
		group.Add(1)

		_, Error := CustomerGrpcClient.CreateRemoteCustomer(
			RemoteCustomerControllers.NewPaymentServiceCustomerCredentials(struct{}{}))

		if Error != nil {
			Cancel()
		} else {
			RequestContext.Done()
		} // If Operation succeeded calling, Done() else Calling Cancel Operation.
		group.Done()

	}(RequestContext)

	group.Wait()

	select { // Processing Response from the remote Transaction...

	case <-RequestContext.Done():
		Updated.Rollback()
		ErrorLogger.Println("Failed to Process Remote Customer Transaction.")
		return false, []string{errors.New("Update Customer Failure.").Error()}

	case <-time.After(time.Duration(5)):
		Updated.Commit()
		InfoLogger.Println("Remote Customer has been created.")
		return true, nil
	}
}

func (this *Customer) DeleteObject(ObjId string) (bool, []string) {
	Customer := Database.Table("customers").Where("id = ?", ObjId)
	if Customer.Error != nil {
		return false, []string{Customer.Error.Error()}
	}
	DeletionTransaction := Database.Table("customers").Delete(&Customer)

	if DeletionTransaction.Error != nil {
		return false, []string{DeletionTransaction.Error.Error()}
	} else {
		DeletionTransaction.SavePoint("pre-deletion")
	}

	group := sync.WaitGroup{}
	RequestContext, CancelMethod := context.WithCancel(context.Background())

	go func(Context context.Context) {
		group.Add(1)

		client := RemoteCustomerControllers.NewPaymentServiceCustomerController()

		Response, Error := client.CreateRemoteCustomer(
			RemoteCustomerControllers.NewPaymentServiceCustomerCredentials(struct{}{}))

		if Response != true || Error != nil {
			CancelMethod()
		} else {
			Context.Done()
		}
		group.Done()
	}(RequestContext)

	group.Wait()
	select {
	case <-RequestContext.Done():
		DeletionTransaction.Rollback()
		return false, []string{errors.New("Failed to Create Customer.").Error()}
	case <-time.After(time.Duration(2 * time.Second)):
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

func (this *CartModelValidator) Validate(ObjectData map[string]string) (map[string]string, []string) {

	ValidationErrors := []string{}
	for Property, Value := range ObjectData {

		if Matches, Error := regexp.MatchString(this.Patterns[Property], Value); Matches != true || Error != nil {
			ValidationErrors = append(
				ValidationErrors, fmt.Sprintf("Invalid Value for Field `%s`", Property))
		}
	}
	if len(ValidationErrors) != 0 {
		return nil,
			ValidationErrors
	} else {
		return ObjectData, []string{}
	}
}

func (this *CartModelValidator) GetPatterns() map[string]string {
	return this.Patterns
}

type Cart struct {
	gorm.Model
	Id         int
	CustomerId string
	ProductId  string
	Owner      Customer `gorm:"foreignKey:CustomerId;references:Id;constraint:OnDelete:Cascade;" json:"Owner"`
	Products   Product  `gorm:"foreignKey:ProductId;references:Id;" json:"Products"`
}

// Cart Create Controller ..

func (this *Cart) CreateObject(Customer *Customer, Products []Product, Validator *BaseModelValidator) (*Cart, []string) {
	// Creating Cart....
	newCart := Cart{Owner: *Customer, Products: Products[0]}
	Saved := Database.Table("carts").Save(&newCart)
	if Saved.Error != nil {
		Saved.Rollback()
		return nil, []string{Saved.Error.Error()}
	} else {
		Saved.Commit()
		return &newCart, []string{}
	}
}

// Cart Update Controller...

func (this *Cart) UpdateObject(ObjId string, UpdatedData struct{ Products []Product }, Validator *BaseModelValidator) (bool, []string) {
	Updated := Database.Table("carts").Where("id = ?", ObjId).Updates(UpdatedData)
	if Updated.Error != nil {
		Updated.Rollback()
		return false, []string{Updated.Error.Error()}
	} else {
		Updated.Commit()
		return true, []string{}
	}
}

// Cart Delete Controller...

func (this *Cart) DeleteObject(ObjId string) (bool, []string) {
	Deleted := Database.Table("carts").Where(
		"id = ?", ObjId).Delete("id = ?", ObjId)
	if Deleted.Error != nil {
		Deleted.Rollback()
		return false, []string{Deleted.Error.Error()}
	} else {
		Deleted.Commit()
		return true, []string{}
	}
}
