package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
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
	values, valuesState, ok := collectStringList(ctx, req, resp, 0, "values")
	if !ok {
		return
	}

	if valuesState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	expected, expectedState, ok := collectStringList(ctx, req, resp, 1, "expected")
	if !ok {
		return
	}

	if expectedState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.NewSetEquals(expected)
	if err := validator.Validate(values); err != nil {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("values"),
			"Set Mismatch",
			err.Error(),
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

func collectStringList(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int, name string) ([]string, valueState, bool) {
	var list types.List
	if err := req.Arguments.GetArgument(ctx, index, &list); err != nil {
		resp.Error = err
		return nil, valueKnown, false
	}

	if list.IsUnknown() {
		return nil, valueUnknown, true
	}

	if list.IsNull() {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root(name),
			"Missing List",
			"List must be provided.",
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return nil, valueKnown, false
	}

	var elements []basetypes.StringValue
	diags := list.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		diags.AddAttributeError(
			path.Root(name),
			"Invalid Elements",
			"List must contain only strings.",
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return nil, valueKnown, false
	}

	values := make([]string, 0, len(elements))
	for _, el := range elements {
		if el.IsUnknown() {
			return nil, valueUnknown, true
		}
		if el.IsNull() {
			diags := diag.Diagnostics{}
			diags.AddAttributeError(
				path.Root(name),
				"Null Element",
				"List must not contain null values.",
			)
			resp.Error = function.FuncErrorFromDiags(ctx, diags)
			return nil, valueKnown, false
		}
		values = append(values, el.ValueString())
	}

	return values, valueKnown, true
}
