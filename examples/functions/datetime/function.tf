terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = ">= 0.1.0"
    }
  }
}

provider "validatefx" {
  datetime_layouts = ["2006-01-02 15:04:05"]
}

locals {
  inputs = {
    default_valid   = "2025-11-02T15:04:05Z"
    default_invalid = "2025-13-02T15:04:05Z"
    custom_valid    = "2025-11-02 15:04:05"
  }

  default_check = provider::validatefx::datetime(local.inputs.default_valid)
  default_fail  = provider::validatefx::datetime(local.inputs.default_invalid)

  custom_check = provider::validatefx::datetime(local.inputs.custom_valid)
}

output "datetime_checks" {
  value = {
    default_valid   = local.default_check
    default_invalid = local.default_fail
    custom_valid    = local.custom_check
  }
}
