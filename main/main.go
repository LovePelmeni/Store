package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LovePelmeni/OnlineStore/StoreService/customers"
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
		AllowMethods:     []string{},
		AllowHeaders:     []string{},
		AllowCredentials: true,
		AllowFiles:       true,
	}))

	// HEALTHCHECK
	router.GET("/healthcheck/", func(context *gin.Context) {
		context.JSON(http.StatusOK, nil)
	})

	// CUSTOMERS

	router.GET("get/profile/:customerId/", customers.GetCustomerProfile)   // Is Authenticated
	router.POST("create/customer/", customers.CreateCustomer)              // AllowAny
	router.PUT("update/customer/:customerId", customers.UpdateCustomer)    // Is Authenticated
	router.DELETE("delete/customer/:customerId", customers.DeleteCustomer) // Is Authenticated

	// PRODUCTS

	router.GET("get/all/products/:productId", products.GetAllProducts)
	router.GET("get/product/:productId", products.GetProduct)          // Is Authenticated
	router.POST("create/product/", products.CreateProduct)             // permission for creating products requires.
	router.PUT("update/product/:productId", products.UpdateProduct)    // permission for own this product
	router.DELETE("delete/product/:productId", products.DeleteProduct) // permissions for own this product.

	// ADDITIONAL PRODUCT URLS:

	router.GET("get/most/popular/week/products", products.TopWeekProducts) // AllowAny
	router.GET("get/discounted/products/", products.GetDiscountedProducts)

	router.Run(fmt.Sprintf(":%s", APPLICATION_PORT))
}



