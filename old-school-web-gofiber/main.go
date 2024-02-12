package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/memory"

	"github.com/gofiber/fiber/v2/utils"
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
	app.Use(compress.New())

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrftoken",
		CookieSameSite: "Lax",
		Expiration:     24 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
	}))

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))

	storage := memory.New()
	store := session.New(session.Config{
		KeyLookup:      "cookie:session",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		Storage:        storage,
	})

	app.Static("/static", "old-school-web-gofiber/static")

	app.Get("/", func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			return err
		}
		if err := session.Save(); err != nil {
			return err
		}
		return c.Render("index", fiber.Map{}, "layout")
	})

	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}
}
