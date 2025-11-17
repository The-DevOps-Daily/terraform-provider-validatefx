locals {
  allowed_roles = ["reader", "writer", "admin"]
  team_roles    = ["writer", "reader"]
}

output "list_subset_example" {
  value = {
    roles_valid = provider::validatefx::list_subset(local.team_roles, local.allowed_roles)
  }
}

