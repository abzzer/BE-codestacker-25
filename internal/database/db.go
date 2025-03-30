package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	var err error
	for attempts := 1; attempts <= 10; attempts++ {
		DB, err = pgx.Connect(context.Background(), currDSN)
		if err == nil {
			log.Println("Connection to the DB was successful.")
			return
		}
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Attempted to connect - but failed: %v", err)
}
