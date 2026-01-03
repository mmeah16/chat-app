package handlers

import (
	"chat/internal/models"
	"chat/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func sendMessage(ctx *gin.Context) {
	senderIDStr := ctx.GetString("user_id")
	senderID, err := uuid.Parse(senderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID."})
		return
	}

	conversationIDStr := ctx.Param("conversation_id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid conversation ID."})
		return
	}

	db, ok := getDB(ctx)
	if !ok { return }

	var req models.MessageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data."})
		zap.S().Warn("invalid message request", "request_id", ctx.GetString("request_id"))
		return
	}

	message := &models.Message{
		ConversationID: conversationID,
		SenderID: senderID,
		Body: req.Body,
	}

	messageDto, err := repository.CreateMessage(ctx.Request.Context(), db, message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create new message."})
		zap.S().Error("error during sendMessage",
			"error", err,
			"request_id", ctx.GetString("request_id"),
		)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Message sent successfully.", "response": messageDto})
}

func getMessagesByConversation(ctx *gin.Context) {

	conversationIDStr := ctx.Param("conversation_id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid conversation ID."})
		return
	}

	db, ok := getDB(ctx)
	if !ok { return }

	messages, err := repository.GetMessagesByConversation(ctx.Request.Context(), db, conversationID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not retrieve messages."})
		zap.S().Warn("could not complete getMessagesByConversation request", "request_id", ctx.GetString("request_id"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Message sent successfully.", "messages": messages})
}