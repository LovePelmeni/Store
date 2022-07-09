package customers

import (
	"log"
	"os"
	"reflection"
	"time"

	"github.com/LovePelmeni/OnlineStore/StoreService/customers/exceptions"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
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
	CreateCustomer(customerData map[string]interface{}) (bool, error)
	UpdateCustomer(customerId string, UpdatedData ...map[string]interface{}) (bool, error)
	DeleteCustomer(customerId string) (bool, error)
}

type CustomerStruct struct{}

func generateJwt(
	customerUsername string, customerEmail string) (string, error) {
	return "", nil
}

func (this *CustomerStruct) createCustomer(RequestContext *gin.Context, customerData struct {
	Username string
	Email    string
	Password string
},
) (*models.Customer, error) {
	// Creates Customer
	newCustomer := models.Customer{
		Username: customerData.Username,
		Password: customerData.Password,
		Email:    customerData.Email,
	}

	SavedModel := models.Database.Save(&newCustomer)
	if SavedModel.Error != nil {
		ErrorLogger.Println("Failed to Create Customer")
	}

	jwtToken, JwtError := generateJwt(customerData.Username, customerData.Email)
	if JwtError != nil {
		ErrorLogger.Println("Failed to generate JWT Token.")
	}

	CookieAgeTime := 10000 * time.Minute
	RequestContext.SetCookie(
		"jwt-token", jwtToken, int(CookieAgeTime.Minutes()),
		"", "", true, false)

	DebugLogger.Println("Customer has been created Successfully.")
	return &newCustomer, nil
}

func (this *CustomerStruct) updateCustomer(customerId string,
	updatedData struct {
		ValidUsername string
		ValidPassword string
	}) (bool, error) {
	// Updates Customer

	MappedItems, error := reflection.Items(updatedData)
	ValidatedData := map[string]string{}

	for element, value := range MappedItems {
		if valid := value != nil && len(value) != 0; valid != false {
			ValidatedData[element] = value
		} else {
			return false, exceptions.ValidationError()
		}
	}

	models.Database.Clauses(clause.Locking{
		Strength: "ON UPDATE",
		Table:    clause.Table{Name: "customers"}})

	customer := models.Database.Table("customers").Where(
		"id = ?", customerId).Updates(ValidatedData)
	if customer.Error != nil {
		ErrorLogger.Println("Failed to Update Customer Profile.")
		return false, exceptions.UpdateFailure()
	}
	return true, nil
}

func (this *CustomerStruct) deleteCustomer(
	RequestContext *gin.Context, CustomerId string) (bool, error) {
	// Deletes Customer

	models.Database.Clauses(clause.Locking{ // Locking Table In Order to prevent some operation while Profile is Deleting...
		Strength: "ON UPDATE",
		Table:    clause.Table{Name: "customers"}})

	deletedCustomer := models.Database.Table(
		"customers").Delete("id", CustomerId)

	if deletedCustomer.Error != nil {
		ErrorLogger.Println("Failed To Delete Customer")
		return false, exceptions.DeleteFailure(
			deletedCustomer.Error.Error())
	}

	if HasJwt, error := RequestContext.Request.Cookie("jwt-token"); HasJwt != nil && error == nil {
		HasJwt.MaxAge = -1 // Forcing Cookie To Expire Right Now...
	} else {
		InfoLogger.Println("No Jwt Token has been found for customer. Looks Like It Expired.")
	}
	return true, nil
}
