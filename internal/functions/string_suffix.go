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

type stringSuffixFunction struct{}

var _ function.Function = (*stringSuffixFunction)(nil)

// NewStringSuffixFunction exposes the string suffix validator as a Terraform function.
func NewStringSuffixFunction() function.Function {
	return &stringSuffixFunction{}
}

func (stringSuffixFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "has_suffix"
}

func (stringSuffixFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string ends with one of the provided suffixes.",
		MarkdownDescription: "Returns true when the input string ends with one of the provided suffixes.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "String value to validate.",
			},
			function.ListParameter{
				Name:                "suffixes",
				AllowNullValue:      false,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "List of suffixes to check.",
			},
		},
	}
}

func (stringSuffixFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	value, valueState, ok := stringArgument(ctx, req, resp, 0)
	if !ok {
		return
	}

	if valueState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	suffixes, suffixState, ok := suffixList(ctx, req, resp)
	if !ok {
		return
	}

	if suffixState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	if len(suffixes) == 0 {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.StringSuffix(suffixes...)
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

func suffixList(ctx context.Context, req function.RunRequest, resp *function.RunResponse) ([]string, valueState, bool) {
	var suffixes types.List
	if err := req.Arguments.GetArgument(ctx, 1, &suffixes); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return nil, valueKnown, false
	}

	values, state, funcErr := prepareSuffixValues(ctx, suffixes)
	if funcErr != nil {
		resp.Error = funcErr
		return nil, valueKnown, false
	}

	return values, state, true
}

func prepareSuffixValues(ctx context.Context, list types.List) ([]string, valueState, *function.FuncError) {
	if list.IsUnknown() {
		return nil, valueUnknown, nil
	}

	if list.IsNull() {
		return nil, valueKnown, function.NewFuncError("suffixes list must be provided")
	}

	var items []basetypes.StringValue
	diags := list.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		diags.AddAttributeError(
			path.Root("suffixes"),
			"Invalid Suffixes",
			"Suffixes must be provided as a list of strings.",
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
