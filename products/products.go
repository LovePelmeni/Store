package products

import (
	"log"
	"net/http"

	"encoding/json"
	"errors"

	"fmt"
	"reflect"
	"sync"

	"github.com/LovePelmeni/OnlineStore/StoreService/authentication"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
)

var (
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
)

var product models.Product

func InitializeLoggers() {}

func init() {}

// CUD Rest Controllers...

func CreateProduct(context *gin.Context) {

	InvalidFieldsErrors := []string{}
	ModelPostValues := reflect.ValueOf(context.PostForm)
	ModelPostNames := reflect.TypeOf(context.PostForm)

	group := sync.WaitGroup{}

	// Validating
	go func() {
		group.Add(1)

		for PropertyIndex := 1; PropertyIndex < reflect.TypeOf(context.PostForm).NumField(); PropertyIndex++ {
			if Valid := ModelPostValues.Field(PropertyIndex).IsValid; Valid() != false {
				if len(ModelPostValues.Field(PropertyIndex).String()) == 0 {
					InvalidFieldsErrors = append(InvalidFieldsErrors, fmt.Sprintf("Invalid Value for Field `%s`",
						ModelPostNames.Field(PropertyIndex).Name))
				}
			} else {
				continue
			}
		}
		group.Done()
	}()
	group.Wait()

	if len(InvalidFieldsErrors) != 0 {
		context.JSON(http.StatusBadRequest, gin.H{"errors": InvalidFieldsErrors})
	}

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
	if newProductCreated == true || len(Errors) != 0 {
		context.JSON(http.StatusOK, gin.H{"Errors": Errors})
	}
}

func UpdateProduct(context *gin.Context) {

	productId := context.Query("productId")
	bodyData, Error := io.ReadAll(context.Request.Body)
	var InvalidFields []error 

	if Error != nil {InfoLogger.Println(
	"Failed to parse Product Form Data."); context.JSON(http.StatusBadRequest, nil)}

	Product := models.Database.Table("products").Where("id = ?", productId) // Receives the object..
	if errors.Is(Product.Error, gorm.ErrRecordNotFound) {context.JSON(http.StatusNotFound, nil)} 


	group := sync.WaitGroup{}

	go func() {

		group.Add(1)
		
			structuredFields := reflect.TypeOf(&bodyData)
			for PropertyIndex := 1; PropertyIndex < structuredFields.NumField(); PropertyIndex ++ {

				if len(reflect.ValueOf(structuredFields.Field(
				PropertyIndex)).String()) == 0 { // Checking if the Field Value is empty...

				InvalidFields = append(InvalidFields, errors.New(fmt.Sprintf("Invalid Value for Field `%s`",
				structuredFields.Field(PropertyIndex).Name)))} 

		group.Done()
	}}()

	group.Wait()

	if len(InvalidFields) != 0 {

	serializedContext, Error := json.Marshal(InvalidFields); _ = Error
	context.JSON(http.StatusBadRequest, gin.H{"InvalidFields": serializedContext})}
	
	context.JSON(http.StatusCreated, nil)
}


func DeleteProduct(context *gin.Context) {

	ProductId := context.Query("ProductId")
	Deleted, Errors := product.DeleteObject(ProductId)
	if Deleted != true || len(Errors) != 0 {
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


func GetProductsCatalog(context *gin.Context) {

	var products []models.Product
	var customer models.Customer

	var CustomerBalance float64
	var AnnotatedProducts []struct {
		Product   models.Product
		Available bool
	}

	jwtToken, error := context.Request.Cookie("jwt-token")
	if error != nil {
	}

	ParsedJwtCredentials, JwtError := authentication.GetCustomerJwtCredentials(jwtToken.String())
	if JwtError != nil {
		DebugLogger.Println("Customer Is not Authenticated.")
	}

	Customer := models.Database.Table("customers").Where( // Receiving Customer.
		"Email = ? AND username = ?",
		ParsedJwtCredentials["email"],
		ParsedJwtCredentials["username"]).First(&customer)

	_ = Customer


	group := sync.WaitGroup{}

	var CustomerPaymentProfileData map[string]string

	go func(){
		group.Add(1)

		client := http.Client{}
		requestUrl := url.URL(fmt.Sprintf("http://%s:%s/get/customer/info/",
	    os.Getenv("PAYMENT_APPLICATION_HOST")))
		requestUrl.Query("CustomerId") = string(customer.ID)
		request, Error := http.NewRequest("GET", requestUrl)
		request.Header.Set("Access-Control-Allow-Origin", "*")
		Response, Error := client.Do(request) 

		SerializedPaymentProfileInfo, Error := io.ReadAll(Response.Body)
		DecodeError := json.Unmarshal(SerializedPaymentProfileInfo, &CustomerPaymentProfileData)
		_ = DecodeError 

	}() // Parsing Customer Balance.. 

	Products := models.Database.Table("products").Find(&products)
	if Products.Error != nil {
		InfoLogger.Println("Failed to Parse Products Query.")
		context.JSON(http.StatusNotImplemented, nil)
	}

	for _, row := range products {

		IsAvailable := row.ProductPrice > CustomerBalance
		AnnotatedProducts = append(AnnotatedProducts, struct {
			Product   models.Product
			Available bool
		}{Product: row, Available: IsAvailable})
	}
	context.JSON(http.StatusOK, gin.H{"products": AnnotatedProducts})
}


func GetProduct(context *gin.Context) {

	productId := context.Query("productId")
	var product *models.Product

	Received := models.Database.Table("products").Where(
		"id = ?", productId).First(&product)
	if errors.Is(Received.Error, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product Not Found"})
	}
	serializedProduct, Error := json.Marshal(product)
	if Error != nil {
		ErrorLogger.Println("Failed to Serialize Product into JSON.")
		context.JSON(http.StatusNotImplemented, nil)
	}
	context.JSON(http.StatusOK, gin.H{"product": serializedProduct})
}
