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

func TestGetPersonById(t *testing.T) {
	// Create a new person
	person, err := json.Marshal(Person{FirstName: "John", LastName: "Smith", Email: "john.smith@email.com"})
	if err != nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(person))
	router.ServeHTTP(w, req)

	// Person exists
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/persons/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Person doesn't exist
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/persons/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreatePersons(t *testing.T) {
	// Nominal case
	person, err := json.Marshal(Person{FirstName: "John", LastName: "Doe", Email: "john.doe@email.com"})
	if err != nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(person))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Email is the same
	person, err = json.Marshal(Person{FirstName: "John", LastName: "Doe", Email: "john.doe@email.com"})
	if err != nil {
		return
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(person))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Email is different
	person, err = json.Marshal(Person{FirstName: "John", LastName: "Doe", Email: "johndoe@email.com"})
	if err != nil {
		return
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(person))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
