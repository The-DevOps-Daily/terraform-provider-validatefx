locals {
  addrs = [
    "8.8.8.8",
    "1.1.1.1",
  ]
}

output "public_ip_example" {
  value = [
    for a in local.addrs : {
      ip    = a
      valid = provider::validatefx::public_ip(a, false, false)
    }
  ]
}
