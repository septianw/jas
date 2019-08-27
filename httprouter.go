package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode("release")

	switch STAGE {
	case "development":
		gin.SetMode("debug")
		break
	case "production":
		gin.SetMode("release")
		break
	case "testing":
		gin.SetMode("debug")
		break
	default:
		gin.SetMode("debug")
	}

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "pong")
	// })

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	return r
}
