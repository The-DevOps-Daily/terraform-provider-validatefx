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
  # Example 1: Only one deployment method specified
  deployment_config1 = ["", "docker", ""] # Only docker is set
  valid1             = provider::validatefx::mutually_exclusive(local.deployment_config1)

  # Example 2: Only one authentication method specified
  auth_config = ["password", "", ""] # Only password is set
  valid2      = provider::validatefx::mutually_exclusive(local.auth_config)

  # Example 3: Only one database backend specified
  db_backends = ["", "", "postgres"] # Only postgres is set
  valid3      = provider::validatefx::mutually_exclusive(local.db_backends)

  # Example 4: Only one storage option specified
  storage_options = ["", "s3", "", ""] # Only s3 is set
  valid4          = provider::validatefx::mutually_exclusive(local.storage_options)
}

output "validation_results" {
  value = {
    deployment_valid = local.valid1
    auth_valid       = local.valid2
    db_valid         = local.valid3
    storage_valid    = local.valid4
  }
}
