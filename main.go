package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/shkh/lastfm-go/lastfm"
)

// LastFMTrack represents the structure of a track from the LastFM API
type LastFMTrack struct {
	NowPlaying string `xml:"nowplaying,attr,omitempty"`
	Artist     struct {
		Name string `xml:",chardata"`
		Mbid string `xml:"mbid,attr"`
	} `xml:"artist"`
	Name       string `xml:"name"`
	Streamable string `xml:"streamable"`
	Mbid       string `xml:"mbid"`
	Album      struct {
		Name string `xml:",chardata"`
		Mbid string `xml:"mbid,attr"`
	} `xml:"album"`
	Url    string `xml:"url"`
	Images []struct {
		Size string `xml:"size,attr"`
		Url  string `xml:",chardata"`
	} `xml:"image"`
	Date struct {
		Uts  string `xml:"uts,attr"`
		Date string `xml:",chardata"`
	} `xml:"date"`
}

// Track represents a single track from LastFM
type Track struct {
	Artist       string
	Name         string
	Album        string
	NowPlaying   bool
	PlayedAt     string
	PlayedAtUnix int64
}

// CustomJSONEncoder creates a JSON encoder that doesn't escape HTML
func CustomJSONEncoder(v any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	// Remove the trailing newline that json.Encoder adds
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}

func main() {
	apiKey := os.Getenv("LASTFM_API_KEY")
	if apiKey == "" {
		log.Fatal("LASTFM_API_KEY is not set")
	}

	client := lastfm.New(apiKey, "")

	app := fiber.New(fiber.Config{
		ServerHeader: "Last.FM Recent",
		JSONEncoder:  CustomJSONEncoder,
	})

	app.Get("/:user", func(c *fiber.Ctx) error {
		return indexHandler(c, client)
	})

	log.Fatal(app.Listen(":3000"))
}

func getFirstTrack(client *lastfm.Api, user string) (Track, error) {
	recentTracks, err := client.User.GetRecentTracks(lastfm.P{"user": user, "limit": "1"})
	if err != nil {
		return Track{}, err
	}

	if len(recentTracks.Tracks) == 0 {
		return Track{}, fmt.Errorf("no tracks found")
	}

	firstTrack := recentTracks.Tracks[0]

	isNowPlaying := firstTrack.NowPlaying == "true"

	var playedAt string
	var playedAtUnix int64

	if !isNowPlaying && firstTrack.Date.Date != "" {
		playedAt = firstTrack.Date.Date
		if firstTrack.Date.Uts != "" {
			if utsVal, err := strconv.ParseInt(firstTrack.Date.Uts, 10, 64); err == nil {
				playedAtUnix = utsVal
			}
		}
	}

	return Track{
		Artist:       firstTrack.Artist.Name,
		Name:         firstTrack.Name,
		Album:        firstTrack.Album.Name,
		NowPlaying:   isNowPlaying,
		PlayedAt:     playedAt,
		PlayedAtUnix: playedAtUnix,
	}, nil
}

func indexHandler(c *fiber.Ctx, lfm *lastfm.Api) error {
	user := c.Params("user")
	track, err := getFirstTrack(lfm, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	format := c.Query("format")

	if format == "shields.io" {
		message := fmt.Sprintf("%s - %s", track.Name, track.Artist)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"schemaVersion": 1,
			"label":         "Currently Playing",
			"message":       message,
			"color":         "green",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"track":        track.Name,
		"artist":       track.Artist,
		"album":        track.Album,
		"nowPlaying":   track.NowPlaying,
		"playedAt":     track.PlayedAt,
		"playedAtUnix": track.PlayedAtUnix,
	})
}
