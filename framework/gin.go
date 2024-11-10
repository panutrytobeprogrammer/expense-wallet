package framework

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Ginapp() *gin.Engine {
	app := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://blaze-dev.panuwat.tech"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"X-Requested-With", "Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	app.Use(cors.New(config))
	return app
}
