terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = "0.0.1"
    }
  }
}

provider "validatefx" {}

locals {
  defaults = [
    "feature_a",
    "feature_b",
    "feature_c",
  ]

  overrides = [
    "feature_c",
    "feature_b",
    "feature_a",
  ]

  validation = provider::validatefx::set_equals(local.defaults, local.overrides)
}

output "set_equals_example" {
  value = local.validation
}

