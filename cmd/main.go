package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("We have a working API!!")
	})

	log.Printf("Server is running")
	log.Fatal(app.Listen(":8080"))
}
