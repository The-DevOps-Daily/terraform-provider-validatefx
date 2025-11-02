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
  checks = [
    provider::validatefx::email("user@example.com"),
    provider::validatefx::uuid("d9428888-122b-11e1-b85c-61cd3cbb3210"),
    provider::validatefx::credit_card("4532015112830366"),
  ]

  all_checks_pass = provider::validatefx::all_valid(local.checks)
}

output "all_pass" {
  value = local.all_checks_pass
}
