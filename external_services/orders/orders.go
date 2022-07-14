package orders

import (
	"context"
	"fmt"
	"log"
	"os"

	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
)

var (
	FIREBASE_DATABASE_NAME = os.Getenv("FIREBASE_DATABASE_NAME")
	FIREBASE_DATABASE_URL  = fmt.Sprintf("https://%s.firebaseio.com")
)

var (
	DebugLogger   *log.Logger
	ErrorLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
)

// Firebase RealTime Database Credentials..

var (
	StorageBucketID  = os.Getenv("STORAGE_BUCKET_ID")
	ProjectID        = os.Getenv("PROJECT_ID")
	ServiceAccountID = os.Getenv("SERVICE_ACCOUNT_ID")
)

func init() {

	LogFile, Exception := os.OpenFile("firebase_order.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if Exception != nil {
		panic(Exception)
	}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Llongfile|log.Ltime)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Llongfile|log.Ltime)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Llongfile|log.Ltime)
}

// Abstractions for the `Orders` Bounded Context...

type FirebaseInitializerInterface interface {
	// Interfaces for Application Intialization..
	// Requires all necessary method to be overridden, in order to provide availability
	InitializeFirebaseApplication()
	InitializeFirebaseDatabase()
}

type FirebaseOrderControllerInterface interface {
	// Interface, that is responsible for Managing `Orders` Real Time Database..
	// Provides Methods for Creating / Deleting Documents.
	CreateFirebaseOrderTransaction(OrderParams OrderCredentialsInterface,
		DatabaseInstance *FirebaseInitializer) (bool, error)
	DeleteFirebaseOrderTransaction(OrderId string) (bool, error)
}

type OrderCredentialsInterface interface {
	// Interface Responsible for Initial Info About the Order..
	Validate(Credentials struct{}) (bool, error)
}

type OrderControllerInterface interface {
	// Interface that represents the Main Controller Responsible For Handling Any Operations,
	// Related to the `Orders`
	CreateOrder(OrderCredentials OrderCredentialsInterface) (bool, error)
}

type FirebaseInitializer struct{}

// Firebase Application Initializer, Contains Method For Initializing Firebase Abstractions...

func (this *FirebaseInitializer) InitializeFirebaseDatabase(
	DatabaseName string, Application *firebase.App) *db.Client {
	// Method Initializes database collection ..
	Context, TimeoutError := context.WithTimeout(context.Background(), 10*time.Second)
	newDatabase, DatabaseError := Application.Database(Context)
	if DatabaseError != nil {
		ErrorLogger.Println("Failed to Initialize Database...")
	}
	defer TimeoutError()
	return newDatabase
}

func (this *FirebaseInitializer) InitializeFirebaseApplication() *firebase.App {

	// Method Initializes Firebase Database ...
	context := context.Background()
	config := &firebase.Config{
		DatabaseURL:      FIREBASE_DATABASE_URL,
		ProjectID:        ProjectID,
		StorageBucket:    StorageBucketID,
		ServiceAccountID: ServiceAccountID,
	}
	newApp, InitializeError := firebase.NewApp(context, config)
	if InitializeError != nil {
		ErrorLogger.Println(
			"Failed to Initialize Application.")
		panic(InitializeError)
	}
	return newApp
}

type FirebaseOrderController struct{}

func (this *FirebaseOrderController) CreateFirebaseOrderTransaction() {}

func (this *FirebaseOrderController) DeleteFirebaseOrderTransaction() {}

type OrderController struct{}

func (this *OrderController) CreateOrder(OrderCredentials OrderCredentialsInterface) (bool, error) {}

func (this *OrderController) CancelOrder(OrderId string) (bool, error) {}

type OrderCredentials struct{}

func (this *OrderCredentials) Validate(Credentials struct{}) (bool, error) {}
