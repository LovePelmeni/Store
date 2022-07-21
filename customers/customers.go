package customers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"sync"

	"github.com/LovePelmeni/OnlineStore/StoreService/authentication"
	models "github.com/LovePelmeni/OnlineStore/StoreService/models"
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

// Validators
var (
	CustomerValidator = models.NewCustomerModelValidator()
)

var customer models.Customer

func CreateCustomerRestController(RequestContext *gin.Context) {
	// Creates Customer

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
	if NewCustomer == nil || Errors != nil {
		serializedErrors, _ := json.Marshal(Errors)

		RequestContext.JSON(http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf(
				"Failed to Create Customer. Error: %v", serializedErrors),
			})
	}

	jwtToken := authentication.CreateJwtToken(
		NewCustomer.Username, NewCustomer.Email)

	CookieAgeTime := 10000 * time.Minute
	RequestContext.SetCookie(
		"jwt-token", jwtToken, int(CookieAgeTime.Minutes()),
		"/", "", true, true)

	DebugLogger.Println("Customer has been created Successfully.")
	RequestContext.JSON(http.StatusOK, gin.H{"customer": NewCustomer})
}

func UpdateCustomerRestController(context *gin.Context) {

	customerId := context.Query("customerId")
	updatedCustomerData := struct{ Password string }{
		Password: context.PostForm("Password"),
	}

	updatedCustomer, Errors := customer.UpdateObject(
		customerId, updatedCustomerData, CustomerValidator)

	if updatedCustomer == false || Errors != nil {
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

	if deleted != true || Errors != nil {
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

	jwtCredentials, JwtError := authentication.GetCustomerJwtCredentials(jwtToken.String())

	if JwtError != nil {
		DebugLogger.Println("Invalid Jwt Token Passed.")
		context.JSON(http.StatusForbidden, nil)
	}

	customer := models.Database.Table("customers").Where(
		"username = ? AND email = ?",
		jwtCredentials["Username"], jwtCredentials["Email"]).First(&customerRef)

	if customer.Error != nil {
		context.JSON(http.StatusForbidden, nil) // If jwt Is Not Valid. Is should return 403..
	}

	// Making Annotations of Purchased Products By User....

	var PurchasedProductsCount int64
	group := sync.WaitGroup{}

	go func(customer *models.Customer) {
		group.Add(1)
		Error := models.Database.Table("customers").Where(
			"id = ?", customer.Id).Select("PurchasedProducts").Count(&PurchasedProductsCount)
		if Error.Error != nil {
			PurchasedProductsCount = 0
		}
		group.Done()
	}(&customerRef)

	// Serializing Customer Profile Data... into JSON...
	jsonSerializedCustomer, EncodeError := json.Marshal(
		struct {
			Username          string
			Email             string
			CreatedAt         string
			PurchasedProducts int64
		}{
			Username: customerRef.Username, Email: customerRef.Email,
			CreatedAt: customerRef.CreatedAt.String(), PurchasedProducts: PurchasedProductsCount})

	if EncodeError != nil {
		context.JSON(http.StatusNotImplemented, nil)
	}
	context.JSON(http.StatusOK, gin.H{"customerProfile": jsonSerializedCustomer})
}
