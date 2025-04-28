package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		key := ctx.Query("key")
		//TO DO
	})

	r.POST("/", func(c *gin.Context) {
		key := c.Query("key")
		value := c.Query("value")

		//TO DO
	})

	r.DELETE("/", func(c *gin.Context) {
		key := c.Query("key")

		//TO DO
	})
}
