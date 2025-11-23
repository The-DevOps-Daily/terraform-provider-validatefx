package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type mapKeysMatchFunction struct{}

var _ function.Function = (*mapKeysMatchFunction)(nil)

// NewMapKeysMatchFunction exposes the map keys matching helper as a Terraform function.
func NewMapKeysMatchFunction() function.Function {
	return &mapKeysMatchFunction{}
}

func (mapKeysMatchFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "map_keys_match"
}

func (mapKeysMatchFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that map keys match allowed or required keys.",
		MarkdownDescription: "Returns true when map keys satisfy the allowed/required constraints. Raises an error when validation fails.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.MapParameter{
				Name:                "values",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "Source map to validate keys.",
			},
			function.ListParameter{
				Name:                "allowed_keys",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "List of allowed keys (empty means all keys allowed).",
			},
			function.ListParameter{
				Name:                "required_keys",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "List of required keys that must be present.",
			},
		},
	}
}

//nolint:cyclop
func (mapKeysMatchFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Get map parameter
	var inputMap types.Map
	if err := req.Arguments.GetArgument(ctx, 0, &inputMap); err != nil {
		resp.Error = err
		return
	}

	if inputMap.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	if inputMap.IsNull() {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("values"),
			"Missing Map",
			"Map must be provided.",
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	// Extract map keys
	mapKeys := make([]string, 0, len(inputMap.Elements()))
	for k := range inputMap.Elements() {
		mapKeys = append(mapKeys, k)
	}

	// Get allowed keys
	allowedKeys, allowedState, ok := stringListArgument(ctx, req, resp, 1, "allowed_keys")
	if !ok {
		return
	}

	if allowedState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	// Get required keys
	requiredKeys, requiredState, ok := stringListArgument(ctx, req, resp, 2, "required_keys")
	if !ok {
		return
	}

	if requiredState == valueUnknown {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	// Validate
	validator := validators.NewMapKeysMatch(allowedKeys, requiredKeys)
	if err := validator.Validate(mapKeys); err != nil {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(
			path.Root("values"),
			"Map Keys Mismatch",
			err.Error(),
		)
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
