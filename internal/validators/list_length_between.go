package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ListLengthBetweenValidator validates that a list has a length between min and max (inclusive).
type ListLengthBetweenValidator struct {
	min int
	max int
}

// Ensure interface compliance.
var _ frameworkvalidator.List = (*ListLengthBetweenValidator)(nil)
var _ frameworkvalidator.Set = (*ListLengthBetweenValidator)(nil)

// NewListLengthBetween creates a new validator that checks list length is between min and max (inclusive).
func NewListLengthBetween(min, max int) *ListLengthBetweenValidator {
	return &ListLengthBetweenValidator{
		min: min,
		max: max,
	}
}

// Description returns a plain text description of the validator.
func (v *ListLengthBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("list must have length between %d and %d (inclusive)", v.min, v.max)
}

// MarkdownDescription returns a markdown formatted description of the validator.
func (v *ListLengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation on the given value.
func (v *ListLengthBetweenValidator) Validate(values []string) error {
	length := len(values)

	if length < v.min {
		return fmt.Errorf("list length %d is less than minimum %d", length, v.min)
	}

	if length > v.max {
		return fmt.Errorf("list length %d is greater than maximum %d", length, v.max)
	}

	return nil
}

// ValidateList validates a list attribute value.
func (v *ListLengthBetweenValidator) ValidateList(_ context.Context, req frameworkvalidator.ListRequest, resp *frameworkvalidator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	length := len(req.ConfigValue.Elements())

	v.validateLength(length, req.Path, &resp.Diagnostics, "List")
}

// ValidateSet validates a set attribute value.
func (v *ListLengthBetweenValidator) ValidateSet(_ context.Context, req frameworkvalidator.SetRequest, resp *frameworkvalidator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	length := len(req.ConfigValue.Elements())

	v.validateLength(length, req.Path, &resp.Diagnostics, "Set")
}

// validateLength is a helper that validates length and adds diagnostics.
func (v *ListLengthBetweenValidator) validateLength(length int, p path.Path, diags *diag.Diagnostics, typeName string) {
	if length < v.min {
		diags.AddAttributeError(
			p,
			fmt.Sprintf("%s Too Short", typeName),
			fmt.Sprintf("%s must have at least %d elements, got %d.", typeName, v.min, length),
		)
		return
	}

	if length > v.max {
		diags.AddAttributeError(
			p,
			fmt.Sprintf("%s Too Long", typeName),
			fmt.Sprintf("%s must have at most %d elements, got %d.", typeName, v.max, length),
		)
	}
}
