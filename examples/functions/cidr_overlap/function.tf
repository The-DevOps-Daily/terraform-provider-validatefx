locals {
  cidr_overlap_checks = {
    non_overlap = provider::validatefx::cidr_overlap(["10.0.0.0/24", "10.0.1.0/24"])
  }
}

output "cidr_overlap" {
  value = local.cidr_overlap_checks
}

