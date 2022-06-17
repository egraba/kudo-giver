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

type kudo struct {
	ID 			int64 	`json:"id"`
	SenderID 	int64 	`json:"senderId"`
	ReceiverID 	int64 	`json:"receiverId"`
	Message		string	`json:"message"`
}

var kudos = []kudo {
	{ID: 1, SenderID: 1, ReceiverID: 2, Message: "Déjà"},
	{ID: 2, SenderID: 2, ReceiverID: 1, Message: "Mira!"},
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

func getKudos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, kudos)
}

func giveKudo(c *gin.Context) {
	var newKudo kudo

	if err := c.BindJSON(&newKudo); err != nil {
		return
	}

	kudos = append(kudos, newKudo)
	c.IndentedJSON(http.StatusCreated, newKudo)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	
	r.GET("/persons", getPersons)
	r.POST("/persons", createPerson)
	
	r.GET("/kudos", getKudos)
	r.POST("/kudos", giveKudo)

	return r
}

func main() {
	router := setupRouter()	
    router.Run("localhost:8080")
}
