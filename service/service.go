package service

import (
	"strconv"

	"github.com/alexraskin/lastfm-now-playing/models"
	"github.com/alexraskin/lastfm-now-playing/utils"
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

func (s *LastFMService) GetFirstTrack(user string) (models.Track, error) {
	recentTracks, err := s.client.User.GetRecentTracks(lastfm.P{"user": user, "limit": "1"})
	if err != nil {
		return models.Track{}, utils.LastFMError{Message: err.Error()}
	}

	if len(recentTracks.Tracks) == 0 {
		return models.Track{}, utils.LastFMError{Message: "no tracks found"}
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

	return models.Track{
		Artist:       firstTrack.Artist.Name,
		Name:         firstTrack.Name,
		Album:        firstTrack.Album.Name,
		NowPlaying:   isNowPlaying,
		Images:       utils.ExtractImageUrls(firstTrack.Images),
		PlayedAt:     playedAt,
		PlayedAtUnix: playedAtUnix,
	}, nil
}
