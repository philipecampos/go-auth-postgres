package server

import (
	"fmt"
	"go-auth-postgres/internal/auth"
	"go-auth-postgres/internal/models"
	"go-auth-postgres/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepository repositories.UsersRepositoryInterface
}

func NewAuthHandler(userRepository repositories.UsersRepositoryInterface) *AuthHandler {
	return &AuthHandler{
		userRepository: userRepository,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the user already exists
	user, err := h.userRepository.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// create user
	user, err = models.NewUser(req.Username, req.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error creating user: %s", err.Error())})
		return
	}

	// save the user to the database
	err = h.userRepository.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userRepository.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid credentials"})
		return
	}

	// hash password & compare with the db password hash
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		15*60, // 15 min
		"/",
		"localhost",
		false,
		true, // não acessível via js
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
