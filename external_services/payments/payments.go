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
	"google.golang.org/grpc"
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
)

func init() {
	LogFile, error := os.OpenFile("payments.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if error != nil {
		return errors.New("")
	}
	DebugLogger = log.New(LogFile, "DEBUG: ")
	InfoLogger = log.New(LogFile, "INFO: ")
	ErrorLogger = log.New(LogFile, "ERROR: ")
	WarningLogger = log.New(LogFile, "WARNING: ")
}

type PaymentInterface interface {
	Pay(PaymentData map[string]string) (bool, error)
	Refund(PaymentId string) (*Refund, error)
}

type gRPCPaymentInterfaceManager interface {
	// Interface for handling gRPC payment requests...

	getGRPCDial() grpc.Dial // returns grpc Server Connection..

	getPaymentIntentClient() *paymentControllers.PaymentIntentClient
	getPaymentSessionClient() *paymentControllers.PaymentSessionClient
	getPaymentRefundClient() *paymentControllers.RefundClient
	getPaymentCheckoutClient() *paymentControllers.PaymentCheckoutClient

	SendPaymentIntentRequest(context *context.Context, params *paymentControllers.PaymentIntentParams) (*paymentControllers.PaymentIntentResponse, error)
	SendPaymentSessionRequest(context *context.Context, params *paymentControllers.PaymentSessionParams) (*paymentControllers.PaymentSessionResponse, error)
	SendPaymentRefundRequest(context *context.Context, params *paymentControllers.RefundParams) (*paymentControllers.PaymentIntentResponse, error)
	SendPaymentCheckoutRequest(context *context.Context, params *paymentControllers.PaymentCheckoutParams) (*paymentControllers.PaymentIntentResponse, error)
}

type PaymentInfoCredentialsInterface interface {
	GetCredentials() map[string]string // returns credentials of the payment.
	Validate() (map[string]string, error)
}

type PaymentSessionInterface interface {
}

type PaymentIntentInterface interface {
}

type PaymentCheckoutInterface interface {
	FormCheckout() //
}

type PaymentInfoCredentials struct {
	mutex sync.RWMutex

	PaymentSessionId string `json:"PaymentSessionId"`
	ProductId        string `json:"ProductId"`
	CustomerId       string `json:"CustomerId"`
	Currency         string `json:"Currency"`
	Price            string `json:"Price"`
}

func (this *PaymentInfoCredentials) GetCredentials() map[string]string {}

func (this *PaymentInfoCredentials) Validate() (bool, error) {
	return true, nil
}

type PaymentCheckoutInfo struct {
	sync.RWMutex

	PaymentId   string            `json:"PaymentId"`
	PurchaserId string            `json:"PurchaserId"`
	Products    []models.Products `json:"Products"`
	TotalPrice  string            `json:"TotalPrice"`
	Currency    string            `json:"Currency"`
	CreatedAt   time.Time         `json:"CreatedAt"`
}

func (this *Payment) Pay(

	Client *grpc.Client,
	PaymentCredentials PaymentInfoCredentialsInterface, // payments credentials interface
	grpcManager gRPCPaymentInterfaceManager, // handler for grpc requests.
) error {

	if valid, error := PaymentCredentials.Validate(); valid != nil || error != nil {
		DebugLogger.Println(
			fmt.Sprintf("Invalid Payment Credentials, Exception: %s"))
	}
	grpcClient := grpcManager.getPaymentIntentClient()
	channel := make(chan struct{}, 10000)

	group := sync.WaitGroup{}

	// Payment Intent Obtaining...

	go func(group *sync.WaitGroup, client paymentControllers.PaymentIntentClient, channel chan struct{}, PaymentValidatedData PaymentInfoCredentials) {

		group.Add(1)

		PaymentRequestParams := paymentClient.PaymentIntentParams{
			ProductId:   PaymentValidatedData.ProductId,
			PurchaserId: PaymentValidatedData.PurchaserId,
			Currency:    PaymentValidatedData.Currency,
			Price:       PaymentValidatedData.Price,
		}

		RequestContext, CancelError := context.WithTimeout(context.Background(), 10*time.Second)
		response, ResponseError := client.CreatePaymentIntent(RequestContext, PaymentRequestParams)

		if NotNone := errors.Is(ResponseError, nil); NotNone != false {
			ErrorLogger.Println(fmt.Sprintf(
				"Payment Failed. Reason: %s", ResponseError))

		} else {
			// Sending Channel
			channel <- struct{ PaymentIntentId }{
				PaymentIntentId: response.PaymentIntentId}
		}

		defer CancelError()
		group.Done()

		// Processing Payment based on Received Payment Intent
		go func(group *sync.WaitGroup, client *grpc.Client) {

			group.Add(1)
			select {

			case data := <-channel:
				// Receives Payment Intent and Processing Payment.
				data := <-channel
				var decodedData struct{ PaymentId string }
				decodedDataError := json.Unmarshal(data, &decodedData)

				if decodedDataError != nil {
					ErrorLogger.Println("Invalid Data Format.")
				}

				checkoutParams := CheckoutRequestParams{PaymentId: decodedData.PaymentId}
				RequestContext, CancelError := context.WithTimeout(context.Background(), 10*time.Second)
				CheckoutResponseconst, error := client.GetPaymentCheckout(checkoutParams)

				if error != nil {
					ErrorLogger.Println(
						fmt.Sprintf("Failed to Receive Checkout Info about the Payment with ID: %s",
							&checkoutParams.PaymentId))
				} else {

					// Processing Sending Email... About the Purchased Order.
					go func(group *sync.WaitGroup, client *grpc.Client, channel chan struct{}, PaymentCredentials *Payment) {
						group.Add(1)
						group.Done()
					}()
				}

				group.Done()
				defer CancelError()

			default:
				InfoLogger.Println("Operation Failed, Aborting Silently...")
				group.Abort()
			}
		}()

		group.Wait()
		return true, nil

	}()
}

func (this *Payment) Refund(CustomerId string, PaymentId string, grpcManager *gRPCPaymentInterfaceManager) {

	customer := models.Database.Table( // Receiving Purchaser Object Info.
		"customers").Where("id = ?", CustomerId).First(&models.Customer)

	if customer.Error != nil { // Check if it's exists...
		InfoLogger.Println("Customer Does Not Exist.")
		return false, exceptions.CustomerDoesNotExist()
	}

	// Initialize Wait Group Add Communication Channel For Goroutines...
	group := sync.WaitGroup{}
	channel := make(chan struct{})

	// Making Refund...
	go func(group *sync.WaitGroup, channel chan bool, refundParams string) {

		group.Add(1)

		// Making Refund Request...
		client := grpcManager.getPaymentRefundClient()
		response, error := client.CreateRefund(RefundParams)
		group.Done()
	}()
}

type RefundInteface interface {
	GetRefundInfo() (map[string]interface{}, error)
}

type CheckoutInterface interface {
	GetCheckoutInfo() (map[string]interface{}, error)
	GetRenderedCheckoutImage() ([]byte, error)
}

type Refund struct{}

type Checkout struct{}
