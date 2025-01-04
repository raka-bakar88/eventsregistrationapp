package middlewares

import (
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Authenticate Middleware function that will be executed before executing other functions
// It protects some APIs, so only authenticated user can execute ther after-APIs
func Authenticate(context *gin.Context) {
	// get token from header
	token := context.Request.Header.Get("Authorization")
	// when token is empty
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is empty"})
		return
	}

	// verify the validity of the token
	userId, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authorized"})
		return
	}
	// userid is set to the context so it can be extracted anytime
	context.Set("userId", userId)
	// execute the next functions
	context.Next()
}
