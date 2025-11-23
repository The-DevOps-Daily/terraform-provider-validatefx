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
  mime_types = [
    "application/json",
    "text/html",
    "image/png",
    "video/mp4",
    "application/vnd.api+json",
    "image/svg+xml",
    "text/plain; charset=utf-8",
  ]

  checked_mime_types = [
    for mime_type in local.mime_types : {
      type  = mime_type
      valid = provider::validatefx::mime_type(mime_type)
    }
  ]
}

output "checked_mime_types" {
  value = local.checked_mime_types
}
