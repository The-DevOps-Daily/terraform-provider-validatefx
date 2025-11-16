locals {
  ips = [
    "10.0.0.1",
    "172.16.0.1",
    "192.168.1.1",
    "fd00::1",
  ]
}

output "private_ip_example" {
  value = [
    for ip in local.ips : {
      ip    = ip
      valid = provider::validatefx::private_ip(ip)
    }
  ]
}

