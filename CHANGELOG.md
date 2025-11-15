# Changelog

## [Unreleased]

### Changes

- None.

### Dependency Updates

- None.

## [0.3.0] - 2025-11-15

### 🚀 Features
- cidr_overlap: add validator and Terraform function to detect overlapping CIDR blocks; examples, tests, integration (PR #351, closes #260)
- semver_range: add validator and function with examples, docs, integration (PR #337)

### 🧰 Maintenance
- make: add fuzz-quick target for short local fuzz runs (PR #349)
- functions: refactor to shared string helper; reuse in semver_range (PR #340, #339)
- ci: add separate fuzz workflow running on PRs in parallel (PR #346)
- validate target: include go vet/test/tidy (PR #342)

### 🧪 Tests
- Fuzz coverage expanded:
  - email, URL, JSON (PR #343)
  - in_list, not_in_list, set_equals (PR #344)
  - uuid, hostname, password_strength (PR #345)
  - ip, cidr, semver (PR #347)
  - credit_card, phone, hex (PR #348)

### Contributors
- @bobbyonmagic


## [0.2.2] - 2025-11-11

### Tests

- tests(datetime): provider defaults path and config Get/Set coverage (PR #335).
- tests(fqdn): add framework-level function tests (PR #334).
- tests(jwt): add framework-level function tests (PR #333).

### Documentation

- Add jwt validator and function; tests, examples, integration; closes #269 (PR #332).

### Contributors

- @bobbyonmagic

## [0.2.1] - 2025-11-11

### Features

- Add fqdn validator and function; tests, examples, integration; closes #253 (PR #329).
- Complete password_strength with framework-compliant function and validator; tests, examples, integration (PR #328).

### Improvements

- fqdn: support punycode (xn--) labels; tests updated; examples/integration unchanged (PR #330).

### Documentation

- docs: add DevOps Daily website reference to provider docs index template (PR #325).
- docs: ensure index.md is generated from template; add DevOps Daily link via template (PR #327).

### Contributors

- @bobbyonmagic
- @smitbhoir20
- @nihalnayak45

- @bobbyonmagic

## [0.2.0] - 2025-11-11

### Features

- Integer validator and Terraform function with tests and examples (PR #322).
- Provider-level datetime defaults scaffolding and wiring to datetime when layouts are null/empty (PR #321).
- Base32 validator and Terraform function with tests, examples, and docs (PR #317).
- SSH public key validator and Terraform function with tests, examples, and docs (PR #318).

### Improvements

- Align function wrappers to use validators consistently (PR #319).
- Refactor matches_regex to rely on validators; remove inline regex logic (PR #320, closes #17).

### Bug Fixes

- Integer validator: trim whitespace and ignore empty-like input to match validator conventions (PR #323).

### Documentation

- Update docs for new validators and provider-level defaults (PRs #316, #317, #318, #321, #323).

### Contributors

- @Ak00005

## [0.1.8] - 2025-11-10

### Features

- Add `not_in_list` validator and Terraform function with examples, docs, and integration coverage (PR #251, closes #158).
- Add `has_prefix` validator and Terraform function with tests and docs (PR #248).
- Add `set_equals` composite function for list set equality (PR #250).
- Add `username` validator and function; align naming for consistency (PR #247).

### Improvements

- Refactor `set_equals` to follow the validator + function wrapper pattern for consistency (ed5b07c).
- Ensure integration scenarios only include successful cases for new functions (4ebe84b).

### Documentation

- Update README CI badge to the unified GitHub Actions badge (PR #241).
- Remove deprecated OS installation document to avoid duplication (PR #240).

### CI / Tooling

- Bump golangci-lint GitHub Action from v8 to v9 (dependabot; 1e90796).

### Reverts

- Revert earlier `list_subset` addition pending follow-up design and coverage (PRs #244, #246).

## [0.1.7] - 2025-11-04

### Features

- Add string contains validator, Terraform function, example, integration scenario, and documentation coverage (#233, #234, #236, #238, #237, #235).
- Add has_suffix string validator with supporting documentation (#176).
- Add exactly_one_valid composite function (#172).

### Improvements

- Run lint as part of make validate (#174).

### Documentation

- Document workflow for adding validators (#173).
- Link to good first issues in README (#175).
- Document has_suffix usage adjustments (#183, #184).
- Add SECURITY policy (#171).

### Contributors

- @bobbyonmagic
- @shanaya-Gupta

## [0.1.6] - 2025-11-04

### Features

- Add the `in_list` validator and Terraform function with documentation, examples, and integration coverage (`0a06356`).
- Expand validation coverage with new helper validators and composite list utilities (`287d1e6`, `2d02eb4`, `b19056f`).

### Improvements

- Increase unit test coverage across validator functions and composites (`330e481`, `5cb498d`, `7cebb5a`).
- Optimize the Docker build pipeline by removing redundant lint stages and caching Go modules (`0c673f7`).

### Documentation

- Document the `in_list` function and add examples demonstrating module usage (`0a06356`).

### CI / Tooling

- Add GitHub issue templates, release drafter updates, and coverage-focused workflows (`4c9af37`, `18306e0`, `ecfd7e8`).


## [0.1.5] - 2025-11-02

### Features

- Add a MAC address validator and expose it via the `validatefx_mac_address` Terraform function with docs, examples, and tests (`1c26b93`, `5448c8c`, `acf160cc`).
- Introduce an RFC 1123 compliant hostname validator and Terraform function coverage (`e25ef04`).

### Improvements

- Add a `make validate` pre-flight target that runs formatting, docs generation, and function coverage checks locally (`25988e3`, `d42e7d6`).

### Documentation

- Publish a Contributor Covenant code of conduct for community guidelines (`c494fe5`, `54f69a7`).

### CI / Tooling

- Refresh pre-commit hook configuration to pick up newer linters (`46e3237`, `f76098a`).
- Update release drafter workflow configuration and token handling (`5a219a8`, `73431fb`, `7a1c0ec`, `7d1a28f`, `803780e`).


## [0.1.4] - 2025-10-30

### Features

- Add string length validator and expose the `validatefx_string_length` Terraform function with docs, examples, and integration coverage (`986d9c4`, `6b504be`, `75379a4`).
- Introduce CIDR validation with Terraform function support, documentation, and scenario updates (`9a3037c`, `53cfc96`, `f1fc4b7`).

### Improvements

- Automate changelog maintenance and grouping via the new update script (`806046b`, `b8ccfe6`, `bf1765f`).
- Add contributor workflow tooling including pre-commit hooks, issue templates, and Dependabot configuration (`4e37890`, `7a2eb13`, `501493a`, `7f165d6`).

### Documentation

- Publish OS-specific installation and troubleshooting guidance and add README contributor highlights (`3165060`, `9e90793`).

### Bug Fixes

- Resolve string length integration regressions uncovered during testing (`be9d1db`, `75379a4`).
- Tidy the CIDR Terraform scenario formatting to keep integration output stable (`0809e49`).

### Dependency Updates

- Bump actions/setup-go from 5 to 6 (`cb51c3e`, `0e44d89`).
- Bump goreleaser/goreleaser-action from 5 to 6 (`a69aba5`, `526fb34`).
- Bump actions/upload-artifact from 4 to 5 (`e74f0ee`, `b9a95d6`).
- Bump actions/checkout from 4 to 5 (`6babd66`, `1e87f4d`).


## [0.1.3] - 2025-10-28

### Features

- Add an HTTP/HTTPS URL validator exposed as `provider::validatefx::url`, including schema tests and Terraform coverage (`faf98d4`, `6a545cf`, `51bef43`).
- Expose provider metadata through the new `provider::validatefx::version` function with integration coverage and documentation updates (`9cdba92`, `84ba24d`, `18dd815`, `81e29af`).

### Improvements

- Expand Terraform integration scenarios to exercise additional validators and the provider version endpoint (`211d656`, `bec4e33`, `c6a6c4f`).
- Add defensive tests ensuring string validation functions surface diagnostics for non-string inputs (`c386eb0`, `e61d50b`).
- Restructure examples and documentation to streamline generation and add a provider quick-start snippet (`7027ef8`, `86db796`, `f67b9b2`, `9472110`).

### Bug Fixes

- Harden URL validation behavior and align imports and formatting (`860cb71`, `6a545cf`).
- Stabilize integration expectations by correcting email/base64 fixtures and handling null inputs (`6992130`, `5a01c2c`).
- Resolve intermittent test failures surfaced during integration expansion (`bfdba96`, `5676adc`).

---

## [0.1.2] - 2025-10-27

### Features

- Add composite validation helpers `all_valid` and `any_valid` for aggregating multiple checks (`a3e1c9a`, `8574455`).
- Expose the existing phone E.164 validator as a Terraform function with docs and examples (`5f62599`).
- Introduce the `matches_regex` Terraform function for pattern validation (`f825340`).

### Bug Fixes

- Cache compiled regular expressions in the `matches_regex` validator to avoid repeated compilation (`db161f7`).

### Misc

- Preserve the provider docs index template during documentation generation (`4171e03`).
- Publish a custom provider index document to improve docs navigation (`337b172`).

---

## [0.1.1] - 2025-10-26

### Features

- Add Terraform functions for JSON structure validation, Semantic Versioning checks, and IP address validation (`1ed7d28`, `ee2e5f3`, `19140c2`).
- Automate regeneration of the README “Available Functions” table to keep documentation in sync (`3bf9caa`, `3c8133a`).

### Bug Fixes

- Correct integration test Docker plugin path, README build badge, and Terraform Registry URLs (`13b6573`, `162c267`, `e3d40a6`).

### Misc

- Remove unused internal function helpers discovered during review (`d397f4d`).

---

## [0.1.0] - 2025-09-28

### Features

- Initial release of the provider scaffold with validators for email, UUID, base64, credit card, domain, and phone numbers plus Terraform examples and unit tests (`046cb51`, `c07ff64`, `0a478c1`, `35497a3`, `211bedc`, `8ce87fd`).
- Add Terraform integration workflows and supporting Makefile targets to validate the provider end to end (`0e74156`, `6944b72`, `2d42556`, `c6845dd`).
- Introduce release automation via GitHub Actions (`58c069f`).

### Bug Fixes

- Iterate on release workflows to resolve checksum, packaging, and pipeline failures (`6d823c7`, `5980981`, `5bb84d9`, `03babef`, `1679763`).
- Fix function parameter naming issues uncovered during early CI automation (`196831d`).

### Misc

- Add contributor guidelines, AGENTS metadata, and README badges to polish the project presentation (`c6845dd`, `7892797`, `9ff1444`).
- Expand validator test coverage with comprehensive table-driven suites (`8b222a5`, `24e67c5`).
