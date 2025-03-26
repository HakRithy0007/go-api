package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a new Fiber app
	app := fiber.New()

	// Define a simple route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber! new world")
	})

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
