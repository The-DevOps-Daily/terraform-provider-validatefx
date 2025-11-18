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

type stringContainsFunction struct{}

var _ function.Function = (*stringContainsFunction)(nil)

// NewStringContainsFunction exposes the string_contains validator as a Terraform function.
func NewStringContainsFunction() function.Function {
	return &stringContainsFunction{}
}

func (stringContainsFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "string_contains"
}

func (stringContainsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string contains at least one of the provided substrings.",
		MarkdownDescription: "Returns true when the input string contains one of the provided substrings.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "String value to validate.",
			},
			function.ListParameter{
				Name:                "substrings",
				AllowNullValue:      false,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "List of substrings to test against.",
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

func (stringContainsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	value, valueState, ok := stringArgument(ctx, req, resp, 0)
	if !ok {
		return
	}

	if valueState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	substrings, substringsState, ok := stringListArgument(ctx, req, resp, 1, "substrings")
	if !ok {
		return
	}

	if substringsState == valueUnknown {
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

	if len(substrings) == 0 {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.StringContains(substrings, ignoreCase)
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
