package server

import (
	"fmt"
	"io"
	"runtime"
	"time"
)

type ExecuteTemplateFunc func(wr io.Writer, name string, data any) error

type Server struct {
	version   string
	tmplFunc  ExecuteTemplateFunc
	lfmclient *LastFMService
}

func NewServer(version string, tmplFunc ExecuteTemplateFunc, lfmclient *LastFMService) *Server {
	return &Server{
		version:   version,
		tmplFunc:  tmplFunc,
		lfmclient: lfmclient,
	}
}

func FormatBuildVersion(version string, commit string, buildTime string) string {
	if len(commit) > 7 {
		commit = commit[:7]
	}

	buildTimeStr := "unknown"
	if buildTime != "unknown" {
		parsedTime, _ := time.Parse(time.RFC3339, buildTime)
		if !parsedTime.IsZero() {
			buildTimeStr = parsedTime.Format(time.ANSIC)
		}
	}
	return fmt.Sprintf("Go Version: %s\nVersion: %s\nCommit: %s\nBuild Time: %s\nOS/Arch: %s/%s\n", runtime.Version(), version, commit, buildTimeStr, runtime.GOOS, runtime.GOARCH)
}
