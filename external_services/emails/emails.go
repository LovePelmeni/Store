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

type grpcEmailSender struct{} // Implementation

// Overriden method for initializing gRPC Server Client.

func (this *grpcEmailSender) getClient() (*grpcControllers.NewEmailClient, error) {

	Connection, Error := grpc.Dial("%s:%s",
		GRPC_SERVER_HOST, GRPC_SERVER_PORT)
	grpcClient := grpcControllers.NewEmailClient(Connection)
	return grpcClient, Error
}

// Method for sending out default Emails without any concrete topic.

func (this *grpcEmailSender) SendDefaultEmail(customerEmail string, message string) (bool, error) {

	EmailRequestCredentials := grpcControllers.EmailDefaultParams{
		CustomerEmail: customerEmail,
		Message:       message,
	}

	client, ClientError := this.getClient()
	if ClientError != nil {
		panic(
			"Failed to Create GRPC Client. Check for the Server Running...")
	}
	response, ResponseError := client.SendEmail(EmailRequestCredentials)

	if ResponseError != nil {
		return false, exceptions.FailedRequest(
			ResponseError)
	}
	return response.Delivered, nil
}

// Method For sending Order Accepted Emails...

func (this *grpcEmailSender) SendAcceptEmail(customerEmail string, message string) (bool, error) {
	client, ClientError := this.getClient()
	if ClientError != nil {
		panic("Failed to Create GRPC Client.")
	}
	RequestParams := grpcControllers.EmailOrderParams{
		Status:        grpcControllers.STATUS_ACCEPTED,
		message:       message,
		CustomerEmail: customerEmail,
	}
	response, gRPCError := client.SendOrderEmail(RequestParams)
	if gRPCError != nil {
		DebugLogger.Println("Failed to send Request")
		return false, exceptions.FailedRequest(gRPCError)
	}
	return response.Delivered, nil
}

// Method for sending Order Rejected Emails...

func (this *grpcEmailSender) SendRejectEmail(customerEmail string, message string) (bool, error) {
	client, ClientError := this.getClient()
	if ClientError != nil {
		panic("Failed to create Client.")
	}
	RequestParams := grpcControllers.EmailOrderStatus{
		Status:        grpcControllers.STATUS_REJECTED,
		message:       message,
		CustomerEmail: customerEmail,
	}
	response, ResponseError := client.SendOrderEmail(RequestParams)
	return response.Delivered, ResponseError
}
