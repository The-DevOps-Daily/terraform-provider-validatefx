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

type betweenFunction struct{}

var _ function.Function = (*betweenFunction)(nil)

func NewBetweenFunction() function.Function {
	return &betweenFunction{}
}

func (betweenFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "between"
}

func (betweenFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a numeric string falls between inclusive minimum and maximum bounds.",
		MarkdownDescription: "Returns true when the input value is a valid decimal greater than or equal to the minimum and less than or equal to the maximum.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "Numeric string to validate.",
				MarkdownDescription: "Numeric string to validate.",
			},
			function.StringParameter{
				Name:                "min",
				Description:         "Inclusive minimum value as a decimal string.",
				MarkdownDescription: "Inclusive minimum value as a decimal string.",
			},
			function.StringParameter{
				Name:                "max",
				Description:         "Inclusive maximum value as a decimal string.",
				MarkdownDescription: "Inclusive maximum value as a decimal string.",
			},
		},
	}
}

func (betweenFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var (
		value types.String
		min   types.String
		max   types.String
	)

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

	validator := validators.Between(stringFrom(min), stringFrom(max))
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

func stringFrom(v types.String) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}

	return v.ValueString()
}
