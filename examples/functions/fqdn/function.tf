locals {
  names = ["example.com", "app.prod.example.com"]
}

output "fqdn_example" {
  value = [
    for n in local.names : {
      value = n
      valid = provider::validatefx::fqdn(n)
    }
  ]
}

