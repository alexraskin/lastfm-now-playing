variable "region" {
  type    = string
  default = "us-east-1"
}

variable "name" {
  type    = string
  default = "lastfm-now-playing"
}

variable "package_path" {
  description = "Zip containing the bootstrap binary, built by `mise run package`."
  type        = string
  default     = "../dist/lambda.zip"
}

variable "domain_name" {
  type    = string
  default = "lastfm.alexraskin.com"
}

variable "cloudflare_zone" {
  type    = string
  default = "alexraskin.com"
}

variable "cloudflare_api_token_parameter" {
  description = "SSM SecureString holding the Cloudflare API token."
  type        = string
  default     = "/infrastructure/cloudflare/api_token"
}

variable "reserved_concurrency" {
  type    = number
  default = 5
}

variable "throttling_rate_limit" {
  type    = number
  default = 10
}

variable "throttling_burst_limit" {
  type    = number
  default = 20
}

variable "github_repository" {
  description = "GitHub repository allowed to assume the deploy role, as owner/repo."
  type        = string
  default     = "alexraskin/lastfm-now-playing"
}

variable "create_oidc_provider" {
  description = "Set to false if the GitHub OIDC provider already exists in this account."
  type        = bool
  default     = false
}

variable "state_bucket" {
  type    = string
  default = "alexraskin-tfstate"
}
