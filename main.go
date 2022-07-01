package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

var dbUrl string
var dbPool *pgxpool.Pool

func initConfig() {
	viper.SetConfigName(fmt.Sprintf("config-%s", os.Getenv("ENVIRONMENT")))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// Logging config
	log.SetPrefix(viper.Get("log.prefix").(string))
	log.SetFlags(viper.Get("log.flags").(int))

	// DB config
	dbUrl = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		viper.Get("postgres.user"),
		viper.Get("postgres.password"),
		viper.Get("postgres.host"),
		viper.Get("postgres.port"),
		viper.Get("postgres.database"))
}

func initDB() {
	var err error

	dbPool, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	_, err = dbPool.Exec(context.Background(), ReadSqlFile("sql/drop_persons_table.sql"))
	if err != nil {
		log.Println(err)
	}
	_, err = dbPool.Exec(context.Background(), ReadSqlFile("sql/create_persons_table.sql"))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	initConfig()
	initDB()
	defer dbPool.Close()

	router := SetupRouter(dbPool)
	router.Run("localhost:8080")
}
