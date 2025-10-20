# 🧩 Terraform Provider - ValidateFX

[![Go Version](https://img.shields.io/github/go-mod/go-version/The-DevOps-Daily/terraform-provider-validatefx?style=flat-square)](https://go.dev/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/The-DevOps-Daily/terraform-provider-validatefx/test.yml?branch=main&style=flat-square)](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/actions)
[![License](https://img.shields.io/github/license/The-DevOps-Daily/terraform-provider-validatefx?style=flat-square)](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/blob/main/LICENSE)
[![Terraform Registry](https://img.shields.io/badge/terraform-registry-623CE4?style=flat-square&logo=terraform)](https://registry.terraform.io/providers/thedevopsdaily/validatefx/latest)

Reusable validation functions for Terraform, built with the latest [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

ValidateFX lets you write cleaner, more expressive validations using functions like `email`, `uuid`, `base64`, and more. Use the `assert` function to validate conditions with custom error messages.

---

## 🚀 Example

```hcl
terraform {
  required_providers {
    validatefx = {
      source  = "thedevopsdaily/validatefx"
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

---

## 🧩 Available Functions

| Function | Description |
| -------------------------- | ------------------------------------------------ |
| `assert(bool, string)` | Validates a condition with a custom error message |
| `email(string)` | Validates email format (RFC 5322) |
| `uuid(string)` | Validates UUID (RFC 4122, versions 1-5) |
| `base64(string)` | Validates Base64 encoding |
| `credit_card(string)` | Validates credit card number (Luhn algorithm) |
| `domain(string)` | Validates domain name (RFC 1123/952) |

---

## 💡 Contributing

Open to PRs! Good first issues include adding new validators like `is_ip`, `is_hostname`, or `matches_regex`.

---

## 📜 License

MIT © 2025 [DevOps Daily](https://github.com/The-DevOps-Daily)
