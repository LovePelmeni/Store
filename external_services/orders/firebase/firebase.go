package firebase

import (
	"context"
	"errors"
	"log"
	"time"

	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/LovePelmeni/StoreService/external_services/orders"
)

// Collection Properties

var (
	OrdersCollectionName = "Orders"
)

// Firebase Credentials

var (
	FIREBASE_DATABASE_NAME = os.Getenv("FIREBASE_DATABASE_NAME")
	FirebaseDatabaseUrl    = fmt.Sprintf("https://%s.firebaseio.com")
)

var (
	StorageBucketID  string
	ProjectID        string
	ServiceAccountID string
)

// Loggers

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
)

// Initializing Package Properties...

func InitializeFirebaseCredentials() bool {

	StorageBucketID = os.Getenv("STORAGE_BUCKET_ID")
	ProjectID = os.Getenv("PROJECT_ID")
	ServiceAccountID = os.Getenv("SERVICE_ACCOUNT_ID")

	return true
}

func InitializeLoggers() bool {
	LogFile, Exception := os.OpenFile("firebase_order.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if Exception != nil {
		panic(Exception)
	}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Llongfile|log.Ltime)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Llongfile|log.Ltime)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Llongfile|log.Ltime)
	return true
}

func init() {

	InitializedLoggers := InitializeLoggers()
	if InitializedLoggers != true {
		panic("Failed to Initialize Loggers in ")
	}

	InitializedFirebaseCredentials := InitializeFirebaseCredentials()
	if InitializedFirebaseCredentials != true {
		panic("Failed to Initialize Firebase Credentials.")
	}
}

// Abstractions...

type FirebaseConfigInterface interface {
	// Configuration of the Firebase Real Time Database...
	// Parameters:

	// {
	//  PROJECT_ID - Id of the Google Cloud Project On Firebase.
	//  Storage Bucket ID - Id of the Firebase Storage Bucket.
	//  Firebase Database Url - Url to the Database. Consists of "https://firebaseio.com/DATABASE_NAME/"
	// 	Service Account ID - Id of the Account, that Serving this Database..
	// }
}

type FirebaseInitializerInterface interface {
	// Interface for Initialization Firebase Entities, such as Database, DatabaseCollection, App
	InitializeFirebaseApplication() (*firebase.App, error)
	InitializeFirebaseDatabase(Application *firebase.App) (*db.Client, error)
	InitializeFirebaseCollection(CollectionName string) (*db.Ref, error)
}

type FirebaseDatabaseOrderManagerInterface interface {
	// Interface for Managing `Order` Collection in Firebase Real Time Database
	CreateOrder(OrderCredentials *orders.OrderCredentialsInterface) (bool, error)
	CancelOrder(OrderId string) (bool, error)
}

// Implementations...

type FirebaseConfig struct {
	// Implementation of the Firebase Configuration Structure.
	// Contains All Necessary Credentials for the Successful Database Connection.
	ProjectID           string
	ServiceAccountID    string
	StorageBucketID     string
	FirebaseDatabaseUrl string
}

func NewFirebaseConfig() *FirebaseConfig {
	// Initializes New Config..
	return &FirebaseConfig{
		ProjectID:           ProjectID,
		ServiceAccountID:    ServiceAccountID,
		StorageBucketID:     StorageBucketID,
		FirebaseDatabaseUrl: FirebaseDatabaseUrl,
	}
}

type FirebaseInitializer struct {
	// Implementation of the `FirebaseInitializerInterface`.

	Config *FirebaseConfig // Firebase Project Configuration.
}

// Firebase Application Initializer, Contains Method For Initializing Firebase Abstractions...

func NewFirebaseInitializer(Config FirebaseConfig) *FirebaseInitializer {
	// Initializes New Instanse of `FirebaseInitializer`
	return &FirebaseInitializer{Config: &Config}
}

func (this *FirebaseInitializer) InitializeFirebaseDatabase(Application *firebase.App) (*db.Client, error) {
	// Method Initializes database collection ..
	Context, TimeoutError := context.WithTimeout(context.Background(), 10*time.Second)
	newDatabase, DatabaseError := Application.Database(Context)
	if DatabaseError != nil {
		ErrorLogger.Println("Failed to Initialize Database...")
		return nil, DatabaseError
	}
	defer TimeoutError()
	return newDatabase, nil
}

func (this *FirebaseInitializer) InitializeFirebaseCollection(CollectionName string) (*db.Ref, error) {

	application, Exception := this.InitializeFirebaseApplication()
	if Exception != nil {
		ErrorLogger.Println("Failed to initialize Firebase Application Instance.")
		return nil, Exception
	}

	database, DatabaseError := this.InitializeFirebaseDatabase(application)
	if DatabaseError != nil {
		ErrorLogger.Println("Failed to Initialize Firebase Database Instance.")
		return nil, DatabaseError
	}

	return database.NewRef(CollectionName), nil // returns initialized Collection...
}

func (this *FirebaseInitializer) InitializeFirebaseApplication() (*firebase.App, error) {

	// Method Initializes Firebase Database ...
	context := context.Background()
	config := &firebase.Config{
		DatabaseURL:      this.Config.FirebaseDatabaseUrl,
		ProjectID:        this.Config.ProjectID,
		StorageBucket:    this.Config.StorageBucketID,
		ServiceAccountID: this.Config.ServiceAccountID,
	}
	newApp, InitializeError := firebase.NewApp(context, config)
	if InitializeError != nil {
		ErrorLogger.Println(
			"Failed to Initialize Application.")
		panic(InitializeError)
	}
	return newApp, nil
}

type FirebaseDatabaseOrderManager struct {
	FirebaseInitializer FirebaseInitializerInterface // Initializer for Getting Access to the Firebase Instances..
}

func NewFirebaseDatabaseOrderManager(FirebaseInitializer FirebaseInitializerInterface) *FirebaseDatabaseOrderManager {
	return &FirebaseDatabaseOrderManager{
		FirebaseInitializer: FirebaseInitializer,
	}
}

func (this *FirebaseDatabaseOrderManager) CreateOrder(OrderCredentials *orders.OrderCredentialsInterface) (bool, error) {

	CollectionReference, Exception := this.FirebaseInitializer.InitializeFirebaseCollection(OrdersCollectionName)
	if Exception != nil {
		ErrorLogger.Println(
			"Failed To Initialize Order Collection.")
		return false, errors.New("Failed to Initialize Collection.")
	}

	DatabaseContext, CancelMethod := context.WithTimeout(context.Background(), time.Second*20)

	TransactionError := CollectionReference.Set(DatabaseContext, OrderCredentials)
	if TransactionError != nil {
		ErrorLogger.Println("Failed to Save")
		return false, errors.New("Failed to Save New Order.")
	}

	defer CancelMethod()
	return true, nil
}

func (this *FirebaseDatabaseOrderManager) CancelOrder(OrderId string) (bool, error) {

	DatabaseContext, CancelMethod := context.WithTimeout(context.Background(), time.Second*20)
	CollectionReference, Error := this.FirebaseInitializer.InitializeFirebaseCollection(OrdersCollectionName)
	if Error != nil {
		ErrorLogger.Println("Failed to Initialize Collection Order Reference.")
		return false, errors.New("Failed to Initialize Collection.")
	}

	CanceledError := CollectionReference.Child(
		OrderId).Delete(DatabaseContext)

	if CanceledError != nil {
		ErrorLogger.Println("Failed to Cancel Order with ID: " + OrderId)
		return false, errors.New("Failed to Cancel Order with ID: " + OrderId)
	}
	defer CancelMethod()
	return true, nil
}
