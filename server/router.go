package server

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	// setup router
	return gin.New()
}
