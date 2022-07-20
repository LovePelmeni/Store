package customers

import (
	"log"
	"os"

	"github.com/mercari/go-circuitbreaker"
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

var (
	grpcServerHost = os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort = os.Getenv("GRPC_SERVER_PORT")
)

func InitializeLoggers() (bool, error) {
	// Initializing Log File For the API Module.
	return true, nil
}

func init() {
	Initialized, Error := InitializeLoggers()
	if Error != nil || Initialized != true {
		panic(Error)
	}
}

// Abstractions...

type PaymentServiceCustomerCredentialsInterface interface {
	// Credentials that Is Necessary For Interacting with Customer Model and making CRUD Operations....
	Validate(Credentials map[string]string) (bool, error)
	GetCredentials(Credentials map[string]string) (PaymentServiceCustomerCredentialsInterface, error)
}

type PaymentServiceCustomerControllerInterface interface {
	// Interface that is responsible for editing Customer Model Instance located in `Payment Service`
	CreateRemoteCustomer(Credentials PaymentServiceCustomerCredentialsInterface) (bool, error)
	DeleteRemoteCustomer(CustomerId string) (bool, error)
}

// Implementation...

type PaymentServiceCustomerCredentials struct {
	Credentials struct{}
}

func NewPaymentServiceCustomerCredentials(Credentials struct{}) *PaymentServiceCustomerCredentials {
	return &PaymentServiceCustomerCredentials{Credentials: Credentials}
}

func (this *PaymentServiceCustomerCredentials) Validate(Credentials map[string]string) (bool, error)

func (this *PaymentServiceCustomerCredentials) GetCredentials(Credentials map[string]string) (PaymentServiceCustomerCredentialsInterface, error)

type PaymentServiceCustomerController struct {
	CircuitBreaker     circuitbreaker.CircuitBreaker
	GrpcCustomerClient interface{} // client for interacting with Payment Service via Grpc.
}

func NewPaymentServiceCustomerController() *PaymentServiceCustomerController {
	return &PaymentServiceCustomerController{CircuitBreaker: *circuitbreaker.New(
		circuitbreaker.WithOpenTimeout(20),
		circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
			if newState == "OPEN" {
				ErrorLogger.Println("Payment Service Customer Controller got failed. CircuitBreaker has blocked any requests.")
			}
			if newState == "CLOSED" {
				InfoLogger.Println("Payment Service Customer Controller recovered..")
			}
		}),
		circuitbreaker.WithHalfOpenMaxSuccesses(10),
	), GrpcCustomerClient: []interface{}{}}
}

func (this *PaymentServiceCustomerController) CreateRemoteCustomer(CustomerCredentials PaymentServiceCustomerCredentialsInterface) (bool, error)

func (this *PaymentServiceCustomerController) DeleteRemoteCustomer(CustomerId string, CustomerEmail string, CustomerUsername string) (bool, error)
