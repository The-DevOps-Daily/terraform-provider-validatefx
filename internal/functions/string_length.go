package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type stringLengthFunction struct{}

var _ function.Function = (*stringLengthFunction)(nil)

// NewStringLengthFunction exposes the string length validator as a Terraform function.
func NewStringLengthFunction() function.Function {
	return &stringLengthFunction{}
}

func (stringLengthFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "string_length"
}

func (stringLengthFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string length falls within optional minimum and maximum bounds.",
		MarkdownDescription: "Returns true when the input string length is within the provided minimum and/or maximum bounds.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "String value to validate.",
			},
			function.Int64Parameter{
				Name:                "min_length",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "Optional minimum length (inclusive).",
			},
			function.Int64Parameter{
				Name:                "max_length",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "Optional maximum length (inclusive).",
			},
		},
	}
}

func (stringLengthFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value types.String
	var min types.Int64
	var max types.Int64

	if err := req.Arguments.GetArgument(ctx, 0, &value); err != nil {
		resp.Error = err
		return
	}

	if err := req.Arguments.GetArgument(ctx, 1, &min); err != nil {
		resp.Error = err
		return
	}

	if err := req.Arguments.GetArgument(ctx, 2, &max); err != nil {
		resp.Error = err
		return
	}

	if value.IsNull() || value.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validation := frameworkvalidator.StringResponse{}
	validateStringLength(ctx, value, min, max, &validation)

	if validation.Diagnostics.HasError() {
		diags := diag.Diagnostics{}
		diags.Append(validation.Diagnostics...)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

func validateStringLength(ctx context.Context, value types.String, min types.Int64, max types.Int64, validation *frameworkvalidator.StringResponse) {
	var minPtr *int
	var maxPtr *int

	if !min.IsNull() && !min.IsUnknown() {
		v := int(min.ValueInt64())
		minPtr = &v
	}

	if !max.IsNull() && !max.IsUnknown() {
		v := int(max.ValueInt64())
		maxPtr = &v
	}

	validator := validators.NewStringLengthValidator(minPtr, maxPtr)
	validator.ValidateString(ctx, frameworkvalidator.StringRequest{
		ConfigValue: value,
		Path:        path.Root("value"),
	}, validation)
}
