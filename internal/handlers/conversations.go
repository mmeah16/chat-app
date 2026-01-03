package handlers

import (
	"chat/internal/models"
	"chat/internal/repository"
	"chat/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func getConversations(ctx *gin.Context) {
	db, ok := getDB(ctx)
	if !ok { return }

	userID, ok := ctx.Get("user_id")
	uid, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID."})
		return
	}

	conversations, err := repository.GetConversationsByUser(ctx.Request.Context(), db, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve conversations."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

func getConversationMessages(ctx *gin.Context) {
	db, ok := getDB(ctx)
	if !ok { return }

	conversationIDStr, ok := ctx.Params.Get("conversation_id")

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing conversation ID."})
		return
	}

	conversationID, _ := uuid.Parse(conversationIDStr)

	messages, err := repository.GetMessagesByConversation(ctx, db, conversationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve messages for specified conversation."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"messages": messages})
}

func createConversationMessage(svc *service.ConversationService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		conversationIDStr, ok := ctx.Params.Get("conversation_id")
		
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing conversation ID."})
			return
		}

		conversationID, err := uuid.Parse(conversationIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid conversation id"})
			zap.S().Error("failed to parse conversationID",
				"error", err.Error(),
				"request_id", ctx.GetString("request_id"),
			)
			return
		}

		senderID := ctx.MustGet("user_id").(uuid.UUID)
		zap.S().Info("senderId",
			"senderId ", senderID,
			"request_id ", ctx.GetString("request_id"),
		)

		var message *models.MessageRequest

		if err := ctx.ShouldBindJSON(&message); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			zap.S().Error("json binding error",
				"error", err.Error(),
				"request_id", ctx.GetString("request_id"),
			)
			return 
		}

		sentMessage, err := svc.CreateMessage(ctx.Request.Context(), senderID, conversationID, message)

		if err != nil {
			switch err {
			case service.ErrUserNotInConversation:
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "user not member of specified conversation"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to send message"})
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"message": sentMessage})

	}
}

func createDMConversation(svc *service.ConversationService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		senderID := ctx.MustGet("user_id").(uuid.UUID)
		zap.S().Info("senderId",
			"senderId ", senderID,
			"request_id ", ctx.GetString("request_id"),
		)

		var req struct {
			UserID string `json:"user_id" binding:"required,uuid"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			zap.S().Error("json binding error",
				"error", err.Error(),
				"request_id", ctx.GetString("request_id"),
			)
			return
		}
		
		// convert User ID into UUID
		receiverID, err := uuid.Parse(req.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid user_id"})
			zap.S().Error("failed to parse UserID",
				"error", err.Error(),
				"request_id", ctx.GetString("request_id"),
			)
			return
		}

		conv, err := svc.CreateOrGetDM(ctx.Request.Context(), senderID, receiverID)
		if err != nil {
			switch err {
			case service.ErrUserNotFound:
				ctx.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			case service.ErrSameUserDM:
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "cannot create dm with self"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create dm"})
			}
			zap.S().Error("create dm error",
				"error", err.Error(),
				"request_id", ctx.GetString("request_id"),
			)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"conversation": conv})
	}
}

