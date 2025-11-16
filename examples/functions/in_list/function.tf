locals {
  allowed_colors = ["red", "green", "blue"]

  color_checks = [
    {
      value       = "green"
      ignore_case = false
      valid       = provider::validatefx::in_list("green", local.allowed_colors, false, null)
    },
    {
      value       = "Green"
      ignore_case = true
      valid       = provider::validatefx::in_list("Green", local.allowed_colors, true, null)
    },
    {
      value       = "purple"
      ignore_case = false
      valid       = provider::validatefx::in_list("purple", local.allowed_colors, false, "Unsupported color; choose one of: ${join(", ", local.allowed_colors)}")
    }
  ]
}

output "validatefx_in_list" {
  value = local.color_checks
}
