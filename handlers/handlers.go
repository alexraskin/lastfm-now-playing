package handlers

import (
	"fmt"

	"github.com/alexraskin/lastfm-now-playing/models"
	"github.com/alexraskin/lastfm-now-playing/service"

	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx) error {
	apiDoc := models.ApiDoc{
		Status: "ok",
		Endpoints: []models.Endpoint{
			{
				Method:      "GET",
				Path:        "/:user",
				Description: "Get the currently playing track for a user",
			},
			{
				Method:      "GET",
				Path:        "/:user/?format=shields.io",
				Description: "Get the currently playing track for a user in Shields.io format",
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(apiDoc)
}

func NowPlayingHandler(c *fiber.Ctx, lfmclient *service.LastFMService) error {
	user := c.Params("user")
	if user == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User is required",
		})
	}

	track, err := lfmclient.GetFirstTrack(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	format := c.Query("format")

	if format == "shields.io" {
		message := fmt.Sprintf("%s - %s", track.Name, track.Artist)
		label := "Currently Playing"
		if !track.NowPlaying {
			label = "Last Played"
			message = fmt.Sprintf("%s - %s", track.Name, track.Artist)
		}

		return c.Status(fiber.StatusOK).JSON(models.ShieldsResponse{
			SchemaVersion: 1,
			Label:         label,
			Message:       message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.TrackResponse{
		Track:      track.Name,
		Artist:     track.Artist,
		Album:      track.Album,
		NowPlaying: track.NowPlaying,
		PlayedAt:   track.PlayedAt,
	})
}
