package polling

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type EC2 struct {
	ID     int
	Status string
}

var db sync.Map

func CreateEC2(id int) {
	fmt.Println("creating EC2 with ID:", id)
	db.Store(id, "todo")
	time.Sleep(10 * time.Second) // Simulate some delay in creation
	db.Store(id, "in-progress")
	fmt.Println("EC2 creation in progress")
	// Simulate some processing time
	time.Sleep(10 * time.Second) // Simulate some delay in processing
	db.Store(id, "done")
	fmt.Println("EC2 creation done")
}

func SetupServer(r *gin.Engine) {
	r.POST("/EC2", func(c *gin.Context) {
		id := rand.IntN(100)
		go CreateEC2(id)
		c.JSON(http.StatusOK, gin.H{
			"id":     id,
			"status": "EC2 creation started"})
	})
	r.GET("/EC2/status/short", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		status, ok := db.Load(id)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "EC2 not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"ID":     id})

	})

	r.GET("/EC2/status/long", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		status := c.Query("status")
		if status == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Status missing in query param",
			})
			return
		}
		for {
			curr_status, ok := db.Load(id)
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "EC2 not found"})
				return
			}
			if curr_status == status {
				// Wait for 1 second before checking again
				time.Sleep(1 * time.Second)
				continue
			}
			c.JSON(http.StatusOK, gin.H{
				"status": curr_status,
				"ID":     id})
			return
		}
	})
}
