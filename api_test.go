package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	initConfig()
	os.Exit(m.Run())
}

func TestGetPersons(t *testing.T) {
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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/persons", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePersons(t *testing.T) {
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

	person, err := json.Marshal(Person{FirstName: "Titi"})
	if err == nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(person))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
