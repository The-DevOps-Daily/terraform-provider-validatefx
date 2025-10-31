terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = "~> 0.1"
    }
  }
}

provider "validatefx" {}

locals {
  mac_addresses = [
    "00:1A:2B:3C:4D:5E",
    "00-1A-2B-3C-4D-5E",
    "001A2B3C4D5E",
    "not-a-mac",
  ]

  mac_validation = [
    for addr in local.mac_addresses : {
      address = addr
      valid   = provider::validatefx::mac_address(addr)
    }
  ]
}

output "mac_validation" {
  value = local.mac_validation
}
