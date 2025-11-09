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

type listSubsetFunction struct{}

var _ function.Function = (*listSubsetFunction)(nil)

// NewListSubsetFunction exposes the list_subset validator as a Terraform function.
func NewListSubsetFunction() function.Function {
	return &listSubsetFunction{}
}

func (listSubsetFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "list_subset"
}

func (listSubsetFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that every element of a list belongs to the allowed collection.",
		MarkdownDescription: "Returns true when all values are present in the allowed list; returns an error otherwise.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "values",
				ElementType:         basetypes.StringType{},
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "List of string values to validate.",
				MarkdownDescription: "List of string values to validate.",
			},
			function.ListParameter{
				Name:                "allowed",
				ElementType:         basetypes.StringType{},
				AllowNullValue:      false,
				AllowUnknownValues:  true,
				Description:         "List of allowed string values.",
				MarkdownDescription: "List of allowed string values.",
			},
		},
	}
}

func (listSubsetFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	values, valuesState, ok := listArg(ctx, req, resp, 0)
	if !ok {
		return
	}

	if valuesState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	allowedList, allowedValues, allowedState, ok := listArgWithStrings(ctx, req, resp, 1)
	if !ok {
		return
	}

	if allowedState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	if len(allowedValues) == 0 {
		resp.Error = function.NewFuncError("allowed values list must contain at least one item")
		return
	}

	validator := validators.ListSubset(allowedValues)
	result := &frameworkvalidator.ListResponse{}

	validator.ValidateList(ctx, frameworkvalidator.ListRequest{
		ConfigValue: values,
		Path:        path.Root("values"),
	}, result)

	if result.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, result.Diagnostics)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

type valueState int

const (
	valueKnown valueState = iota
	valueUnknown
)

func listArg(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int) (types.List, valueState, bool) {
	var list types.List
	if err := req.Arguments.GetArgument(ctx, index, &list); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return types.List{}, valueKnown, false
	}

	if list.IsUnknown() {
		return list, valueUnknown, true
	}

	return list, valueKnown, true
}

func listArgWithStrings(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int) (types.List, []string, valueState, bool) {
	list, state, ok := listArg(ctx, req, resp, index)
	if !ok {
		return types.List{}, nil, valueKnown, false
	}

	if state == valueUnknown {
		return list, nil, valueUnknown, true
	}

	if list.IsNull() {
		return list, nil, valueKnown, true
	}

	var values []string
	if diags := list.ElementsAs(ctx, &values, false); diags.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return types.List{}, nil, valueKnown, false
	}

	return list, values, valueKnown, true
}
