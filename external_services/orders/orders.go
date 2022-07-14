package orders

import (
	"log"
	"os"

	"errors"
	"sync"

	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/orders/firebase"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
)

var (
	DebugLogger   *log.Logger
	ErrorLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
)

// Firebase RealTime Database Credentials..

var (
	StorageBucketID  = os.Getenv("STORAGE_BUCKET_ID")
	ProjectID        = os.Getenv("PROJECT_ID")
	ServiceAccountID = os.Getenv("SERVICE_ACCOUNT_ID")
)

func InitializeLoggers() (bool, error) {

	LogFile, Exception := os.OpenFile("order.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if Exception != nil {
		return false, Exception
	}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Llongfile|log.Ltime)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Llongfile|log.Ltime)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Llongfile|log.Ltime)
	return true, nil
}

func init() {

	Initialized, Error := InitializeLoggers()
	if !Initialized && Error != nil {
		panic(Error)
	}

}

// Abstractions...

//go:generate mockgen generate -destination=Storeservice/mocks/orders.go . OrderCredentialsInterface
type OrderCredentialsInterface interface {

	// Interface Responsible for Initial Info About the Order..
	// Requires Following Structure To Be passed as a Order Credentials:

	Validate() (bool, error)
	GetCredentials() (OrderCredentialsInterface, error)
}

//go:generate -destination=StoreService/mocks/orders.go . OrderControllerInterface
type OrderControllerInterface interface {
	// Interface that represents the Main Controller Responsible For Handling Any Operations,
	// Related to the `Orders`
	CreateOrder(OrderCredentials OrderCredentialsInterface) (bool, error)
	CancelOrder(OrderId string) (bool, error)
}

// Implementations...

type OrderCredentials struct {
	mutex sync.RWMutex

	Credentials struct {
		mutex            sync.RWMutex
		OrderName        string
		OrderDescription string

		PurchasersInfo struct {
			Purchaser *models.Customer
		}

		ProductsInfo struct {
			Products         []*models.Product
			TotalPrice       string
			Currency         string
			ProductsQuantity string
		}
	}
}

func NewOrderCredentials(Credentials struct {
	Credentials struct {
		mutex            sync.RWMutex
		OrderName        string
		OrderDescription string

		PurchasersInfo struct {
			Purchaser *models.Customer
		}

		ProductsInfo struct {
			Products         []*models.Product
			TotalPrice       string
			Currency         string
			ProductsQuantity string
		}
	}
}) *OrderCredentials {
	return &OrderCredentials{Credentials: Credentials}
}

func (this *OrderCredentials) Validate() error {

	var ValidatedCustomersInfo interface{}
	var ValidatedProductsInfo interface{}

	_ = ValidatedCustomersInfo
	_ = ValidatedProductsInfo

	group := sync.WaitGroup{}

	// Validating Customer Info..

	go func(CustomerCredentials ...this.Credentials) {
		// If Data is valid it will put it into `ValidatedCustomersInfo`
		group.Add(1)

		group.Done()
	}(this.Credentials)

	// Validating Products Info..
	go func(ProductsCredentials ...this.Credentials) {

		// If Data is Valid it will put it into `ValidatedProductsInfo`
		group.Add(1)
		group.Done()
	}(this.Credentials)
	return nil
}

func (this *OrderCredentials) GetCredentials() (*OrderCredentials, error) {
	Error := this.Validate()
	if Error != nil {
		InfoLogger.Println(
			"Invalid Order Credentials has been passed.")
		return nil, errors.New("Invalid Credentials")
	}
	return this, nil
}

type OrderController struct {
	FirebaseManager *firebase.FirebaseDatabaseOrderManager // for managing orders...
}

func NewOrderController(FirebaseOrderManager *firebase.FirebaseDatabaseOrderManager) *OrderController {
	return &OrderController{FirebaseManager: FirebaseOrderManager}
}

func (this *OrderController) CreateOrder(OrderCredentials OrderCredentialsInterface) (bool, error) {

	OrderCredentials, Error := OrderCredentials.GetCredentials() // Validating Credentials First..
	if Error != nil {
		return false, Error
	}

	if Created, Error := this.FirebaseManager.CreateOrder(OrderCredentials); Created != true || Error != nil {
		ErrorLogger.Println("Failed to Create Order.. Reason: " + Error.Error())
		return false, Error
	}
	return true, nil
}

func (this *OrderController) CancelOrder(OrderId string) (bool, error) {

	if Deleted, Error := this.FirebaseManager.CancelOrder(OrderId); Deleted != true || Error != nil {
		ErrorLogger.Println("Failed to Cancel Order.. Reason: " + Error.Error())
		return false, Error
	} else {
		return true, nil
	}
}
