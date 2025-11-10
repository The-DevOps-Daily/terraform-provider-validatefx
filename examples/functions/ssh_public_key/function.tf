locals {
  keys = [
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKJf0N0nH7kz5Zr4xkz0GWWJrPq9uO2m6sR3j0s8v2QG test@example",
  ]
}

output "ssh_public_key_example" {
  value = [
    for k in local.keys : {
      value = k
      valid = provider::validatefx::ssh_public_key(k)
    }
  ]
}

