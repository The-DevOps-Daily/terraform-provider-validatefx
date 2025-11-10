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

type hexFunction struct{}

var _ function.Function = (*hexFunction)(nil)

func NewHexFunction() function.Function { return &hexFunction{} }

func (hexFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "hex"
}

func (hexFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string contains only hexadecimal characters.",
		MarkdownDescription: "Returns true when the input string contains only hexadecimal characters.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:               "value",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				Description:        "String value to validate.",
			},
		},
	}
}

func (hexFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	value, state, ok := stringArgument(ctx, req, resp, 0)
	if !ok {
		return
	}
	if state == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.Hex()
	res := frameworkvalidator.StringResponse{}
	validator.ValidateString(ctx, frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: value}, &res)
	if res.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, res.Diagnostics)
		return
	}
	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
