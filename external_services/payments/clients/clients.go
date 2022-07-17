package grpc_clients

import (
	"log"

	"fmt"
	"os"

	"context"

	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/exceptions"
	grpcControllers "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
	"github.com/mercari/go-circuitbreaker"
	"google.golang.org/grpc"
)

var (
	PAYMENT_GRPC_SERVER_HOST = os.Getenv("PAYMENT_GRPC_SERVER_HOST")
	PAYMENT_GRPC_SERVER_PORT = os.Getenv("PAYMENT_GRPC_SERVER_PORT")
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
)

func InitializeLoggers() (bool, error) {

	LogFile, Error := os.OpenFile("Main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	WarningLogger = log.New(LogFile, "WARNING: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)

	if Error != nil {
		return false, Error
	}
	return true, nil
}

func init() {
	Initialized, Errors := InitializeLoggers()
	if Initialized != true || Errors != nil {
		panic(
			"Failed to Initialize Loggers for Grpc Payment Clients...")
	}
}

// Abstraction Clients.

type GrpcServerConnectionInterface interface {
	// Interface  represents the Connection Layer for gRPC Client.
	// Required Params:
	// - grpc server host
	// - grpc server port
	GetConnection() (*grpc.ClientConn, error)
}

type PaymentIntentClientInterface interface {
	GetClient() (grpcControllers.PaymentIntentClient, error)
}

type PaymentSessionClientInterface interface {
	GetClient() (grpcControllers.PaymentSessionClient, error)
}

type PaymentRefundClientInterface interface {
	GetClient() (grpcControllers.RefundClient, error)
}

// Implementations

type GrpcServerConnection struct {
	ServerHost     string
	ServerPort     string
	CircuitBreaker circuitbreaker.CircuitBreaker
}

func NewGrpcServerConnection(grpcServerHost string, grpcServerPort string) *GrpcServerConnection {
	return &GrpcServerConnection{ServerHost: grpcServerHost, ServerPort: grpcServerPort}
}

func (this *GrpcServerConnection) GetConnection() (*grpc.ClientConn, error) {

	var ConnectionObj *grpc.ClientConn
	Response, Error := this.CircuitBreaker.Do(

		context.Background(),
		func() (interface{}, error) {

			Connection, Error := grpc.Dial(
				fmt.Sprintf("%s:%s", this.ServerHost, this.ServerPort))
			if Connection == nil || Error != nil {
				this.CircuitBreaker.FailWithContext(nil)
				return nil, exceptions.ServiceUnavailable()
			} else {
				ConnectionObj = Connection
				return nil, nil
			}
		})

	_ = Response

	if Error != nil {
		return nil, Error
	}
	return ConnectionObj, nil
}

// Clients Goes There....

type PaymentIntentClient struct {
	Connection GrpcServerConnectionInterface
}

func NewPaymentIntentClient(Connection GrpcServerConnectionInterface) *PaymentIntentClient {
	return &PaymentIntentClient{Connection: Connection}
}

func (this *PaymentIntentClient) GetClient() (grpcControllers.PaymentIntentClient, error) {
	Connection, Error := this.Connection.GetConnection()
	if Error != nil {
		return nil, exceptions.ServiceUnavailable()
	}
	return grpcControllers.NewPaymentIntentClient(Connection), nil
}

type PaymentSessionClient struct {
	Connection GrpcServerConnectionInterface
}

func NewPaymentSessionClient(Connection GrpcServerConnectionInterface) *PaymentSessionClient {
	return &PaymentSessionClient{Connection: Connection}
}

func (this *PaymentSessionClient) GetClient() (grpcControllers.PaymentSessionClient, error) {
	Connection, Error := this.Connection.GetConnection()
	if Error != nil {
		return nil, Error
	}
	return grpcControllers.NewPaymentSessionClient(Connection), nil
}

type PaymentRefundClient struct {
	Connection GrpcServerConnectionInterface
}

func NewPaymentRefundClient(Connection GrpcServerConnectionInterface) *PaymentRefundClient {
	return &PaymentRefundClient{Connection: Connection}
}

func (this *PaymentRefundClient) GetClient() grpcControllers.RefundClient {
	ServerConnection, Error := this.Connection.GetConnection()
	if Error != nil {
		return nil
	}
	return grpcControllers.NewRefundClient(ServerConnection)
}
