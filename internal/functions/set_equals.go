package functions

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type setEqualsFunction struct{}

var _ function.Function = (*setEqualsFunction)(nil)

// NewSetEqualsFunction exposes the set equivalence helper as a Terraform function.
func NewSetEqualsFunction() function.Function {
	return &setEqualsFunction{}
}

func (setEqualsFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "set_equals"
}

func (setEqualsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that two string lists contain the same elements regardless of order.",
		MarkdownDescription: "Returns true when both string lists contain the same unique elements. Raises an error when they differ.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "values",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "Source list of strings to evaluate.",
			},
			function.ListParameter{
				Name:                "expected",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "List of strings that should match the source list.",
			},
		},
	}
}

func (setEqualsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	values, valuesState, ok := collectStringSet(ctx, req, resp, 0, "values")
	if !ok {
		return
	}

	if valuesState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	expected, expectedState, ok := collectStringSet(ctx, req, resp, 1, "expected")
	if !ok {
		return
	}

	if expectedState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	if len(values.items) != len(expected.items) || !setsEqual(values.items, expected.items) {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("values"),
			"Set Mismatch",
			fmt.Sprintf("Values %v must match expected %v (order independent).", values.ordered, expected.ordered),
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

type stringSet struct {
	items   map[string]struct{}
	ordered []string
}

func collectStringSet(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int, name string) (stringSet, valueState, bool) {
	var list types.List
	if err := req.Arguments.GetArgument(ctx, index, &list); err != nil {
		resp.Error = err
		return stringSet{}, valueKnown, false
	}

	if list.IsUnknown() {
		return stringSet{}, valueUnknown, true
	}

	if list.IsNull() {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root(name),
			"Missing List",
			fmt.Sprintf("Parameter %q must be provided.", name),
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return stringSet{}, valueKnown, false
	}

	var elements []basetypes.StringValue
	diags := list.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		diags.AddAttributeError(
			path.Root(name),
			"Invalid Elements",
			fmt.Sprintf("Parameter %q must be a list of strings.", name),
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return stringSet{}, valueKnown, false
	}

	set := stringSet{
		items:   make(map[string]struct{}, len(elements)),
		ordered: make([]string, 0, len(elements)),
	}

	for _, element := range elements {
		if element.IsUnknown() {
			return stringSet{}, valueUnknown, true
		}

		if element.IsNull() {
			diags := diag.Diagnostics{}
			diags.AddAttributeError(
				path.Root(name),
				"Null Element",
				fmt.Sprintf("Parameter %q must not contain null values.", name),
			)
			resp.Error = function.FuncErrorFromDiags(ctx, diags)
			return stringSet{}, valueKnown, false
		}

		value := element.ValueString()
		if _, exists := set.items[value]; !exists {
			set.items[value] = struct{}{}
			set.ordered = append(set.ordered, value)
		}
	}

	sort.Strings(set.ordered)

	return set, valueKnown, true
}

func setsEqual(left, right map[string]struct{}) bool {
	for key := range left {
		if _, ok := right[key]; !ok {
			return false
		}
	}

	for key := range right {
		if _, ok := left[key]; !ok {
			return false
		}
	}

	return true
}
