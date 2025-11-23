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
  sizes = [
    { value = "5", min = "1", max = "10" },
    { value = "1", min = "1", max = "10" },
    { value = "10", min = "1", max = "10" },
    { value = "0.5", min = "0", max = "1" },
  ]

  checked_sizes = [
    for size in local.sizes : {
      value = size.value
      min   = size.min
      max   = size.max
      valid = provider::validatefx::size_between(size.value, size.min, size.max)
    }
  ]
}

output "checked_sizes" {
  value = local.checked_sizes
}
