package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	initConfig()
	initDB()
	defer dbPool.Close()

	router = SetupRouter(dbPool)

	os.Exit(m.Run())
}

func TestGetPersons(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/persons", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var persons []Person
	err := json.Unmarshal([]byte(w.Body.String()), persons)
	if err != nil {
		return
	}
	assert.IsType(t, Person{}, persons[0])
}

func TestCreatePersons(t *testing.T) {
	person, err := json.Marshal(Person{FirstName: "Titi"})
	if err != nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(person))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
