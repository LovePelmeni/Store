package middlewares 

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func JwtAuthenticationMiddleware(context *gin.Context) *gin.HandlerFunc{
	jwtToken, error := context.Request.Cookie("")
	if error != nil || jwtToken == nil {return } 
	if strings.HasPrefix(jwtToken, "Bearer") || !len(strings.Split(jwtToken, " ")) == 2 {return }
	decodedData := 
}