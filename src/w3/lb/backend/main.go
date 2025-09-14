package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()
	port := flag.Arg(0)
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		// time.Sleep(10 * time.Millisecond) // Simulate some processing delay
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})
	router.Run(fmt.Sprintf(":%s", port))
}
