output "function_name" {
  value = module.lambda.lambda_function_name
}

output "function_arn" {
  value = module.lambda.lambda_function_arn
}

output "ssm_parameter_name" {
  value = aws_ssm_parameter.lastfm_api_key.name
}
