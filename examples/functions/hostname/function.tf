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
  hostnames = [
    "service.internal",
    "xn--bcher-kva.example",
    "example.com.",
    "bad_name",
  ]

  hostname_checks = [
    for host in local.hostnames : {
      hostname = host
      valid    = provider::validatefx::hostname(host)
    }
  ]
}

output "hostname_validation" {
  value = local.hostname_checks
}
