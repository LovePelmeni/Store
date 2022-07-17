package emails

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"fmt"

	"github.com/LovePelmeni/EmailService/emails/proto/grpcControllers"
	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/exceptions"
	"github.com/mercari/go-circuitbreaker"
	"google.golang.org/grpc"
)

var (
	GRPC_SERVER_HOST = os.Getenv("GRPC_SERVER_HOST")
	GRPC_SERVER_PORT = os.Getenv("GRPC_SERVER_PORT")
)

var (
	DebugLogger  *log.Logger
	InfoLogger   *log.Logger
	ErrorLogger  *log.Logger
	WarningError *log.Logger
)

func InitializeLoggers() (bool, error) {

	LogFile, Error := os.OpenFile("EmailLogger.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Llongfile|log.Ltime)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Llongfile|log.Ltime)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return true, Error
}

func init() {

	Initialized, Error := InitializeLoggers()
	if Initialized != true || Error != nil {
		panic("Failed to Initialize Loggers in Emails.go ")
	}
}

type grpcEmailSenderInterface interface {
	// grpc Email Sender Interface.
	// Is Used For Comminucating with `Email Service` API
	// using gRPC. Provides Following API Methods for handling emails behaviour.

	SendDefaultEmail(
		customerEmail string,
		message string,
		backgroundImage ...[]byte,
	) (bool, error)

	SendAcceptEmail(
		customerEmail string,
		message string) (bool, error)

	SendRejectEmail(
		customerEmail string,
		message string,
	) (bool, error)
}

type grpcEmailClientInterface interface {
	// Email Interface represents Email grpc Client...
	getClient() (grpcControllers.EmailSenderClient, error)
}

// Implementation....

type grpcEmailClient struct {
	CircuitBreaker *circuitbreaker.CircuitBreaker
}

func NewGrpcEmailClient() *grpcEmailClient {
	return &grpcEmailClient{CircuitBreaker: circuitbreaker.New(

		circuitbreaker.WithOpenTimeout(20),
		circuitbreaker.WithOnStateChangeHookFn(

			func(old_state, new_state circuitbreaker.State) {
				if hasPrefix := strings.HasPrefix(
					strings.ToLower(string(new_state)), "cl"); hasPrefix == true {
					DebugLogger.Println(
						"Failed to Connect to Grpc Server of `Email` Service.",
					)
				}
			}),
	)}
}

func (this *grpcEmailClient) getClient() (grpcControllers.EmailSenderClient, error) {

	var EmailClient grpcControllers.EmailSenderClient
	if this.CircuitBreaker.Ready() != true {
		return nil, errors.New("Email Service Is Unavailable.")
	} // Checks for Circuit Breaker Availability..
	context, CancelError := context.WithTimeout(context.Background(), time.Second*10)

	Client, Error := this.CircuitBreaker.Do(context,

		func() (interface{}, error) {

			Connection, Error := grpc.Dial(fmt.Sprintf("%s:%s",
				GRPC_SERVER_HOST, GRPC_SERVER_PORT))

			grpcClient := grpcControllers.NewEmailSenderClient(Connection)
			defer CancelError()
			switch Error {

			case nil:
				EmailClient = grpcClient
				this.CircuitBreaker.Done(context, nil)
				return nil, Error

			default:
				this.CircuitBreaker.FailWithContext(context)
				return nil, Error
			}
		})

	_ = Client
	if Error != nil {
		ErrorLogger.Println("Failed to Initialize Grpc Email Client.")
	}
	return EmailClient, nil
}

type GrpcEmailSender struct {
	Client         grpcEmailClientInterface
	CircuitBreaker circuitbreaker.CircuitBreaker
} // Implementation

func NewGrpcEmailSender() *GrpcEmailSender {

	Client := NewGrpcEmailClient()           // Initializing Grpc Client
	NewCircuitBreaker := circuitbreaker.New( // Initializing new Circuit breaker....
		circuitbreaker.WithOpenTimeout(20),
		circuitbreaker.WithFailOnContextCancel(true),
		circuitbreaker.WithOnStateChangeHookFn(func(old_state, new_state circuitbreaker.State) {
			if new_state == "OPEN" {
				ErrorLogger.Println("CircuitBreaker is UP, Failure in the Email Service. State: " + new_state)
			} else {
				InfoLogger.Println("CircuitBreaker is closed.. Issue Fixed... State: " + new_state)
			}
		}),
	)
	if Client == nil {
		panic("Failed To Intialize Grpc Client for Email Service. Seems Like it did not respond.")
	}
	return &GrpcEmailSender{Client: Client, CircuitBreaker: *NewCircuitBreaker}
}

// Method for sending out default Emails without any concrete topic.

func (this *GrpcEmailSender) SendDefaultEmail(customerEmail string, message string) (bool, error) {

	if !this.CircuitBreaker.Ready() {
		return false, exceptions.ServiceUnavailable()
	} // Checking for Circuit Breaker Availability..

	grpcClient, Error := this.Client.getClient()
	if Error != nil {
		return false, errors.New("Connection GRPC Server Error.")
	}

	EmailRequestCredentials := grpcControllers.DefaultEmailParams{
		CustomerEmail: customerEmail,
		EmailMessage:  message,
	}

	var DeliveredResponse grpcControllers.EmailResponse // Variable that the Response is going to be stored in...
	Response, ResponseError := this.CircuitBreaker.Do(context.Background(),

		func() (interface{}, error) {

			RequestContext, CancelError := context.WithTimeout(context.Background(), time.Second*10)
			Response, Exception := grpcClient.SendEmail(RequestContext, &EmailRequestCredentials)

			defer CancelError()
			switch Exception {

			case nil:
				defer this.CircuitBreaker.Done(RequestContext, nil)
				DebugLogger.Println("Notification Has been Sended...")
				return Response.Delivered, nil

			default:
				this.CircuitBreaker.FailWithContext(RequestContext)
				return false, Exception
			}
		})

	_ = Response

	if ResponseError != nil {
		return false, exceptions.ServiceUnavailable()
	}

	return DeliveredResponse.Delivered, nil
}

// Method For sending Order Accepted Emails...

func (this *GrpcEmailSender) SendAcceptEmail(customerEmail string, message string) (bool, error) {

	if !this.CircuitBreaker.Ready() {
		return false, exceptions.ServiceUnavailable()
	} // Checking If Circuit Breaker Is Available...

	grpcClient, Error := this.Client.getClient()
	if Error != nil {
		return false, errors.New("Connection GRPC Server Error.")
	}

	var DeliveredResponse *grpcControllers.EmailResponse
	RequestParams := grpcControllers.OrderEmailParams{
		Status:        grpcControllers.OrderStatus_ACCEPTED,
		Message:       message,
		CustomerEmail: customerEmail,
	}

	Response, gRPCError := this.CircuitBreaker.Do(

		context.Background(),
		func() (interface{}, error) {

			RequestContext, CancelError := context.WithTimeout(context.Background(), time.Second*10)
			Response, Error := grpcClient.SendOrderEmail(context.Background(), &RequestParams)

			defer CancelError()
			switch Error {
			case nil:
				this.CircuitBreaker.Done(RequestContext, nil)
				DeliveredResponse = Response
				return true, nil
			default:
				DeliveredResponse = Response
				return false, exceptions.ServiceUnavailable()
			}
		})

	_ = Response

	if gRPCError != nil {
		DebugLogger.Println("Failed to send Request")
		return false, gRPCError
	}
	return DeliveredResponse.Delivered, nil
}

// Method for sending Order Rejected Emails...

func (this *GrpcEmailSender) SendRejectEmail(customerEmail string, message string) (bool, error) {

	if !this.CircuitBreaker.Ready() {
		return false, exceptions.ServiceUnavailable()
	} // Checks for Circuit Breaker Availability.

	var DeliveredResponse *grpcControllers.EmailResponse
	grpcClient, Error := this.Client.getClient()

	if Error != nil {
		return false, errors.New("Connection GRPC Server Error.")
	}

	RequestParams := grpcControllers.OrderEmailParams{
		Status:        grpcControllers.OrderStatus_REJECTED,
		Message:       message,
		CustomerEmail: customerEmail,
	}

	Response, ResponseError := this.CircuitBreaker.Do(

		context.Background(),
		func() (interface{}, error) {
			context, CancelError := context.WithTimeout(context.Background(), time.Second*10)
			Response, Error := grpcClient.SendOrderEmail(context, &RequestParams)

			defer CancelError()
			switch Error {
			case nil:
				DeliveredResponse = Response
				return nil, nil

			default:
				return nil, exceptions.ServiceUnavailable()
			}
		})
	_ = Response
	if ResponseError != nil {
		return false, ResponseError
	}
	return DeliveredResponse.Delivered, ResponseError
}
