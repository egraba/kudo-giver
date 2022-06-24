package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

const dbUrl = "postgres://devuser:devpwd@localhost:5432/kudo-giver"

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
