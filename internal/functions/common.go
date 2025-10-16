package functions

import (
	"context"

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
		resp.Result = function.NewResultData(basetypes.NewBoolValue(false))
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
