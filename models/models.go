package models

import (
	"fmt"
	"log"
	"os"
	_ "strings"
	"time"

	"reflection"
	"regexp"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
)
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

var ValidationError error

type BaseValidator interface {
	// Base Validator Interface ...
	Validate() (string, ValidationError)
}

var (
	Validators = []BaseValidator{
		PriceValidator{},
		CurrencyValidator{},
	}
)

type CurrencyValidator struct {
	Value string
}

func (this CurrencyValidator) Validate() (string, ValidationError) {
	CurrencyPattern := ""
	if Matched, Error := regexp.MatchString(CurrencyPattern, this.Value); Matched == false {
		return "", Error
	} else {
		return this.Value, nil
	}
}

type PriceValidator struct {
	Value string
}

func (this PriceValidator) Validate() (string, ValidationError) {
	PricePattern := ""
	if Matched, Error := regexp.MatchString(PricePattern, this.Value); Matched == false {
		return "", Error
	} else {
		return this.Value, nil
	}
}

/// Models ...

var customer *Customer
var product *Product
var cart *Cart

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

func NewBaseModel(Model *gorm.DB) {
	return
}

type BaseModel interface {
	// ORM Model Interface with base Methods that every model need to have.
	ApplyRestrictedFields() bool
	GetRestrictedFields() []string
	Create(ObjectData map[string]interface{}) interface{}
	Update(UpdatedData map[string]interface{}) bool
	Delete(ObjectId string)
}

func init() {
	// Applying Tables Constraints...
	Models := []BaseModel{} // List of the Models...
	for _, model := range Models {
		if applied := model.ApplyRestrictedFields(); applied != true {
			ErrorLogger.Println("Failed to Apply Orm Table Restrict Dependencies.")
			panic(fmt.Sprintf("Orm Restriction Error, Model: %s", model))
		}
	}
	DebugLogger.Println("Constraints has been applied successfully.")
}

type Product struct {
	gorm.Model

	OwnerId            string
	Owner              *Customer `gorm:"foreignKey: Customer;references:OwnerId"`
	ProductName        string    `gorm:"VARCHAR(100) NOT NULL"`
	ProductDescription string    `gorm:"VARCHAR(100) NOT NULL DEFAULT 'This Product Has No Description'"`
	ProductPrice       string    `gorm:"NUMERIC(10, 5) NOT NULL"`
	Currency           string    `gorm:"VARCHAR(10) NOT NULL"`
}

// Create Controller...

func (this *Product) Create(

	ObjectData struct {
		ProductName        string
		ProductDescription string
		OwnerId            string
	},

	PriceCredentials struct {
		ProductPrice string
		Currency     string
	},

) *Product {

	Validators := []BaseValidator{
		PriceValidator{Value: PriceCredentials.ProductPrice},
		CurrencyValidator{Value: PriceCredentials.Currency},
	}

	// Validated Price Credentials...

	ValidatedPrice, PriceError := Validators[0].Validate()
	ValidatedCurrency, CurrencyError := Validators[1].Validate()

	if PriceError != nil || CurrencyError != nil {
		return nil
	}

	ValidatedPriceCredentials := map[string]string{
		"Price":    ValidatedPrice,
		"Currency": ValidatedCurrency,
	}

	// Validating Other String Product Params...

	for element, value := range reflection.Items(ObjectData) {
		if len(element) == 0 || element == nil {
			return nil
		}
	}

	// Creates New Object...

	newProduct := Product{
		ProductName:        ObjectData.ProductName,
		OwnerId:            ObjectData.OwnerId,
		ProductDescription: ObjectData.ProductDescription,
		ProductPrice:       ValidatedPriceCredentials["Price"],
		Currency:           ValidatedPriceCredentials["Currency"],
	}

	// Saving to the Database...

	Saved := Database.Table("products").Save(&newProduct)
	if Saved.Error != nil {
		ErrorLogger.Println(fmt.Sprintf(
			"Failed To Create Product. Reason: %s", Saved.Error))
		return nil
	}
	return &newProduct
}

