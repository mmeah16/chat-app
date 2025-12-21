package handlers

import (
	"github.com/gin-gonic/gin"
)

// referenced functions are accessible because they are in the same package
func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)
}