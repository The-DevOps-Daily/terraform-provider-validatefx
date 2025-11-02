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
    "bad-email",
  ]

  uuids = [
    "d9428888-122b-11e1-b85c-61cd3cbb3210",
    "not-a-uuid",
  ]

  base64_values = [
    "U29sdmVkIQ==",
    "invalid base64",
  ]

  credit_cards = [
    "4532015112830366",
    "4532015112830367",
  ]

  phone_numbers = [
    "+14155552671",
    "14155552671",
  ]

  mac_addresses = [
    "00:1A:2B:3C:4D:5E",
    "00-1A-2B-3C-4D-5E",
    "001A2B3C4D5E",
    "AA:BB:CC:DD:EE",
  ]

  url_values = [
    "https://example.com",
    "ftp://example.com",
  ]

  domains = [
    "example.com",
    "invalid..domain",
  ]

  hostnames = [
    "service.internal",
    "xn--bcher-kva.example",
    "example.com.",
    "bad_name",
  ]

  json_payloads = [
    "{\"key\": \"value\"}",
    "{\"invalid\":",
    "[]",
  ]

  semver_values = [
    "1.0.0",
    "v1.0.0",
    "1.0",
  ]

  ip_values = [
    "127.0.0.1",
    "::1",
    "999.999.999.999",
  ]

  regex_samples = [
    {
      value   = "user_123"
      pattern = "^[a-z0-9_]+$"
    },
    {
      value   = "Invalid-User"
      pattern = "^[a-z0-9_]+$"
    },
  ]

  cidr_values = [
    "10.0.0.0/24",
    "2001:db8::/48",
    "bad-cidr",
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

  ip_results = [
    for value in local.ip_values : {
      value = value
      valid = provider::validatefx::ip(value)
    }
  ]

  matches_regex_results = [
    for item in local.regex_samples : {
      value   = item.value
      pattern = item.pattern
      valid   = provider::validatefx::matches_regex(item.value, item.pattern)
    }
  ]

  string_length_values = [
    {
      value      = "short"
      min_length = 3
      max_length = 10
    },
    {
      value      = "extremely-long-string"
      min_length = 3
      max_length = 10
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

output "validatefx_hostname" {
  value = local.hostname_results
}

output "validatefx_json" {
  value = local.json_results
}

output "validatefx_semver" {
  value = local.semver_results
}

output "validatefx_ip" {
  value = local.ip_results
}

output "validatefx_matches_regex" {
  value = local.matches_regex_results
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

output "validatefx_string_length" {
  value = local.string_length_results
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
