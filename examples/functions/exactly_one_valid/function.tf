locals {
  checks = [
    provider::validatefx::email("alice@example.com"),
    provider::validatefx::email("invalid"),
    provider::validatefx::email("bob@example.com"),
  ]

  exactly_one_valid = provider::validatefx::exactly_one_valid(local.checks)
}

output "exactly_one_valid_example" {
  value = local.exactly_one_valid
}
