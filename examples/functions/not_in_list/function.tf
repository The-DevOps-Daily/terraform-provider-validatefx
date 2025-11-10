locals {
  disallowed_colors = ["red", "blue"]
  sample_values     = ["green", "yellow"]
}

output "not_in_list_example" {
  value = [
    for v in local.sample_values : {
      value = v
      valid = provider::validatefx::not_in_list(v, local.disallowed_colors, false)
    }
  ]
}

