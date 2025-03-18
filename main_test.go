package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestTrackResponseJSON tests the JSON response format
func TestTrackResponseJSON(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		track := Track{
			Artist:     "Test Artist",
			Name:       "Test Track",
			Album:      "Test Album",
			NowPlaying: true,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"track":        track.Name,
			"artist":       track.Artist,
			"album":        track.Album,
			"nowPlaying":   track.NowPlaying,
			"playedAt":     track.PlayedAt,
			"playedAtUnix": track.PlayedAtUnix,
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, "Test Track", result["track"])
	assert.Equal(t, "Test Artist", result["artist"])
	assert.Equal(t, "Test Album", result["album"])
	assert.Equal(t, true, result["nowPlaying"])
}

// TestShieldsIOFormat tests the shields.io format response
func TestShieldsIOFormat(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"schemaVersion": 1,
			"label":         "Currently Playing",
			"message":       "Test Artist - Test Track",
			"color":         "green",
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), result["schemaVersion"])
	assert.Equal(t, "Currently Playing", result["label"])
	assert.Equal(t, "Test Artist - Test Track", result["message"])
	assert.Equal(t, "green", result["color"])
}

// TestError tests the error response
func TestError(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "API error",
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Contains(t, result, "error")
	assert.Equal(t, "API error", result["error"])
}
