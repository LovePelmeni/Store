package middlewares

import (
	"fmt"
	"strings"

	"log"
	"net/http"

	"github.com/LovePelmeni/OnlineStore/StoreService/authentication"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/gin-gonic/gin"
)

// Loggers

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
)

func SetAuthHeaderMiddleware() gin.HandlerFunc {

	return func(context *gin.Context) {

		switch {
		case len(context.Request.Header.Get("AUTHORIZATION")) == 0:
			AuthCookie, error := context.Request.Cookie("jwt-token")
			if error != nil || AuthCookie == nil {
				context.Request.Header.Set("AUTHORIZATION", fmt.Sprintf("Bearer %s", AuthCookie.String()))
			}
			context.Next()
		default:
			context.Next()
		}
	}
}

func JwtAuthenticationMiddleware() gin.HandlerFunc {

	return func(context *gin.Context) {

		jwtToken := context.Request.Header.Get("AUTHORIZATION")
		switch {
		case len(jwtToken) == 0:
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Not Authorized."})
			return
		case authentication.CheckValidJwtToken(jwtToken) != nil:
			context.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{"Error": "Not Authorized."})
			return
		default:
			context.Next()
		}
	}
}

// Product Middlewares....

func IsProductOwnerMiddleware() gin.HandlerFunc {

	return gin.HandlerFunc(func(context *gin.Context) {

		var product models.Product
		productId := context.Query("productId")

		token, CookieError := context.Request.Cookie("jwt-token")
		jwtParams, jwtError := authentication.GetCustomerJwtCredentials(token.String()) // customer that makes request...
		models.Database.Table(
			"products").Where("id = ?", productId).First(&product)

		switch {
		case jwtError != nil || CookieError != nil:
			DebugLogger.Println(
				"Exceptions: " + CookieError.Error() + fmt.Sprintf("%s", jwtError))
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Does Not Have Enough Permissions to remove the  Product."})
			return

		case strings.ToLower(
			product.OwnerEmail) != strings.ToLower(jwtParams["email"]): // if customer is not an owner...
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Error": "Does not have enough permissions to Remove The Product."})
			return

		default:
			context.Next()
		}
	})
}
