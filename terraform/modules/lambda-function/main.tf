# Standard-tier SecureString parameters with the aws/ssm key are free, unlike Secrets Manager.
# Terraform only creates the parameter; set the real value out of band (see the README).
resource "aws_ssm_parameter" "lastfm_api_key" {
  name  = "/${var.name}/lastfm-api-key"
  type  = "SecureString"
  value = "changeme"

  lifecycle {
    ignore_changes = [value]
  }
}

module "lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 8.0"

  function_name = var.name
  description   = var.description
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  architectures = ["arm64"]
  memory_size   = 128
  timeout       = 10

  reserved_concurrent_executions = var.reserved_concurrency

  create_package         = false
  local_existing_package = var.package_path

  environment_variables = {
    LASTFM_API_KEY_SSM_PARAMETER = aws_ssm_parameter.lastfm_api_key.name
  }

  attach_policy_statements = true
  policy_statements = {
    ssm = {
      effect    = "Allow"
      actions   = ["ssm:GetParameter"]
      resources = [aws_ssm_parameter.lastfm_api_key.arn]
    }
  }

  cloudwatch_logs_retention_in_days = 14
}
