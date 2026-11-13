variable "name" {
  description = "Prefix of the resources the deploy role may manage."
  type        = string
}

variable "github_repository" {
  description = "GitHub repository allowed to assume the deploy role, as owner/repo."
  type        = string
}

variable "create_oidc_provider" {
  description = "Set to false if the GitHub OIDC provider already exists in this account."
  type        = bool
  default     = true
}

variable "state_bucket" {
  type = string
}

variable "state_key" {
  type = string
}

variable "read_parameter_arns" {
  description = "SSM parameters outside the name prefix the deploy role may read."
  type        = list(string)
  default     = []
}
