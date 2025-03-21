package main

import (
	"log"
	"os"
	"time"

	"github.com/alexraskin/lastfm-now-playing/handlers"
	"github.com/alexraskin/lastfm-now-playing/service"
	"github.com/alexraskin/lastfm-now-playing/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config := struct {
		Port          string
		LastFMAPIKey  string
		RateLimit     int
		RateLimitTime time.Duration
	}{
		Port:          "3000",
		RateLimit:     10,
		RateLimitTime: 1 * time.Minute,
	}
	if envPort := os.Getenv("PORT"); envPort != "" {
		config.Port = envPort
	}

	config.LastFMAPIKey = os.Getenv("LASTFM_API_KEY")
	if config.LastFMAPIKey == "" {
		log.Fatal("LASTFM_API_KEY is not set")
	}

	client := service.NewLastFMService(config.LastFMAPIKey)

	app := fiber.New(fiber.Config{
		ServerHeader: "Last.FM Recent",
		JSONEncoder:  utils.CustomJSONEncoder,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, OPTIONS",
	}))

	// don't ddos last.fm lol
	app.Use(limiter.New(limiter.Config{
		Max:        config.RateLimit,
		Expiration: config.RateLimitTime,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded, please try again later",
			})
		},
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] - ${method} ${path} - ${status} ${latency}\n",
	}))

	app.Get("/", handlers.IndexHandler)

	app.Get("/:user", func(c *fiber.Ctx) error {
		return handlers.NowPlayingHandler(c, client)
	})

	log.Fatal(app.Listen(":" + config.Port))
}
