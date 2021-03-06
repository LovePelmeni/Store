package orders

import (
	"log"
	"os"
	"strings"

	"errors"
	"sync"

	"fmt"
	"reflect"
	"regexp"

	"context"

	"github.com/LovePelmeni/Store/external_services/exceptions"
	"github.com/LovePelmeni/Store/external_services/firebase"
	"github.com/LovePelmeni/Store/models"
	"github.com/mercari/go-circuitbreaker"
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
	GetCredentials() (OrderCredentialsInterface, []error)
}

type OrderSessionInterface interface {
	// Interface that represents Order Session,
	GetSession() (map[string]string, error) // Returns Json Structure, Possible errors: NotFound, Timeout
	SetParams() (bool, error)               // Returns bool if params has been updated successfully.
}

//go:generate -destination=StoreService/mocks/orders.go . OrderControllerInterface
type OrderControllerInterface interface {
	// Interface that represents the Main Controller Responsible For Handling Any Operations,
	// Related to the `Orders`
	CreateOrder(OrderCredentials OrderCredentialsInterface) (bool, []error)
	CancelOrder(OrderId string) (bool, error)
}

type OrderCurcuitBreakerInterface interface {
	// Interface that implements curcuit breaker for the Order Interface gRPC Calls.
	Open() (bool, error)
	HalfOpen() (bool, error)
	Close() (bool, error)
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

func (this *OrderCredentials) Validate() (*OrderCredentials, []error) {

	var ValidationErrors struct {
		Mutex  sync.RWMutex
		Errors []error
	}

	// If goroutines Run Successfully, `Errors` field will be equals to emtpy list.

	var ValidatedCustomersInfo struct {
		mutex     sync.RWMutex
		Purchaser *models.Customer
	}

	var ValidatedProductsInfo struct {
		mutex            sync.RWMutex
		Products         []*models.Product
		TotalPrice       string
		Currency         string
		ProductsQuantity string
	}

	group := sync.WaitGroup{}

	// Validating Customer Info..

	go func(CustomerCredentials struct{ Purchaser *models.Customer }) {
		// If Data is valid it will put it into `ValidatedCustomersInfo`
		group.Add(1)

		customer := &models.Customer{}

		if Exists := models.Database.Table("customers").Where(
			"Username = ?", CustomerCredentials.Purchaser.Username).First(&customer); Exists.Error != nil {

			ValidationErrors.Mutex.Lock() // Writing Exception to the Errors List and Locking in order to avoid Race Condition.
			ValidationErrors.Errors = append(ValidationErrors.Errors,
				errors.New(fmt.Sprintf("User Specific in the Order"+
					"As Purchaser With Username: %s Not Found.", CustomerCredentials.Purchaser.Username)))
			ValidationErrors.Mutex.Unlock()

		} else {
			ValidatedCustomersInfo.mutex.Lock()
			ValidatedCustomersInfo.Purchaser = customer
		} // Updated the Customer Purchaser Field...

		// Matching Customer's ORM Model Object Properties, to that, that has been Passed in the Orders Form.

		StructuredCustomerCredentialsValue := reflect.ValueOf(CustomerCredentials).Elem() // Data, that has been Passed from the Order Form Credentials.
		StructuredCustomersProfileValue := reflect.ValueOf(customer).Elem()               // Data That has been retrieved from the ORM Model Object of the Customer.

		for PropertyIndex := 1; PropertyIndex < reflect.TypeOf(&CustomerCredentials).NumField(); PropertyIndex++ {

			CustomerFieldValue := StructuredCustomerCredentialsValue.Field(PropertyIndex).String() // Receiving the Customer Field

			if Equals := CustomerFieldValue ==
				StructuredCustomersProfileValue.Field(PropertyIndex).String(); Equals != true { // If not Matches to the Original Value, Writing Exception...

				ValidationErrors.Mutex.Lock()
				ValidationErrors.Errors = append(ValidationErrors.Errors,
					errors.New(fmt.Sprintf("Invalid Value for Field `%s` Passed to the Order,"+
						" Does not Match Purchaser Confidential Information.", reflect.TypeOf(
						StructuredCustomerCredentialsValue).Field(PropertyIndex).Name)))

				ValidationErrors.Mutex.Unlock()
			}
		}

		group.Done()
	}(this.Credentials.PurchasersInfo)

	// Validating Products Info..
	go func(ProductsCredentials struct {
		Products         []*models.Product
		TotalPrice       string
		Currency         string
		ProductsQuantity string
	}) {

		// If Data is Valid it will put it into `ValidatedProductsInfo`
		group.Add(1)

		// matching value patterns
		matchingPatterns := map[string]string{
			"Currency":         "regex-pattern-for-currency", // Regex for Currency, Explanation: ....
			"TotalPrice":       "regex-pattern-for-price",    // Regex pattern For Price, Explanation: ....
			"ProductsQuantity": "[0-100]",                    // regex pattern for Products Quantity, Explanation: Just checks if the passed string is a number from 0 to 100.
		}
		ProductProperties := reflect.ValueOf(ProductsCredentials).Elem()

		for PropertyValueIndex := 1; PropertyValueIndex < reflect.TypeOf(ProductsCredentials).NumField(); PropertyValueIndex++ {

			// Receiving Properties Info From the Structure....
			Field := reflect.TypeOf(ProductProperties).Field(PropertyValueIndex)
			FieldValue := ProductProperties.Field(PropertyValueIndex)

			// Matching The Regex
			if Matches, Error := regexp.MatchString(matchingPatterns[Field.Name], FieldValue.String()); Matches != true { // Checking For the Regex Value match to the Presented Value..

				_ = Error
				ValidationErrors.Mutex.Lock() // Locking Mutex in order to avoid Race Condition...
				ValidationErrors.Errors = append(ValidationErrors.Errors, errors.New(
					fmt.Sprintf("Invalid Value for field")))
				ValidationErrors.Mutex.Unlock()

			} else {
				_ = Error
				ValidatedProductsInfo.mutex.Lock()

				// Checking If it's possible to replace the value, and setting up the right one.
				if CanSet := ProductProperties.Field(PropertyValueIndex).CanSet(); CanSet == true {
					ProductProperties.Field(PropertyValueIndex).Set(FieldValue)
					DebugLogger.Println("Failed to Set Valid Property to Structure.")
					ValidatedProductsInfo.mutex.Unlock()
				}
			}
		}
		group.Done()
	}(this.Credentials.ProductsInfo)

	group.Wait()

	Response, Errors := func() (*OrderCredentials, []error) {
		if len(ValidationErrors.Errors) != 0 {
			return nil, ValidationErrors.Errors
		} else {

			return NewOrderCredentials(

				struct {
					mutex            sync.RWMutex
					OrderName        string
					OrderDescription string

					Credentials struct {
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
				}{
					OrderName:        this.Credentials.OrderName,
					OrderDescription: this.Credentials.OrderDescription,
					Credentials: {
						PurchasersInfo: ValidatedCustomersInfo,
						ProductsInfo:   ValidatedProductsInfo,
					},
				}), nil
		}
	}()
	return Response, Errors
}

func (this *OrderCredentials) GetCredentials() (*OrderCredentials, []error) {
	OrderCredentials, Errors := this.Validate()
	if len(Errors) != 0 {
		InfoLogger.Println(
			"Invalid Order Credentials has been passed.")
		return nil, Errors
	}
	return OrderCredentials, nil
}

type OrderController struct {
	FirebaseManager firebase.FirebaseDatabaseOrderManagerInterface // for managing orders...
	CircuitBreaker  *circuitbreaker.CircuitBreaker
}

func NewOrderController(FirebaseOrderManager *firebase.FirebaseDatabaseOrderManager) *OrderController {

	return &OrderController{
		FirebaseManager: FirebaseOrderManager,

		CircuitBreaker: circuitbreaker.New(

			circuitbreaker.WithOpenTimeout(20),
			circuitbreaker.WithHalfOpenMaxSuccesses(10),
			circuitbreaker.WithOnStateChangeHookFn(func(old_state, new_state circuitbreaker.State) {
				if hasPrefix := strings.HasPrefix(strings.ToLower(
					string(new_state)), "cl"); hasPrefix == true {
				} else { // check if state has been changed to close..

					ErrorLogger.Println(fmt.Sprintf(
						"CircuitBreaker Is Opened. New State: %s", new_state))
				}
			}),
		)}
}

func (this *OrderController) CreateOrder(OrderCredentials OrderCredentialsInterface) (bool, []error) {

	if this.CircuitBreaker.Ready() != true {
		return false, []error{errors.New("Service Is Unavailable")}
	}
	OrderCredentials, Errors := OrderCredentials.GetCredentials() // Validating Credentials First..
	if len(Errors) != 0 {
		return false, Errors
	}

	var Created bool

	// This is copy of the OrderCredentials Data.
	orderFirebaseDataStruct := struct {
		mutex            sync.RWMutex
		OrderName        string
		OrderDescription string

		Credentials struct {
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
	}{
		OrderName:        this.Credentials.OrderName,
		OrderDescription: this.Credentials.OrderDescription,
		Credentials: {
			PurchasersInfo: ValidatedCustomersInfo,
			ProductsInfo:   ValidatedProductsInfo,
		},
	}

	_, Error := this.CircuitBreaker.Do(context.Background(), func() (interface{}, error) {
		Response, Error := this.FirebaseManager.CreateOrder()

		switch Error {

		case nil:
			this.CircuitBreaker.FailWithContext(context.Background())
			return false, exceptions.ServiceUnavailable()

		default:
			Created = Response
			this.CircuitBreaker.Done(context.Background(), nil)
			return Created, Error
		}
	})
	return false, []error{Error}
}

func (this *OrderController) CancelOrder(OrderId string) (bool, error) {

	// Checking if circuit breaker works properly..
	if this.CircuitBreaker.Ready() != true {
		return false, errors.New("Service is Unavailable.")
	}
	var Canceled bool
	_, Error := this.CircuitBreaker.Do(

		context.Background(),
		func() (interface{}, error) {

			Deleted, Error := this.FirebaseManager.CancelOrder(OrderId)
			switch Error {

			case nil:
				Canceled = true
				this.CircuitBreaker.Done(context.Background(), nil)
				return true, nil

			default:
				Canceled = Deleted
				ErrorLogger.Println("Failed to Cancel Order.. Reason: " + Error.Error())
				return false, Error
			}
		})
	return Canceled, Error
}
