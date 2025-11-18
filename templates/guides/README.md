# Provider Guides

This directory contains user-facing guides for the validatefx provider.

## Adding a New Guide

1. **Create the guide file** in this directory:
   ```bash
   templates/guides/your-guide-name.md.tmpl
   ```

2. **Add frontmatter** (required for Terraform Registry):
   ```markdown
   ---
   page_title: "Your Guide Title"
   subcategory: "Guides"
   description: |-
     Brief description of what this guide covers.
   ---
   ```

3. **Update the index** to link to your guide:
   Edit `templates/index.md.tmpl` and add your guide under the `## Guides` section:
   ```markdown
   ## Guides
   
   - [List Validators: Usage Patterns and Tips](guides/list-validators.md)
   - [Your Guide Title](guides/your-guide-name.md)
   ```

4. **Regenerate docs**:
   ```bash
   make docs
   ```

5. **Verify the output**:
   - Check that `docs/guides/your-guide-name.md` was created
   - Verify the link appears in `docs/index.md`

## Guide Writing Tips

### Structure
- Start with a clear introduction explaining what the guide covers
- Use clear section headers (##, ###)
- Include practical code examples
- End with troubleshooting or common pitfalls

### Code Examples
Use fenced code blocks with `terraform` syntax:

\`\`\`terraform
locals {
  example = provider::validatefx::email("user@example.com")
}
\`\`\`

### Best Practices
- Keep examples realistic and runnable
- Explain the "why" not just the "how"
- Link to related functions when relevant
- Use consistent terminology

## Current Guides

- **list-validators.md.tmpl** - Usage patterns for list-focused validators (in_list, not_in_list, set_equals, list_subset)

## Guide Ideas

Future guides could cover:
- **Network Validators** - IP, CIDR, subnet, port validation patterns
- **String Validation Strategies** - When to use regex vs specific validators
- **Composite Validation** - Using all_valid, any_valid, exactly_one_valid
- **Date/Time Validation** - Working with datetime formats and timezones
- **AWS Resource Validation** - ARN validation patterns
- **Testing Your Validations** - How to test Terraform configurations using validatefx
