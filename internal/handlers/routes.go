package handlers

import (
	"chat/internal/middleware"
	"chat/internal/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

// referenced functions are accessible because they are in the same package
func RegisterRoutes(server *gin.Engine, db *sql.DB) {

	conversationService := service.NewConversationService(db)

	server.POST("/signup", signup)
	server.POST("/login", login)

	
	authenticated := server.Group("/")
	authenticated.Use(middleware.Authenticate)
	authenticated.GET("/conversations", getConversations)
	authenticated.GET("/conversations/:conversation_id/messages", getConversationMessages)
	authenticated.POST("/conversations/:conversation_id/messages", createConversationMessage(conversationService))
	authenticated.POST("/conversations/dm", createDMConversation(conversationService))
}