terraform {
  required_providers {
    validatefx = {
      source = "the-devops-daily/validatefx"
    }
  }
}

locals {
  # Valid Azure locations
  valid_us_east   = provider::validatefx::azure_location("eastus")
  valid_europe    = provider::validatefx::azure_location("westeurope")
  valid_asia      = provider::validatefx::azure_location("southeastasia")
  valid_australia = provider::validatefx::azure_location("australiaeast")
  valid_gov_cloud = provider::validatefx::azure_location("usgovvirginia")
}

output "azure_location_checks" {
  value = {
    us_east   = local.valid_us_east
    europe    = local.valid_europe
    asia      = local.valid_asia
    australia = local.valid_australia
    gov_cloud = local.valid_gov_cloud
  }
}
