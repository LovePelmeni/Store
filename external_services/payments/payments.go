package payments

import (
	"errors"
	"log"
	"os"
	"sync"

	paymentClients "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/clients"
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
)

func InitializeLoggers() (bool, error) {
	LogFile, error := os.OpenFile("payments.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if error != nil {
		return false, errors.New("Failed to Set Up Logs for Payments Package.")
	}
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
	WarningLogger = log.New(LogFile, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile)
	return true, nil
}
func init() {
	Initialized, Errors := InitializeLoggers()
	if Errors != nil || !Initialized {
		panic("Failed to Initialize Payment Logs.")
	}
}

// Abstractions...

//go:generate -destination=StoreService/mocks/payments.go --build_flags=--mod=mod . PaymentIntentCredentialsInterface

type PaymentIntentCredentialsInterface interface {
	// Payment Intent Credentials Interface, Describes the Payment Intent Document..
	// Key-Required Parameters:
	// PaymentSessionId string - `Identifier` of the Payment Session, that was returned from the PaymentSessionControllerInterface after Calling `CreatePaymentSession` Method.

	Validate() (*PaymentIntentCredentialsInterface, []error)
	GetCredentials() (*PaymentIntentCredentialsInterface, []error)
}

//go:generate -destination=StoreService/mocks/payments.go --build_flags=--mod=mod . PaymentSessionCredentialsInterface

type PaymentSessionCredentialsInterface interface {
	// Payment Session Credentials, Describe the Payment Session Document...
	// Key-Required Params:
	// -  PurchasersInfo
	//		- Purchaser *models.Customer
	// 		- PurchaserEmail string

	// - ProductsInfo
	//		- TotalPrice string
	// 		- Currency string
	// 		-
	Validate() (*PaymentSessionCredentialsInterface, []error)
	GetCredentials() (*PaymentSessionCredentialsInterface, []error)
}

//go:generate -destination=StoreService/mocks/payments.go --build_flags=--mod=mod . PaymentRefundCredentialsInterface

type PaymentRefundCredentialsInterface interface {
	// Controller Interface, represents Payment Refund Model, Requires Following Params.
	// - Payment Id of ORM Model Object the Object That Was Created, during Successful Payment.
	Validate() (*PaymentRefundCredentialsInterface, []error)
	GetCredentials() (*PaymentRefundCredentialsInterface, []error)
}

// Controllers Interfaces...

//go:generate -destination=StoreService/mocks/payments.go --build_flags=--mod=mod . PaymentIntentControllerInterface

type PaymentIntentControllerInterface interface {
	// Controller Interface Responsible for handling Payment Intents.
	// Requires Following Params to be Passed as a structure inside the implementation.
	// - grpc Payment Intent Client
	CreatePaymentIntent(PaymentIntentCredentials *PaymentIntentCredentialsInterface) (map[string]string, error)
}

//go:generate -destination=StoreService/mocks/payments.go --build_flags=--mod=mod . PaymentSessionControllerInterface

type PaymentSessionControllerInterface interface {
	// Controller Interface Responsible for handling Payment Sessions.
	// Requires Following Params to be Passed as a Structure Inside the Implementation.
	// - grpc Payment Session Client.
	CreatePaymentSession(PaymentSessionCredentials *PaymentSessionCredentialsInterface) (map[string]string, error)
}

//go:generate -destination=StoreService/mocks/payments.go --build_flags=--mod=mod . PaymentRefundControllerInterface

type PaymentRefundControllerInterface interface {
	// Controller Interface Responsible for handling Payment Refunds..
	// Requires Following Params to be Passed as a Structure Inside the Implementation
	// - grpc Payment Refund Client.
	CreatePaymentRefund(PaymentRefundCredentials *PaymentRefundCredentialsInterface) (map[string]string, error)
}

// Implementations..

// Credentials Implementations...

type PaymentIntentCredentials struct {
	Mutex       sync.RWMutex
	Credentials struct {
	}
}

func NewPaymentIntentCredentials(Credentials struct{}) *PaymentIntentCredentials {
	return &PaymentIntentCredentials{Credentials: Credentials}
}

func (this *PaymentIntentCredentials) Validate() (*PaymentIntentCredentials, []error)

func (this *PaymentIntentCredentials) GetCredentials()

type PaymentSessionCredentials struct {
	Mutex       sync.RWMutex
	Credentials struct {
	}
}

func NewPaymentSessionCredentials(Credentials struct{}) *PaymentSessionCredentials {
	return &PaymentSessionCredentials{Credentials: Credentials}
}

func (this *PaymentSessionCredentials) Validate() (*PaymentSessionCredentials, []error)

func (this *PaymentSessionCredentials) GetCredentials() (*PaymentSessionCredentials, []error)

type PaymentRefundCredentials struct {
	Mutex       sync.RWMutex
	Credentials struct {
	}
}

func NewPaymentRefundCredentials(Credentials struct{}) *PaymentRefundCredentials {
	return &PaymentRefundCredentials{Credentials: Credentials}
}

func (this *PaymentRefundCredentials) Validate() (*PaymentRefundCredentials, []error)

func (this *PaymentRefundCredentials) GetCredentials() (*PaymentRefundCredentials, []error)

// Controller Implementations

type PaymentIntentController struct {
	Client *paymentClients.PaymentIntentClientInterface
}

func NewPaymentIntentController(Client *paymentClients.PaymentIntentClientInterface) *PaymentIntentController {
	return &PaymentIntentController{Client: Client}
}

func (this *PaymentIntentController) CreatePaymentIntent(Credentials *PaymentIntentCredentials) (map[string]string, error)

type PaymentSessionController struct {
	Client *paymentClients.PaymentSessionClientInterface
}

func NewPaymentSessionController(Client *paymentClients.PaymentSessionClientInterface) *PaymentSessionController {
	return &PaymentSessionController{Client: Client}
}

func (this *PaymentSessionController) CreatePaymentSession(Credentials *PaymentSessionCredentials) (map[string]string, error)

type PaymentRefundController struct {
	Client *paymentClients.PaymentRefundClientInterface
}

func NewPaymentRefundController(Client *paymentClients.PaymentRefundClientInterface) *PaymentRefundController {
	return &PaymentRefundController{Client: Client}
}

func (this *PaymentRefundController) CreateRefundIntent(Credentials *PaymentRefundCredentials) (map[string]string, error)
