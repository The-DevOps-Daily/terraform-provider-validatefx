package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type listUniqueFunction struct{}

var _ function.Function = (*listUniqueFunction)(nil)

// NewListUniqueFunction exposes the list_unique validator as a Terraform function.
func NewListUniqueFunction() function.Function {
	return &listUniqueFunction{}
}

func (listUniqueFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "list_unique"
}

func (listUniqueFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that all list elements are unique.",
		MarkdownDescription: "Returns true when all elements in the list are unique (no duplicates).",
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

func (listUniqueFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var values basetypes.ListValue

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &values))
	if resp.Error != nil {
		return
	}

	if values.IsNull() || values.IsUnknown() {
		resp.Result = function.NewResultData(basetypes.NewBoolUnknown())
		return
	}

	validator := validators.NewListUnique()
	validation := frameworkvalidator.ListResponse{}
	validator.ValidateList(ctx, frameworkvalidator.ListRequest{ConfigValue: values, Path: path.Root("values")}, &validation)

	if validation.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, validation.Diagnostics)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
