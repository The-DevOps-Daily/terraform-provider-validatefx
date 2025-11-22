terraform {
  required_providers {
    validatefx = {
      source = "the-devops-daily/validatefx"
    }
  }
}

locals {
  # Valid GCP regions
  valid_us_central    = provider::validatefx::gcp_region("us-central1")
  valid_europe        = provider::validatefx::gcp_region("europe-west1")
  valid_asia          = provider::validatefx::gcp_region("asia-southeast1")
  valid_australia     = provider::validatefx::gcp_region("australia-southeast1")
  valid_south_america = provider::validatefx::gcp_region("southamerica-east1")
}

output "gcp_region_checks" {
  value = {
    us_central    = local.valid_us_central
    europe        = local.valid_europe
    asia          = local.valid_asia
    australia     = local.valid_australia
    south_america = local.valid_south_america
  }
}
