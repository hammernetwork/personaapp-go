package server

import (
	"log"
	"persona/middlewares"
	"persona/models"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	authMiddleware, err := middlewares.AuthMiddleware([]models.Role{models.EmployeeRole})
	if err != nil {
		log.Fatalf("Auth middleware error: %v", err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/register", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"success": "true",
		})
	})

	// TODO: handle different roles endpoints
	authGroup := r.Group("/auth")
	authGroup.GET("refresh_token", authMiddleware.RefreshHandler)
	authGroup.Use(authMiddleware.MiddlewareFunc())
	{
		// add authorized only routes
	}

	return r
}
