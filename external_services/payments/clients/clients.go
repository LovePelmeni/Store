package grpc_clients

import (
	"log"

	"fmt"
	"os"

	grpcControllers "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
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

func InitializeLoggers() {

}

func init() {

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
	GetClient() (*grpcControllers.PaymentIntentClient, error)
}

type PaymentSessionClientInterface interface {
	GetClient() (*grpcControllers.PaymentSessionClient, error)
}

type PaymentRefundClientInterface interface {
	GetClient() (*grpcControllers.RefundClient, error)
}

// Implementations

type GrpcServerConnection struct {
	ServerHost string
	ServerPort string
}

func NewGrpcServerConnection(grpcServerHost string, grpcServerPort string) *GrpcServerConnection {
	return &GrpcServerConnection{ServerHost: grpcServerHost, ServerPort: grpcServerPort}
}
func (this *GrpcServerConnection) GetConnection() (*grpc.ClientConn, error) {
	Connection, Error := grpc.Dial(
		fmt.Sprintf("%s:%s", this.ServerHost, this.ServerPort))
	if Connection == nil || Error != nil {
		panic("Failed to Connect To Payment Grpc Server.")
	}
	return Connection, nil
}

type PaymentIntentClient struct {
	Connection *GrpcServerConnectionInterface
}

func NewPaymentIntentClient(Connection *GrpcServerConnectionInterface) *PaymentIntentClient {
	return &PaymentIntentClient{Connection: Connection}
}

func (this *PaymentIntentClient) GetClient() (*grpcControllers.PaymentIntentClient, error)

type PaymentSessionClient struct {
	Connection *GrpcServerConnectionInterface
}

func NewPaymentSessionClient(Connection *GrpcServerConnectionInterface) *PaymentSessionClient {
	return &PaymentSessionClient{Connection: Connection}
}

func (this *PaymentSessionClient) GetClient() (*grpcControllers.PaymentSessionClient, error)

type PaymentRefundClient struct {
	Connection *GrpcServerConnectionInterface
}

func NewPaymentRefundClient(Connection *GrpcServerConnectionInterface) *PaymentRefundClient {
	return &PaymentRefundClient{Connection: Connection}
}

func (this *PaymentRefundClient) GetClient() *grpcControllers.RefundClient
