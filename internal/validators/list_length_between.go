package validators

import (
	"context"
	"fmt"

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

	if length < v.min {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"List Too Short",
			fmt.Sprintf("List must have at least %d elements, got %d.", v.min, length),
		)
		return
	}

	if length > v.max {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"List Too Long",
			fmt.Sprintf("List must have at most %d elements, got %d.", v.max, length),
		)
	}
}

// ValidateSet validates a set attribute value.
func (v *ListLengthBetweenValidator) ValidateSet(_ context.Context, req frameworkvalidator.SetRequest, resp *frameworkvalidator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	length := len(req.ConfigValue.Elements())

	if length < v.min {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Set Too Small",
			fmt.Sprintf("Set must have at least %d elements, got %d.", v.min, length),
		)
		return
	}

	if length > v.max {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Set Too Large",
			fmt.Sprintf("Set must have at most %d elements, got %d.", v.max, length),
		)
	}
}
