package models

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
	Images       []string
	PlayedAt     string
	PlayedAtUnix int64
}

// Endpoint represents a single endpoint in the API
type Endpoint struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

// ApiDoc represents the API documentation
type ApiDoc struct {
	Status    string     `json:"status"`
	Endpoints []Endpoint `json:"endpoints"`
}

type TrackResponse struct {
	Track        string   `json:"track"`
	Artist       string   `json:"artist"`
	Album        string   `json:"album"`
	NowPlaying   bool     `json:"nowPlaying"`
	Images       []string `json:"image,omitempty"`
	PlayedAt     string   `json:"playedAt,omitempty"`
	PlayedAtUnix int64    `json:"playedAtUnix,omitempty"`
}

type ShieldsResponse struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
}
