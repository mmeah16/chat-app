package handlers

import (
	"chat/internal/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func bindAuth(ctx *gin.Context) (*models.AuthRequest, bool) {
	var req models.AuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data."})
		zap.S().Warn("invalid auth request", "request_id", ctx.GetString("request_id"))
		return nil, false
	}
	return &req, true
}

func getDB(ctx *gin.Context) (*sql.DB, bool) {
	dbAny, ok := ctx.Get("db")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database not available"})
		zap.S().Error("db error during login",
			"request_id", ctx.GetString("request_id"),
		)
		return nil, false 
	}
	return dbAny.(*sql.DB), true
}

func authFail(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials."})
	zap.S().Warn("auth failed",
		"reason", reason,
		"request_id", ctx.GetString("request_id"),
	)
}