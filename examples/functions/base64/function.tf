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
  # Valid Base64 encoded strings
  test_strings = [
    "U29sdmVkIQ==",     # "Solved!"
    "SGVsbG8gV29ybGQh", # "Hello World!"
    "VGVycmFmb3Jt",     # "Terraform"
  ]

  # Validate each string
  validation_results = [
    for str in local.test_strings : {
      value = str
      valid = provider::validatefx::base64(str)
    }
  ]
}

output "base64_results" {
  description = "Base64 validation results for test strings"
  value       = local.validation_results
}
