module "api_gateway" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "~> 6.0"

  name          = var.name
  description   = var.description
  protocol_type = "HTTP"

  # The module's built-in certificate and records assume Route53; DNS is on Cloudflare.
  # Setting a domain also disables the default execute-api endpoint.
  domain_name                 = var.domain_name
  domain_name_certificate_arn = aws_acm_certificate_validation.this.certificate_arn
  create_certificate          = false
  create_domain_records       = false

  routes = {
    "$default" = {
      integration = {
        uri                    = var.function_arn
        payload_format_version = "2.0"
      }
    }
  }

  stage_access_log_settings = null
  stage_default_route_settings = {
    detailed_metrics_enabled = false
    throttling_burst_limit   = var.throttling_burst_limit
    throttling_rate_limit    = var.throttling_rate_limit
  }
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${module.api_gateway.api_execution_arn}/*/*"
}

resource "aws_acm_certificate" "this" {
  domain_name       = var.domain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

data "cloudflare_zone" "this" {
  filter = {
    name = var.cloudflare_zone
  }
}

resource "cloudflare_dns_record" "acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.this.domain_validation_options : dvo.domain_name => dvo
  }

  zone_id = data.cloudflare_zone.this.zone_id
  name    = trimsuffix(each.value.resource_record_name, ".")
  type    = each.value.resource_record_type
  content = trimsuffix(each.value.resource_record_value, ".")
  ttl     = 60
  proxied = false
  comment = "ACM validation for ${var.name}"
}

resource "aws_acm_certificate_validation" "this" {
  certificate_arn         = aws_acm_certificate.this.arn
  validation_record_fqdns = [for r in cloudflare_dns_record.acm_validation : r.name]
}

resource "cloudflare_dns_record" "api" {
  zone_id = data.cloudflare_zone.this.zone_id
  name    = var.domain_name
  type    = "CNAME"
  content = module.api_gateway.domain_name_target_domain_name
  ttl     = 1
  proxied = true
  comment = "API Gateway for ${var.name}"
}
