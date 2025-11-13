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
  ranges = {
    single_ge      = ">=1.2.3"
    bounded_open   = ">=1.2.3, <2.0.0"
    with_v_prefix  = ">=v1.0.0,<=v1.5.0"
  }

  results = {
    for name, expr in local.ranges : name => provider::validatefx::semver_range(expr)
  }
}

output "semver_range_validation" {
  value = local.results
}

