# Contributing Guide

Thanks for your interest in contributing to **terraform-provider-validatefx**! This document outlines the recommended local development workflow and expectations for pull requests.

## Prerequisites

- Go 1.25.2 or newer (matching the version configured in `go.mod`)
- Docker and Docker Compose (required for integration tests)
- [`golangci-lint`](https://golangci-lint.run/) v1.61 or newer available on your `PATH`
- [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs) v0.19.x for documentation generation (`go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.2`)

## Getting Started

```bash
git clone https://github.com/The-DevOps-Daily/terraform-provider-validatefx.git
cd terraform-provider-validatefx
make deps
```

## Common Tasks

| Command | Description |
| --- | --- |
| `make fmt` | Format Go source files using `go fmt`. |
| `make build` | Compile the provider to verify it builds. |
| `make test` | Run unit tests. |
| `make lint` | Execute `golangci-lint` using the local installation. |
| `make docs` | Regenerate the Markdown docs under `docs/` via `tfplugindocs`. |
| `make integration` | Build the Docker image and run the Terraform integration scenario end-to-end. |
| `make coverage` | Generate test coverage report for functions and validators. |
| `make coverage-html` | Generate and open HTML coverage report in browser. |
| `make pre-push` | Run complete pre-push checklist (format, test, lint, coverage, docs). |
| `make clean` | Remove build artifacts and reset local integration state. |

> **Tip:** The `make help` command lists all available targets and descriptions.

## Adding a New Validator

Use the following checklist when introducing a new validator so every layer stays in sync:

1. **Author the core validator**
   - Create or update `internal/validators/<name>.go` with the validation logic.
   - Add comprehensive table-driven tests in `internal/validators/<name>_test.go` covering success, failure, and null/unknown inputs.
2. **Expose the Terraform function**
   - Implement the wrapper in `internal/functions/<name>.go` using `newStringValidationFunction` or the appropriate helper.
   - Add a focused unit test file `internal/functions/<name>_test.go` validating success, failure diagnostics, and null/unknown propagation.
   - Register the function in `internal/functions/registry.go` so Terraform can discover it.
3. **Exercise real scenarios**
   - Add or update example usage under `examples/functions/<name>/function.tf`.
   - Extend the integration suite in `integration/main.tf` with passing scenarios that call the new function.
4. **Document the addition**
   - Generate a doc stub `docs/functions/<name>.md` (or update an existing one) so `tfplugindocs` has content to render.
   - If the README lists available validators, update the relevant table or section.

### Diagnostics and Edge Cases

- Prefer `resp.Diagnostics.AddAttributeError` (or the helper plumbing already used in the project) so Terraform users receive actionable messages.
- Treat null or unknown inputs as indeterminate: return `types.BoolUnknown()` from functions and bail out early in validators when `ConfigValue` is null or unknown.
- Keep error wording consistent with existing validators—short summary plus a detail that includes the problematic value when safe.

### Automation Checklist

Run the project automation after making changes to guarantee consistency:

- `go fmt ./...` — ensures Go sources stay formatted.
- `go test ./...` — covers both validator and function unit tests.
- `make docs` — refreshes `docs/functions/*.md` using `tfplugindocs` (required after schema/function changes).
- `make validate` — runs the full validation pipeline (fmt, docs, and coverage checks) invoked by CI.

## Fuzz Testing

Go ships with native fuzzing support starting from Go 1.18. We include fuzz tests for all validators to harden them against edge cases and unexpected inputs.

### Running Fuzz Tests

Run a targeted fuzz test (10 seconds):

```bash
go test ./internal/validators -run FuzzEmail -fuzz FuzzEmail -fuzztime=10s
```

Run all fuzz tests in the validators package (1 minute):

```bash
go test ./internal/validators -fuzz Fuzz -fuzztime=1m
```

Run fuzz tests for specific validators:

```bash
# Email validator
go test ./internal/validators -run FuzzEmail -fuzz FuzzEmail -fuzztime=30s

# URL validator
go test ./internal/validators -run FuzzURL -fuzz FuzzURL -fuzztime=30s

# JSON validator
go test ./internal/validators -run FuzzJSON -fuzz FuzzJSON -fuzztime=30s
```

### Fuzz Test Requirements

When adding a new validator, you must include a corresponding fuzz test:

1. Create `internal/validators/<name>_fuzz_test.go`
2. Implement a `Fuzz<Name>` function that exercises the validator with random inputs
3. Seed the fuzzer with both valid and invalid examples
4. The `make validate` command enforces this requirement via `scripts/check-fuzz-coverage.go`

### Understanding Fuzz Failures

When the fuzzer finds a failure, it:

1. Minimizes the input to the smallest failing case
2. Writes it to `testdata/fuzz/<FuzzFunction>/` as a corpus entry
3. Re-runs that corpus entry in future test runs to prevent regressions

If a fuzz test fails:

1. Examine the failing input in the corpus file
2. Fix the validator to handle that case correctly
3. Add a unit test for the specific case to document the fix
4. Re-run the fuzz test to verify the fix

## Pull Request Checklist

Before opening a PR, please ensure:

- Code is formatted (`make fmt`).
- Unit tests pass (`make test`).
- Linting reports no issues (`make lint`).
- Integration tests succeed (`make integration`) when applicable.
- Documentation is regenerated (`make docs`) when schema changes affect the published docs.
- Commits are descriptive and scoped to a logical change.
- PR includes a brief summary explaining motivation and testing performed.


### Pre-commit Hooks
This repository uses [pre-commit](https://pre-commit.com/) to automatically format and lint code before it is committed. This helps maintain code quality and consistency.
#### Installation
To use the hooks, you must install them locally:
1.  Install the `pre-commit` tool. A common way is with Python's package manager:
    ```bash
    pip install pre-commit
    ```
2.  Install the hooks in this repository. From the root directory, run:
    ```bash
    pre-commit install
    ```
After this, the hooks (including `terraform fmt`, `go fmt`, and `make lint`) will run automatically on every `git commit`. If they find an issue, they may fix it for you or stop the commit so you can fix it manually.
