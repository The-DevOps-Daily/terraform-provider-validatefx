# Development Guide

This guide provides detailed information for developers working on the ValidateFX provider.

## Architecture Overview

The provider is organized into three main layers:

### 1. Validators (`internal/validators/`)

Core validation logic implementing the `validator.String` interface.

**Structure:**
```
internal/validators/
├── <name>.go          # Validator implementation
├── <name>_test.go     # Unit tests
└── <name>_fuzz_test.go # Fuzz tests (required)
```

**Pattern:**
```go
func ValidatorName() frameworkvalidator.String {
    return validatorNameValidator{}
}

type validatorNameValidator struct{}

func (v validatorNameValidator) Description(_ context.Context) string {
    return "short description"
}

func (v validatorNameValidator) MarkdownDescription(_ context.Context) string {
    return "markdown description"
}

func (v validatorNameValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
    // Validation logic
}
```

### 2. Functions (`internal/functions/`)

Terraform function wrappers exposing validators as provider functions.

**Structure:**
```
internal/functions/
├── <name>.go          # Function wrapper
├── <name>_test.go     # Function tests
├── common.go          # Shared helpers
└── registry.go        # Function registration
```

**Simple validators use:**
```go
func NewValidatorNameFunction() function.Function {
    return newStringValidationFunction(
        "validator_name",
        "Short summary",
        "Detailed description",
        validators.ValidatorName(),
    )
}
```

### 3. Provider (`internal/provider/`)

Provider configuration and initialization.

## Adding a New Validator

### Step-by-Step Checklist

- [ ] **1. Implement validator** in `internal/validators/<name>.go`
  - Export function returning `frameworkvalidator.String`
  - Implement Description, MarkdownDescription, ValidateString
  - Handle null/unknown values (return early)
  - Use clear error messages with `resp.Diagnostics.AddAttributeError`

- [ ] **2. Add unit tests** in `internal/validators/<name>_test.go`
  - Test valid inputs (should pass)
  - Test invalid inputs (should error)
  - Test null/unknown (should pass without error)
  - Test empty string behavior
  - Use table-driven tests

- [ ] **3. Add fuzz test** in `internal/validators/<name>_fuzz_test.go`
  - Required for `make validate` to pass
  - Seed with representative valid and invalid inputs
  - Test for panics and unexpected behavior

- [ ] **4. Create function wrapper** in `internal/functions/<name>.go`
  - Use `newStringValidationFunction` helper for simple cases
  - Complex validators may need custom implementation

- [ ] **5. Add function tests** in `internal/functions/<name>_test.go`
  - Test function execution
  - Test null/unknown propagation
  - Test error handling

- [ ] **6. Register function** in `internal/functions/registry.go`
  - Add to `ProviderFunctionFactories()` slice
  - Alphabetical order preferred

- [ ] **7. Create example** in `examples/functions/<name>/function.tf`
  - Show practical usage
  - **Use only valid inputs** (examples should work)
  - Include comments explaining the validation

- [ ] **8. Add integration test** in `integration/main.tf`
  - Add passing scenario only
  - Use the function with valid input
  - Add output to verify

- [ ] **9. Generate documentation**
  - Run `make docs` to generate `docs/functions/<name>.md`
  - README function table updates automatically

- [ ] **10. Validate**
  - Run `make validate` (all checks must pass)
  - Run `make coverage` (check coverage didn't drop)

## Code Style Guidelines

### Error Messages

**Structure:**
```go
resp.Diagnostics.AddAttributeError(
    req.Path,
    "Error Title",          // Short, capitalized
    "Detailed message",     // Include value in quotes
)
```

**Examples:**
- Good: `"Value \"abc\" is not a valid email address: missing '@'"`
- Bad: `"invalid email"`

### Null/Unknown Handling

Always return early for null or unknown values:
```go
if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
    return
}
```

### Empty String

Decide if empty string is valid for your validator:
```go
value := req.ConfigValue.ValueString()
if value == "" {
    return  // or error, depending on requirements
}
```

## Testing Philosophy

### Unit Tests
- Comprehensive coverage of edge cases
- Test both success and failure paths
- Fast execution

### Fuzz Tests
- Robustness testing with random inputs
- Catch panics and unexpected behavior
- Required for all validators

### Integration Tests
- End-to-end validation
- **Only passing scenarios** (errors tested in unit tests)
- Ensures functions work in Terraform context

### Examples
- Demonstrate real-world usage
- **Must use valid inputs only**
- Should be copy-paste ready

## Common Patterns

### Using Helpers

**String list parsing:**
```go
values, state, ok := stringListArgument(ctx, req, resp, 1, "items")
if !ok {
    return
}
if state == valueUnknown {
    resp.Result = function.NewResultData(types.BoolUnknown())
    return
}
```

**Boolean flags:**
```go
ignoreCase, state, ok := ignoreCaseFlag(ctx, req, resp, 2)
if !ok {
    return
}
```

### Embedding Validation Logic

```go
validator := validators.MyValidator()
validationResp := frameworkvalidator.StringResponse{}
validator.ValidateString(ctx, frameworkvalidator.StringRequest{
    ConfigValue: value,
    Path:        path.Root("value"),
}, &validationResp)

if validationResp.Diagnostics.HasError() {
    resp.Error = function.FuncErrorFromDiags(ctx, validationResp.Diagnostics)
    return
}
```

## Performance Considerations

- Validators are called frequently - keep them fast
- Compile regexes once (package-level variables)
- Avoid unnecessary allocations
- Use early returns to skip unnecessary work

## Release Process

1. Update CHANGELOG.md
2. Create version tag: `git tag v0.X.Y`
3. Push tag: `git push origin v0.X.Y`
4. GitHub Actions builds and publishes automatically
5. Registry updates within ~30 minutes

## Useful Commands

```bash
# Full validation suite
make validate

# Quick checks before push
make pre-push

# Coverage report
make coverage
make coverage-html

# Run specific tests
go test ./internal/validators -run TestEmail
go test ./internal/functions -run TestEmailFunction

# Fuzz specific validator
go test ./internal/validators -run FuzzEmail -fuzz FuzzEmail -fuzztime=30s

# Integration tests
make integration
```

## Troubleshooting

### "missing fuzz test" error
- Every validator needs a `<name>_fuzz_test.go` file
- Run `scripts/check-fuzz-coverage.go` to identify missing tests

### "missing function example" error
- Every function needs `examples/functions/<name>/function.tf`
- Run `scripts/check-function-coverage.go examples` to identify missing examples

### "missing integration test" error
- Every function must be called in `integration/main.tf`
- Run `scripts/check-function-coverage.go integration` to identify missing coverage

### Tests fail in CI but pass locally
- Run `go fmt ./...` and commit
- Ensure `go.mod` and `go.sum` are up to date (`go mod tidy`)
- Check that all files are committed

## Getting Help

- Check existing validators for patterns
- Review CONTRIBUTING.md for guidelines
- Open a draft PR for early feedback
- Tag maintainers for specific questions
