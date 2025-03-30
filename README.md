# Last.fm Now Playing

A simple API that fetches the currently playing or most recently played track from Last.fm for a given user.

[![Image from Gyazo](https://i.gyazo.com/5632a2462e3cee91a25d1824a45f318d.png)](https://gyazo.com/5632a2462e3cee91a25d1824a45f318d)

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

3. Take a look at the `config.yaml.example` file and create a `config.yaml` file in the root of the project.


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
   docker run -p 8000:8000 -e LASTFM_API_KEY={your lastfm api key} lastfm-now-playing
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
[![Last.FM Last Played Song](https://img.shields.io/endpoint?color=blueviolet&url=https://lastfm.alexraskin.com/twizycat?format=shields.io)](https://github.com/alexraskin/lastfm-now-playing)
```

#### Glance Widget Format:

```
GET /widget
```
```yaml
- type: extension
  url: https://lastfm.alexraskin.com/widget
  allow-potentially-dangerous-html: true
  cache: 30s
  parameters:
    user: twizycat
```
