package main

import (
	"log"
	"os"
)

func ReadSqlFile(fileName string) string {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}