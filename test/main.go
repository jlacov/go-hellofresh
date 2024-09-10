package main

import (
	"go-hellofresh/test/controllers"
	"go-hellofresh/test/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	r := gin.Default()
	r.GET("/stats", controllers.GetStats)
	r.POST("/event", controllers.PostEvents)
	r.Run()
}
