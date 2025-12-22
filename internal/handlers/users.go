package handlers

import (
	"chat/internal/models"
	"chat/internal/repository"
	"chat/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func signup(ctx *gin.Context) {
	req, ok := bindAuth(ctx)
	if !ok { return }

	db, ok := getDB(ctx)
	if !ok { return }

	hashedPassword, err := utils.HashPassword(req.Password) // hash the password before saving
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create new user."})
		zap.S().Error("hashing error",
			"error", err,
			"request_id", ctx.GetString("request_id"),
		)
		return
	}

	user := &models.User{ 
		Email: req.Email,
		PasswordHash: hashedPassword,
	}

	err = repository.CreateUser(db, user)
	if err != nil {
		if utils.IsUniqueViolation(err) {
			ctx.JSON(http.StatusConflict, gin.H{"message": "User already exists."})
			zap.S().Warn("signup failed",
				"reason", "user_exists",
				"request_id", ctx.GetString("request_id"),
			)
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create new user."})
		zap.S().Error("db error during signup",
			"error", err,
			"request_id", ctx.GetString("request_id"),
		)
		return
	}


	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}

func login(ctx *gin.Context) {
	req, ok := bindAuth(ctx)
	if !ok { return }

	db, ok := getDB(ctx)
	if !ok { return }

	user, err := repository.GetUserByEmail(db, req.Email)
	if err != nil {
		authFail(ctx, "user_not_found")
		ctx.Set("auth_result", "failure")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		authFail(ctx, "invalid_password")
		ctx.Set("auth_result", "failure")
		return
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		zap.S().Warn("could not generate token",
			"reason", "user_not_found",
			"request_id", ctx.GetString("request_id"),
		)
		ctx.Set("auth_result", "failure")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful.", "token": token})
	ctx.Set("auth_result", "success")
}

