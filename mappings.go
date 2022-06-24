package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupRouter(dbPool *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(connectDb(dbPool))

	router.GET("/persons", getPersons)
	router.POST("/persons", createPerson)

	return router
}
