locals {
  allowed_colors = ["red", "green", "blue"]

  color_checks = [
    {
      value       = "green"
      ignore_case = false
      valid       = provider::validatefx::in_list("green", local.allowed_colors, false)
    },
    {
      value       = "Green"
      ignore_case = true
      valid       = provider::validatefx::in_list("Green", local.allowed_colors, true)
    }
  ]
}

output "validatefx_in_list" {
  value = local.color_checks
}
