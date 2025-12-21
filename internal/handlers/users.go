package handlers

import (
	"chat/internal/models"
	"chat/internal/repository"
	"chat/internal/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signup(ctx *gin.Context) {
	var req models.AuthRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data."})
		return
	}

	dbAny, ok := ctx.Get("db")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database not available"})
		return
	}

	db := dbAny.(*sql.DB)

	hashedPassword, err := utils.HashPassword(req.Password) // hash the password before saving

	user := &models.User{ 
		Email: req.Email,
		PasswordHash: hashedPassword,
	}

	err = repository.CreateUser(db, user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create new user."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully.", "user_id": user.ID})
}

func login(ctx *gin.Context) {
	var req models.AuthRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data."})
		return
	}

	dbAny, ok := ctx.Get("db")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Database not available"})
		return
	}

	db := dbAny.(*sql.DB)

	existingUser, err := repository.GetUserByEmail(db, req.Email)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "User with email and password not found."})
		return
	}

	validPassword := utils.CheckPasswordHash(req.Password, existingUser.PasswordHash)

	if !validPassword {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Either email or password are incorrect."})
		return
	}

	token, err := utils.GenerateToken(existingUser.Email, existingUser.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful.", "token": token})
}

