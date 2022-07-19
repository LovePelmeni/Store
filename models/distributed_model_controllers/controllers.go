package distributed_transaction_controllers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mercari/go-circuitbreaker"
	"google.golang.org/grpc"
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

type PaymentServiceCustomerCredentialsInterface interface {
	// Credentials that Is Necessary For Interacting with Customer Model and making CRUD Operations....
	Validate(Credentials map[string]string) PaymentServiceCustomerCredentialsInterface
	GetCredentials(Credentials map[string]string) (PaymentServiceCustomerCredentialsInterface, error)
}

type PaymentServiceProductCredentials interface {
	Validate()
	GetCredentials()
}

type PaymentServiceCustomerControllerInterface interface {
	// Interface that is responsible for editing Customer Model Instance located in `Payment Service`
	CreateRemoteCustomer() (bool, error)
	UpdateRemoteCustomer() (bool, error)
	DeleteRemoteCustomer() (bool, error)
}

type PaymentServiceProductControllerInterface interface {
	// Interface that is responsible for editing `Product` Model Instance located in `Payment` Service
	CreateRemoteProduct() (bool, error)
	UpdateRemoteProduct() (bool, error)
	DeleteRemoteProduct() (bool, error)
}

type PaymentServiceCustomerController struct {
	CircuitBreaker     circuitbreaker.CircuitBreaker
	GrpcCustomerClient *model_clients.GrpcCustomerClientInterface // client for interacting with Payment Service via Grpc.
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
	), GrpcCustomerClient: newGrpcCustomerClient()}
}

func (this *PaymentServiceCustomerController) CreateRemoteCustomer(CustomerCredentials *PaymentServiceCustomerCredentialsInterface) (bool, error)

type PaymentServiceProductController struct {
	CircuitBreaker     circuitbreaker.CircuitBreaker
	GrpcCustomerClient *model_clients.GrpcProductClientInterface
}

func NewPaymentServiceProductController(GrpcPaymentClient *model_clients.GrpcPaymentClientInterface) *PaymentServiceProductController {

	grpcServerConnection, ConnectionError := grpc.Dial(fmt.Sprintf("%s:%s", grpcServerHost, grpcServerPort))
	if ConnectionError != nil {
		ErrorLogger.Println("Failed to Connect To Grpc Payment Service Server.. Time: " + time.Now().String())
	}

	return &PaymentServiceProductController{

		CircuitBreaker: *circuitbreaker.New(
			circuitbreaker.WithOpenTimeout(20),
			circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
				if newState == "OPEN" {
					ErrorLogger.Println("Payment Service Product Controller got failed. CircuitBreaker has blocked any requests.")
				}
				if newState == "CLOSED" {
					InfoLogger.Println("Payment Service Product Controller recovered..")
				}
			}),
			circuitbreaker.WithHalfOpenMaxSuccesses(10),
		),
		GrpcProductClient: model_clients.NewPaymentServiceProductClient(
			grpcServerConnection),
	}
}
