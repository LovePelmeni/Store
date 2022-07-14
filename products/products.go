package products

import (
	"log"
	"net/http"

	"github.com/LovePelmeni/StoreService/models"
	"github.com/gin-gonic/gin"
)

var (
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
)

// CUD Rest Controllers...

func CreateProduct(context *gin.Context) {}

func UpdateProduct(context *gin.Context) {}

func DeleteProduct(context *gin.Context) {}

// Getter Rest Controllers...

func GetTopWeekProducts(context *gin.Context) {}

func GetProductsCatalog(context *gin.Context) {}

func GetMostLikedProducts(context *gin.Context) {

	var ProductsMergedQuery struct {
		Product    models.Product `json:"Product"`
		LikedUsers struct {
			Username string `json:"Username"`
			Email    string `json:"Email"`
		}
	}

	var Products []models.Product
	var ResponseQuery []ProductsMergedQuery

	products, RowsError := models.Database.Table("products").Order(
		"COUNT(liked_users) desc").Find(&Products).Limit(10).Rows()

	if RowsError != nil {
		DebugLogger.Println("Failed to Convert Query to Rows..")
	}
	defer products.Close()

	for products.Next() {

		var AnnotatedProduct models.Product

		AnnotatedLikedCustomers := products.Association("LikedUsers").Select("Username, Email").Find(
			&ProductsMergedQuery.LikedUsers)

		models.Database.ScanRows(products, &AnnotatedProduct)

		mergedProductInfo := ProductsMergedQuery{
			Product:    AnnotatedProduct,
			LikedUsers: AnnotatedLikedCustomers,
		}
		ResponseQuery = append(ResponseQuery, mergedProductInfo)
	}
	context.JSON(http.StatusOK, gin.H{"query": ResponseQuery})
}

func GetProduct(context *gin.Context) {

}
