package server

import (
	"go-auth-postgres/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepository repositories.UsersRepositoryInterface
}

func newUserHandler(userRepository repositories.UsersRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
	}

	user, err := h.userRepository.FindByID(c.Request.Context(), userID.(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}
