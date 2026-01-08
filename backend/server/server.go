package server

import (
	"fmt"
	"go-auth-postgres/internal/database"
	"go-auth-postgres/internal/middlewares"
	"go-auth-postgres/internal/repositories"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewServer() *http.Server {
	port := os.Getenv("PORT")

	db := database.New()
	userRepository := repositories.NewUserRepository(db.Database)

	server := &Server{
		usersRepository: userRepository,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      server.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	authHandler := NewAuthHandler(s.usersRepository)
	userHandler := newUserHandler(s.usersRepository)

	// Setup routes
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})
		authRoutes.POST("/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})
	}

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware(s.usersRepository))
	{
		protected.GET("/user", userHandler.GetUser)
	}

	return r
}
