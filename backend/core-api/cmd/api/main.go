package main

import (
	"log"

	"github.com/Poted/raitometer/backend/core-api/internal/server"
)

func main() {
	dbConnection := connectToDB()
	defer dbConnection.Close()

	srv := server.New(dbConnection)

	port := ":8080"
	err := srv.Start(port)
	if err != nil {
		log.Fatalf("fatal server error: %v", err)
	}
}
