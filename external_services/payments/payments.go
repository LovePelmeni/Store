package payments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/LovePelmeni/OnlineStore/StoreService/customers/exceptions"
	paymentControllers "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
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
		return false, errors.New("")
	}
	DebugLogger = log.New(LogFile, "DEBUG: ")
	InfoLogger = log.New(LogFile, "INFO: ")
	ErrorLogger = log.New(LogFile, "ERROR: ")
	WarningLogger = log.New(LogFile, "WARNING: ")
	return true, nil 
}
func init() {
	Initialized, Errors := InitializeLoggers()
	if Errors.Error != nil || !Initialized {panic("Failed to Initialize Payment Logs.")}
}


// Abstractions...


// Credentials Interfaces...

type PaymentIntentCredentialsInterface interface {
	// Payment Intent Credentials Interface, Describes the Payment Intent Document..
	// Key-Required Parameters:	
		// PaymentSessionId string - `Identifier` of the Payment Session, that was returned from the PaymentSessionControllerInterface after Calling `CreatePaymentSession` Method.
		
	Validate() (*PaymentIntentCredentialsInterface, []error)
	GetCredentials() (*PaymentIntentCredentialsInterface, []error)
}

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

type PaymentRefundCredentialsInterface interface {
	// Controller Interface, represents Payment Refund Model, Requires Following Params.
		// - Payment Id of ORM Model Object the Object That Was Created, during Successful Payment.
	Validate() (*PaymentRefundCredentialsInterface, []error)
	GetCredentials() (*PaymentRefundCredentialsInterface, []error)
}




// Controllers Interfaces...

type PaymentIntentControllerInterface interface {
	// Controller Interface Responsible for handling Payment Intents.
	// Requires Following Params to be Passed as a structure inside the implementation.
		// - grpc Payment Intent Client 
	CreatePaymentIntent(PaymentIntentCredentials *PaymentIntentCredentialsInterface) (map[string]string, error)
}

type PaymentSessionControllerInterface interface {
	// Controller Interface Responsible for handling Payment Sessions.
	// Requires Following Params to be Passed as a Structure Inside the Implementation.
		// - grpc Payment Session Client.
	CreatePaymentSession(PaymentSessionCredentials *PaymentSessionCredentialsInterface) (map[string]string, error)
}

type PaymentRefundControllerInterface interface {
	// Controller Interface Responsible for handling Payment Refunds..
	// Requires Following Params to be Passed as a Structure Inside the Implementation 
		// - grpc Payment Refund Client.
	CreatePaymentRefund(PaymentRefundCredentials *PaymentRefundCredentialsInterface) (map[string]string, error)
}



// Implementations..





// Credentials Implementations...


type PaymentIntentCredentials struct {
	Mutex sync.RWMutex 
	Credentials struct{

	}
}

func NewPaymentIntentCredentials(Credentials struct{}) (*PaymentIntentCredentials) {
	return &PaymentIntentCredentials{Credentials: Credentials}
}

func (this *PaymentIntentCredentials) Validate() (*PaymentIntentCredentials, []error)

func (this *PaymentIntentCredentials) GetCredentials()



type PaymentSessionCredentials struct {
	Mutex sync.RWMutex
	Credentials struct{

	}
}

func NewPaymentSessionCredentials(Credentials struct{}) (*PaymentSessionCredentials) {
	return &PaymentSessionCredentials{Credentials: Credentials}
}

func (this *PaymentSessionCredentials) Validate() (*PaymentSessionCredentials, []error)

func (this *PaymentSessionCredentials) GetCredentials() (*PaymentSessionCredentials, []error)



type PaymentRefundCredentials struct {
	Mutex sync.RWMutex
	Credentials struct{

	}
}

func NewPaymentRefundCredentials(Credentials struct{}) (*PaymentRefundCredentials) {
	return &PaymentRefundCredentials{Credentials: Credentials}
}

func (this *PaymentRefundCredentials) Validate() (*PaymentRefundCredentials, []error)  

func (this *PaymentRefundCredentials) GetCredentials() (*PaymentRefundCredentials, []error)



// Controller Implementations 




type PaymentIntentController struct {
	Client *paymentClients.PaymentIntentClientInterface 
}

func NewPaymentIntentController(Client *paymentClients.PaymentIntentClientInterface) (*PaymentIntentController) {
	return &PaymentIntentController{Client: Client}
}

func (this *PaymentIntentController) CreatePaymentIntent(Credentials *PaymentIntentCredentials) (map[string]string, error)  




type PaymentSessionController struct {
	Client *paymentClients.PaymentIntentClientInterface 
}

func NewPaymentSessionController(Client *paymentClients.PaymentIntentClientInterface) (*PaymentIntentController) {
	return &PaymentIntentController{Client: Client}
}

func (this *PaymentSessionController) CreatePaymentSession(Credentials *PaymentSessionCredentials) (map[string]string, error)  




type PaymentRefundController struct {
	Client *paymentClients.PaymentIntentClientInterface 
}

func NewPaymentRefundController(Client *paymentClients.PaymentIntentClientInterface) (*PaymentIntentController) {
	return &PaymentIntentController{Client: Client}
}

func (this *PaymentRefundController) CreateRefundIntent(Credentials *PaymentRefundCredentials) (map[string]string, error)  



