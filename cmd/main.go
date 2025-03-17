package main

import (
	"log"

	"github.com/abzzer/BE-codestacker-25/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("you need a .env to run the project - please see the docs")
	}

	database.ConnectPostgres()
	database.ConnectMinIO()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("We have a working API with databases!!")
	})

	log.Println("The server is running")
	log.Fatal(app.Listen(":8080"))
}
