# Changelog

## [Unreleased]

### Changes

- None.

### Dependency Updates

- None.


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
