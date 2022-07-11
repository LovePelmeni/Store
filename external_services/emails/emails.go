package emails

import (
	"os"

	"log"

	"github.com/LovePelmeni/OnlineStore/EmailService/emails/proto/grpcControllers"
	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/exceptions"
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

func init() {
	LogFile, error := os.OpenFile("EmailGRPClog.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if error != nil {
		panic("Failed to Create Log File At Emails.go")
	}
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime)
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

type GrpcEmailSender struct {
	Client *grpcControllers.NewEmailClient
} // Implementation

// Overriden method for initializing gRPC Server Client.

func getClient() (*grpcControllers.NewEmailClient, error) {

	Connection, Error := grpc.Dial("%s:%s",
		GRPC_SERVER_HOST, GRPC_SERVER_PORT)
	grpcClient := grpcControllers.NewEmailClient(Connection)
	return grpcClient, Error
}

// Method for sending out default Emails without any concrete topic.

func (this *GrpcEmailSender) SendDefaultEmail(customerEmail string, message string) (bool, error) {

	EmailRequestCredentials := grpcControllers.EmailDefaultParams{
		CustomerEmail: customerEmail,
		Message:       message,
	}

	response, ResponseError := this.Client.SendEmail(EmailRequestCredentials)

	if ResponseError != nil {
		return false, exceptions.FailedRequest(
			ResponseError)
	}
	return response.Delivered, nil
}

// Method For sending Order Accepted Emails...

func (this *GrpcEmailSender) SendAcceptEmail(customerEmail string, message string) (bool, error) {

	RequestParams := grpcControllers.EmailOrderParams{
		Status:        grpcControllers.STATUS_ACCEPTED,
		message:       message,
		CustomerEmail: customerEmail,
	}
	response, gRPCError := this.Client.SendOrderEmail(RequestParams)
	if gRPCError != nil {
		DebugLogger.Println("Failed to send Request")
		return false, exceptions.FailedRequest(gRPCError)
	}
	return response.Delivered, nil
}

// Method for sending Order Rejected Emails...

func (this *GrpcEmailSender) SendRejectEmail(customerEmail string, message string) (bool, error) {

	RequestParams := grpcControllers.EmailOrderStatus{
		Status:        grpcControllers.STATUS_REJECTED,
		message:       message,
		CustomerEmail: customerEmail,
	}
	response, ResponseError := this.Client.SendOrderEmail(RequestParams)
	return response.Delivered, ResponseError
}
