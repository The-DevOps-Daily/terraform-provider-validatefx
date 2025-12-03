terraform {
  required_providers {
    validatefx = {
      source = "the-devops-daily/validatefx"
    }
  }
}

# Valid Kubernetes label values
output "valid_simple_value" {
  value = provider::validatefx::k8s_label_value("production")
}

output "valid_empty_value" {
  value = provider::validatefx::k8s_label_value("")
}

output "valid_value_with_dash" {
  value = provider::validatefx::k8s_label_value("prod-env")
}

output "valid_value_with_underscore" {
  value = provider::validatefx::k8s_label_value("app_v1")
}

output "valid_value_with_mixed_case" {
  value = provider::validatefx::k8s_label_value("Production")
}

# Validation will fail for invalid values:
# - provider::validatefx::k8s_label_value(strings.repeat("a", 64))  # Too long
# - provider::validatefx::k8s_label_value("prod@env")              # Invalid characters
# - provider::validatefx::k8s_label_value("-production")           # Starts with dash
# - provider::validatefx::k8s_label_value("production-")           # Ends with dash
