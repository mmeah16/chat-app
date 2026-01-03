package middleware

import (
	"net/http"
	"strings"

	"chat/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Authenticate(ctx *gin.Context) {
	// Middleware logic to authenticate requests
	token := ctx.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization token is required."})
		return
	}

	userIdStr, err := utils.VerifyToken(tokenString)
	userUUID, err := uuid.Parse(userIdStr)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization token."})
		return
	}
	
	// Allows us to add data to the context for handlers to use
	ctx.Set("user_id", userUUID)
	// Ensures next event handler can execute correctly
	ctx.Next()
}