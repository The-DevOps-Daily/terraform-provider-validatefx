package functions

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type dateTimeFunction struct{}

var _ function.Function = (*dateTimeFunction)(nil)

// NewDateTimeFunction exposes the datetime validator as a Terraform function.
func NewDateTimeFunction() function.Function {
	return &dateTimeFunction{}
}

func (dateTimeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "datetime"
}

func (dateTimeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string is an ISO 8601 / RFC 3339 datetime.",
		MarkdownDescription: "Returns true when the input string matches RFC 3339 (default) or caller-provided layouts.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "value",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Description:         "Datetime string to validate.",
				MarkdownDescription: "Datetime string to validate.",
			},
			function.ListParameter{
				Name:                "layouts",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				ElementType:         basetypes.StringType{},
				Description:         "Optional list of Go time layouts to accept in addition to RFC 3339.",
				MarkdownDescription: "Optional list of Go time layouts to accept in addition to RFC 3339.",
			},
		},
	}
}

func (dateTimeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value types.String
	var layouts types.List

	if err := req.Arguments.GetArgument(ctx, 0, &value); err != nil {
		resp.Error = err
		return
	}

	if err := req.Arguments.GetArgument(ctx, 1, &layouts); err != nil {
		var funcErr *function.FuncError
		if errors.As(err, &funcErr) && funcErr.FunctionArgument != nil && *funcErr.FunctionArgument == 1 {
			layouts = basetypes.NewListNull(basetypes.StringType{})
		} else {
			resp.Error = err
			return
		}
	}

	if value.IsNull() || value.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	layoutStrings, diags := extractLayouts(ctx, layouts)
	if diags.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	validation := frameworkvalidator.StringResponse{}
	validator := validators.DateTime(layoutStrings)
	if len(layoutStrings) == 0 {
		cfg := GetProviderConfiguration()
		validator = validators.DateTimeWithLocation(cfg.DatetimeLayouts, cfg.Timezone)
	}

	validator.ValidateString(ctx, frameworkvalidator.StringRequest{
		ConfigValue: value,
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

func extractLayouts(ctx context.Context, layouts types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if layouts.IsNull() || layouts.IsUnknown() {
		return nil, diags
	}

	var elements []basetypes.StringValue

	if d := layouts.ElementsAs(ctx, &elements, false); d.HasError() {
		d.AddAttributeError(
			path.Root("layouts"),
			"Invalid Layouts",
			"Layouts must be a list of string values.",
		)
		diags.Append(d...)
		return nil, diags
	}

	result := make([]string, 0, len(elements))
	for _, elem := range elements {
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}
		result = append(result, elem.ValueString())
	}

	return result, diags
}
