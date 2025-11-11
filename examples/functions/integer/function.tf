locals {
  values = ["0", "-42", "+7", "3.14"]
}

output "integer_example" {
  value = [
    for v in local.values : {
      value = v
      valid = provider::validatefx::integer(v)
    }
  ]
}
