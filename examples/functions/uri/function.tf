locals {
  uris = [
    "http://example.com",
    "https://example.com/path?x=1",
    "ssh://user@host",
    "urn:isbn:0451450523",
  ]
}

output "uri_example" {
  value = [
    for u in local.uris : {
      uri   = u
      valid = provider::validatefx::uri(u)
    }
  ]
}

