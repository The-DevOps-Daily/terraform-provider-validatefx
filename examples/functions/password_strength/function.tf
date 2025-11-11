locals {
  passwords = ["Abc@1234", "Abcdef1!"]
}

output "password_strength_example" {
  value = [
    for p in local.passwords : {
      value = p
      valid = provider::validatefx::password_strength(p)
    }
  ]
}

