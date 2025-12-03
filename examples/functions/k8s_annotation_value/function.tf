terraform {
  required_providers {
    validatefx = {
      source = "the-devops-daily/validatefx"
    }
  }
}

# Valid Kubernetes annotation values
output "valid_empty_annotation" {
  value = provider::validatefx::k8s_annotation_value("")
}

output "valid_simple_annotation" {
  value = provider::validatefx::k8s_annotation_value("This is a simple annotation")
}

output "valid_annotation_with_special_chars" {
  value = provider::validatefx::k8s_annotation_value("annotation@example.com: value!")
}

output "valid_annotation_with_newlines" {
  value = provider::validatefx::k8s_annotation_value("line1\nline2\nline3")
}

# Validation will fail for invalid annotations:
# - provider::validatefx::k8s_annotation_value(strings.repeat("a", 262145))  # Exceeds 256KB
