locals {
  subnets = [
    "192.168.1.0/24",
    "10.0.0.0/8",
    "2001:db8::/64",
  ]
}

output "subnet_example" {
  value = [
    for s in local.subnets : {
      subnet = s
      valid  = provider::validatefx::subnet(s)
    }
  ]
}

