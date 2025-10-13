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
}

locals {
  email_results = [
    for email in local.emails : {
      value = email
      valid = validatefx_email(email)
    }
  ]

  uuid_results = [
    for id in local.uuids : {
      value = id
      valid = validatefx_uuid(id)
    }
  ]

  base64_results = [
    for value in local.base64_values : {
      value = value
      valid = validatefx_base64(value)
    }
  ]

  credit_card_results = [
    for number in local.credit_cards : {
      value = number
      valid = validatefx_credit_card(number)
    }
  ]
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
