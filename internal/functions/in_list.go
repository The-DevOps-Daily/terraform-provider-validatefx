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

type inListFunction struct{}

var _ function.Function = (*inListFunction)(nil)

// NewInListFunction exposes the in_list validator as a Terraform function.
func NewInListFunction() function.Function {
	return &inListFunction{}
}

func (inListFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "in_list"
}

func (inListFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string matches one of the allowed values.",
		MarkdownDescription: "Returns true when the input string equals one of the supplied allowed values.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "String value to validate.",
				MarkdownDescription: "String value to validate.",
			},
			function.ListParameter{
				Name:                "allowed",
				AllowNullValue:      false,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				Description:         "List of allowed string values.",
				MarkdownDescription: "List of allowed string values.",
			},
			function.BoolParameter{
				Name:                "ignore_case",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "Whether comparisons are case-insensitive.",
				MarkdownDescription: "Whether comparisons are case-insensitive.",
			},
			function.StringParameter{
				Name:                "message",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "Optional custom diagnostic message to surface on validation failure.",
				MarkdownDescription: "Optional custom diagnostic message to surface on validation failure.",
			},
		},
	}
}

func (inListFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	value, vState, ok := stringArgument(ctx, req, resp, 0)
	if !ok {
		return
	}

	allowedValues, aState, ok := allowedList(ctx, req, resp, 1)
	if !ok {
		return
	}

	ignore, iState, ok := ignoreCaseFlag(ctx, req, resp, 2)
	if !ok {
		return
	}

	// Optional custom message
	message, mState, ok := messageArgument(ctx, req, resp, 3)
	if !ok {
		return
	}

	if unknownIf(resp, vState, aState, iState, mState) {
		return
	}

	if len(allowedValues) == 0 {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	validator := selectInListValidator(allowedValues, ignore, message)

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

func selectInListValidator(allowed []string, ignore bool, message types.String) frameworkvalidator.String {
	if message.IsNull() || message.IsUnknown() || message.ValueString() == "" {
		return validators.NewInListValidator(allowed, ignore)
	}
	return validators.NewInListValidatorWithMessage(allowed, ignore, message.ValueString())
}

func stringArgument(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int) (types.String, valueState, bool) {
	var value types.String
	if err := req.Arguments.GetArgument(ctx, index, &value); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return types.String{}, valueKnown, false
	}

	if value.IsNull() || value.IsUnknown() {
		return value, valueUnknown, true
	}

	return value, valueKnown, true
}

func allowedList(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int) ([]string, valueState, bool) {
	var allowed types.List
	if err := req.Arguments.GetArgument(ctx, index, &allowed); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return nil, valueKnown, false
	}

	values, state, err := prepareAllowedValues(ctx, allowed)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return nil, valueKnown, false
	}

	return values, state, true
}

func ignoreCaseFlag(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int) (bool, valueState, bool) {
	var flag types.Bool
	if err := req.Arguments.GetArgument(ctx, index, &flag); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return false, valueKnown, false
	}

	value, state := boolFromOptional(flag)
	return value, state, true
}

type valueState int

const (
	valueKnown valueState = iota
	valueUnknown
)

func boolFromOptional(v types.Bool) (bool, valueState) {
	if v.IsUnknown() {
		return false, valueUnknown
	}

	if v.IsNull() {
		return false, valueKnown
	}

	return v.ValueBool(), valueKnown
}

func prepareAllowedValues(ctx context.Context, list types.List) ([]string, valueState, error) {
	if list.IsUnknown() {
		return nil, valueUnknown, nil
	}

	if list.IsNull() {
		return nil, valueKnown, function.NewFuncError("allowed values list must be provided")
	}

	var items []basetypes.StringValue
	diags := list.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		diags.AddAttributeError(
			path.Root("allowed"),
			"Invalid Allowed Values",
			"Allowed values must be provided as a list of strings.",
		)
		return nil, valueKnown, function.FuncErrorFromDiags(ctx, diags)
	}

	values := make([]string, 0, len(items))
	for _, item := range items {
		if item.IsNull() || item.IsUnknown() {
			continue
		}
		values = append(values, item.ValueString())
	}

	return values, valueKnown, nil
}

// messageArgument fetches an optional string argument and returns its state.
func messageArgument(ctx context.Context, req function.RunRequest, resp *function.RunResponse, index int) (types.String, valueState, bool) {
	var msg types.String
	if err := req.Arguments.GetArgument(ctx, index, &msg); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return types.String{}, valueKnown, false
	}
	if msg.IsUnknown() {
		return msg, valueUnknown, true
	}
	return msg, valueKnown, true
}

// unknownIf returns unknown result when any provided state is unknown.
func unknownIf(resp *function.RunResponse, states ...valueState) bool {
	for _, s := range states {
		if s == valueUnknown {
			resp.Result = function.NewResultData(types.BoolUnknown())
			return true
		}
	}
	return false
}
