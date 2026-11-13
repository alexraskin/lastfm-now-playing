output "url" {
  value = "https://${var.domain_name}"
}

output "api_gateway_target" {
  value = module.api_gateway.domain_name_target_domain_name
}
