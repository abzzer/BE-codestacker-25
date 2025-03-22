package main

import (
	"log"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/abzzer/BE-codestacker-25/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("You need .env for this project to run -> Please check the docker-compose for the required env fields")
	}

	database.ConnectPostgres()
	database.ConnectMinIO()

	app := fiber.New()

	routes.RegisterRoutes(app)

	log.Println("The server is running")
	log.Fatal(app.Listen(":8080"))
}
