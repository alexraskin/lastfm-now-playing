module "lambda_function" {
  source = "./modules/lambda-function"

  name                 = var.name
  description          = "Last.fm now playing API"
  package_path         = var.package_path
  reserved_concurrency = var.reserved_concurrency
}

module "http_api" {
  source = "./modules/http-api"

  name                   = var.name
  description            = "Last.fm now playing API"
  function_name          = module.lambda_function.function_name
  function_arn           = module.lambda_function.function_arn
  domain_name            = var.domain_name
  cloudflare_zone        = var.cloudflare_zone
  throttling_rate_limit  = var.throttling_rate_limit
  throttling_burst_limit = var.throttling_burst_limit
}

module "github_oidc" {
  source = "./modules/github-oidc"

  name                 = var.name
  github_repository    = var.github_repository
  create_oidc_provider = var.create_oidc_provider
  state_bucket         = var.state_bucket
  state_key            = "${var.name}/terraform.tfstate"
  read_parameter_arns  = [data.aws_ssm_parameter.cloudflare_api_token.arn]
}
