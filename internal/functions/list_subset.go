package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type listSubsetFunction struct{}

var _ function.Function = (*listSubsetFunction)(nil)

func NewListSubsetFunction() function.Function { return &listSubsetFunction{} }

func (listSubsetFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "list_subset"
}

func (listSubsetFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that all elements of a list/set are contained in a reference list.",
		MarkdownDescription: "Returns true when every element of the input collection exists in the allowed collection.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:               "values",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				ElementType:        basetypes.StringType{},
			},
			function.ListParameter{
				Name:               "allowed",
				AllowNullValue:     false,
				AllowUnknownValues: true,
				ElementType:        basetypes.StringType{},
			},
		},
	}
}

func (listSubsetFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	values, vState, ok := getListArg(ctx, req, resp, 0, "values")
	if !ok {
		return
	}
	allowed, aState, ok := getListArg(ctx, req, resp, 1, "allowed")
	if !ok {
		return
	}

	if vState || aState { // unknown states
		resp.Result = function.NewResultData(basetypes.NewBoolUnknown())
		return
	}

	if allowed.IsNull() {
		resp.Error = function.NewFuncError("allowed list must be provided")
		return
	}

	var allowedItems []basetypes.StringValue
	d := allowed.ElementsAs(ctx, &allowedItems, false)
	if d.HasError() {
		d.AddAttributeError(path.Root("allowed"), "Invalid Allowed Values", "Allowed must be a list of strings.")
		resp.Error = function.FuncErrorFromDiags(ctx, d)
		return
	}
	arr := collectStrings(allowedItems)

	validator := validators.NewListSubset(arr)
	validation := frameworkvalidator.ListResponse{}
	validator.ValidateList(ctx, frameworkvalidator.ListRequest{ConfigValue: values, Path: path.Root("values")}, &validation)
	if validation.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, validation.Diagnostics)
		return
	}
	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

func getListArg(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int, name string) (basetypes.ListValue, bool, bool) {
	var v basetypes.ListValue
	if err := req.Arguments.GetArgument(ctx, index, &v); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return basetypes.ListValue{}, false, false
	}
	if v.IsUnknown() {
		return v, true, true
	}
	return v, false, true
}

func collectStrings(items []basetypes.StringValue) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		if it.IsNull() || it.IsUnknown() {
			continue
		}
		out = append(out, it.ValueString())
	}
	return out
}
