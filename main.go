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

func getPersons(c *gin.Context) {
	var persons = []Person{}
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
			Email:     values[3].(string),
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

	c.Status(http.StatusCreated)
}

func connectDb(dbPool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbConnection", dbPool)
		c.Next()
	}
}

func main() {
	log.SetPrefix("[kudo-giver] ")
	log.SetFlags(7)

	dbPool, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	_, err = dbPool.Exec(context.Background(), ReadSqlFile("sql/create_persons_table.sql"))
	if err != nil {
		log.Println(err)
	}

	router := SetupRouter(dbPool)
	router.Run("localhost:8080")
}
