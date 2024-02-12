package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("old-school-web-gofiber/templates/", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(requestid.New())
	app.Use(helmet.New())
	app.Use(etag.New())

	app.Static("/static", "old-school-web-gofiber/static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{}, "layout")
	})

	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}
}
