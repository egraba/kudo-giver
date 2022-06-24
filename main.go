package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

const dbUrl = "postgres://devuser:devpwd@localhost:5432/kudo-giver"

const PersonEndpoint = "/persons"
const KudoEndpoint = "/kudos"

type Person struct {
	ID        int32  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

var persons = []Person{}

type Kudo struct {
	ID         int64  `json:"id"`
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	Message    string `json:"message"`
}

var kudos = []Kudo{}

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
			ID:        values[0].(int32),
			FirstName: values[1].(string),
			LastName:  values[2].(string),
		}

		persons = append(persons, person)
	}

	c.IndentedJSON(http.StatusOK, persons)
}

func createPerson(c *gin.Context) {
	var person Person

	if err := c.BindJSON(&person); err != nil {
		log.Println(err)
		return
	}

	dbPool := c.MustGet("dbConnection").(*pgxpool.Pool)
	values := fmt.Sprintf("VALUES ('%s', '%s', '%s');", person.FirstName, person.LastName, person.Email)
	sqlStr := "INSERT INTO persons (first_name, last_name, email)" + values
	log.Println(sqlStr)
	_, err := dbPool.Exec(context.Background(), sqlStr)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, err)
	}

	persons = append(persons, person)
	c.Status(http.StatusCreated)
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

	_, err = dbPool.Exec(context.Background(), `CREATE TABLE persons (
													id SERIAL PRIMARY KEY NOT NULL,
													first_name VARCHAR(32) NOT NULL,
													last_name VARCHAR(32) NOT NULL,
													email VARCHAR(254) UNIQUE NOT NULL
												);`)
	if err != nil {
		log.Fatal(err)
	}

	router := setupRouter(dbPool)
	router.Run("localhost:8080")
}
