terraform {
  required_providers {
    validatefx = {
      source = "the-devops-daily/validatefx"
    }
  }
}

locals {
  # Valid AWS regions
  valid_us_east    = provider::validatefx::aws_region("us-east-1")
  valid_eu_west    = provider::validatefx::aws_region("eu-west-1")
  valid_ap_south   = provider::validatefx::aws_region("ap-southeast-1")
  valid_gov_region = provider::validatefx::aws_region("us-gov-west-1")
  valid_cn_region  = provider::validatefx::aws_region("cn-north-1")
}

output "aws_region_checks" {
  value = {
    us_east    = local.valid_us_east
    eu_west    = local.valid_eu_west
    ap_south   = local.valid_ap_south
    gov_region = local.valid_gov_region
    cn_region  = local.valid_cn_region
  }
}
