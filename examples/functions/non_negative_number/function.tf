terraform {
  required_providers {
    validatefx = {
      source  = "The-DevOps-Daily/validatefx"
      version = ">= 0.0.1"
    }
  }
}

provider "validatefx" {}

locals {
  numbers = [
    "0",
    "42",
    "3.14",
    "0.001",
    "100.5",
  ]

  checked_numbers = [
    for number in local.numbers : {
      number = number
      valid  = provider::validatefx::non_negative_number(number)
    }
  ]
}

output "checked_numbers" {
  value = local.checked_numbers
}
