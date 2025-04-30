package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sanjay-vasudeva/ioutil"
)

func main() {
	r := gin.Default()
	db := ioutil.NewConn("3306", "root", "password", "delivery")

	r.POST("/agent/reserve", func(c *gin.Context) {

	})

	r.POST("/agent/book", func(c *gin.Context) {
	})

	r.Run(":8081")
}
