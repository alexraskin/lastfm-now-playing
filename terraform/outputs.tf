output "url" {
  value = module.http_api.url
}

output "api_gateway_target" {
  value = module.http_api.api_gateway_target
}

output "ssm_parameter_name" {
  value = module.lambda_function.ssm_parameter_name
}

output "deploy_role_arn" {
  description = "Set as the AWS_ROLE_ARN GitHub Actions variable."
  value       = module.github_oidc.role_arn
}
