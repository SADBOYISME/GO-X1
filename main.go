package main

import (
	"GO-X/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// creat type for user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func main() {

	app := fiber.New()

	// * set cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, HEAD",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	users = append(users, User{"admin", "admin"})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusBadRequest)
	})

	app.Get("/uuid", func(c *fiber.Ctx) error {
		return c.SendString(auth.GenUuid())
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(users)
	})

	app.Listen(":8080")

}
