package main

import (
	"embed"
	"html/template"
	"log/slog"
	"os"

	"github.com/alexraskin/lastfm-now-playing/server"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
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
	slog.Info("Starting lastfm-now-playing...", slog.Any("version", version), slog.Any("commit", commit), slog.Any("buildTime", buildTime))

	tmpl, err := template.New("").ParseFS(Templates, "templates/*.gohtml")
	if err != nil {
		slog.Error("failed to parse templates", slog.Any("error", err))
		os.Exit(-1)
	}

	client := server.NewLastFMService(server.LoadLastFMAPIKey())

	server := server.NewServer(
		server.FormatBuildVersion(version, commit, buildTime),
		tmpl.ExecuteTemplate,
		client,
	)

	// Function URLs use the API Gateway v2 payload format.
	lambda.Start(httpadapter.NewV2(server.Routes()).ProxyWithContext)
}
