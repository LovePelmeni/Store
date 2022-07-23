package main

import (
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/LovePelmeni/Store/customers"
	"github.com/LovePelmeni/Store/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/LovePelmeni/Store/models"
	"github.com/LovePelmeni/Store/products"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/utrack/gin-csrf"
)

var (
	APPLICATION_HOST = os.Getenv("APPLICATION_HOST")
	APPLICATION_PORT = os.Getenv("APPLICATION_PORT")

	EMAIL_APPLICATION_HOST = os.Getenv("EMAIL_APPLICATION_HOST")
	EMAIL_APPLICATION_PORT = os.Getenv("EMAIL_APPLICATION_PORT")

	ORDER_APPLICATION_HOST = os.Getenv("ORDER_APPLICATION_HOST")
	ORDER_APPLICATION_PORT = os.Getenv("ORDER_APPLICATION_PORT")

	PAYMENT_APPLICATION_HOST = os.Getenv("PAYMENT_APPLICATION_HOST")
	PAYMENT_APPLICATION_PORT = os.Getenv("PAYMENT_APPLICATION_PORT")
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

	// Migrating into Postgresql
	Failed := models.Database.AutoMigrate(
		&models.Cart{},
		&models.Customer{},
		&models.Product{},
	)

	if Failed != nil {
		panic(fmt.Sprintf("Failed To Auto Migrate, Error: %s", Failed.Error()))
	}
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

		AllowOrigins:     []string{fmt.Sprintf("http://%s:%s", EMAIL_APPLICATION_HOST, EMAIL_APPLICATION_PORT)},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowFiles:       true,
	}))

	// router.Use(csrf.Middleware(csrf.Options{
	// 	Secret: os.Getenv("CSRF_SECRET_KEY"),
	// 	ErrorFunc: func(context *gin.Context) {
	// 		context.String(400, "CSRF token mismatch")
	// 		context.Abort()
	// 	},
	// }))

	
	// HEALTHCHECK

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, nil)
	})

	// Prometheus metrics... 

	http.Handle("/metrics", promhttp.Handler())

	router.Group("customer/")
	{
	router.POST("create/", customers.CreateCustomerRestController) // AllowAny
	router.Use(middlewares.JwtAuthenticationMiddleware())
	{
		router.GET("get/profile/", customers.GetCustomerProfileRestController)    // Is Authenticated   // AllowAny
		router.PUT("update/", customers.UpdateCustomerRestController)    // Is Authenticated
		router.DELETE("delete/", customers.DeleteCustomerRestController) // Is Authenticated
	}}

	// PRODUCTS

	// Getter Endpoints.
	router.Group("retrieve/").Use(middlewares.SetAuthHeaderMiddleware())
	{
		router.GET("all/products/", products.GetProductsCatalog)
		router.GET("product/", products.GetProduct)
	}

	// CUD Endpoints.
	router.Group("product/").Use(middlewares.SetAuthHeaderMiddleware(),
		middlewares.JwtAuthenticationMiddleware(), middlewares.IsProductOwnerMiddleware())
	{ // Is Authenticated
		router.POST("create/", products.CreateProduct)   // permission for creating products requires.
		router.PUT("update/", products.UpdateProduct)    // permission for own this product
		router.DELETE("delete/", products.DeleteProduct) // permissions for own this product.
	}

	// Most ... Products.
	router.Group("get/most").Use(middlewares.SetAuthHeaderMiddleware())
	{
		router.GET("/popular/week/products", products.GetTopWeekProducts) // AllowAny
	}
	DebugLogger.Println("Running HTTP Server...")
	http.ListenAndServe(fmt.Sprintf("%s:%s", APPLICATION_HOST, APPLICATION_PORT), router)
}
