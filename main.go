package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexraskin/lastfm-now-playing/server"
)

var (
	version   = "unknown"
	commit    = "unknown"
	buildTime = "unknown"
)

var (
	//go:embed templates/**
	embeddedTemplates embed.FS
)

func main() {

	config := server.LoadConfig("config.yaml")

	slog.Info("Starting lastfm-now-playing...", slog.Any("version", version), slog.Any("commit", commit), slog.Any("buildTime", buildTime))

	client := server.NewLastFMService(config.LastFMAPIKey)

	server := server.NewServer(
		server.FormatBuildVersion(version, commit, buildTime),
		http.DefaultClient,
		embeddedTemplates,
		client,
		config,
	)

	go server.Start()
	defer server.Close()

	slog.Info("started lastfm-now-playing", slog.Any("listen_addr", config.Port))
	si := make(chan os.Signal, 1)
	signal.Notify(si, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-si
}
