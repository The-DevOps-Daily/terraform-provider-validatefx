package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type dependentValueFunction struct{}

var _ function.Function = (*dependentValueFunction)(nil)

// NewDependentValueFunction exposes the dependent_value validator as a Terraform function.
func NewDependentValueFunction() function.Function {
	return &dependentValueFunction{}
}

func (dependentValueFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "dependent_value"
}

func (dependentValueFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate a dependent relationship between two values.",
		MarkdownDescription: "Returns true when the dependent value is set if the condition value is set. If condition is empty, dependent can be empty or set.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "condition",
				MarkdownDescription: "The condition value to check.",
			},
			function.StringParameter{
				Name:                "dependent",
				MarkdownDescription: "The dependent value that must be set if condition is set.",
			},
		},
	}
}

func (dependentValueFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var condition basetypes.StringValue
	var dependent basetypes.StringValue

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &condition, &dependent))
	if resp.Error != nil {
		return
	}

	if condition.IsUnknown() || dependent.IsUnknown() {
		resp.Result = function.NewResultData(basetypes.NewBoolUnknown())
		return
	}

	// Treat null as empty string
	conditionVal := ""
	if !condition.IsNull() {
		conditionVal = condition.ValueString()
	}

	dependentVal := ""
	if !dependent.IsNull() {
		dependentVal = dependent.ValueString()
	}

	validator := validators.NewDependentValue()
	if err := validator.Validate(conditionVal, dependentVal); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
