package main

import (
	"servicedependency/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Entrypoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "server is running",
	})
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // The React Port
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Cache-Control"}, // ✅ Added more headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", Entrypoint)
	r.GET("/service", controllers.JsonMarshalling)

	r.Run()
}
