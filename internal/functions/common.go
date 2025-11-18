package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemavalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type stringValidationFunction struct {
	name        string
	summary     string
	description string
	validator   schemavalidator.String
}

var _ function.Function = (*stringValidationFunction)(nil)

func newStringValidationFunction(name, summary, description string, v schemavalidator.String) function.Function {
	return &stringValidationFunction{
		name:        name,
		summary:     summary,
		description: description,
		validator:   v,
	}
}

func (f *stringValidationFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
}

func (f *stringValidationFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             f.summary,
		MarkdownDescription: f.description,
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "String value to validate.",
				MarkdownDescription: "String value to validate.",
			},
		},
	}
}

func (f *stringValidationFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input types.String

	if err := req.Arguments.GetArgument(ctx, 0, &input); err != nil {
		resp.Error = err
		return
	}

	if input.IsNull() || input.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validation := schemavalidator.StringResponse{}

	f.validator.ValidateString(ctx, schemavalidator.StringRequest{
		ConfigValue: input,
		Path:        path.Root("value"),
	}, &validation)

	if validation.Diagnostics.HasError() {
		diags := diag.Diagnostics{}
		diags.Append(validation.Diagnostics...)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

// stringFrom safely converts a Terraform framework String to a plain string,
// returning an empty string for null/unknown values. Useful for optional
// parameters passed as strings.
func stringFrom(v types.String) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}
	return v.ValueString()
}

// stringListArgument extracts a list of strings from function arguments at the specified index.
// It handles null, unknown, and invalid list scenarios consistently.
// Returns the string slice, value state, and success boolean.
func stringListArgument(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
	index int,
	paramName string,
) ([]string, valueState, bool) {
	var list types.List
	if err := req.Arguments.GetArgument(ctx, index, &list); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return nil, valueKnown, false
	}

	if list.IsUnknown() {
		return nil, valueUnknown, true
	}

	if list.IsNull() {
		resp.Error = function.NewFuncError(paramName + " list must be provided")
		return nil, valueKnown, false
	}

	var items []basetypes.StringValue
	diags := list.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		diags.AddAttributeError(
			path.Root(paramName),
			"Invalid "+paramName,
			"Must be provided as a list of strings.",
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return nil, valueKnown, false
	}

	values := make([]string, 0, len(items))
	for _, item := range items {
		if item.IsNull() || item.IsUnknown() {
			continue
		}
		values = append(values, item.ValueString())
	}

	return values, valueKnown, true
}
