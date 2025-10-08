package main

import (
	"log"
	"os"

	"github.com/Poted/raitometer/backend/core-api/internal/server"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := connectToDB(connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("database connection established successfully")

	s := server.New(db)
	err = s.Start(":8080")
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