// Updating Controller...

func (this *Product) Update(ObjId string,
	UpdatedData struct {
		ProductName        string
		ProductDescription string
	}) bool {

	for element, value := range reflection.Items(UpdatedData) {
		if len(value) == 0 || value == nil {
			reflection.Remove(UpdatedData, element)
		}
	}

	Updated := Database.Table("products").Where("id = ?", ObjId).Updates(UpdatedData)
	if Updated.Error != nil {
		return false
	} else {
		return true
	}
}

// Deleting Controller...

func (this *Product) Delete(ObjId string) bool {
	Deleted := Database.Table("products").Delete("id = ?", ObjId)
	if Deleted.Error != nil {
		return false
	} else {
		return true
	}
}

type Customer struct {
	gorm.Model

	Username          string `gorm:"VARCHAR(100) NOT NULL"`
	Password          string `gorm:"VARCHAR(100) NOT NULL"`
	Email             string `gorm:"VARCHAR(100) NOT NULL"`
	ProductId         string
	PurchasedProducts []Product `gorm:"foreignKey:Product;references:ProductId"`
	CreatedAt         time.Time `gorm:"DATE DEFAULT CURRENT DATE"`
}

func (this *Customer) Create(ObjectData map[string]interface{}) *Customer {

}

func (this *Customer) Update(ObjId string, UpdatedData map[string]interface{}) bool {

}

func (this *Customer) Delete(ObjId string) bool {

}

type Cart struct {
	gorm.Model

	CustomerId string
	ProductId  string
	Owner      Customer  `gorm:"foreignKey:Customer;references:CustomerId"`
	Products   []Product `gorm:"foreignKey:Customer;references:ProductId"`
}

// Cart Create Controller ..

func (this *Cart) Create(Customer *Customer, Products []Product) *Cart {
	// Creating Cart....
	newCart := Cart{Owner: *Customer, Products: Products}
	Saved := Database.Table("carts").Save(&newCart)
	if Saved.Error != nil {
		return nil
	} else {
		return &newCart
	}
}

// Cart Update Controller...

func (this *Cart) Update(ObjId string, UpdatedData struct{ Products []Product }) bool {
	Updated := Database.Table("carts").Where("id = ?", ObjId).Updates(UpdatedData)
	if Updated.Error != nil {
		return false
	} else {
		return true
	}
}

// Cart Delete Controller...

func (this *Cart) Delete(ObjId string) bool {
	Deleted := Database.Table("carts").Where(
		"id = ?", ObjId).Delete("id = ?", ObjId)
	if Deleted.Error != nil {
		return false
	} else {
		return true
	}
}

func CartOneOwnerConstraintTrigger() {
	// Adds trigger constraint that allows to have only one Owner Per Cart, In avoid of merging Orders.
	command := fmt.Sprintf(`CREATE FUNCTION public.check_one_cart_owner() 
	RETURNS TRIGGER 
	LANGUAGE 'plpgsql'
	AS $BODY$ 
	DECLARE updated integer;
	BEGIN UPDATE %s SET owner = owner + 1 WHERE cart.id = NEW.cart_id AND cart.owner < 1;
	GET DIAGNOSTICS addedOwners = ROW_COUNT 
	IF addedOwners = 0 THEN 
	RAISE EXCEPTION 'Cart can have only one owner.'
	END IF;
	RETURN NEW;
	END; 
	$BODY$;
	
	CREATE TRIGGER OwnerCartConstraintTrigger 
	BEFORE INSERT ON public.cart
	FOR EACH ROW EXECUTE PROCEDURE public.check_one_cart_owner();`, "cart")
	Database.Exec(command)
	DebugLogger.Println("Unique Constraint Has Been Integrated.")
}
