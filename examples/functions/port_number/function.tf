locals {
  ports = ["80", "443", "65535"]
}

output "port_number_example" {
  value = [
    for p in local.ports : {
      port  = p
      valid = provider::validatefx::port_number(p)
    }
  ]
}

