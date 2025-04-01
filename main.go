package main

import (
	"embed"
	"flag"
	"html/template"
	"io"
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
	Templates embed.FS
)

func main() {

	port := flag.String("port", "3000", "port to listen on")
	devMode := flag.Bool("dev", false, "run in dev mode")
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	config := server.LoadConfig(*configPath)
	config.Port = *port

	var (
		tmplFunc server.ExecuteTemplateFunc
	)

	slog.Info("Starting lastfm-now-playing...", slog.Any("version", version), slog.Any("commit", commit), slog.Any("buildTime", buildTime))

	if *devMode {
		slog.Info("running in dev mode")
		tmplFunc = func(wr io.Writer, name string, data any) error {
			tmpl, err := template.New("").ParseGlob("templates/*.gohtml")
			if err != nil {
				return err
			}
			return tmpl.ExecuteTemplate(wr, name, data)
		}
	} else {
		tmpl, err := template.New("").ParseFS(Templates, "templates/*.gohtml")
		if err != nil {
			slog.Error("failed to parse templates", slog.Any("error", err))
			os.Exit(-1)
		}
		tmplFunc = tmpl.ExecuteTemplate
	}

	client := server.NewLastFMService(config.LastFMAPIKey)

	server := server.NewServer(
		server.FormatBuildVersion(version, commit, buildTime),
		http.DefaultClient,
		tmplFunc,
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
