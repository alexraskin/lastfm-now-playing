# Last.fm Recent Tracks API

A simple API that fetches the currently playing or most recently played track from Last.fm for a given user.

## Try it out:

```bash
curl -s https://playing.alexraskin.com/{your lastfm username}
```

## Setup

### Prerequisites

- Go 1.24 or later
- Last.fm API key (get one from [Last.fm API](https://www.last.fm/api/))

### Local Development

1. Clone this repository
2. Install dependencies:
   ```
   go mod download
   ```
3. Run the application:
   ```
   go run main.go
   ```

### Using Docker

1. Build the Docker image:
   ```
   docker build -t lastfm-api .
   ```
2. Run the container:
   ```
   docker run -p 3000:3000 lastfm-api
   ```

## API Usage

### Get a user's currently playing track

```
GET /:username
```

Example:
```
GET /twizycat
```

### Response Formats

#### Default JSON response:

```json
{
  "album": "Album Name",
  "artist": "Artist Name",
  "nowPlaying": true,
  "playedAt": "",
  "playedAtUnix": 0,
  "track": "Track Name"
}
```

#### Shields.io format:

Add `?format=shields.io` to get a response compatible with Shields.io dynamic badges:

```
GET /twizycat?format=shields.io
```

Response:
```json
{
  "schemaVersion": 1,
  "label": "Currently Playing",
  "message": "Artist Name - Track Name",
  "color": "green"
}
```

## License

MIT 