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

func main() {
	router := gin.Default()
	router.GET("/persons", getPersons)

    router.Run("localhost:8080")
}