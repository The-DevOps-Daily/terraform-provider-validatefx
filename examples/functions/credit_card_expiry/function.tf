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
  expiry_dates = [
    "12/25",      # Valid MM/YY format
    "01/2026",    # Valid MM/YYYY format
    "06/30",      # Valid future date
    "12/99",      # Valid far future
    "12/2099",    # Valid far future (4-digit)
    "1/25",       # Invalid (single digit month)
    "13/25",      # Invalid (month > 12)
    "00/25",      # Invalid (month < 01)
    "12/2020",    # Invalid (past date)
    "01-25",      # Invalid (wrong separator)
    "0125",       # Invalid (no separator)
  ]

  checked = [
    for expiry in local.expiry_dates : {
      expiry_date = expiry
      is_valid    = provider::validatefx::credit_card_expiry(expiry)
    }
  ]
}

output "credit_card_expiry_validation" {
  value = local.checked
}
