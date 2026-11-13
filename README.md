# Last.fm Now Playing

A simple API that fetches the currently playing or most recently played track from Last.fm for a given user.

[![Image from Gyazo](https://i.gyazo.com/5632a2462e3cee91a25d1824a45f318d.png)](https://gyazo.com/5632a2462e3cee91a25d1824a45f318d)

[![Last.FM Last Played Song](https://img.shields.io/endpoint?color=blueviolet&url=https://lastfm.alexraskin.com/alexraskin?format=shields.io)](https://github.com/alexraskin/lastfm-now-playing)

```bash
curl -s https://lastfm.alexraskin.com/{your lastfm username}
```

## Usage

### Get a user's currently playing track

```
GET /:username
```

Example:
```
GET /twizycat
```

### Response formats

#### JSON (default)

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

#### Shields.io

```
GET /twizycat?format=shields.io
```
```
[![Last.FM Last Played Song](https://img.shields.io/endpoint?color=blueviolet&url=https://lastfm.alexraskin.com/twizycat?format=shields.io)](https://github.com/alexraskin/lastfm-now-playing)
```

#### Glance widget

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

## Self-hosting

The API runs on AWS Lambda (`provided.al2023`, arm64) behind an API Gateway HTTP API, with a custom domain on Cloudflare DNS. Everything is deployed with Terraform.

### Requirements

- [mise](https://mise.jdx.dev), which installs the pinned Go and Terraform versions from `mise.toml`
- An AWS account and a Cloudflare zone
- A Last.fm API key from [last.fm/api](https://www.last.fm/api/). Only the API key is used, not the shared secret.

### Layout

```
terraform/                  root stack
  modules/lambda-function/  Lambda + SSM parameter for the Last.fm API key
  modules/http-api/         API Gateway + ACM certificate + Cloudflare DNS
  modules/github-oidc/      IAM role GitHub Actions deploys with
```

### Configuration

Before deploying a fork, change these to match your own setup:

- The `backend "s3"` block in `terraform/versions.tf`, which holds the state bucket
- In `terraform/variables.tf`:
  - `domain_name` and `cloudflare_zone`
  - `github_repository`
  - `state_bucket`
  - `cloudflare_api_token_parameter`: the SSM `SecureString` that holds a Cloudflare API token with `Zone:Read` and `DNS:Edit`
  - `create_oidc_provider`: set it to `true` if your account doesn't have the GitHub OIDC provider yet

In Cloudflare, set the zone's SSL/TLS mode to **Full (strict)**, because API Gateway only serves HTTPS.

### Deploy

1. Run the first apply locally with your own AWS credentials:
   ```
   mise install
   mise run deploy   # builds dist/lambda.zip and applies terraform/
   ```

2. Set the Last.fm API key. Terraform creates the parameter with a placeholder value and ignores later changes to it:
   ```
   aws ssm put-parameter --name /lastfm-now-playing/lastfm-api-key --type SecureString --overwrite --value <your api key>
   ```
   Lambda environments that are already running keep the cached key until they're recycled.

3. Optionally, let GitHub Actions deploy every push to `main`. It assumes the deploy role over OIDC, so no AWS keys are stored in GitHub:
   ```
   gh variable set AWS_ROLE_ARN --body "$(terraform -chdir=terraform output -raw deploy_role_arn)"
   ```
   The deploy role can't modify itself, so changes to `modules/github-oidc` have to be applied locally.

Requests are capped by API Gateway throttling (10 req/s, burst 20) and Lambda reserved concurrency (5). Both settings are variables in `terraform/variables.tf`.
