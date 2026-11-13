variable "name" {
  type = string
}

variable "description" {
  type    = string
  default = null
}

variable "function_name" {
  type = string
}

variable "function_arn" {
  type = string
}

variable "domain_name" {
  type = string
}

variable "cloudflare_zone" {
  type = string
}

variable "throttling_rate_limit" {
  description = "Steady-state requests per second."
  type        = number
}

variable "throttling_burst_limit" {
  type = number
}
