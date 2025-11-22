terraform {
  required_providers {
    validatefx = {
      source = "the-devops-daily/validatefx"
    }
  }
}

locals {
  # Valid GCP zones
  valid_us_central_a    = provider::validatefx::gcp_zone("us-central1-a")
  valid_us_east_b       = provider::validatefx::gcp_zone("us-east1-b")
  valid_europe_west_b   = provider::validatefx::gcp_zone("europe-west1-b")
  valid_asia_east_a     = provider::validatefx::gcp_zone("asia-east1-a")
  valid_australia_a     = provider::validatefx::gcp_zone("australia-southeast1-a")
  valid_africa_a        = provider::validatefx::gcp_zone("africa-south1-a")
}

output "gcp_zone_checks" {
  value = {
    us_central_a    = local.valid_us_central_a
    us_east_b       = local.valid_us_east_b
    europe_west_b   = local.valid_europe_west_b
    asia_east_a     = local.valid_asia_east_a
    australia_a     = local.valid_australia_a
    africa_a        = local.valid_africa_a
  }
}
