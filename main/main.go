package main

import (
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/LovePelmeni/OnlineStore/StoreService/customers"
	"github.com/LovePelmeni/OnlineStore/StoreService/middlewares"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/LovePelmeni/OnlineStore/StoreService/products"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	APPLICATION_HOST = os.Getenv("APPLICATION_HOST")
	APPLICATION_PORT = os.Getenv("APPLICATION_PORT")

	EMAIL_APPLICATION_HOST = os.Getenv("EMAIL_APPLICATION_HOST")
	EMAIL_APPLICATION_PORT = os.Getenv("EMAIL_APPLICATION_PORT")

	ORDER_APPLICATION_HOST = os.Getenv("ORDER_APPLICATION_HOST")
	ORDER_APPLICATION_PORT = os.Getenv("ORDER_APPLICATION_PORT")
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
)

func init() {
	// Initializing API Servers...

	// Migrating into Postgresql
	Failed := models.Database.AutoMigrate(
		&models.Product{},
		&models.Customer{},
		&models.Cart{},
	)
	if Failed != nil {
		panic(fmt.Sprintf("Failed To Auto Migrate, Error: %s"))
	}
	LogFile, error := os.OpenFile("Main.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if error != nil {
		panic("Failed To Create Log file For Main.go File")
	}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	WarnLogger = log.New(LogFile, "WARNING: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
}

func main() {
	// ROUTER CONFIGURATION
	router := gin.Default()

	// TRUSTED PROXIES.
	router.SetTrustedProxies([]string{
		fmt.Sprintf("%s:%s", EMAIL_APPLICATION_HOST, EMAIL_APPLICATION_PORT),
		fmt.Sprintf("%s:%s", ORDER_APPLICATION_HOST, ORDER_APPLICATION_PORT),
	})

	// CORS CONFIGURATION
	router.Use(cors.New(cors.Config{

		AllowOrigins:     []string{},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowFiles:       true,
	}))

	// HEALTHCHECK

	router.GET("/healthcheck/", func(context *gin.Context) {
		context.JSON(http.StatusOK, nil)
	})

	// CUSTOMERS
	router.Use(middlewares.SetAuthHeaderMiddleware, middlewares.JwtAuthenticationMiddleware)
	{
		router.GET("get/profile/:customerId/", customers.GetCustomerProfile)   // Is Authenticated
		router.POST("create/customer/", customers.CreateCustomer)              // AllowAny
		router.PUT("update/customer/:customerId", customers.UpdateCustomer)    // Is Authenticated
		router.DELETE("delete/customer/:customerId", customers.DeleteCustomer) // Is Authenticated
	}

	// PRODUCTS

	// Getter Endpoints.
	router.Group("retrieve/").Use(middlewares.SetAuthHeaderMiddleware)
	{
		router.GET("all/products/:productId", products.GetProductsCatalog)
		router.GET("product/:productId", products.GetProduct)
	}

	// CUD Endpoints.
	router.Group("product/").Use(middlewares.SetAuthHeaderMiddleware, 
	middlewares.JwtAuthenticationMiddleware, middlewares.IsProductOwnerMiddleware)
	{ // Is Authenticated
		router.POST("create/", products.CreateProduct)             // permission for creating products requires.
		router.PUT("update/:productId", products.UpdateProduct)    // permission for own this product
		router.DELETE("delete/:productId", products.DeleteProduct) // permissions for own this product.
	}

	// Most ... Products.
	router.Group("get/most/").Use(middlewares.SetAuthHeaderMiddleware)
	{
		router.GET("/popular/week/products", products.GetTopWeekProducts) // AllowAny
		router.GET("/liked/products/", products.GetMostLikedProducts)
	}


	DebugLogger.Println("Running HTTP Server.")
	router.Run(fmt.Sprintf(":%s", APPLICATION_PORT))
}
