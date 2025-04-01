package server

import (
	"errors"
	"strconv"

	"github.com/shkh/lastfm-go/lastfm"
)

type LastFMService struct {
	client *lastfm.Api
}

func NewLastFMService(apiKey string) *LastFMService {
	return &LastFMService{
		client: lastfm.New(apiKey, ""),
	}
}

func (s *LastFMService) getFirstTrack(user string) (Track, error) {
	recentTracks, err := s.client.User.GetRecentTracks(lastfm.P{"user": user, "limit": "1"})
	if err != nil {
		return Track{}, errors.New(err.Error())
	}

	if len(recentTracks.Tracks) == 0 {
		return Track{}, errors.New("no tracks found")
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
		Images:       extractImageUrls(firstTrack.Images),
		PlayedAt:     playedAt,
		PlayedAtUnix: playedAtUnix,
	}, nil
}

func extractImageUrls(images []struct {
	Size string `xml:"size,attr"`
	Url  string `xml:",chardata"`
}) []string {
	urls := make([]string, len(images))
	for i, img := range images {
		urls[i] = img.Url
	}
	return urls
}
