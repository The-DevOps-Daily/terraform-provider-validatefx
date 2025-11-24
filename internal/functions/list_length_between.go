package functions

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type listLengthBetweenFunction struct{}

var _ function.Function = (*listLengthBetweenFunction)(nil)

// NewListLengthBetweenFunction exposes the list_length_between validator as a Terraform function.
func NewListLengthBetweenFunction() function.Function {
	return &listLengthBetweenFunction{}
}

func (listLengthBetweenFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "list_length_between"
}

func (listLengthBetweenFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a list has a length between minimum and maximum bounds.",
		MarkdownDescription: "Returns true when the list length is within the specified range (inclusive).",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:               "values",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				ElementType:        basetypes.StringType{},
				MarkdownDescription: "List to validate.",
			},
			function.StringParameter{
				Name:                "min",
				MarkdownDescription: "Minimum length (inclusive).",
			},
			function.StringParameter{
				Name:                "max",
				MarkdownDescription: "Maximum length (inclusive).",
			},
		},
	}
}

func (listLengthBetweenFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var values basetypes.ListValue
	var minStr basetypes.StringValue
	var maxStr basetypes.StringValue

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &values, &minStr, &maxStr))
	if resp.Error != nil {
		return
	}

	if values.IsNull() || values.IsUnknown() || minStr.IsUnknown() || maxStr.IsUnknown() {
		resp.Result = function.NewResultData(basetypes.NewBoolUnknown())
		return
	}

	min, err := strconv.Atoi(minStr.ValueString())
	if err != nil {
		resp.Error = function.NewFuncError("min must be a valid integer")
		return
	}

	max, err := strconv.Atoi(maxStr.ValueString())
	if err != nil {
		resp.Error = function.NewFuncError("max must be a valid integer")
		return
	}

	if min < 0 || max < 0 {
		resp.Error = function.NewFuncError("min and max must be non-negative")
		return
	}

	if min > max {
		resp.Error = function.NewFuncError("min must be less than or equal to max")
		return
	}

	validator := validators.NewListLengthBetween(min, max)
	validation := frameworkvalidator.ListResponse{}
	validator.ValidateList(ctx, frameworkvalidator.ListRequest{ConfigValue: values, Path: path.Root("values")}, &validation)

	if validation.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, validation.Diagnostics)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
