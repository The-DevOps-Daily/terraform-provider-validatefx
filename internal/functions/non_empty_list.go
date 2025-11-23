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

type nonEmptyListFunction struct{}

var _ function.Function = (*nonEmptyListFunction)(nil)

// NewNonEmptyListFunction exposes the non-empty list validator as a Terraform function.
func NewNonEmptyListFunction() function.Function {
	return &nonEmptyListFunction{}
}

func (nonEmptyListFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "non_empty_list"
}

func (nonEmptyListFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a list is not empty.",
		MarkdownDescription: "Returns true when the list contains at least one element. Raises an error when the list is empty.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "values",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "List to validate.",
			},
		},
	}
}

func (nonEmptyListFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	values, valuesState, ok := stringListArgument(ctx, req, resp, 0, "values")
	if !ok {
		return
	}

	if valuesState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.NewNonEmptyList()
	if err := validator.Validate(values); err != nil {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("values"),
			"Empty List",
			err.Error(),
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
