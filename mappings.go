package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupRouter(dbPool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(connectDb(dbPool))

	v1 := r.Group("api/v1")
	{
		v1.POST("/persons", CreatePerson)
		v1.GET("/persons", GetPersons)
		v1.GET("/persons/:id", GetPersonById)

		v1.POST("/kudos", GiveKudo)
		v1.GET("/kudos", GetKudos)
		v1.GET("/kudos/:id", GetKudoById)
	}

	return r
}
