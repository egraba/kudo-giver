package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	HOST     = "localhost"
	PORT     = 5432
	DATABASE = "kudo-giver"
	USER     = "devuser"
	PASSWORD = "devpwd"
)

var dbUrl = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", USER, PASSWORD, HOST, PORT, DATABASE)

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
