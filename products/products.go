package products

import (
	"log"
	"net/http"

	"encoding/json"

	"fmt"
	"sync"

	"io"
	"os"

	"strconv"

	"github.com/LovePelmeni/OnlineStore/StoreService/authentication"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/gin-gonic/gin"
)

var (
	DebugLogger   *log.Logger
	ErrorLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
)

var product models.Product
var products []models.Product
var customer models.Customer
var ProductValidator = models.NewProductModelValidator()

type CustomerBalance string

func InitializeLoggers() (bool, error) {

	LogFile, Error := os.OpenFile("Products.Log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if Error != nil {
		return false, Error
	}
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Llongfile|log.Ltime)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Llongfile|log.Ltime)
	WarningLogger = log.New(LogFile, "WARNING: ", log.Ldate|log.Llongfile|log.Ltime)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Llongfile|log.Ltime)
	return true, nil
}

func init() {
	Initialized, Error := InitializeLoggers()
	if Initialized != true || Error != nil {
		panic("Failed To Initialize Loggers.")
	}
}

// CUD Rest Controllers...

func CreateProduct(context *gin.Context) {

	ProductValidator := models.NewProductModelValidator()
	newProductCreated, Errors := product.CreateObject(
		map[string]string{
			"ProductName":        context.PostForm("ProductName"),
			"ProductDescription": context.PostForm("ProductDescription"),
			"ProductPrice":       context.PostForm("ProductPrice"),
			"Currency":           context.PostForm("Currency"),
		},
		ProductValidator,
	)
	if newProductCreated == false || Errors != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Errors": Errors})
	} else {
		context.JSON(http.StatusCreated, nil)
	}
}

func UpdateProduct(context *gin.Context) {

	var decodedData struct {
		ProductName        string
		ProductDescription string
		ProductPrice       float64
	}
	productId := context.Query("productId")
	bodyData, Error := io.ReadAll(context.Request.Body)
	json.Unmarshal(bodyData, &decodedData)

	if Error != nil {
		InfoLogger.Println(
			"Failed to parse Product Form Data.")
		context.JSON(http.StatusBadRequest, nil)
	}

	Product := models.Database.Table("products").Where("id = ?", productId) // Receives the object..
	Updated, Errors := product.UpdateObject(productId, decodedData, ProductValidator)
	if Product.Error != nil || !Updated || Errors != nil {
		context.JSON(http.StatusBadRequest, nil)
	}
	context.JSON(http.StatusCreated, nil)
}

func DeleteProduct(context *gin.Context) {

	ProductId := context.Query("ProductId")
	Deleted, Errors := product.DeleteObject(ProductId)
	if Deleted != true || Errors != nil {
		context.JSON(
			http.StatusNotImplemented, gin.H{"errors": Errors})
	} else {
		context.JSON(http.StatusOK, nil)
	}
}

// Getter Rest Controllers...

func GetTopWeekProducts(context *gin.Context) {
	productsQuery := models.Database.Table(
		"products").Order("likes desc").Limit(10)
	context.JSON(http.StatusOK, gin.H{"query": productsQuery})
}

func CustomerBalanceReceiver(ResponseChannel chan CustomerBalance, CustomerUsername string, CustomerEmail string) {

	Customer := models.Database.Table("customers").Where( // Receiving Customer.
		"Email = ? AND username = ?",
		CustomerEmail,
		CustomerUsername).First(&customer)

	if Customer.Error != nil {

		var CustomerPaymentBalance = struct{ Balance string }{}

		client := http.Client{}
		requestUrl := fmt.Sprintf("http://%s:%s/get/customer/info/",
			os.Getenv("PAYMENT_APPLICATION_HOST"), os.Getenv("PAYMENT_APPLICATION_PORT"))

		requestUrl += fmt.Sprintf("?customerId=%s", fmt.Sprint(customer.ID))
		request, Error := http.NewRequest("GET", requestUrl, nil)
		if Error != nil {
			ErrorLogger.Println("Failed To Initialize Request.")
		}

		request.Header.Set("Access-Control-Allow-Origin", "*")
		Response, Error := client.Do(request)

		SerializedPaymentProfileInfo, Error := io.ReadAll(Response.Body)
		json.Unmarshal(SerializedPaymentProfileInfo, &CustomerPaymentBalance)

		ResponseChannel <- CustomerBalance(CustomerPaymentBalance.Balance)
		// returning Eventual Value... to the Channel
	}
}

func GetProductsCatalog(context *gin.Context) {

	var AnnotatedProducts []struct {
		Product   models.Product
		Available bool
	}

	jwtToken, _ := context.Request.Cookie("jwt-token")
	ParsedJwtCredentials, JwtError := authentication.GetCustomerJwtCredentials(jwtToken.String())
	if JwtError != nil {
		DebugLogger.Println("Customer Is not Authenticated.")
	}

	ResponseChannel := make(chan CustomerBalance, 10)                 // Creating A Channel to let goroutine to send Response.
	go CustomerBalanceReceiver(ResponseChannel, ParsedJwtCredentials[ // Running goroutine that returns
	"Username"], ParsedJwtCredentials["Email"])                       // the Customer Balance..

	models.Database.Table("products").Find(&products) // Receiving Products from DB.

	var serializedProducts []byte // Once the Response is returning,
	//  the Products Serialized Context is getting put in this var
	// after that defering response

	// Method that returns Eventual HTTP Response to the Client...
	defer func(products []byte) {
		context.JSON(http.StatusOK, gin.H{"products": products})
	}(serializedProducts)

	select {
	case <-ResponseChannel:

		Balance := <-ResponseChannel

		if len(string(Balance)) != 0 { // If Balance Is Not None,
			// Making Product Annotations.. and returns Products Query.

			group := sync.WaitGroup{}
			convertedBalance, _ := strconv.ParseFloat(string(Balance), 5)

			go func() {
				group.Add(1)
				for _, row := range products {
					IsAvailable := row.ProductPrice > convertedBalance
					AnnotatedProducts = append(AnnotatedProducts, struct {
						Product   models.Product
						Available bool
					}{Product: row, Available: IsAvailable})
				}
				group.Done()
			}()
			group.Wait()

			serializedProducts, _ = json.Marshal(products)
		}
	default:
		serializedProducts, _ = json.Marshal(products)
	}
	defer close(ResponseChannel) // Closing Channel
}

func GetProduct(context *gin.Context) {

	productId := context.Query("productId")
	var product *models.Product

	models.Database.Table("products").Where(
		"id = ?", productId).First(&product)

	serializedProduct, _ := json.Marshal(product)
	context.JSON(http.StatusOK, gin.H{"product": serializedProduct})
}
