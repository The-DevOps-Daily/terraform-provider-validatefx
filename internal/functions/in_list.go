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

type inListFunction struct{}

var _ function.Function = (*inListFunction)(nil)

// NewInListFunction exposes the in_list validator as a Terraform function.
func NewInListFunction() function.Function {
	return &inListFunction{}
}

func (inListFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "in_list"
}

func (inListFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string matches one of the allowed values.",
		MarkdownDescription: "Returns true when the input string equals one of the supplied allowed values.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "String value to validate.",
				MarkdownDescription: "String value to validate.",
			},
			function.ListParameter{
				Name:                "allowed",
				AllowNullValue:      false,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				Description:         "List of allowed string values.",
				MarkdownDescription: "List of allowed string values.",
			},
			function.BoolParameter{
				Name:                "ignore_case",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "Whether comparisons are case-insensitive.",
				MarkdownDescription: "Whether comparisons are case-insensitive.",
			},
		},
	}
}

func (inListFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value types.String
	var allowed types.List
	var ignoreCase types.Bool

	if err := req.Arguments.GetArgument(ctx, 0, &value); err != nil {
		resp.Error = err
		return
	}

	if err := req.Arguments.GetArgument(ctx, 1, &allowed); err != nil {
		resp.Error = err
		return
	}

	if err := req.Arguments.GetArgument(ctx, 2, &ignoreCase); err != nil {
		resp.Error = err
		return
	}

	if value.IsNull() || value.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	if allowed.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	if allowed.IsNull() {
		resp.Error = function.NewFuncError("allowed values list must be provided")
		return
	}

	if ignoreCase.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	var ignore bool
	if ignoreCase.IsNull() {
		ignore = false
	} else {
		ignore = ignoreCase.ValueBool()
	}

	var allowedValues []basetypes.StringValue
	if diags := allowed.ElementsAs(ctx, &allowedValues, false); diags.HasError() {
		diags.AddAttributeError(
			path.Root("allowed"),
			"Invalid Allowed Values",
			"Allowed values must be provided as a list of strings.",
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	candidates := make([]string, 0, len(allowedValues))
	for _, item := range allowedValues {
		if item.IsNull() || item.IsUnknown() {
			continue
		}
		candidates = append(candidates, item.ValueString())
	}

	if len(candidates) == 0 {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.NewInListValidator(candidates, ignore)

	validation := frameworkvalidator.StringResponse{}
	validator.ValidateString(ctx, frameworkvalidator.StringRequest{
		Path:        path.Root("value"),
		ConfigValue: value,
	}, &validation)

	if validation.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, validation.Diagnostics)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
