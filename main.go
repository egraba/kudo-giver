package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

const dbUrl = "postgres://devuser:devpwd@localhost:5432/kudo-giver"

const PersonEndpoint = "/persons"
const KudoEndpoint = "/kudos"

type Person struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
}

var persons = []Person{
	{ID: 1, FirstName: "Eric"},
	{ID: 2, FirstName: "Yadi"},
}

type Kudo struct {
	ID         int64  `json:"id"`
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	Message    string `json:"message"`
}

var kudos = []Kudo{
	{ID: 1, SenderID: 1, ReceiverID: 2, Message: "Déjà"},
	{ID: 2, SenderID: 2, ReceiverID: 1, Message: "Mira!"},
}

func getPersons(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, persons)
}

func createPerson(c *gin.Context) {
	var person Person

	if err := c.BindJSON(&person); err != nil {
		return
	}

	persons = append(persons, person)
	c.IndentedJSON(http.StatusCreated, person)
}

func getKudos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, kudos)
}

func giveKudo(c *gin.Context) {
	var kudo Kudo

	if err := c.BindJSON(&kudo); err != nil {
		log.Fatal(err)
	}

	kudos = append(kudos, kudo)
	c.IndentedJSON(http.StatusCreated, kudo)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET(PersonEndpoint, getPersons)
	r.POST(PersonEndpoint, createPerson)

	r.GET(KudoEndpoint, getKudos)
	r.POST(KudoEndpoint, giveKudo)

	return r
}

func main() {
	log.SetPrefix("[kudo-giver] ")
	log.SetFlags(7)

	dbpool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sucessfully connected to the database!")
	defer dbpool.Close()

	router := setupRouter()
	router.Run("localhost:8080")
}
