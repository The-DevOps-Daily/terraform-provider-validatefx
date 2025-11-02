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
  versions = {
    stable     = "1.2.3"
    prefixed   = "v2.0.0"
    prerelease = "3.1.0-beta.1"
    build      = "1.0.0+build.5"
    invalid    = "1.0"
  }

  results = {
    for name, value in local.versions : name => provider::validatefx::semver(value)
  }

  enforce_release = provider::validatefx::assert(
    provider::validatefx::semver(local.versions.stable),
    "release version must be valid semver"
  )
}

output "semver_validation" {
  value = local.results
}

output "assert_semver" {
  value = local.enforce_release
}
