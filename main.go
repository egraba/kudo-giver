package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

const dbUrl = "postgres://devuser:devpwd@localhost:5432/kudo-giver"

const PersonEndpoint = "/persons"
const KudoEndpoint = "/kudos"

type Person struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate time.Time `json:"birthDate"`
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
	dbPool := c.MustGet("dbConnection").(*pgxpool.Pool)

	rows, err := dbPool.Query(context.Background(), "SELECT * FROM persons ORDER BY first_name;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}

		var person = Person{
			ID:        values[0].(int64),
			FirstName: values[1].(string),
			LastName:  values[2].(string),
			BirthDate: values[3].(time.Time),
		}

		persons = append(persons, person)
	}

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

func connectDb(dbPool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbConnection", dbPool)
		c.Next()
	}
}

func setupRouter(dbPool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(connectDb(dbPool))

	r.GET(PersonEndpoint, getPersons)
	r.POST(PersonEndpoint, createPerson)

	r.GET(KudoEndpoint, getKudos)
	r.POST(KudoEndpoint, giveKudo)

	return r
}

func main() {
	log.SetPrefix("[kudo-giver] ")
	log.SetFlags(7)

	dbPool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	_, err = dbPool.Query(context.Background(), `CREATE TABLE persons (
														id INTEGER NOT NULL PRIMARY KEY,
														first_name VARCHAR(32) NOT NULL,
														last_name VARCHAR(32) NOT NULL,
														birth_date DATE NOT NULL
													);`)
	if err != nil {
		log.Fatal(err)
	}

	router := setupRouter(dbPool)
	router.Run("localhost:8080")
}
