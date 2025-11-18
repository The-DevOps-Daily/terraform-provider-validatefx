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

type arnFunction struct{}

var _ function.Function = (*arnFunction)(nil)

func NewARNFunction() function.Function { return &arnFunction{} }

func (arnFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "arn"
}

func (arnFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string is an AWS ARN.",
		MarkdownDescription: "Returns true when the string is a valid AWS ARN.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:               "value",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				Description:        "String value to validate as AWS ARN.",
			},
		},
	}
}

func (arnFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input types.String
	if err := req.Arguments.GetArgument(ctx, 0, &input); err != nil {
		resp.Error = err
		return
	}
	if input.IsNull() || input.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	v := validators.ARN()
	r := frameworkvalidator.StringResponse{}
	v.ValidateString(ctx, frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: input}, &r)
	if r.Diagnostics.HasError() {
		diags := diag.Diagnostics{}
		diags.Append(r.Diagnostics...)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}
	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
