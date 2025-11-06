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
  phrases = [
    "Hello, Terraform!",
    "ValidateFX is great",
    "Infrastructure as Code",
  ]

  results = [
    for phrase in local.phrases : {
      value = phrase
      valid = provider::validatefx::string_contains(phrase, ["Terraform", "ValidateFX"], true)
    }
  ]
}

output "string_contains_example" {
  value = local.results
}
