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
  # Example 1: List with length at minimum
  list1  = ["item1", "item2"]
  valid1 = provider::validatefx::list_length_between(local.list1, "2", "5")

  # Example 2: List with length at maximum
  list2  = ["a", "b", "c", "d", "e"]
  valid2 = provider::validatefx::list_length_between(local.list2, "2", "5")

  # Example 3: List with length in middle of range
  list3  = ["apple", "banana", "cherry"]
  valid3 = provider::validatefx::list_length_between(local.list3, "2", "5")

  # Example 4: Exact length (min equals max)
  list4  = ["one", "two", "three"]
  valid4 = provider::validatefx::list_length_between(local.list4, "3", "3")

  # Example 5: Empty list with zero minimum
  list5  = []
  valid5 = provider::validatefx::list_length_between(local.list5, "0", "10")

  # Example 6: Validating security group rules (1-50 rules)
  security_rules = ["rule1", "rule2", "rule3"]
  valid6         = provider::validatefx::list_length_between(local.security_rules, "1", "50")
}

output "validation_results" {
  value = {
    at_minimum     = local.valid1
    at_maximum     = local.valid2
    in_middle      = local.valid3
    exact_length   = local.valid4
    empty_list     = local.valid5
    security_rules = local.valid6
  }
}
