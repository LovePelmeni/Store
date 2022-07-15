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

func UpdateProduct(context *gin.Context) {}

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

func GetTopWeekProducts(context *gin.Context) {}

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

	// if customer != nil {
	// 	CustomerBalance, ParserError := strconv.ParseFloat("10000000.00", 5)
	// 	if ParserError != nil {
	// 		DebugLogger.Println("Invalid Customer Balance Format.")
	// 	}
	// }

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

func GetMostLikedProducts(context *gin.Context) {
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
