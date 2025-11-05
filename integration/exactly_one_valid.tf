locals {
  email_checks = [
    provider::validatefx::email("alice@example.com"),
    provider::validatefx::email("bob@example.com"),
    provider::validatefx::email("not-an-email"),
  ]

  exactly_one_valid_results = {
    valid_when_one_true = provider::validatefx::exactly_one_valid([
      provider::validatefx::email("alice@example.com"),
      provider::validatefx::email("invalid"),
      provider::validatefx::email("bob@example.com"),
    ])

    false_when_multiple_true = provider::validatefx::exactly_one_valid([
      provider::validatefx::email("alice@example.com"),
      provider::validatefx::email("bob@example.com"),
    ])

    false_when_none_true = provider::validatefx::exactly_one_valid([
      provider::validatefx::email("invalid"),
      provider::validatefx::email("bad"),
    ])

    unknown_when_unknown_present = provider::validatefx::exactly_one_valid([
      provider::validatefx::email("alice@example.com"),
      provider::validatefx::email(var.unknown_email),
      provider::validatefx::email("invalid"),
    ])
  }
}

variable "unknown_email" {
  description = "Optional email value that may be unknown"
  type        = string
  default     = null
}

output "validatefx_exactly_one_valid" {
  value = local.exactly_one_valid_results
}
