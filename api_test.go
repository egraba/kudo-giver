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
	req, _ := http.NewRequest(http.MethodGet, PersonEndpoint, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePersons(t *testing.T) {
	router := setupRouter()

	person, err := json.Marshal(Person{FirstName: "Titi"})
	if err == nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, PersonEndpoint, bytes.NewBuffer(person))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetKudos(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, KudoEndpoint, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGiveKudos(t *testing.T) {
	router := setupRouter()

	kudo, err := json.Marshal(Kudo{SenderID: 1, ReceiverID: 2, Message: "Bedesi"})
	if err == nil {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, KudoEndpoint, bytes.NewBuffer(kudo))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
