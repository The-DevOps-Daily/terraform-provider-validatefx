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
  valid_base64   = provider::validatefx::base64("U29sdmVkIQ==")
  invalid_base64 = provider::validatefx::base64("invalid base64")
}

output "base64_results" {
  value = {
    valid   = local.valid_base64
    invalid = local.invalid_base64
  }
}
