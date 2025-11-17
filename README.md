# 🧩 Terraform Provider - ValidateFX

[![Go Version](https://img.shields.io/github/go-mod/go-version/The-DevOps-Daily/terraform-provider-validatefx?style=flat-square)](https://go.dev/)
[![CI](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/actions/workflows/ci.yml/badge.svg)](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/actions/workflows/ci.yml)
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


---

## 🧩 Available Functions

| Function | Description |
| -------------------------- | ------------------------------------------------ |
| `all_valid` | Return true when all provided validation checks evaluate to true. |
| `any_valid` | Return true when any provided validation check evaluates to true. |
| `assert` | Assert a condition with a custom error message. |
| `base32` | Validate that a string is Base32 encoded. |
| `base64` | Validate that a string is Base64 encoded. |
| `between` | Validate that a numeric string falls between inclusive minimum and maximum bounds. |
| `cidr` | Validate that a string is an IPv4 or IPv6 CIDR block. |
| `cidr_overlap` | Validate that provided CIDR blocks do not overlap. |
| `credit_card` | Validate that a string is a credit card number using the Luhn algorithm. |
| `datetime` | Validate that a string is an ISO 8601 / RFC 3339 datetime. |
| `domain` | Validate that a string is a compliant domain name. |
| `email` | Validate that a string is an RFC 5322 compliant email address. |
| `exactly_one_valid` | Return true when exactly one validation check evaluates to true. |
| `fqdn` | Validate that a string is a fully qualified domain name (FQDN). |
| `has_prefix` | Validate that a string starts with one of the provided prefixes. |
| `has_suffix` | Validate that a string ends with one of the provided suffixes. |
| `hex` | Validate that a string contains only hexadecimal characters. |
| `hostname` | Validate that a string is a hostname compliant with RFC 1123. |
| `in_list` | Validate that a string matches one of the allowed values. |
| `integer` | Validate that a string represents a valid integer. |
| `ip` | Validate that a string is a valid IPv4 or IPv6 address. |
| `json` | Validate that a string decodes to a JSON object. |
| `jwt` | Validate that a string is a well-formed JSON Web Token (JWT). |
| `list_subset` | Validate that all elements of a list/set are contained in a reference list. |
| `mac_address` | Validate that a string is a MAC address in colon, dash, or compact format. |
| `matches_regex` | Validate that a string matches a provided regular expression. |
| `not_in_list` | Validate that a string does not match any of the provided disallowed values. |
| `password_strength` | Checks if a password meets strength requirements |
| `phone` | Validate that a string is an E.164 compliant phone number. |
| `port_number` | Validate that a string is a valid TCP/UDP port number (1..65535). |
| `port_range` | Validate that a string is a valid port range (start-end). |
| `private_ip` | Validate that an IP address is private (RFC1918 / IPv6 ULA). |
| `semver` | Validate that a string follows Semantic Versioning (SemVer 2.0.0). |
| `semver_range` | Validate that a string is a valid semantic version range expression. |
| `set_equals` | Validate that two string lists contain the same elements regardless of order. |
| `ssh_public_key` | Validate that a string is a valid SSH public key. |
| `string_contains` | Validate that a string contains at least one of the provided substrings. |
| `string_length` | Validate that a string length falls within optional minimum and maximum bounds. |
| `subnet` | Validate that a string is a subnet address (IP equals network) in CIDR notation. |
| `uri` | Validate that a string is a URI. |
| `url` | Validate that a string is an HTTP(S) URL. |
| `username` | Validate that a string is a valid username. |
| `uuid` | Validate that a string is an RFC 4122 UUID (versions 1-5). |
| `version` | Return the provider version string. |


---

## 💡 Contributing

Open to PRs! Browse our [good first issues](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/issues?q=is%3Aopen+label%3A"good+first+issue") to get started.

### Fuzz Tests

Go ships with native fuzzing support. We include fuzz tests for core string validators (email, URL, JSON) to harden them against edge cases.

- Run a targeted fuzz (10s) for email:
  - `go test ./internal/validators -run FuzzEmailValidator -fuzz FuzzEmailValidator -fuzztime=10s`
- For URL or JSON, replace the function name accordingly (e.g., `FuzzURLValidator`, `FuzzJSONValidator`).
- To fuzz all in the package for 1 minute:
  - `go test ./internal/validators -fuzz Fuzz -fuzztime=1m`

When the fuzzer finds a failure, it writes a minimizing corpus entry that becomes part of future test runs.

---

## 📜 License

MIT © 2025 [DevOps Daily](https://github.com/The-DevOps-Daily)

## Thanks to all contributors ❤

[![Contributors](https://contrib.rocks/image?repo=The-DevOps-Daily/terraform-provider-validatefx)](https://github.com/The-DevOps-Daily/terraform-provider-validatefx/graphs/contributors)
