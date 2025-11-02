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
    provider::validatefx::credit_card("0000000000000000"),
  ]

  some_checks_pass = provider::validatefx::any_valid(local.checks)
}

output "any_pass" {
  value = local.some_checks_pass
}
