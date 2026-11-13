variable "name" {
  type = string
}

variable "description" {
  type    = string
  default = null
}

variable "package_path" {
  description = "Zip containing the bootstrap binary."
  type        = string
}

variable "reserved_concurrency" {
  description = "Maximum concurrent executions; requests beyond this are throttled. null for no limit."
  type        = number
  default     = null
}
