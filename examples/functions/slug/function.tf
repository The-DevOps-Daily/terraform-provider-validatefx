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
  slugs = [
    "my-application",
    "web-server-01",
    "api-v2",
    "user-profile",
    "data-pipeline-prod",
  ]

  checked_slugs = [
    for slug in local.slugs : {
      slug  = slug
      valid = provider::validatefx::slug(slug)
    }
  ]
}

output "checked_slugs" {
  value = local.checked_slugs
}
