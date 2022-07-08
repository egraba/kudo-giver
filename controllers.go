package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func connectDb(dbPool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbConnection", dbPool)
		c.Next()
	}
}

func CreatePerson(c *gin.Context) {
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

func GetPersons(c *gin.Context) {
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

func GetPersonById(c *gin.Context) {
	var person Person
	dbPool := c.MustGet("dbConnection").(*pgxpool.Pool)

	personId := c.Param("id")
	sqlStr := fmt.Sprintf("SELECT * FROM persons WHERE id = %s;", personId)
	log.Println(sqlStr)
	rows, err := dbPool.Query(context.Background(), sqlStr)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}

		person = Person{
			ID:        values[0].(int32),
			FirstName: values[1].(string),
			LastName:  values[2].(string),
			Email:     values[3].(string),
		}
		c.IndentedJSON(http.StatusOK, person)
	}
	c.Status(http.StatusNotFound)
}

func GiveKudo(c *gin.Context) {
	var kudo Kudo

	if err := c.BindJSON(&kudo); err != nil {
		log.Println(err)
		return
	}

	if kudo.SenderID == kudo.ReceiverID {
		c.Status(http.StatusForbidden)
		return
	}

	dbPool := c.MustGet("dbConnection").(*pgxpool.Pool)
	values := fmt.Sprintf("VALUES ('%d', '%d', '%s');", kudo.SenderID, kudo.ReceiverID, kudo.Message)
	sqlStr := "INSERT INTO kudos (sender_id, receiver_id, message)" + values
	log.Println(sqlStr)
	_, err := dbPool.Exec(context.Background(), sqlStr)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, err)
	}

	c.Status(http.StatusCreated)
}

func GetKudos(c *gin.Context) {
	var kudos = []Kudo{}
	dbPool := c.MustGet("dbConnection").(*pgxpool.Pool)

	rows, err := dbPool.Query(context.Background(), "SELECT * FROM kudos ORDER BY id;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}

		var kudo = Kudo{
			ID:         values[0].(int32),
			SenderID:   values[1].(int32),
			ReceiverID: values[2].(int32),
			Message:    values[3].(string),
		}

		kudos = append(kudos, kudo)
	}

	c.IndentedJSON(http.StatusOK, kudos)
}

func GetKudoById(c *gin.Context) {
	var kudo Kudo
	dbPool := c.MustGet("dbConnection").(*pgxpool.Pool)

	kudoId := c.Param("id")
	sqlStr := fmt.Sprintf("SELECT * FROM kudos WHERE id = %s;", kudoId)
	log.Println(sqlStr)
	rows, err := dbPool.Query(context.Background(), sqlStr)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal(err)
		}

		kudo = Kudo{
			ID:         values[0].(int32),
			SenderID:   values[1].(int32),
			ReceiverID: values[2].(int32),
			Message:    values[3].(string),
		}
		c.IndentedJSON(http.StatusOK, kudo)
	}
	c.Status(http.StatusNotFound)
}
