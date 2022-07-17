package payments

import (
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"context"
	"time"

	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/exceptions"
	paymentClients "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/clients"
	paymentGrpcControllers "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
	"github.com/mercari/go-circuitbreaker"
	curcuitbreaker "github.com/mercari/go-circuitbreaker"
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

	Validate() (PaymentIntentCredentialsInterface, []error)
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
	Validate() (PaymentSessionCredentialsInterface, []error)
	GetCredentials() (*PaymentSessionCredentialsInterface, []error)
}

//go:generate -destination=StoreService/mocks/payments.go --build_flags=--mod=mod . PaymentCheckoutStructInterface
type PaymentCheckoutStructRendererInterface interface {
	// Payment Checkout Interface, that is responsible for Processing Checkout Content, 
	GetImage(CheckoutContent map[string]string) ([]byte, error)
	GetJson(CheckoutContent map[string]string) ([]byte, error)
	GetXml(CheckoutContent map[string]string) ([]byte, error)
}

//go:generate -destination=StoreService/mocks/payments.go --build_flags=--mod=mod . PaymentRefundCredentialsInterface

type PaymentRefundCredentialsInterface interface {
	// Controller Interface, represents Payment Refund Model, Requires Following Params.
	// - Payment Id of ORM Model Object the Object That Was Created, during Successful Payment.
	Validate() (PaymentRefundCredentialsInterface, []error)
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

type PaymentCheckoutControllerInterface interface {
	// Payment Checkout Interface, that is responsible for obtaining Payment Checkout Info About the Products, that has been purchased...
	GetPaymentCheckout() (*PaymentCheckoutStructInterface, error)
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

		ProductId string 
		PurchaserId string 
		PurchaserUsername string 
		TotalPrice string 
		Currency string
	}
}

func NewPaymentSessionCredentials(Credentials struct{}) *PaymentSessionCredentials {
	return &PaymentSessionCredentials{Credentials: Credentials}
}

func (this *PaymentSessionCredentials) Validate() (PaymentSessionCredentials, []error)

func (this *PaymentSessionCredentials) GetCredentials() (*PaymentSessionCredentials, []error)

type PaymentRefundCredentials struct {

	Mutex       sync.RWMutex
	Credentials struct {
		PaymentId string 
	}
}

func NewPaymentRefundCredentials(Credentials struct{}) *PaymentRefundCredentials {
	return &PaymentRefundCredentials{Credentials: Credentials}
}

func (this *PaymentRefundCredentials) Validate() (PaymentRefundCredentials, []error)

func (this *PaymentRefundCredentials) GetCredentials() (*PaymentRefundCredentials, []error)










// Controller Implementations

type PaymentIntentController struct {
	// Represents Interface of the Payment Intent was implemented for communication between this app and `Payment Service`, to allow making Payment Intents...
	// Requires Attributes:
		// Client - grpc Client that represents Payment Grpc Service as described in proto file: https://github.com/LovePelemeni/Payment-Service/API/grpc/proto/payment.proto 
		// CircuitBreaker - For Preventing Failure Operations, while the `Payment Service` is not available...

	Client         paymentClients.PaymentIntentClientInterface
	CircuitBreaker curcuitbreaker.CircuitBreaker
}

func NewPaymentIntentController(Client *paymentClients.PaymentIntentClientInterface) *PaymentIntentController {
	return &PaymentIntentController{Client: *Client,
		CircuitBreaker: *curcuitbreaker.New(
			
			circuitbreaker.WithOpenTimeout(20),
			circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State){

				if newState == "OPEN" {ErrorLogger.Fatal("Payment Service Is Down, And Does Not Respond On Any Requests.")} 
				if newState == "CLOSED" {InfoLogger.Fatal("Payment Service Is Now Recorved. Time: " + time.Now().String())}
		}),
	)}
}


