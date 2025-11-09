package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type stringPrefixFunction struct{}

var _ function.Function = (*stringPrefixFunction)(nil)

// NewStringPrefixFunction exposes the string prefix validator as a Terraform function.
func NewStringPrefixFunction() function.Function {
	return &stringPrefixFunction{}
}

func (stringPrefixFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "has_prefix"
}

func (stringPrefixFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string starts with one of the provided prefixes.",
		MarkdownDescription: "Returns true when the input string starts with one of the provided prefixes.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "String value to validate.",
			},
			function.ListParameter{
				Name:                "prefixes",
				AllowNullValue:      false,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "List of prefixes to test against.",
			},
			function.BoolParameter{
				Name:                "ignore_case",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "Whether comparisons should ignore case.",
			},
		},
	}
}

func (stringPrefixFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	value, valueState, ok := stringArgument(ctx, req, resp, 0)
	if !ok {
		return
	}

	if valueState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	prefixes, prefixesState, ok := prefixesList(ctx, req, resp)
	if !ok {
		return
	}

	if prefixesState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	ignoreCase, ignoreState, ok := ignoreCaseFlag(ctx, req, resp, 2)
	if !ok {
		return
	}

	if ignoreState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	if len(prefixes) == 0 {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.StringPrefix(prefixes, ignoreCase)
	validation := frameworkvalidator.StringResponse{}
	validator.ValidateString(ctx, frameworkvalidator.StringRequest{
		ConfigValue: value,
		Path:        path.Root("value"),
	}, &validation)

	if validation.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, validation.Diagnostics)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

func prefixesList(ctx context.Context, req function.RunRequest, resp *function.RunResponse) ([]string, valueState, bool) {
	var prefixes types.List
	if err := req.Arguments.GetArgument(ctx, 1, &prefixes); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return nil, valueKnown, false
	}

	values, state, funcErr := preparePrefixValues(ctx, prefixes)
	if funcErr != nil {
		resp.Error = funcErr
		return nil, valueKnown, false
	}

	return values, state, true
}

func preparePrefixValues(ctx context.Context, list types.List) ([]string, valueState, *function.FuncError) {
	if list.IsUnknown() {
		return nil, valueUnknown, nil
	}

	if list.IsNull() {
		return nil, valueKnown, function.NewFuncError("prefixes list must be provided")
	}

	var items []basetypes.StringValue
	diags := list.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		diags.AddAttributeError(
			path.Root("prefixes"),
			"Invalid Prefixes",
			"Prefixes must be provided as a list of strings.",
		)
		return nil, valueKnown, function.FuncErrorFromDiags(ctx, diags)
	}

	values := make([]string, 0, len(items))
	for _, item := range items {
		if item.IsNull() || item.IsUnknown() {
			continue
		}

		values = append(values, item.ValueString())
	}

	return values, valueKnown, nil
}
