package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ frameworkvalidator.List = (*listSubsetValidator)(nil)
var _ frameworkvalidator.Set = (*listSubsetValidator)(nil)

// ListSubset returns a list validator ensuring all configured values are contained within the allowed collection.
func ListSubset(allowed []string) frameworkvalidator.List {
	return newListSubsetValidator(allowed)
}

// SetSubset returns a set validator ensuring all configured values are contained within the allowed collection.
func SetSubset(allowed []string) frameworkvalidator.Set {
	return newListSubsetValidator(allowed)
}

type listSubsetValidator struct {
	allowed []string
	lookup  map[string]struct{}
}

func newListSubsetValidator(allowed []string) *listSubsetValidator {
	lookup := make(map[string]struct{}, len(allowed))
	normalized := make([]string, 0, len(allowed))

	for _, candidate := range allowed {
		trimmed := strings.TrimSpace(candidate)
		if trimmed == "" {
			continue
		}

		if _, exists := lookup[trimmed]; exists {
			continue
		}

		lookup[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	return &listSubsetValidator{
		allowed: normalized,
		lookup:  lookup,
	}
}

func (v *listSubsetValidator) Description(_ context.Context) string {
	if len(v.allowed) == 0 {
		return "values must belong to the allowed set"
	}

	return fmt.Sprintf("values must be a subset of: %s", strings.Join(v.allowed, ", "))
}

func (v *listSubsetValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *listSubsetValidator) ValidateList(ctx context.Context, req frameworkvalidator.ListRequest, resp *frameworkvalidator.ListResponse) {
	v.validateCollection(ctx, req.ConfigValue, &resp.Diagnostics, req.Path)
}

func (v *listSubsetValidator) ValidateSet(ctx context.Context, req frameworkvalidator.SetRequest, resp *frameworkvalidator.SetResponse) {
	v.validateCollection(ctx, req.ConfigValue, &resp.Diagnostics, req.Path)
}

func (v *listSubsetValidator) validateCollection(ctx context.Context, value types.List, diags *diag.Diagnostics, attrPath path.Path) {
	if value.IsNull() || value.IsUnknown() {
		return
	}

	if len(v.lookup) == 0 {
		diags.AddAttributeError(
			attrPath,
			"Allowed Set Empty",
			"The allowed values collection must contain at least one entry for subset validation.",
		)
		return
	}

	var elements []string
	if d := value.ElementsAs(ctx, &elements, false); d.HasError() {
		diags.Append(d...)
		return
	}

	for _, element := range elements {
		if _, ok := v.lookup[element]; ok {
			continue
		}

		diags.AddAttributeError(
			attrPath,
			"Value Not Allowed",
			fmt.Sprintf("Value %q is not part of the allowed set (%s)", element, strings.Join(v.allowed, ", ")),
		)
	}
}
