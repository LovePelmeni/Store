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
	"github.com/LovePelmeni/OnlineStore/StoreService/external_services/exceptions"

	// "github.com/LovePelmeni/OnlineStore/StoreService/external_services/orders"
	"github.com/mercari/go-circuitbreaker"
)

// Collection Properties

var (
	OrdersCollectionName = "Orders"
)

// Firebase Credentials

var (
	FIREBASE_DATABASE_NAME = os.Getenv("FIREBASE_DATABASE_NAME")
	FirebaseDatabaseUrl    = fmt.Sprintf("https://%s.firebaseio.com", FIREBASE_DATABASE_NAME)
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

	StorageBucketID = os.Getenv("FIREBASE_STORAGE_BUCKET_ID")
	ProjectID = os.Getenv("FIREBASE_PROJECT_ID")
	ServiceAccountID = os.Getenv("FIREBASE_SERVICE_ACCOUNT_ID")

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

//go:generate mockgen generate -destination=StoreService/mocks/firebase.go --build_flags=--mod=mod . FirebaseConfigInterface
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

//go:generate mockgen generate -destination=StoreService/mocks/firebase.go --build_flags=--mod=mod . FirebaseInitializerInteface

type FirebaseInitializerInterface interface {
	// Interface for Initialization Firebase Entities, such as Database, DatabaseCollection, App
	InitializeFirebaseApplication() (*firebase.App, error)
	InitializeFirebaseDatabase(Application *firebase.App) (*db.Client, error)
	InitializeFirebaseCollection(CollectionName string) (*db.Ref, error)
}

//go:generate mockgen generate -destination=StoreService/mocks/firebase.go --build_flags=--mod=mod . FirebaseDatabaseOrderManagerInterface

type FirebaseDatabaseOrderManagerInterface interface {
	// Interface for Managing `Order` Collection in Firebase Real Time Database
	CreateOrder(OrderCredentials struct{}) (bool, error)
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

	Config         *FirebaseConfig // Firebase Project Configuration.
	CircuitBreaker circuitbreaker.CircuitBreaker
}

// Firebase Application Initializer, Contains Method For Initializing Firebase Abstractions...

func NewFirebaseInitializer(Config FirebaseConfig) *FirebaseInitializer {
	// Initializes New Instanse of `FirebaseInitializer`
	return &FirebaseInitializer{Config: &Config,
		CircuitBreaker: *circuitbreaker.New(
			circuitbreaker.WithOpenTimeout(20),
			circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
				if newState == "OPEN" {
					ErrorLogger.Println("Firebase database is not available. Any Operations Not Allowed. Time: " + time.Now().String())
				}
				if newState == "CLOSED" {
					InfoLogger.Println("Firebase database is Available and Ready For New Requests... Time: " + time.Now().String())
				}
			}),
		)}
}

func (this *FirebaseInitializer) InitializeFirebaseDatabase(Application *firebase.App) (*db.Client, error) {
	// Method Initializes database collection ..

	if !this.CircuitBreaker.Ready() {
		return nil, exceptions.ServiceUnavailable()
	}

	var databaseClientInstance *db.Client
	_, DatabaseError := this.CircuitBreaker.Do(

		context.Background(),
		func() (interface{}, error) {

			Context, TimeoutError := context.WithTimeout(context.Background(), 10*time.Second)
			newDatabase, Error := Application.Database(Context)

			defer TimeoutError()
			switch Error {

			case nil:
				databaseClientInstance = newDatabase
				this.CircuitBreaker.Done(Context, nil)
				return true, Error
			default:
				this.CircuitBreaker.FailWithContext(Context)
				return false, exceptions.ServiceUnavailable()
			}
		},
	)
	if DatabaseError != nil {
		return nil, DatabaseError
	}
	return databaseClientInstance, nil
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
	CircuitBreaker      circuitbreaker.CircuitBreaker
}

func NewFirebaseDatabaseOrderManager(FirebaseInitializer FirebaseInitializerInterface) *FirebaseDatabaseOrderManager {
	return &FirebaseDatabaseOrderManager{
		FirebaseInitializer: FirebaseInitializer,
		CircuitBreaker: *circuitbreaker.New(
			circuitbreaker.WithHalfOpenMaxSuccesses(10),
			circuitbreaker.WithOpenTimeout(20),
			circuitbreaker.WithOnStateChangeHookFn(func(oldState, newState circuitbreaker.State) {
				if newState == "OPEN" {
					ErrorLogger.Println("Firebase database is not available. Any Operations Not Allowed. Time: " + time.Now().String())
				}
				if newState == "CLOSED" {
					InfoLogger.Println("Firebase database is Available and Ready For New Requests... Time: " + time.Now().String())
				}
			}),
		)}
}

func (this *FirebaseDatabaseOrderManager) CreateOrder(OrderCredentials struct{}) (bool, error) {

	if !this.CircuitBreaker.Ready() {
		return false, exceptions.ServiceUnavailable()
	} // Checking for CircuitBreaker State..

	CollectionReference, Exception := this.FirebaseInitializer.InitializeFirebaseCollection(OrdersCollectionName)
	if Exception != nil {
		ErrorLogger.Println(
			"Failed To Initialize Order Collection.")
		return false, errors.New("Failed to Initialize Collection.")
	}

	var Created bool

	_, Error := this.CircuitBreaker.Do(context.Background(), func() (interface{}, error) {

		DatabaseContext, CancelMethod := context.WithTimeout(context.Background(), time.Second*10)
		TransactionError := CollectionReference.Set(DatabaseContext, OrderCredentials)

		defer CancelMethod()
		switch TransactionError {
		case nil:
			this.CircuitBreaker.Done(DatabaseContext, context.Canceled)
			Created = true
			return nil, nil
		default:
			this.CircuitBreaker.FailWithContext(DatabaseContext)
			Created = false
			return nil, exceptions.ServiceUnavailable()
		}
	})
	if Error != nil {
		return false, exceptions.ServiceUnavailable()
	}
	return Created, Error
}

func (this *FirebaseDatabaseOrderManager) CancelOrder(OrderId string) (bool, error) {

	if !this.CircuitBreaker.Ready() {
		return false, exceptions.ServiceUnavailable()
	} // Checking CircuitBreaker State..

	CollectionReference, Error := this.FirebaseInitializer.InitializeFirebaseCollection(OrdersCollectionName)
	if Error != nil {
		ErrorLogger.Println("Failed to Initialize Collection Order Reference.")
		return false, errors.New("Failed to Initialize Collection.")
	}

	var Canceled bool

	_, TransactionError := this.CircuitBreaker.Do(context.Background(), func() (interface{}, error) {

		DatabaseContext, CancelMethod := context.WithTimeout(context.Background(), time.Second*20)
		TransactionError := CollectionReference.Child(
			OrderId).Delete(DatabaseContext)

		defer CancelMethod()
		switch TransactionError {

		case nil:
			this.CircuitBreaker.Done(DatabaseContext, nil)
			return nil, nil
		default:
			this.CircuitBreaker.FailWithContext(DatabaseContext)
			return nil, exceptions.ServiceUnavailable()
		}
	})
	return Canceled, TransactionError
}
