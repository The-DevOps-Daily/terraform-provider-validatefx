# 🧩 Terraform Provider - ValidateFX

[![Go Version](https://img.shields.io/github/go-mod/go-version/The-DevOps-Daily/terraform-provider-validatefx?style=flat-square)](https://go.dev/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/The-DevOps-Daily/terraform-provider-validatefx/ci.yml?branch=main&style=flat-square)](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/actions)
[![License](https://img.shields.io/github/license/The-DevOps-Daily/terraform-provider-validatefx?style=flat-square)](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/blob/main/LICENSE)
[![Terraform Registry](https://img.shields.io/badge/terraform-registry-623CE4?style=flat-square&logo=terraform)](https://registry.terraform.io/providers/The-DevOps-Daily/validatefx/latest)

Reusable validation functions for Terraform, built with the latest [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

ValidateFX lets you write cleaner, more expressive validations using functions like `email`, `uuid`, `base64`, and more. Use the `assert` function to validate conditions with custom error messages.

---

## 🚀 Example

```hcl
terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = "0.1.0"
    }
  }
}

provider "validatefx" {}

variable "email" {
  type = string
}

locals {
  # Validate email with custom error message
  email_check = provider::validatefx::assert(
    provider::validatefx::email(var.email),
    "Invalid email address provided!"
  )

  # Or use in variable validation
  age_validation = provider::validatefx::assert(
    var.user_age >= 18,
    "User must be at least 18 years old!"
  )
}
```

---

## ⚙️ Development

```bash
git clone https://github.com/The-DevOps-Daily/terraform-provider-validatefx.git
cd terraform-provider-validatefx
go mod tidy
make build
make install
make dev
```

Example usage in `examples/basic/main.tf`.

- [OS-specific installation & troubleshooting](docs/os-installation.md) — instructions for Windows, macOS (Intel/ARM), and Linux.


---

## 🧩 Available Functions

| Function | Description |
| -------------------------- | ------------------------------------------------ |
| `all_valid` | Return true when all provided validation checks evaluate to true. |
| `any_valid` | Return true when any provided validation check evaluates to true. |
| `assert` | Assert a condition with a custom error message. |
| `base64` | Validate that a string is Base64 encoded. |
| `between` | Validate that a numeric string falls between inclusive minimum and maximum bounds. |
| `cidr` | Validate that a string is an IPv4 or IPv6 CIDR block. |
| `credit_card` | Validate that a string is a credit card number using the Luhn algorithm. |
| `datetime` | Validate that a string is an ISO 8601 / RFC 3339 datetime. |
| `domain` | Validate that a string is a compliant domain name. |
| `email` | Validate that a string is an RFC 5322 compliant email address. |
| `exactly_one_valid` | Return true when exactly one validation check evaluates to true. |
| `hostname` | Validate that a string is a hostname compliant with RFC 1123. |
| `in_list` | Validate that a string matches one of the allowed values. |
| `ip` | Validate that a string is a valid IPv4 or IPv6 address. |
| `json` | Validate that a string decodes to a JSON object. |
| `mac_address` | Validate that a string is a MAC address in colon, dash, or compact format. |
| `matches_regex` | Validate that a string matches a provided regular expression. |
| `phone` | Validate that a string is an E.164 compliant phone number. |
| `semver` | Validate that a string follows Semantic Versioning (SemVer 2.0.0). |
| `string_length` | Validate that a string length falls within optional minimum and maximum bounds. |
| `url` | Validate that a string is an HTTP(S) URL. |
| `uuid` | Validate that a string is an RFC 4122 UUID (versions 1-5). |
| `version` | Return the provider version string. |


---

## 💡 Contributing

Open to PRs! Browse our [good first issues](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/issues?q=is%3Aopen+label%3A"good+first+issue") to get started.

---

## 📜 License

MIT © 2025 [DevOps Daily](https://github.com/The-DevOps-Daily)

## Thanks to all contributors ❤

[![Contributors](https://contrib.rocks/image?repo=The-DevOps-Daily/terraform-provider-validatefx)](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/graphs/contributors)
