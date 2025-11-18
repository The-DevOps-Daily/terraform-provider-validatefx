# ValidateFX Function Examples

This directory contains usage examples for all ValidateFX provider functions. Each example demonstrates practical use cases with **valid inputs only**, following Terraform best practices.

## Example Structure

Each function has its own directory under `examples/functions/<function_name>/` containing:

- **function.tf** - Terraform configuration demonstrating the function
- Provider configuration with version constraints  
- Local variables showing typical usage patterns
- Output values for verification

## Important Notes

### ✅ Valid Inputs Only

Examples in this directory contain **only valid inputs** that will pass validation. This ensures:
- Examples can be run successfully with `terraform plan` or `terraform apply`
- Users see realistic, working configurations
- CI integration tests pass consistently

### ❌ Testing Invalid Inputs

To test validation failures, refer to:
- **Unit tests**: `internal/validators/*_test.go` - comprehensive test cases including failure scenarios
- **Function tests**: `internal/functions/*_test.go` - function-level validation including error handling
- **Integration tests**: `integration/main.tf` - end-to-end validation with passing scenarios only

## Running Examples

To try an example:

```bash
cd examples/functions/<function_name>
terraform init
terraform plan
```

## Contributing

When adding new function examples:

1. Create a directory: `examples/functions/<function_name>/`
2. Add `function.tf` with valid inputs only
3. Include helpful comments explaining the validation
4. Add output showing the validation result
5. Test the example runs successfully

See [CONTRIBUTING.md](../CONTRIBUTING.md) for complete guidelines.
