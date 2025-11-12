locals {
  tokens = [
    // header: {"alg":"HS256"} payload: {"sub":"1234"}
    "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  ]
}

output "jwt_example" {
  value = [
    for t in local.tokens : {
      value = t
      valid = provider::validatefx::jwt(t)
    }
  ]
}

