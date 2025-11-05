package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type compositeValidationFunction struct {
	name        string
	summary     string
	description string
	allMustPass bool
}

func newCompositeValidationFunction(name, summary, description string, allMustPass bool) function.Function {
	return &compositeValidationFunction{
		name:        name,
		summary:     summary,
		description: description,
		allMustPass: allMustPass,
	}
}

func (f *compositeValidationFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
}

func (f *compositeValidationFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             f.summary,
		MarkdownDescription: f.description,
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "checks",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.BoolType{},
				Description:         "List of boolean validation results to evaluate.",
				MarkdownDescription: "List of boolean validation results to evaluate.",
			},
		},
	}
}

func (f *compositeValidationFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var checks types.List

	if err := req.Arguments.GetArgument(ctx, 0, &checks); err != nil {
		resp.Error = err
		return
	}

	if checks.IsNull() || checks.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	var bools []basetypes.BoolValue
	if diags := checks.ElementsAs(ctx, &bools, false); diags.HasError() {
		diags.AddAttributeError(
			path.Root("checks"),
			"Invalid Boolean",
			"List elements must be boolean validation results.",
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	if len(bools) == 0 {
		resp.Result = function.NewResultData(basetypes.NewBoolValue(f.allMustPass))
		return
	}

	eval := summarizeBoolValues(bools)
	result := f.resolveCompositeOutcome(eval)

	resp.Result = function.NewResultData(result)
}

type boolEvaluation struct {
	anyTrue    bool
	anyFalse   bool
	anyUnknown bool
}

func summarizeBoolValues(values []basetypes.BoolValue) boolEvaluation {
	var eval boolEvaluation

	for _, value := range values {
		switch {
		case value.IsUnknown():
			eval.anyUnknown = true
		case value.IsNull():
			eval.anyFalse = true
		case value.ValueBool():
			eval.anyTrue = true
		default:
			eval.anyFalse = true
		}
	}

	return eval
}

func (f *compositeValidationFunction) resolveCompositeOutcome(eval boolEvaluation) basetypes.BoolValue {
	if f.allMustPass {
		if eval.anyFalse {
			return basetypes.NewBoolValue(false)
		}
		if eval.anyUnknown {
			return basetypes.NewBoolUnknown()
		}
		return basetypes.NewBoolValue(true)
	}

	if eval.anyTrue {
		return basetypes.NewBoolValue(true)
	}

	if eval.anyUnknown {
		return basetypes.NewBoolUnknown()
	}

	return basetypes.NewBoolValue(false)
}

type exactlyOneEvaluation struct {
	trueCount int
	unknown   bool
}

func evaluateExactlyOne(values []basetypes.BoolValue) exactlyOneEvaluation {
	var eval exactlyOneEvaluation

	for _, value := range values {
		switch {
		case value.IsUnknown():
			eval.unknown = true
			if eval.trueCount > 0 {
				return eval
			}
		case value.IsNull():
			continue
		case value.ValueBool():
			eval.trueCount++
			if eval.trueCount > 1 {
				return eval
			}
		}
	}

	return eval
}

func NewAllValidFunction() function.Function {
	return newCompositeValidationFunction(
		"all_valid",
		"Return true when all provided validation checks evaluate to true.",
		"Accepts a list of boolean validation results and returns true only when every element is true.",
		true,
	)
}

func NewAnyValidFunction() function.Function {
	return newCompositeValidationFunction(
		"any_valid",
		"Return true when any provided validation check evaluates to true.",
		"Accepts a list of boolean validation results and returns true when at least one element is true.",
		false,
	)
}

func NewExactlyOneValidFunction() function.Function {
	return &exactlyOneValidFunction{}
}

type exactlyOneValidFunction struct{}

var _ function.Function = (*exactlyOneValidFunction)(nil)

func (exactlyOneValidFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "exactly_one_valid"
}

func (exactlyOneValidFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Return true when exactly one validation check evaluates to true.",
		MarkdownDescription: "Accepts a list of boolean validation results and returns true only when exactly one element is true.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "checks",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.BoolType{},
				Description:         "List of boolean validation results to evaluate.",
				MarkdownDescription: "List of boolean validation results to evaluate.",
			},
		},
	}
}

func (exactlyOneValidFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var checks types.List

	if err := req.Arguments.GetArgument(ctx, 0, &checks); err != nil {
		resp.Error = err
		return
	}

	if checks.IsNull() || checks.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	var bools []basetypes.BoolValue
	if diags := checks.ElementsAs(ctx, &bools, false); diags.HasError() {
		diags.AddAttributeError(
			path.Root("checks"),
			"Invalid Boolean",
			"List elements must be boolean validation results.",
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	if len(bools) == 0 {
		resp.Result = function.NewResultData(basetypes.NewBoolValue(false))
		return
	}

	eval := evaluateExactlyOne(bools)

	switch {
	case eval.trueCount > 1:
		resp.Result = function.NewResultData(basetypes.NewBoolValue(false))
	case eval.trueCount == 1 && !eval.unknown:
		resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
	case eval.trueCount == 1 && eval.unknown:
		resp.Result = function.NewResultData(basetypes.NewBoolValue(false))
	case eval.trueCount == 0 && eval.unknown:
		resp.Result = function.NewResultData(types.BoolUnknown())
	default:
		resp.Result = function.NewResultData(basetypes.NewBoolValue(false))
	}
}
