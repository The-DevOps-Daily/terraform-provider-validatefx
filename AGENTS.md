# Repository Guidelines

## Project Structure & Module Organization
- `internal/provider/` exposes the Terraform provider entry point and now registers all validatefx functions.
- `internal/functions/` wraps string validators as Terraform functions; reuse these helpers when adding new callable utilities.
- `internal/validators/` holds the core Go validators and their unit tests.
- `integration/` contains HCL scenarios that exercise every exported function end to end.
- `examples/` provides practitioner-facing usage samples grouped by validator type.

## Build, Test, and Development Commands
- `go fmt ./...` — format Go sources across the repository; run before committing.
- `go test ./...` — execute unit tests, ensuring validators and helper packages stay green.
- `make lint-shell` — run shellcheck on bash scripts in scripts/ directory (requires shellcheck).
- `make build` — compile the provider binary into `bin/terraform-provider-validatefx`.
- `make install` — install the provider into the local Terraform plugin directory for manual validation.

### Local prerequisites
- Go 1.22+
- Docker and Docker Compose
- `golangci-lint` in PATH (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0`)
- `shellcheck` (optional, for shell script linting - https://github.com/koalaman/shellcheck#installing)
- `tfplugindocs` in PATH (`go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.2`)
- GitHub CLI (`gh`) for PRs, issues and releases

### Developer workflow
- Use feature branches with prefixes: `feat/`, `fix/`, `docs/`, `chore/`, `refactor/`, `test/`.
- Before pushing, always run `make validate` (formats, vets, lints, shell lints, tests, docs, function coverage, fuzz coverage).
- Integration scenarios are expected to be success-only; remove or adjust any failing examples before committing.

## Coding Style & Naming Conventions
- Go files must remain `gofmt` clean with tabs for indentation.
- Exported Go symbols follow mixed-case naming (`ValidateSomething`), while Terraform function names use the `validatefx_<name>` pattern.
- HCL examples and integration configs prefer two-space indentation and descriptive local variable names.

### Terraform Functions Behavior
- Provider functions should signal validation failures via diagnostics (errors), not by returning `false`.
- Wrappers typically return `true` when validation succeeds; invalid inputs must produce a diagnostic error.
- Terraform functions have fixed arity: emulate optional parameters by accepting `null` and handling it internally.

## Testing Guidelines
- Unit tests live beside their targets (e.g., `internal/validators/base64_test.go`) and should cover success, failure, and null/unknown cases.
- Integration tests rely on Docker (`make dev`) and the `integration/` directory; update scenarios whenever adding new provider functions.
- Use table-driven tests for new validators and ensure `go test ./...` passes prior to submission.

### Fuzz and Coverage Gates
- Every validator in `internal/validators/` must have at least one fuzz test: `*_fuzz_test.go`.
- `make validate` enforces:
  - All functions are referenced in examples/integration (`scripts/check-function-coverage.go`).
  - All validators have fuzz tests (`scripts/check-fuzz-coverage.go`).
- Keep fuzz tests robust and fast; CI also runs a dedicated fuzz workflow separately from main checks.

### Function wrapper tests
- Add tests under `internal/functions/` mirroring validator coverage (valid, invalid, null/unknown).
- Ensure wrappers propagate diagnostics as errors and yield `unknown` for null/unknown inputs.

## Commit & Pull Request Guidelines
- Commit messages typically use the imperative mood (`Add credit-card validator`, `Refactor functions helper`). Keep them scoped and reference issues where applicable.
- Pull requests should describe the validator/function changes, list testing evidence (`go test`, Terraform run), and link to any related integration updates.
- Include screenshots or logs only when troubleshooting Terraform apply output; otherwise prefer concise bullet summaries.

### Branch, PR, and Issue Hygiene
- Reference issues in commits/PRs (e.g., "feat: add subnet validator (#257)").
- For multi-part work, prefer one commit per acceptance criterion when reasonable.
- Use `gh pr create --fill` (or provide a clear title/body) after pushing a branch.

## Security & Configuration Tips
- Avoid hardcoding secrets in examples; rely on Terraform variables when demonstrating configurable values.
- When adding new validators, ensure external dependencies are vetted and present in `go.mod` with minimal version bump.

## Docs Generation
- Run `make docs` or `make validate` to regenerate provider docs via `tfplugindocs`.
- Do not hand-edit files under `docs/`; they are generated. Update templates/examples instead.
- The README validator table is updated by `scripts/update-readme-functions-table.go` during docs generation.

## Adding a New Validator (Checklist)
1. Implement validator in `internal/validators/<name>.go`.
2. Add unit tests: table-driven + null/unknown + error cases.
3. Add fuzz test: `internal/validators/<name>_fuzz_test.go`.
4. Add Terraform function wrapper in `internal/functions/<name>.go` (use `newStringValidationFunction`/common helpers when applicable).
5. Add wrapper unit tests under `internal/functions/`.
6. Register the function in `internal/functions/registry.go`.
7. Add example at `examples/functions/<name>/function.tf`.
8. Add a passing-only scenario to `integration/main.tf`.
9. Run `make validate` and address any failures (lint, coverage, docs).
10. Commit and open a PR linking the parent issue and subtasks.

## Adding a New Guide

Provider guides appear on the Terraform Registry and help users understand usage patterns.

1. Create guide file: `templates/guides/your-guide-name.md.tmpl`
2. Add frontmatter with `page_title`, `subcategory: \"Guides\"`, and `description`
3. Write guide content with examples and best practices
4. Update `templates/index.md.tmpl` to link to your guide under `## Guides`
5. Run `make docs` to generate `docs/guides/your-guide-name.md`
6. Verify guide appears in `docs/index.md`
7. Commit and open PR

See `templates/guides/README.md` for detailed guide writing tips and conventions.

## Integration Tests Rules
- Keep only passing scenarios in `integration/main.tf`; failures are covered by unit tests.
- When a function requires extra parameters, pass `null` for optional ones.
- Use descriptive local names and small, readable samples.

## CI and Releases
- Linting and tests run in GitHub Actions on PRs. Fuzz tests run in a separate workflow to avoid blocking.
- Docker images are built and pushed on tags (v*); required secrets in the repo:
  - `DOCKERHUB_USERNAME`
  - `DOCKERHUB_TOKEN`
- Releases and changelogs are prepared from draft release notes; tags follow semver (e.g., v0.6.0).

## Common Pitfalls
- Forgetting to register a new function in `internal/functions/registry.go` (docs/README table will also miss it).
- Adding failing cases to integration; these must be removed to keep CI green.
- Missing fuzz tests for a new validator; `make validate` will fail.
- Terraform function arity mismatches: supply `null` for optional parameters.

## Commands Cheat Sheet
- `make validate` — end-to-end local checks (fmt, tidy, vet, lint, test, docs, coverage gates).
- `make integration` — run Dockerized Terraform scenario from `integration/`.
- `go test ./...` — unit tests only.
- `$(go env GOPATH)/bin/tfplugindocs generate` — regenerate docs manually.
