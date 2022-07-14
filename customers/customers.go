package customers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

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

func init() {
	LogFile, error := os.OpenFile("CustomerLog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if error != nil {
		panic("Failed to Create Log file for Customers API.")
	}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Llongfile|log.Ltime)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Llongfile|log.Ltime)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Llongfile|log.Ltime)
}

//go:generate mockgen -destination=mocks/customer.go --build_flags=--mod=mod . CustomerInterface
type CustomerInterface interface {
	// Interface for Managing Customer Model
	CreateCustomer(customerData map[string]interface{})
	UpdateCustomer(customerId string, UpdatedData ...map[string]interface{})
	DeleteCustomer(customerId string)
}

var customer models.Customer

func CreateCustomerRestController(RequestContext *gin.Context) {
	// Creates Customer

	NonSpecifiedFields := []string{}
	for Property, Value := range reflection.Items(models.Customer) {
		if Value := RequestContext.PostForm(Property); len(Value) == 0 {
			result := append(NonSpecifiedFields, Property)
			NonSpecifiedFields = result
		} else {
			continue
		}
	}

	if len(NonSpecifiedFields) != 0 {
		ErrorStatus := http.StatusBadRequest
		RequestContext.JSON(ErrorStatus,
			gin.H{"NonSpecifiedField": NonSpecifiedFields})
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

	NewCustomer := customer.CreateObject(newCustomerData, nil, []models.Product{})
	if NewCustomer == nil {
		RequestContext.JSON(http.StatusNotImplemented,
			gin.H{"error": "Failed to Create Customer."})
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
	// Updates Customer

	customerId := context.Query("customerId")
	MappedItems, error := reflection.Items(customer)

	for element, value := range MappedItems {
		if EmptyValue := context.PostForm(element); len(EmptyValue) == 0 {
			context.Request.PostForm.Del(element)
		} else {
			continue
		}
	}

	updatedCustomerData := struct{ Password string }{
		Password: context.PostForm("Password"),
	}
	updatedCustomer := customer.UpdateObject(customerId, updatedCustomerData)

	if updatedCustomer == false {
		context.JSON(
			http.StatusNotImplemented, gin.H{"error": error})
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
	deleted := customer.DeleteObject(customerId)

	if deleted != true {
		RequestContext.JSON(http.StatusNotImplemented, nil)
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

	context.JSON(http.StatusOK, gin.H{"customerProfile": string(
		jsonSerializedCustomer[:len(jsonSerializedCustomer)])})
}
