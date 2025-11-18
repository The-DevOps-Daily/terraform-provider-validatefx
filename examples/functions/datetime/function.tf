terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = ">= 0.1.0"
    }
  }
}

provider "validatefx" {}

locals {
  # Valid datetime strings in different formats
  inputs = {
    default_valid1 = "2025-11-02T15:04:05Z"
    default_valid2 = "2025-01-15T10:30:00Z"
    custom_valid = {
      value  = "2025-11-02 15:04:05"
      layout = "2006-01-02 15:04:05"
    }
  }

  # Validate with default RFC3339 format
  default_check1 = provider::validatefx::datetime(local.inputs.default_valid1)
  default_check2 = provider::validatefx::datetime(local.inputs.default_valid2)

  # Validate with custom layout
  custom_check = provider::validatefx::datetime(
    local.inputs.custom_valid.value,
    [local.inputs.custom_valid.layout],
  )
}

output "datetime_checks" {
  description = "Datetime validation results"
  value = {
    default_valid1 = local.default_check1
    default_valid2 = local.default_check2
    custom_valid   = local.custom_check
  }
}
