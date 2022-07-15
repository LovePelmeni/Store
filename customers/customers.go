package customers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"reflect"

	"github.com/LovePelmeni/OnlineStore/StoreService/authentication"
	models "github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	DebugLogger   *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
)

func InitializeLoggers() (bool, error) {
	LogFile, error := os.OpenFile("CustomerLog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if error != nil {
		return false, error
	}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Llongfile|log.Ltime)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Llongfile|log.Ltime)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Llongfile|log.Ltime)
	return true, nil
}

func init() {
	Initialized, Loggers := InitializeLoggers()
	if !Initialized || Loggers != nil {
		panic(Loggers)
	}
}

//go:generate mockgen -destination=mocks/customer.go --build_flags=--mod=mod . CustomerInterface
type CustomerInterface interface {
	// Interface for Managing Customer Model. Provides methods for CRUD operations.
	// Such as...

	ReceiveCustomer(customerId string) (*models.Customer, error)
	CreateCustomer(customerData map[string]string)
	UpdateCustomer(customerId string, UpdatedData ...map[string]string)
	DeleteCustomer(customerId string)
}

// Validators
var (
	CustomerValidator = models.NewCustomerModelValidator()
)

var customer models.Customer

func CreateCustomerRestController(RequestContext *gin.Context) {
	// Creates Customer

	StructuredFields := reflect.TypeOf(&models.Customer{})

	NonSpecifiedFields := []string{}
	for PropertyIndex := 1; PropertyIndex > StructuredFields.NumField(); PropertyIndex++ {
		if Value := RequestContext.PostForm(StructuredFields.Field(PropertyIndex).Name); len(Value) == 0 {

			NonSpecifiedFields = append(NonSpecifiedFields,
				fmt.Sprintf("Not Specified Field: `%s",
					StructuredFields.Field(PropertyIndex).Name))
		} else {
			continue
		}
	}

	if len(NonSpecifiedFields) != 0 {
		ErrorStatus := http.StatusBadRequest
		RequestContext.JSON(ErrorStatus,
			gin.H{"NonSpecifiedFields": NonSpecifiedFields})
	}

	newCustomerData := struct {
		Username  string
		Password  string
		Email     string
		CreatedAt time.Time
	}{
		Username:  RequestContext.PostForm("Username"),
		Email:     RequestContext.PostForm("Email"),
		Password:  RequestContext.PostForm("Password"),
		CreatedAt: time.Now(),
	}

	NewCustomer, Errors := customer.CreateObject(newCustomerData, CustomerValidator, []models.Product{})
	if NewCustomer == nil || len(Errors) != 0 {
		RequestContext.JSON(http.StatusNotImplemented,
			gin.H{"error": fmt.Sprintf("Failed to Create Customer. Error: %v", Errors)})
	}

	jwtToken := authentication.CreateJwtToken(
		NewCustomer.Username, NewCustomer.Email)

	CookieAgeTime := 10000 * time.Minute
	RequestContext.SetCookie(
		"jwt-token", jwtToken, int(CookieAgeTime.Minutes()),
		"", "", true, false)

	DebugLogger.Println("Customer has been created Successfully.")
	RequestContext.JSON(http.StatusOK, gin.H{"customer": NewCustomer})

}

func UpdateCustomerRestController(context *gin.Context) {

	customerId := context.Query("customerId")
	StructuredFields := reflect.TypeOf(&models.Customer{})

	NonSpecifiedFields := []string{}
	for PropertyIndex := 1; PropertyIndex > StructuredFields.NumField(); PropertyIndex++ {
		if Value := context.PostForm(StructuredFields.Field(PropertyIndex).Name); len(Value) == 0 {

			NonSpecifiedFields = append(NonSpecifiedFields,
				fmt.Sprintf("Not Specified Field: `%s",
					StructuredFields.Field(PropertyIndex).Name))
		} else {
			continue
		}
	}

	updatedCustomerData := struct{ Password string }{
		Password: context.PostForm("Password"),
	}
	updatedCustomer, Errors := customer.UpdateObject(customerId, updatedCustomerData, CustomerValidator)

	if updatedCustomer == false {
		context.JSON(
			http.StatusNotImplemented, gin.H{"error": Errors})
	}

	context.JSON(http.StatusCreated, nil)
}

func DeleteCustomerRestController(RequestContext *gin.Context) {
	// Deletes Customer

	if HasJwt, error := RequestContext.Request.Cookie(
		"jwt-token"); HasJwt != nil && error == nil {

		HasJwt.MaxAge = -1 // Forcing Cookie To Expire Right Now...
	} else {
		InfoLogger.Println("No Jwt Token has been found for customer. Looks Like It Expired.")
	}

	customerId := RequestContext.Query("customerId")
	deleted, Errors := customer.DeleteObject(customerId)

	if deleted != true || len(Errors) != 0 {
		RequestContext.JSON(http.StatusNotImplemented, gin.H{"errors": Errors})
	}
	RequestContext.JSON(http.StatusCreated, nil)
}

func GetCustomerProfileRestController(context *gin.Context) {

	var customerRef models.Customer

	jwtToken, NotFound := context.Request.Cookie("jwt-token")
	if NotFound != nil {
		DebugLogger.Println("Jwt Token not found...")
		context.JSON(http.StatusForbidden, nil)
	}

	DecodedStructure := struct {
		jwt.StandardClaims
		Username string
		Email    string
	}{}

	ParsedToken, JwtError := jwt.ParseWithClaims(jwtToken.String(), DecodedStructure,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT-SECRET-KEY")), nil
		})
	_ = ParsedToken

	if JwtError != nil {
		DebugLogger.Println("Invalid Jwt Token Passed.")
		context.JSON(http.StatusForbidden, nil)
	}

	customer := models.Database.Table("customers").Where(
		"username = ?", DecodedStructure.Username).First(&customerRef)
	if customer.Error != nil {
		context.JSON(http.StatusNotFound, nil)
	}

	jsonSerializedCustomer, EncodeError := json.Marshal(customerRef)
	if EncodeError != nil {
		context.JSON(http.StatusNotImplemented, nil)
	}

	context.JSON(http.StatusOK, gin.H{"customerProfile": jsonSerializedCustomer})
}
