package middlewares

import (
	"fmt"
	"strings"

	"github.com/LovePelmeni/OnlineStore/StoreService/authentication"
	"github.com/LovePelmeni/OnlineStore/StoreService/models"
	"github.com/gin-gonic/gin"
)

func SetAuthHeaderMiddleware(context *gin.Context) gin.HandlerFunc {

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

func JwtAuthenticationMiddleware(context *gin.Context) gin.HandlerFunc {

	return func(context *gin.Context) {

		jwtToken := context.Request.Header.Get("AUTHORIZATION")
		if len(jwtToken) == 0 {
			return
		}

		if valid := authentication.CheckValidJwtToken(jwtToken); valid != nil {
			return
		}
		context.Next()
	}
}

// Product Middlewares....

func IsProductOwnerMiddleware(context *gin.Context) gin.HandlerFunc {

	return func(context *gin.Context) {

		var product models.Product
		productId := context.Query("productId")
		requestCustomerId := context.Query("customerId") // customer that makes request...

		models.Database.Table(
			"products").Where("id = ?", productId).First(&product)

		if isOwner := strings.ToLower(product.Owner.Username); isOwner == strings.ToLower(requestCustomerId) { // if customer is not an owner...
			return
		}
		context.Next()
	}
}
