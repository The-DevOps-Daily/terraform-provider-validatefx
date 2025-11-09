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
  usernames = [
    "alice",
    "bob_123",
    "invalid-user",
  ]

  results = [
    for username in local.usernames : {
      value = username
      valid = provider::validatefx::username(username)
    }
  ]
}

output "username_example" {
  value = local.results
}
