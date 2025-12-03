terraform {
  required_providers {
    validatefx = {
      source = "the-devops-daily/validatefx"
    }
  }
}

# Valid Kubernetes label keys
output "valid_simple_key" {
  value = provider::validatefx::k8s_label_key("app")
}

output "valid_key_with_prefix" {
  value = provider::validatefx::k8s_label_key("kubernetes.io/name")
}

output "valid_key_with_subdomain" {
  value = provider::validatefx::k8s_label_key("example.com/app-name")
}

output "valid_key_with_dashes" {
  value = provider::validatefx::k8s_label_key("app-name")
}

# Validation will fail for invalid keys:
# - provider::validatefx::k8s_label_key("")                    # Empty key
# - provider::validatefx::k8s_label_key("a/b/c")               # Too many slashes
# - provider::validatefx::k8s_label_key("Example.com/app")     # Uppercase in prefix
# - provider::validatefx::k8s_label_key("app@name")            # Invalid characters
