terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = "0.0.1"
    }
  }
}

provider "validatefx" {}

locals {
  emails = [
    "alice@example.com",
  ]

  uuids = [
    "d9428888-122b-11e1-b85c-61cd3cbb3210",
  ]

  base64_values = [
    "U29sdmVkIQ==",
  ]

  credit_cards = [
    "4532015112830366",
  ]

  phone_numbers = [
    "+14155552671",
  ]

  mac_addresses = [
    "00:1A:2B:3C:4D:5E",
    "00-1A-2B-3C-4D-5E",
    "001A2B3C4D5E",
  ]

  url_values = [
    "https://example.com",
  ]

  domains = [
    "example.com",
  ]

  hostnames = [
    "service.internal",
    "xn--bcher-kva.example",
  ]

  fqdn_idna_checks = [
    {
      label = "punycode label"
      value = "xn--bcher-kva.example"
      valid = provider::validatefx::fqdn("xn--bcher-kva.example")
    }
  ]

  json_payloads = [
    "{\"key\": \"value\"}",
  ]

  semver_values = [
    "1.0.0",
  ]

  semver_range_values = [
    ">=1.2.3",
    ">=1.2.3, <2.0.0",
  ]

  datetime_values = [
    {
      value   = "2025-11-02T15:04:05Z"
      layouts = []
    },
    {
      value   = "2025-11-02 15:04:05"
      layouts = ["2006-01-02 15:04:05"]
    },
  ]

  ip_values = [
    "127.0.0.1",
    "::1",
  ]

  private_ip_checks = [
    {
      label = "rfc1918"
      value = "10.0.0.1"
      valid = provider::validatefx::private_ip("10.0.0.1")
    },
    {
      label = "ipv6-ula"
      value = "fd00::1"
      valid = provider::validatefx::private_ip("fd00::1")
    }
  ]

  regex_samples = [
    {
      value   = "user_123"
      pattern = "^[a-z0-9_]+$"
    },
  ]

  string_contains_samples = [
    {
      label       = "matches Terraform"
      value       = "Hello Terraform"
      substrings  = ["Terraform", "ValidateFX"]
      ignore_case = false
    },
    {
      label       = "matches ValidateFX case-insensitive"
      value       = "I love validatefx"
      substrings  = ["Terraform", "ValidateFX"]
      ignore_case = true
    }
  ]

  cidr_values = [
    "10.0.0.0/24",
    "2001:db8::/48",
  ]

  ip_range_size_checks = [
    {
      label = "within bounds"
      value = "10.0.0.0/16"
      valid = provider::validatefx::ip_range_size("10.0.0.0/16", 8, 28)
    }
  ]

  subnet_checks = [
    {
      label = "ipv4 subnet"
      value = "192.168.1.0/24"
      valid = provider::validatefx::subnet("192.168.1.0/24")
    }
  ]

  email_results = [
    for value in local.emails : {
      value = value
      valid = provider::validatefx::email(value)
    }
  ]

  uuid_results = [
    for value in local.uuids : {
      value = value
      valid = provider::validatefx::uuid(value)
    }
  ]

  base64_results = [
    for value in local.base64_values : {
      value = value
      valid = provider::validatefx::base64(value)
    }
  ]

  credit_card_results = [
    for value in local.credit_cards : {
      value = value
      valid = provider::validatefx::credit_card(value)
    }
  ]

  domain_results = [
    for value in local.domains : {
      value = value
      valid = provider::validatefx::domain(value)
    }
  ]

  fqdn_checks = [
    {
      label = "valid fqdn"
      value = "app.example.com"
      valid = provider::validatefx::fqdn("app.example.com")
    }
  ]

  hostname_results = [
    for host in local.hostnames : {
      hostname = host
      valid    = provider::validatefx::hostname(host)
    }
  ]

  json_results = [
    for value in local.json_payloads : {
      value = value
      valid = provider::validatefx::json(value)
    }
  ]

  semver_results = [
    for value in local.semver_values : {
      value = value
      valid = provider::validatefx::semver(value)
    }
  ]

  semver_range_results = [
    for expr in local.semver_range_values : {
      value = expr
      valid = provider::validatefx::semver_range(expr)
    }
  ]

  datetime_results = [
    for item in local.datetime_values : {
      value   = item.value
      layouts = item.layouts
      valid   = provider::validatefx::datetime(item.value, item.layouts)
    }
  ]

  ip_results = [
    for value in local.ip_values : {
      value = value
      valid = provider::validatefx::ip(value)
    }
  ]

  public_ip_checks = [
    {
      label = "google"
      value = "8.8.8.8"
      valid = provider::validatefx::public_ip("8.8.8.8", false, false)
    }
  ]

  matches_regex_results = [
    for item in local.regex_samples : {
      value   = item.value
      pattern = item.pattern
      valid   = provider::validatefx::matches_regex(item.value, item.pattern)
    }
  ]

  uri_checks = [
    {
      label = "http"
      value = "http://example.com"
      valid = provider::validatefx::uri("http://example.com")
    },
    {
      label = "urn"
      value = "urn:isbn:0451450523"
      valid = provider::validatefx::uri("urn:isbn:0451450523")
    }
  ]

  arn_checks = [
    {
      label = "iam role"
      value = "arn:aws:iam::123456789012:role/Admin"
      valid = provider::validatefx::arn("arn:aws:iam::123456789012:role/Admin")
    }
  ]

  username_values = [
    "alice",
    "bob_123",
  ]

  username_results = [
    for value in local.username_values : {
      value = value
      valid = provider::validatefx::username(value)
    }
  ]

  string_contains_results = [
    for sample in local.string_contains_samples : {
      label       = sample.label
      value       = sample.value
      substrings  = sample.substrings
      ignore_case = sample.ignore_case
      valid       = provider::validatefx::string_contains(sample.value, sample.substrings, sample.ignore_case)
    }
  ]

  port_range_checks = [
    {
      label = "common web range"
      value = "80-8080"
      valid = provider::validatefx::port_range("80-8080")
    },
    {
      label = "full range"
      value = "0-65535"
      valid = provider::validatefx::port_range("0-65535")
    }
  ]

  port_number_checks = [
    {
      label = "https"
      value = "443"
      valid = provider::validatefx::port_number("443")
    }
  ]

  hex_checks = [
    {
      label = "lowercase"
      value = "deadbeef"
      valid = provider::validatefx::hex("deadbeef")
    },
    {
      label = "uppercase"
      value = "CAFE1234"
      valid = provider::validatefx::hex("CAFE1234")
    }
  ]

  integer_checks = [
    {
      label = "zero"
      value = "0"
      valid = provider::validatefx::integer("0")
    },
    {
      label = "negative"
      value = "-42"
      valid = provider::validatefx::integer("-42")
    },
    {
      label = "plus sign"
      value = "+7"
      valid = provider::validatefx::integer("+7")
    }
  ]

  base32_checks = [
    {
      label = "hello"
      value = "NBSWY3DP"
      valid = provider::validatefx::base32("NBSWY3DP")
    },
    {
      label = "validate"
      value = "OZQWY2LEMF2GK==="
      valid = provider::validatefx::base32("OZQWY2LEMF2GK===")
    }
  ]

  ssh_public_key_checks = [
    {
      label = "ed25519"
      value = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKJf0N0nH7kz5Zr4xkz0GWWJrPq9uO2m6sR3j0s8v2QG test@example"
      valid = provider::validatefx::ssh_public_key("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKJf0N0nH7kz5Zr4xkz0GWWJrPq9uO2m6sR3j0s8v2QG test@example")
    }
  ]

  jwt_checks = [
    {
      label = "valid jwt"
      value = "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
      valid = provider::validatefx::jwt("eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
    }
  ]

  password_strength_checks = [
    {
      label = "strong"
      value = "Abc@1234"
      valid = provider::validatefx::password_strength("Abc@1234")
    }
  ]

  in_list_checks = [
    {
      label       = "valid"
      value       = "draft"
      allowed     = ["draft", "review", "published"]
      ignore_case = false
      valid       = provider::validatefx::in_list("draft", ["draft", "review", "published"], false, null)
    },
    {
      label       = "case-insensitive"
      value       = "Published"
      allowed     = ["draft", "review", "published"]
      ignore_case = true
      valid       = provider::validatefx::in_list("Published", ["draft", "review", "published"], true, null)
    }
  ]

  not_in_list_checks = [
    {
      label       = "not disallowed"
      value       = "green"
      disallowed  = ["red", "blue"]
      ignore_case = false
      valid       = provider::validatefx::not_in_list("green", ["red", "blue"], false)
    },
    {
      label       = "case-insensitive not in list"
      value       = "Green"
      disallowed  = ["red", "blue"]
      ignore_case = true
      valid       = provider::validatefx::not_in_list("Green", ["red", "blue"], true)
    }
  ]

  exactly_one_valid_checks = {
    valid_when_one_true = provider::validatefx::exactly_one_valid([
      provider::validatefx::all_valid([true, true]),
      provider::validatefx::any_valid([false, false]),
      false,
    ])

    false_when_multiple_true = provider::validatefx::exactly_one_valid([
      provider::validatefx::all_valid([true, true]),
      provider::validatefx::any_valid([true]),
    ])

    false_when_none_true = provider::validatefx::exactly_one_valid([
      provider::validatefx::any_valid([]),
      provider::validatefx::all_valid([true, false]),
      false,
    ])

    unknown_when_unknown_present = provider::validatefx::exactly_one_valid([
      provider::validatefx::email(var.integration_unknown_email),
      false,
    ])
  }

  in_list_integration_checks = [
    {
      label       = "valid"
      value       = "draft"
      ignore_case = false
      valid       = provider::validatefx::in_list("draft", ["draft", "review", "published"], false, null)
    },
    {
      label       = "case-insensitive"
      value       = "Published"
      ignore_case = true
      valid       = provider::validatefx::in_list("Published", ["draft", "review", "published"], true, null)
    }
  ]

  has_suffix_checks = [
    {
      label    = "yaml configuration"
      value    = "config.yaml"
      suffixes = [".yaml", ".yml"]
      valid    = provider::validatefx::has_suffix("config.yaml", [".yaml", ".yml"])
    },
    {
      label    = "text notes"
      value    = "notes.txt"
      suffixes = [".log", ".txt"]
      valid    = provider::validatefx::has_suffix("notes.txt", [".log", ".txt"])
    }
  ]

  list_subset_checks = [
    {
      label = "roles subset"
      roles = ["reader", "writer"]
      allow = ["reader", "writer", "admin"]
      valid = provider::validatefx::list_subset(["reader", "writer"], ["reader", "writer", "admin"])
    }
  ]

  has_prefix_checks = [
    {
      label       = "terraform prefix"
      value       = "tf-backend"
      prefixes    = ["tf-", "iac-"]
      ignore_case = false
      valid       = provider::validatefx::has_prefix("tf-backend", ["tf-", "iac-"], false)
    },
    {
      label       = "case insensitive"
      value       = "TF-service"
      prefixes    = ["tf-", "iac-"]
      ignore_case = true
      valid       = provider::validatefx::has_prefix("TF-service", ["tf-", "iac-"], true)
    }
  ]

  set_equals_checks = [
    {
      label    = "matching sets"
      values   = ["feature_a", "feature_b", "feature_c"]
      expected = ["feature_c", "feature_a", "feature_b"]
      valid = provider::validatefx::set_equals(
        ["feature_a", "feature_b", "feature_c"],
        ["feature_c", "feature_a", "feature_b"],
      )
    }
  ]

  string_length_values = [
    {
      value      = "short"
      min_length = 3
      max_length = 10
    },
  ]

  between_checks = [
    {
      label = "within"
      value = "7.5"
      min   = "5"
      max   = "10"
    },
  ]

  phone_results = [
    for value in local.phone_numbers : {
      value = value
      valid = provider::validatefx::phone(value)
    }
  ]

  mac_address_results = [
    for value in local.mac_addresses : {
      value = value
      valid = provider::validatefx::mac_address(value)
    }
  ]

  url_results = [
    for value in local.url_values : {
      value = value
      valid = provider::validatefx::url(value)
    }
  ]

  cidr_results = [
    for value in local.cidr_values : {
      value = value
      valid = provider::validatefx::cidr(value)
    }
  ]

  string_length_results = [
    for item in local.string_length_values : {
      value = item.value
      valid = provider::validatefx::string_length(item.value, item.min_length, item.max_length)
    }
  ]

  between_results = [
    for sample in local.between_checks : {
      label = sample.label
      value = sample.value
      valid = provider::validatefx::between(sample.value, sample.min, sample.max)
    }
  ]

  all_valid_results = [
    for values in [
      [true, true, true],
      [true, false],
      [true, null],
      ] : {
      checks = values
      result = provider::validatefx::all_valid(values)
    }
  ]

  any_valid_results = [
    for values in [
      [false, false],
      [false, true],
      [false, null, false],
      ] : {
      checks = values
      result = provider::validatefx::any_valid(values)
    }
  ]

  # Assert function tests
  assert_email_valid = provider::validatefx::assert(
    provider::validatefx::email("alice@example.com"),
    "Email validation failed!"
  )

  assert_uuid_valid = provider::validatefx::assert(
    provider::validatefx::uuid("d9428888-122b-11e1-b85c-61cd3cbb3210"),
    "UUID validation failed!"
  )

  assert_custom_condition = provider::validatefx::assert(
    length("test") == 4,
    "String length assertion failed!"
  )

  provider_version = provider::validatefx::version()
  # AWS Region checks
  aws_region_checks = [
    {
      description = "Valid US East region"
      value       = "us-east-1"
      valid       = provider::validatefx::aws_region("us-east-1")
    },
    {
      description = "Valid EU West region"
      value       = "eu-west-1"
      valid       = provider::validatefx::aws_region("eu-west-1")
    },
    {
      description = "Valid Asia Pacific region"
      value       = "ap-southeast-1"
      valid       = provider::validatefx::aws_region("ap-southeast-1")
    },
    {
      description = "Valid US GovCloud region"
      value       = "us-gov-west-1"
      valid       = provider::validatefx::aws_region("us-gov-west-1")
    },
    {
      description = "Valid China region"
      value       = "cn-north-1"
      valid       = provider::validatefx::aws_region("cn-north-1")
    },
  ]
}

locals {
  cidr_overlap_checks = {
    non_overlap = provider::validatefx::cidr_overlap(["10.0.0.0/24", "10.0.1.0/24"])
  }
}

output "validatefx_cidr_overlap" {
  value = local.cidr_overlap_checks
}

output "validatefx_subnet" {
  value = local.subnet_checks
}

output "validatefx_email" {
  value = local.email_results
}

output "validatefx_uuid" {
  value = local.uuid_results
}

output "validatefx_base64" {
  value = local.base64_results
}

output "validatefx_credit_card" {
  value = local.credit_card_results
}

output "validatefx_domain" {
  value = local.domain_results
}

output "validatefx_fqdn" {
  value = local.fqdn_checks
}

output "validatefx_hostname" {
  value = local.hostname_results
}

output "validatefx_fqdn_idna" {
  value = local.fqdn_idna_checks
}

output "validatefx_json" {
  value = local.json_results
}

output "validatefx_semver" {
  value = local.semver_results
}

output "validatefx_semver_range" {
  value = local.semver_range_results
}

output "validatefx_datetime" {
  value = local.datetime_results
}

output "validatefx_ip" {
  value = local.ip_results
}

output "validatefx_public_ip" {
  value = local.public_ip_checks
}

output "validatefx_matches_regex" {
  value = local.matches_regex_results
}

output "validatefx_uri" {
  value = local.uri_checks
}

output "validatefx_arn" {
  value = local.arn_checks
}

output "validatefx_private_ip" {
  value = local.private_ip_checks
}

output "validatefx_username" {
  value = local.username_results
}

output "validatefx_string_contains" {
  value = local.string_contains_results
}

output "validatefx_hex" {
  value = local.hex_checks
}

output "validatefx_base32" {
  value = local.base32_checks
}

output "validatefx_integer" {
  value = local.integer_checks
}

output "validatefx_ssh_public_key" {
  value = local.ssh_public_key_checks
}

output "validatefx_jwt" {
  value = local.jwt_checks
}

output "validatefx_password_strength" {
  value = local.password_strength_checks
}

output "validatefx_in_list" {
  value = local.in_list_checks
}

output "validatefx_not_in_list" {
  value = local.not_in_list_checks
}

output "validatefx_exactly_one_valid" {
  value = local.exactly_one_valid_checks
}

variable "integration_unknown_email" {
  description = "Optional email value for exactly_one_valid integration tests"
  type        = string
  default     = null
}

output "validatefx_in_list_integration" {
  value = local.in_list_integration_checks
}

output "validatefx_phone" {
  value = local.phone_results
}

output "validatefx_mac_address" {
  value = local.mac_address_results
}

output "validatefx_url" {
  value = local.url_results
}

output "validatefx_cidr" {
  value = local.cidr_results
}

output "validatefx_ip_range_size" {
  value = local.ip_range_size_checks
}

output "validatefx_string_length" {
  value = local.string_length_results
}

output "validatefx_between" {
  value = local.between_results
}

output "validatefx_all_valid" {
  value = local.all_valid_results
}

output "validatefx_any_valid" {
  value = local.any_valid_results
}

output "validatefx_assert" {
  value = {
    email_check      = local.assert_email_valid
    uuid_check       = local.assert_uuid_valid
    custom_condition = local.assert_custom_condition
  }
}

output "validatefx_version" {
  value = local.provider_version
}

output "validatefx_has_suffix" {
  value = local.has_suffix_checks
}

output "validatefx_has_prefix" {
  value = local.has_prefix_checks
}

output "validatefx_set_equals" {
  value = local.set_equals_checks
}

output "validatefx_list_subset" {
  value = local.list_subset_checks
}

output "validatefx_port_range" {
  value = local.port_range_checks
}

output "validatefx_port_number" {
  value = local.port_number_checks
}


output "validatefx_aws_region" {
  value = local.aws_region_checks
}
