terraform {
  required_version = ">= 1.10"

  backend "s3" {
    bucket       = "alexraskin-tfstate"
    key          = "lastfm-now-playing/terraform.tfstate"
    region       = "us-east-1"
    encrypt      = true
    use_lockfile = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

# Shared with the infrastructure repo, which owns the parameter (00-global).
data "aws_ssm_parameter" "cloudflare_api_token" {
  name = var.cloudflare_api_token_parameter
}

provider "cloudflare" {
  api_token = data.aws_ssm_parameter.cloudflare_api_token.value
}
