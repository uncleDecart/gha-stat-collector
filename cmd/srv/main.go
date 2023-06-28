package main

import (
	"github.com/gin-gonic/gin"
	"github.com/uncleDecart/gha-stat-collector/controllers"
	"github.com/uncleDecart/gha-stat-collector/middlewares"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", controllers.Ping)

	protected := r.Group("/api/v1")
	protected.Use(middlewares.TokenAuthMiddleware())
	protected.POST("/publish", controllers.Publish)

	return r
}

func main() {
	r := setupRouter()
	r.Run()
}
