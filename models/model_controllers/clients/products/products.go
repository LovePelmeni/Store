package products

import (
	"log"
	"os"

	"github.com/LovePelmeni/OnlineStore/StoreService/exceptions"
	"github.com/mercari/go-circuitbreaker"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
)

func InitializeLoggers() (bool, error) {

	LogFile, Error := os.OpenFile("Main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	WarnLogger = log.New(LogFile, "WARNING: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)

	if Error != nil {
		return false, Error
	}
	return true, nil
}

func init() {

	// Initializing Loggers...

	Initialized, Error := InitializeLoggers()
	if Error != nil || Initialized == false {
		panic(Error)
	}
}

type ProductRemoteTransactionParamsInterface interface {
	Validate(struct{}) (bool, error)
	GetCredentials(struct{}) (bool, error)
}

type ProductRemoteTransactionControllerInterface interface {
	CreateRemoteProduct(Params *ProductRemoteTransactionParamsInterface) (bool, error)
	UpdateRemoteProduct(Params *ProductRemoteTransactionParamsInterface) (bool, error)
	DeleteRemoteProduct(Params *ProductRemoteTransactionParamsInterface) (bool, error)
}

// Implementations...

type ProductRemoteTransactionParams struct {
	Credentials struct{}
}

func NewProductRemoteTransactionParams(Credentials struct{}) *ProductRemoteTransactionParams {
	return &ProductRemoteTransactionParams{Credentials: Credentials}
}

func (this *ProductRemoteTransactionParams) Validate(
	Credentials *ProductRemoteTransactionParamsInterface) (bool, error) {
	return true, nil
}

func (this *ProductRemoteTransactionParams) GetCredentials(Credentials *ProductRemoteTransactionParamsInterface) (
	*ProductRemoteTransactionParamsInterface, error) {

	if ValidatedCredentials, ValidationError := this.Validate(
		Credentials); !ValidatedCredentials || ValidationError != nil {
		DebugLogger.Println("Invalid Credentials to Create a Product.")
		return false, exceptions.ValidationError()
	}
	return Credentials, nil
}

type ProductRemoteTransactionController struct {
	GrpcProductClient GrpcProductControllers.ProductClient
	CircuitBreaker    circuitbreaker.CircuitBreaker
}

func NewProductRemoteTransactionController() *ProductRemoteTransactionController {
	return &ProductRemoteTransactionController{}
}

func (this *ProductRemoteTransactionController) CreateRemoteProduct() {

}

func (this *ProductRemoteTransactionController) UpdateRemoteProduct() {

}
func (this *ProductRemoteTransactionController) DeleteRemoteProduct() {

}
