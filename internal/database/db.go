package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectPostgres() {
	currDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	estConnection, err := pgx.Connect(context.Background(), currDSN)
	if err != nil {
		log.Fatalf("Couldn't connect to your DB: %v", err)
	}

	log.Println("Connection to the db was succcessful")
	DB = estConnection
}
