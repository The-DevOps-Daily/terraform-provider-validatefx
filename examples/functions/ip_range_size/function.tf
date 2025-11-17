locals {
  # Allow between /8 and /28 for IPv4
  ranges = [
    {
      cidr = "10.0.0.0/16"
      min  = 8
      max  = 28
      ok   = provider::validatefx::ip_range_size("10.0.0.0/16", 8, 28)
    },
  ]
}

output "ip_range_size_examples" {
  value = local.ranges
}

