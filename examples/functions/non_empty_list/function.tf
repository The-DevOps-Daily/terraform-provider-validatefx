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
  # Example 1: Single element list
  list1  = ["item"]
  valid1 = provider::validatefx::non_empty_list(local.list1)

  # Example 2: Multiple elements
  list2  = ["apple", "banana", "cherry"]
  valid2 = provider::validatefx::non_empty_list(local.list2)

  # Example 3: Many elements
  list3  = ["one", "two", "three", "four", "five"]
  valid3 = provider::validatefx::non_empty_list(local.list3)
}

output "validation_results" {
  value = {
    list1_valid = local.valid1
    list2_valid = local.valid2
    list3_valid = local.valid3
  }
}
