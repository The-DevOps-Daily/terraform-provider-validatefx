locals {
  arns = [
    {
      label = "iam role"
      value = "arn:aws:iam::123456789012:role/Admin"
      valid = provider::validatefx::arn("arn:aws:iam::123456789012:role/Admin")
    },
  ]
}

output "arn_examples" {
  value = local.arns
}

