# Last.fm Now Playing

A simple API that fetches the currently playing or most recently played track from Last.fm for a given user.

[![Last.FM Last Played Song](https://img.shields.io/endpoint?color=purple&url=https://lastfm.alexraskin.com/twizycat?format=shields.io)](https://github.com/alexraskin/lastfm-now-playing)


```bash
curl -s https://lastfm.alexraskin.com/{your lastfm username}
```

## Setup

- Go 1.24 or later
- Last.fm API key (get one from [Last.fm API](https://www.last.fm/api/))

### Local Dev

1. Clone this repository
   ```
   git clone https://github.com/alexraskin/lastfm-now-playing.git
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set the `LASTFM_API_KEY` environment variable:
   ```
   export LASTFM_API_KEY={your lastfm api key}
   ```

4. Run the application:
   ```
   go run main.go
   ```

### Docker

1. Build the Docker image:
   ```
   docker build -t lastfm-now-playing .
   ```
2. Run the container:
   ```
   docker run -p 3000:3000 -e LASTFM_API_KEY={your lastfm api key} lastfm-now-playing
   ```

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

#### Shields.io Format:

```
GET /twizycat?format=shields.io
```
```
[![Last.FM Last Played Song](https://img.shields.io/endpoint?color=blueviolet&url=https://playing.alexraskin.com/twizycat?format=shields.io)](https://github.com/alexraskin/lastfm-now-playing)
```

#### Glance Widget Format:

```
GET /widget/twizycat
```
```yaml
- type: extension
  url: https://lastfm.alexraskin.com/widget/twizycat
  allow-potentially-dangerous-html: true
  cache: 30s
```
