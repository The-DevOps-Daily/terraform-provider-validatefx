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
  # Example 1: All keys allowed from a list
  config1 = {
    name = "server-01"
    port = "8080"
  }
  valid1 = provider::validatefx::map_keys_match(
    local.config1,
    ["name", "port", "protocol"],
    []
  )

  # Example 2: Required keys must be present
  config2 = {
    name = "server-02"
    port = "9090"
  }
  valid2 = provider::validatefx::map_keys_match(
    local.config2,
    ["name", "port", "protocol"],
    ["name", "port"]
  )

  # Example 3: Empty allowed list means all keys allowed
  config3 = {
    name  = "server-03"
    port  = "3000"
    extra = "value"
  }
  valid3 = provider::validatefx::map_keys_match(
    local.config3,
    [],
    ["name"]
  )
}

output "validation_results" {
  value = {
    config1_valid = local.valid1
    config2_valid = local.valid2
    config3_valid = local.valid3
  }
}
