package functions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type matchesRegexFunction struct{}

var _ function.Function = (*matchesRegexFunction)(nil)

// NewMatchesRegexFunction exposes a regex matching validator.
func NewMatchesRegexFunction() function.Function { return &matchesRegexFunction{} }

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
				Name:               "value",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				Description:        "String value to validate.",
			},
			function.StringParameter{
				Name:        "pattern",
				Description: "Regular expression pattern to apply.",
			},
		},
	}
}

func (matchesRegexFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Pull the pattern argument first and handle unknown/null
	var pattern basetypes.StringValue
	if err := req.Arguments.GetArgument(ctx, 1, &pattern); err != nil {
		resp.Error = err
		return
	}

	if pattern.IsNull() || pattern.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	// Validate the regex pattern and surface errors on the pattern attribute
	if _, err := regexp.Compile(pattern.ValueString()); err != nil {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("pattern"),
			"Invalid Regex Pattern",
			err.Error(),
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	// Use the shared single-string validation wrapper with a dynamic validator
	fn := newStringValidationFunction(
		"matches_regex",
		"Validate that a string matches a provided regular expression.",
		"Returns true when the input string matches the supplied regular expression pattern.",
		validators.MatchesRegex(pattern.ValueString()),
	)

	// Delegate to the common single string validator wrapper.
	fn.Run(ctx, req, resp)
}
