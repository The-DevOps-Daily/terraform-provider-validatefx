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
  # Example 1: All unique elements
  list1  = ["apple", "banana", "cherry"]
  valid1 = provider::validatefx::list_unique(local.list1)

  # Example 2: Single element (always unique)
  list2  = ["single"]
  valid2 = provider::validatefx::list_unique(local.list2)

  # Example 3: Empty list (always unique)
  list3  = []
  valid3 = provider::validatefx::list_unique(local.list3)

  # Example 4: Unique subnet CIDRs
  subnet_cidrs = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  valid4       = provider::validatefx::list_unique(local.subnet_cidrs)

  # Example 5: Unique resource names
  resource_names = ["web-server-1", "web-server-2", "db-server-1"]
  valid5         = provider::validatefx::list_unique(local.resource_names)

  # Example 6: Unique environment tags
  env_tags = ["production", "staging", "development"]
  valid6   = provider::validatefx::list_unique(local.env_tags)
}

output "validation_results" {
  value = {
    all_unique       = local.valid1
    single_element   = local.valid2
    empty_list       = local.valid3
    subnet_cidrs     = local.valid4
    resource_names   = local.valid5
    environment_tags = local.valid6
  }
}
