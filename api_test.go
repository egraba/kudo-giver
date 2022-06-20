package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPersons(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/persons", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePersons(t *testing.T) {
	router := setupRouter()

	newPerson, err := json.Marshal(person{FirstName: "Titi"})
	if err == nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(newPerson))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetKudos(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/kudos", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGiveKudos(t *testing.T) {
	router := setupRouter()

	newKudo, err := json.Marshal(kudo{SenderID: 1, ReceiverID: 2, Message: "Bedesi"})
	if err == nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/kudos", bytes.NewBuffer(newKudo))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
