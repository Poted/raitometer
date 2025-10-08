package main

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func connectToDB() *sqlx.DB {
	dsn := "host=localhost port=5432 user=raitometer_user password=raitometer_password dbname=raitometer_db sslmode=disable"

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	log.Println("database connection established successfully")

	err = db.Ping()
	if err != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	return db
}
