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
  credit_cards = [
    "4532015112830366",    # Valid Visa
    "5555555555554444",    # Valid MasterCard
    "378282246310005",     # Valid American Express
    "4532 0151 1283 0366", # Valid Visa with spaces
    "4532-0151-1283-0366", # Valid Visa with hyphens
    "4532015112830367",    # Invalid (wrong checksum)
    "123456789012",        # Invalid (too short)
    "not-a-credit-card",   # Invalid (contains letters)
  ]

  checked = [
    for card in local.credit_cards : {
      card_number = card
      is_valid    = provider::validatefx::credit_card(card)
    }
  ]
}

output "credit_card_validation" {
  value = local.checked
}
