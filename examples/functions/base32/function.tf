locals {
  values = ["NBSWY3DP", "OZQWY2LEMF2GK==="]
}

output "base32_example" {
  value = [
    for v in local.values : {
      value = v
      valid = provider::validatefx::base32(v)
    }
  ]
}

