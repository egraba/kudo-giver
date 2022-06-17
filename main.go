package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type person struct {
	ID 			int64 	`json:"id"`
	FirstName 	string	`json:"firstName"`
}

var persons = []person {
	{ID: 1, FirstName: "Eric"},
	{ID: 2, FirstName: "Yadi"},
}

func getPersons(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, persons)
}

func createPerson(c *gin.Context) {
	var newPerson person

	if err := c.BindJSON(&newPerson); err != nil {
		return
	}

	persons = append(persons, newPerson)
	c.IndentedJSON(http.StatusCreated, newPerson)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/persons", getPersons)
	r.POST("/persons", createPerson)
	return r
}

func main() {
	router := setupRouter()
    router.Run("localhost:8080")
}