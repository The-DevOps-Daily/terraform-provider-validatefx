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

type notInListFunction struct{}

var _ function.Function = (*notInListFunction)(nil)

// NewNotInListFunction exposes the not_in_list validator as a Terraform function.
func NewNotInListFunction() function.Function { return &notInListFunction{} }

func (notInListFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "not_in_list"
}

func (notInListFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string does not match any of the provided disallowed values.",
		MarkdownDescription: "Returns true when the input string is not present in the list of disallowed values.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:               "value",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				Description:        "String value to validate.",
			},
			function.ListParameter{
				Name:               "disallowed",
				AllowNullValue:     false,
				AllowUnknownValues: true,
				ElementType:        basetypes.StringType{},
				Description:        "List of disallowed string values.",
			},
			function.BoolParameter{
				Name:               "ignore_case",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				Description:        "Whether comparisons are case-insensitive.",
			},
		},
	}
}

func (notInListFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	value, state, ok := stringArgument(ctx, req, resp, 0)
	if !ok {
		return
	}
	if state == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	disallowed, state, ok := allowedList(ctx, req, resp, 1)
	if !ok {
		return
	}
	if state == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	ignore, state, ok := ignoreCaseFlag(ctx, req, resp, 2)
	if !ok {
		return
	}
	if state == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	if len(disallowed) == 0 {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := validators.NewNotInListValidator(disallowed, ignore)
	validation := frameworkvalidator.StringResponse{}
	validator.ValidateString(ctx, frameworkvalidator.StringRequest{
		Path:        path.Root("value"),
		ConfigValue: value,
	}, &validation)

	if validation.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, validation.Diagnostics)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
