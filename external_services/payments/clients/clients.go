package grpc_clients

import (
	"log"

	grpcControllers "github.com/LovePelmeni/OnlineStore/StoreService/external_services/payments/proto"
	"google.golang.org/grpc"
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

func NewGrpcServerConnection(Host string, Port string) *GrpcServerConnection {
	return &GrpcServerConnection{ServerHost: Host, ServerPort: Port}
}
func GetConnection() (*grpc.ClientConn, error)

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
