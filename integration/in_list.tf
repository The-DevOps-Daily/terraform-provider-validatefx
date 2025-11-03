locals {
  allowed_status = ["draft", "review", "published"]

  status_checks = [
    {
      label       = "valid"
      value       = "draft"
      ignore_case = false
      valid       = provider::validatefx::in_list("draft", local.allowed_status, false)
    },
    {
      label       = "case-insensitive"
      value       = "Published"
      ignore_case = true
      valid       = provider::validatefx::in_list("Published", local.allowed_status, true)
    }
  ]
}

output "validatefx_in_list_integration" {
  value = local.status_checks
}
