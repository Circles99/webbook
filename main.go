package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好")
	})

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
