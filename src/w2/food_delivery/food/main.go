package main

import (
	ioutil "github.com/sanjay-vasudeva/ioutil"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db := ioutil.NewConn("3306", "root", "password", "delivery")
	r.POST("/food/reserve", func(c *gin.Context) {
	})

	r.POST("/food/book", func(c *gin.Context) {
	})

	r.Run(":8082")
}
