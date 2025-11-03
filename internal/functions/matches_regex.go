package functions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type matchesRegexFunction struct{}

var _ function.Function = (*matchesRegexFunction)(nil)

// NewMatchesRegexFunction exposes a regex matching helper.
func NewMatchesRegexFunction() function.Function {
	return &matchesRegexFunction{}
}

func (matchesRegexFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "matches_regex"
}

func (matchesRegexFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string matches a provided regular expression.",
		MarkdownDescription: "Returns true when the input string matches the supplied regular expression pattern.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "String value to validate.",
				MarkdownDescription: "String value to validate.",
			},
			function.StringParameter{
				Name:                "pattern",
				Description:         "Regular expression pattern to apply.",
				MarkdownDescription: "Regular expression pattern to apply.",
			},
		},
	}
}

func (matchesRegexFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value types.String
	var pattern types.String

	if err := req.Arguments.GetArgument(ctx, 0, &value); err != nil {
		resp.Error = err
		return
	}

	if err := req.Arguments.GetArgument(ctx, 1, &pattern); err != nil {
		resp.Error = err
		return
	}

	if value.IsNull() || value.IsUnknown() || pattern.IsNull() || pattern.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	reCompiled, err := regexp.Compile(pattern.ValueString())
	if err != nil {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("pattern"),
			"Invalid Regex Pattern",
			err.Error(),
		)

		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	if !reCompiled.MatchString(value.ValueString()) {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("value"),
			"Pattern Mismatch",
			"Value must match the provided regular expression pattern.",
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
