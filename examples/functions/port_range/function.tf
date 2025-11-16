locals {
  ranges = [
    "80-8080",
    "0-65535",
  ]
}

output "port_range_example" {
  value = [
    for r in local.ranges : {
      range = r
      valid = provider::validatefx::port_range(r)
    }
  ]
}

