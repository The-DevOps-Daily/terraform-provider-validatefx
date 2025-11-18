terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = ">= 0.0.1"
    }
  }
}

provider "validatefx" {}

locals {
  emails = [
    "alice@example.com",
    "bob@example.org",
    "charlie+tag@test.com",
  ]

  checked_emails = [
    for email in local.emails : {
      address = email
      valid   = provider::validatefx::email(email)
    }
  ]
}

output "checked_emails" {
  value = local.checked_emails
}
