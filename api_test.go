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

func TestGetPersons(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/persons", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPersonById(t *testing.T) {
	// Person exists
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/persons/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Person doesn't exist
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/persons/10000", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGiveKudo(t *testing.T) {
	// Nominal case
	kudo, err := json.Marshal(Kudo{SenderID: 1, ReceiverID: 3, Message: "Awesome!"})
	if err != nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/kudos", bytes.NewBuffer(kudo))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Receiver doesn't exist
	kudo, err = json.Marshal(Kudo{SenderID: 1, ReceiverID: 2, Message: "Awesome!"})
	if err != nil {
		return
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/kudos", bytes.NewBuffer(kudo))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Sender doesn't exist
	kudo, err = json.Marshal(Kudo{SenderID: 1000, ReceiverID: 1, Message: "Awesome!"})
	if err != nil {
		return
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/kudos", bytes.NewBuffer(kudo))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Self kudos are not allowed
	kudo, err = json.Marshal(Kudo{SenderID: 1, ReceiverID: 1, Message: "Awesome!"})
	if err != nil {
		return
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/kudos", bytes.NewBuffer(kudo))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

}

func TestGetKudos(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/kudos", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
