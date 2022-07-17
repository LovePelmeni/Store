package emails

import (
	"os"
	"log"
	"github.com/LovePelmeni/OnlineStore/EmailService/emails/proto/grpcControllers"
	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/exceptions"
	"github.com/mercari/go-circuitbreaker"
	"google.golang.org/grpc"
	"context"
	"errors"
	"time"
	"strings"
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
	if Initialized != true || Error != nil {panic("Failed to Initialize Loggers in Emails.go ")}
}


type grpcEmailSenderInterface interface {
	// grpc Email Sender Interface.
	// Is Used For Comminucating with `Email Service` API
	// using gRPC. Provides Following API Methods for handling emails behaviour.

	getClient() (*grpcControllers.NewEmailClient, error)

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
	getClient() (*grpcControllers.NewEmailClient, error)
}


type grpcEmailClient struct {
	CircuitBreaker *circuitbreaker.CircuitBreaker
}

func NewGrpcEmailClient() (*grpcEmailClient) {
	return &grpcEmailClient{CircuitBreaker: circuitbreaker.New(

		circuitbreaker.WithOpenTimeout(20), 
		circuitbreaker.WithOnStateChangeHookFn(

		func(old_state, new_state circuitbreaker.State){
			if hasPrefix := strings.HasPrefix(
			strings.ToLower(string(new_state)), "cl"); hasPrefix == true {
				DebugLogger.Println(
					"Failed to Connect to Grpc Server of `Email` Service.",
				)
			}
		}),
	)}
}



func (this *grpcEmailClient) getClient() (*grpcControllers.NewEmailClient, error) {

	if this.CircuitBreaker.Ready() != true {return nil, errors.New("Email Service Is Unavailable.")}
	context, CancelError := context.WithTimeout(context.Background(), time.Second * 10)
	defer CancelError()

	return this.CircuitBreaker.Do(context, func() (interface{}, error){

	Connection, Error := grpc.Dial("%s:%s",
		GRPC_SERVER_HOST, GRPC_SERVER_PORT)
	grpcClient := grpcControllers.NewEmailClient(Connection)
	if Error != nil {ErrorLogger.Println(
	"Failed to Connect To Email Grpc Server: Error " + Error.Error());

    this.CircuitBreaker.FailWithContext(context)}else{
	this.CircuitBreaker.Done(context, nil)}
	return grpcClient, Error
	})
}


type GrpcEmailSender struct {
	Client grpcEmailClientInterface
	CircuitBreaker circuitbreaker.CircuitBreaker 
} // Implementation


func NewGrpcEmailSender() (*GrpcEmailSender){

	Client := NewGrpcEmailClient() // Initializing Grpc Client 
	NewCircuitBreaker := circuitbreaker.New( // Initializing new Circuit breaker....
		circuitbreaker.WithOpenTimeout(20),
		circuitbreaker.WithFailOnContextCancel(true),
		circuitbreaker.WithOnStateChangeHookFn(func(old_state, new_state circuitbreaker.State){
			if new_state == "OPEN" {ErrorLogger.Println("CircuitBreaker is UP, Failure in the Email Service. State: " + new_state)}else{
				InfoLogger.Println("CircuitBreaker is closed.. Issue Fixed... State: " + new_state)
			}
		}),
	)
	if Client == nil {panic("Failed To Intialize Grpc Client for Email Service. Seems Like it did not respond.")}
	return &GrpcEmailSender{Client: Client, CircuitBreaker: *NewCircuitBreaker}
}

// Method for sending out default Emails without any concrete topic.

func (this *GrpcEmailSender) SendDefaultEmail(customerEmail string, message string) (bool, error) {

	grpcClient, Error := this.Client.getClient()
	if Error != nil {return false, errors.New("Connection GRPC Server Error.")}

	EmailRequestCredentials := grpcControllers.EmailDefaultParams{
		CustomerEmail: customerEmail,
		Message:       message,
	}

	RequestContext, CancelError := context.WithTimeout(context.Background(), time.Second * 10)
	Delivered, ResponseError := this.CircuitBreaker.Do(RequestContext,
	   
		
		func() (interface{}, error){ 
			Response, Exception := grpcClient.SendEmail(EmailRequestCredentials);
			if Exception != nil {this.CircuitBreaker.FailWithContext(RequestContext);
		    return false, Exception} else{

			defer this.CircuitBreaker.Done(RequestContext, nil)	
			DebugLogger.Println("Notification Has been Sended...")
			return Response.Delivered, nil 
			}
		})
	

	if ResponseError != nil {
		return false, exceptions.FailedRequest(
		ResponseError)
	}

	defer CancelError()
	return Delivered, nil
}

// Method For sending Order Accepted Emails...

func (this *GrpcEmailSender) SendAcceptEmail(customerEmail string, message string) (bool, error) {

	grpcClient, Error := this.Client.getClient()
	if Error != nil {return false, errors.New("Connection GRPC Server Error.")}

	RequestParams := grpcControllers.EmailOrderParams{
		Status:        grpcControllers.STATUS_ACCEPTED,
		message:       message,
		CustomerEmail: customerEmail,
	}
	response, gRPCError := grpcClient.SendOrderEmail(RequestParams)
	if gRPCError != nil {
		DebugLogger.Println("Failed to send Request")
		return false, exceptions.FailedRequest(gRPCError)
	}
	return response.Delivered, nil
}

// Method for sending Order Rejected Emails...

func (this *GrpcEmailSender) SendRejectEmail(customerEmail string, message string) (bool, error) {

	grpcClient, Error := this.Client.getClient()
	if Error != nil {return false, errors.New("Connection GRPC Server Error.")}

	RequestParams := grpcControllers.EmailOrderStatus{
		Status:        grpcControllers.STATUS_REJECTED,
		message:       message,
		CustomerEmail: customerEmail,
	}
	response, ResponseError := grpcClient.SendOrderEmail(RequestParams)
	return response.Delivered, ResponseError
}
