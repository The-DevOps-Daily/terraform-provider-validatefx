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
  environments = [
    "tf-production",
    "TF-staging",
    "dev",
  ]

  results = [
    for name in local.environments : {
      value       = name
      ignore_case = true
      valid       = provider::validatefx::has_prefix(name, ["tf-", "iac-"], true)
    }
  ]
}

output "has_prefix_example" {
  value = local.results
}

