locals {
  values = ["deadbeef", "CAFE1234"]
}

output "hex_example" {
  value = [
    for v in local.values : {
      value = v
      valid = provider::validatefx::hex(v)
    }
  ]
}

