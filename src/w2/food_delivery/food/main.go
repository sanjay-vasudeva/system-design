package main

import (
	ioutil "github.com/sanjay-vasudeva/ioutil"

	"food/src"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db := ioutil.NewConn("3308", "root", "password", "delivery")
	r.POST("/food/reserve", func(c *gin.Context) {
		stock, err := src.Reserve(db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, stock)
	})

	r.POST("/food/book", func(c *gin.Context) {
		orderID := c.Query("order_id")
		if orderID == "" {
			c.JSON(400, gin.H{"error": "order_id is required"})
			return
		}

		stock, err := src.Book(orderID, db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, stock)
	})

	r.Run(":8082")
}
