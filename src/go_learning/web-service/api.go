package main

import (
	"github.com/gin-gonic/gin"
)

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int16  `json:"age"`
}

var people []Person = []Person{
	{ID: "1", Name: "Sanjay", Age: 28},
	{ID: "2", Name: "Aarthi", Age: 25},
}

func main() {
	r := gin.Default()
	r.GET("/people", GetPeople)
	r.GET("/people/:id", GetPerson)
	r.POST("/people", CreatePerson)

	r.Run(":8080")
}

func GetPeople(c *gin.Context) {
	c.JSON(200, people)
}

func GetPerson(c *gin.Context) {
	id := c.Param("id")
	for _, person := range people {
		if person.ID == id {
			c.JSON(200, person)
			return
		}
	}
	c.AbortWithStatus(404)
}

func CreatePerson(c *gin.Context) {
	var p Person
	if err := c.Bind(&p); err != nil {
		c.AbortWithStatus(400)
		return
	}
	people = append(people, p)
	c.JSON(201, p)
}
