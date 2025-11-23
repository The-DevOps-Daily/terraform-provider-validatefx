package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type mutuallyExclusiveFunction struct{}

var _ function.Function = (*mutuallyExclusiveFunction)(nil)

// NewMutuallyExclusiveFunction exposes the mutually exclusive validator as a Terraform function.
func NewMutuallyExclusiveFunction() function.Function {
	return &mutuallyExclusiveFunction{}
}

func (mutuallyExclusiveFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "mutually_exclusive"
}

func (mutuallyExclusiveFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that exactly one value is set.",
		MarkdownDescription: "Returns true when exactly one value from the list is set (non-empty). Raises an error when zero or multiple values are set.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "values",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "List of values to validate.",
			},
		},
	}
}

func (mutuallyExclusiveFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	values, valuesState, ok := stringListArgument(ctx, req, resp, 0, "values")
	if !ok {
		return
	}

	if valuesState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.NewMutuallyExclusive()
	if err := validator.Validate(values); err != nil {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("values"),
			"Mutually Exclusive Validation Failed",
			err.Error(),
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
