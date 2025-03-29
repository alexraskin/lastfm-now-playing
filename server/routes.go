package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Origin, Content-Type, Accept"},
		AllowedMethods: []string{"GET", "OPTIONS"},
	}))

	r.Use(httprate.Limit(
		100,
		time.Minute,
		httprate.WithLimitHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
			}),
		),
	))

	r.Get("/", s.index)
	r.Get("/widget", s.nowPlayingWidget)
	r.Get("/{user}", s.nowPlaying)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	return r
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	apiDoc := ApiDoc{
		Status: "ok",
		Endpoints: []Endpoint{
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
			{
				Method:      "GET",
				Path:        "/widget/:user",
				Description: "Get the currently playing track that supports Glance Widgets",
			},
		},
	}

	_ = json.NewEncoder(w).Encode(apiDoc)
}

func (s *Server) nowPlaying(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	if user == "" {
		http.Error(w, "User is required", http.StatusBadRequest)
	}

	track, err := s.lfmclient.GetFirstTrack(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	format := r.URL.Query().Get("format")

	if format == "shields.io" {
		message := fmt.Sprintf("%s - %s", track.Name, track.Artist)
		label := "Currently Playing"
		if !track.NowPlaying {
			label = "Last Played"
			message = fmt.Sprintf("%s - %s", track.Name, track.Artist)
		}

		_ = json.NewEncoder(w).Encode(ShieldsResponse{
			SchemaVersion: 1,
			Label:         label,
			Message:       message,
		})
		return
	}

	_ = json.NewEncoder(w).Encode(TrackResponse{
		Track:      track.Name,
		Artist:     track.Artist,
		Album:      track.Album,
		NowPlaying: track.NowPlaying,
		Images:     track.Images,
		PlayedAt:   track.PlayedAt,
	})
}

func (s *Server) nowPlayingWidget(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "Missing 'user' parameter", http.StatusBadRequest)
	}

	track, err := s.lfmclient.GetFirstTrack(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageURL := ""
	if len(track.Images) > 0 {
		for i := len(track.Images) - 1; i >= 0; i-- {
			if track.Images[i] != "" {
				imageURL = track.Images[i]
				break
			}
		}
	}

	file, err := s.templates.Open("widget.gohtml")
	if err != nil {
		http.Error(w, "Template not found: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("widget").Parse(string(data))
	if err != nil {
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, struct {
		Track
		ImageURL string
	}{
		Track:    track,
		ImageURL: imageURL,
	}); err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Widget-Title", "Last.fm")
	w.Header().Set("Widget-Content-Type", "html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(out.Bytes())
}
