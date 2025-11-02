terraform {
  required_providers {
    validatefx = {
      source = "The-DevOps-Daily/validatefx"
    }
  }
}

provider "validatefx" {}

locals {
  urls = [
    "https://example.com",
    "ftp://example.com",
  ]

  results = [
    for value in local.urls : {
      url   = value
      valid = provider::validatefx::url(value)
    }
  ]
}

output "url_validation_results" {
  value = local.results
}