func (this *PaymentIntentController) CreatePaymentIntent(Credentials *PaymentIntentCredentials) (struct{PaymentIntentId string}, []error) {

	if !this.CircuitBreaker.Ready() {return struct{PaymentIntentId string}{}, []error{exceptions.ServiceUnavailable()}}

	PaymentGrpcClient, Error := this.Client.GetClient()
	paymentCredentials, ValidationErrors := Credentials.Validate()

	if errors.Is(Error, exceptions.ServiceUnavailable()){
	return struct{PaymentIntentId string}{}, []error{Error}}


	if len(ValidationErrors) != 0 {
		return struct{PaymentIntentId string}{}, ValidationErrors
	}

	var PaymentIntentResponse *paymentGrpcControllers.PaymentIntentResponse 
	RequestContext, CancelError := context.WithTimeout(context.Background(), time.Second*10) // Initializing Request Context..
	PaymentResponseIntentId, Error := this.CircuitBreaker.Do(
		RequestContext,
	   
		func() (interface{}, error) {

			// TODO: Sending Grpc Request to the grpc Endpoints, which is located in `Payment Service.`
			// For more info check out the `https://github.com/LovePelmeni/Payment-Service/README.md`

			PaymentResponse, Error := PaymentGrpcClient.CreatePaymentIntent( 
				RequestContext,
				&paymentGrpcControllers.PaymentIntentParams{
					ProductId: paymentCredentials.Credentials.ProductId,
					PurchaserId: paymentCredentials.Credentials.PurchaserId,
					Currency: paymentCredentials.Credentials.Currency,
					Price: paymentCredentials.TotalPrice,
				},
			)

			if Error != nil {
				this.CircuitBreaker.FailWithContext(RequestContext)
				InfoLogger.Println(
					"Failure Response from Payment Grpc Server..")
				return nil, Error 
					
			}else{
				PaymentIntentResponse = PaymentResponse
				this.CircuitBreaker.Done(RequestContext, nil)
				return nil, nil 
		}})

		if errors.Is(Error, exceptions.ServiceUnavailable()) {
			return struct{PaymentIntentId string}{}, []error{Error}
		}

		// Opening Curcuit Breaker In order to prevent any Potential Errors.
		return struct{PaymentIntentId string}{PaymentIntentId: PaymentIntentResponse.PaymentIntentId}, nil

	_ = PaymentResponseIntentId

	defer CancelError()
	return struct{PaymentIntentId string}{PaymentIntentId: PaymentIntentResponse.PaymentIntentId}, []error{Error}
}




type PaymentSessionController struct {

	// Interface Represents Entity of the Payment Session...
	// Requires Following Params... 

	// Client - Grpc Client, That Represents Communication Layer for making `Payment Sessions`,
	// between this application and `Payment Service.` for more info read: `https://github.com/LovePelmeni/Payment-Service/README.md`

	// Circuit Breaker - Circuit Breaker Object for handling Request State.
	Client *paymentClients.PaymentSessionClientInterface
	CircuitBreaker curcuitbreaker.CircuitBreaker
}

func NewPaymentSessionController(Client *paymentClients.PaymentSessionClientInterface) *PaymentSessionController {
	return &PaymentSessionController{Client: Client,
	 CircuitBreaker: *circuitbreaker.New(
		 circuitbreaker.WithOpenTimeout(20),
		 circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State){
			 if newState == "OPEN" {ErrorLogger.Println("Payment Service Is Not Available And Not Allowed to Start Any Payment Sessions.")}
			 if newState == "CLOSED" {ErrorLogger.Println("Payment Service is Available now. Time: " + time.Now().String())
		 }},
	)),
}

func (this *PaymentSessionController) CreatePaymentSession(Credentials *PaymentSessionCredentials) (struct{PaymentSessionId string}, []error) {

	if !this.CircuitBreaker.Ready() {return struct{PaymentSessionId string}{},
    []error{exceptions.ServiceUnavailable()}} // Checks for Circuit Breaker Status...

	grpcClient := this.Client.GetClient() // obtaining grpc Client for `PaymentSession` Service.
	credentials, Error := Credentials.Validate()
	if Error != nil {return struct{PaymentSessionId string}{},
    []error{exceptions.ValidationError()}}

	var PaymentSession *paymentGrpcControllers.PaymentSessionResponse
	_, Error := this.CircuitBreaker.Do(

		context.Background(),
		func() (interface{}, error){

			paymentSessionCredentials := paymentGrpcControllers.PaymentSessionParams{
				ProductId: Credentials.ProductId, 
				PurchaserId: Credentials.PurchaserId,
			}

			context, CancelError := context.WithTimeout(context.Background(), time.Second * 10)
			grpcResponse, Error := grpcClient.CreatePaymentSession(context, paymentSessionCredentials)

			defer CancelError()
			switch Error {
				case nil:
					this.CircuitBreaker.Done(context)
					return nil, nil

				default:
					this.CircuitBreaker.FailWithContext(context)
					return nil, exceptions.ServiceUnavailable()
			}
	})
	return struct{PaymentSessionId string}{
	PaymentSessionId: PaymentSession.PaymentSessionId}, []error{}
}




type PaymentRefundController struct {
	Client *paymentClients.PaymentRefundClientInterface
	CircuitBreaker curcuitbreaker.CircuitBreaker
}

func NewPaymentRefundController(Client *paymentClients.PaymentRefundClientInterface) *PaymentRefundController {
	return &PaymentRefundController{Client: Client}
}

func (this *PaymentRefundController) CreateRefundIntent(Credentials *PaymentRefundCredentials) (map[string]string, []error)
