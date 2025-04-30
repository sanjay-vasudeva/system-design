package order

import (
	ioutil "github.com/sanjay-vasudeva/ioutil"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	ioutil.NewConn("3306", "root", "password", "delivery")
	r.POST("/order", func(c *gin.Context) {

		// reserve food
		// reserve agent

		// commit food
		// commit agent

		// place order
	})

	r.Run(":8080")
}
