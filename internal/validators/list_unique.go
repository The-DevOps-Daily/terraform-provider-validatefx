package validators

import (
	"context"
	"fmt"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ListUniqueValidator validates that all elements in a list are unique.
type ListUniqueValidator struct{}

// Ensure interface compliance.
var _ frameworkvalidator.List = (*ListUniqueValidator)(nil)
var _ frameworkvalidator.Set = (*ListUniqueValidator)(nil)

// NewListUnique creates a new validator that checks all list elements are unique.
func NewListUnique() *ListUniqueValidator {
	return &ListUniqueValidator{}
}

// Description returns a plain text description of the validator.
func (ListUniqueValidator) Description(_ context.Context) string {
	return "all list elements must be unique"
}

// MarkdownDescription returns a markdown formatted description of the validator.
func (v ListUniqueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation on the given values.
func (ListUniqueValidator) Validate(values []string) error {
	seen := make(map[string]struct{}, len(values))
	duplicateSet := make(map[string]struct{})
	var duplicates []string

	for _, value := range values {
		if _, exists := seen[value]; exists {
			if _, alreadyRecorded := duplicateSet[value]; !alreadyRecorded {
				duplicates = append(duplicates, value)
				duplicateSet[value] = struct{}{}
			}
		} else {
			seen[value] = struct{}{}
		}
	}

	if len(duplicates) > 0 {
		return fmt.Errorf("duplicate elements found: %s", strings.Join(duplicates, ", "))
	}

	return nil
}

// ValidateList validates a list attribute value.
//
//nolint:cyclop
func (ListUniqueValidator) ValidateList(ctx context.Context, req frameworkvalidator.ListRequest, resp *frameworkvalidator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var items []basetypes.StringValue
	if err := req.ConfigValue.ElementsAs(ctx, &items, false); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Collection", "Expected a list of strings.")
		return
	}

	seen := make(map[string]struct{}, len(items))
	duplicateSet := make(map[string]struct{})
	var duplicates []string

	for _, item := range items {
		if item.IsNull() || item.IsUnknown() {
			continue
		}
		value := item.ValueString()
		if _, exists := seen[value]; exists {
			if _, alreadyRecorded := duplicateSet[value]; !alreadyRecorded {
				duplicates = append(duplicates, value)
				duplicateSet[value] = struct{}{}
			}
		} else {
			seen[value] = struct{}{}
		}
	}

	if len(duplicates) > 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Duplicate Elements",
			fmt.Sprintf("List contains duplicate elements: %s", strings.Join(duplicates, ", ")),
		)
	}
}

// ValidateSet validates a set attribute value.
// Note: Sets are inherently unique, so this always passes.
func (ListUniqueValidator) ValidateSet(_ context.Context, req frameworkvalidator.SetRequest, resp *frameworkvalidator.SetResponse) {
	// Sets are always unique by definition, so no validation needed
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	// Always passes for sets
}
