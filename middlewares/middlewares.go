package middlewares

import (
	"fmt"
	"strings"

	"log"

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

		if hasAuthHeader := context.Request.Header.Get("AUTHORIZATION"); len(hasAuthHeader) == 0 {
			AuthCookie, error := context.Request.Cookie("jwt-token")
			if error != nil || AuthCookie == nil {
				context.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AuthCookie))
			}
		}
		context.Next()
	}
}

func JwtAuthenticationMiddleware() gin.HandlerFunc {

	return func(context *gin.Context) {

		jwtToken := context.Request.Header.Get("AUTHORIZATION")
		if len(jwtToken) == 0 {
			return
		}

		if validationError := authentication.CheckValidJwtToken(jwtToken); validationError != nil {
			return
		}
		context.Next()
	}
}

// Product Middlewares....

func IsProductOwnerMiddleware() gin.HandlerFunc {

	return func(context *gin.Context) {

		var product models.Product
		productId := context.Query("productId")

		token, CookieError := context.Request.Cookie("jwt-token")
		jwtParams, jwtError := authentication.GetCustomerJwtCredentials(token.String()) // customer that makes request...

		if jwtError != nil || CookieError != nil {
			DebugLogger.Println(
				"Exceptions: " + CookieError.Error() + fmt.Sprintf("%s", jwtError))
			return
		}
		models.Database.Table(
			"products").Where("id = ?", productId).First(&product)

		if isOwner := strings.ToLower(
			product.OwnerEmail); isOwner == strings.ToLower(jwtParams["email"]) { // if customer is not an owner...
			return
		}
		context.Next()
	}
}
