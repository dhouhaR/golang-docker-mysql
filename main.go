package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang-docker-todo/api"
)

func main() {
	r := gin.New()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	//r.Run()
	api.Run()
}