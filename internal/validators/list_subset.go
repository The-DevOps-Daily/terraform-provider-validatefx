package validators

import (
	"context"
	"fmt"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ListSubsetValidator validates that all elements of a provided list/set are contained in a reference allowlist.
// It accepts list(string) or set(string) inputs from Terraform callers.
type ListSubsetValidator struct {
	allowed map[string]struct{}
	display []string
}

// NewListSubset constructs a ListSubsetValidator with the allowed elements.
func NewListSubset(allowed []string) *ListSubsetValidator {
	m := make(map[string]struct{}, len(allowed))
	disp := make([]string, 0, len(allowed))
	for _, v := range allowed {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		disp = append(disp, v)
	}
	return &ListSubsetValidator{allowed: m, display: disp}
}

// Ensure interface compliance for string collection validation through framework.
var _ frameworkvalidator.List = (*ListSubsetValidator)(nil)
var _ frameworkvalidator.Set = (*ListSubsetValidator)(nil)

func (v *ListSubsetValidator) Description(_ context.Context) string {
	return "all elements must be contained in the allowed set"
}

func (v *ListSubsetValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *ListSubsetValidator) ValidateList(ctx context.Context, req frameworkvalidator.ListRequest, resp *frameworkvalidator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var items []basetypes.StringValue
	if err := req.ConfigValue.ElementsAs(ctx, &items, false); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Collection", "Expected a list of strings.")
		return
	}

	var offenders []string
	for _, it := range items {
		if it.IsNull() || it.IsUnknown() {
			continue
		}
		if _, ok := v.allowed[it.ValueString()]; !ok {
			offenders = append(offenders, it.ValueString())
		}
	}
	if len(offenders) > 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Disallowed Elements", fmt.Sprintf("Elements not allowed: %v. Allowed: %v", offenders, v.display))
	}
}

func (v *ListSubsetValidator) ValidateSet(ctx context.Context, req frameworkvalidator.SetRequest, resp *frameworkvalidator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var items []basetypes.StringValue
	if err := req.ConfigValue.ElementsAs(ctx, &items, false); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Collection", "Expected a set of strings.")
		return
	}

	var offenders []string
	for _, it := range items {
		if it.IsNull() || it.IsUnknown() {
			continue
		}
		if _, ok := v.allowed[it.ValueString()]; !ok {
			offenders = append(offenders, it.ValueString())
		}
	}
	if len(offenders) > 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Disallowed Elements", fmt.Sprintf("Elements not allowed: %v. Allowed: %v", offenders, v.display))
	}
}
