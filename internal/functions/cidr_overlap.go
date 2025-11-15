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

type cidrOverlapFunction struct{}

var _ function.Function = (*cidrOverlapFunction)(nil)

// NewCIDROverlapFunction exposes the CIDR overlap validator.
func NewCIDROverlapFunction() function.Function { return &cidrOverlapFunction{} }

func (cidrOverlapFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "cidr_overlap"
}

func (cidrOverlapFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that provided CIDR blocks do not overlap.",
		MarkdownDescription: "Returns true when none of the provided CIDR blocks overlap. Fails with an error when overlap is detected.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "cidrs",
				Description:         "List of CIDR blocks to check for overlap.",
				MarkdownDescription: "List of CIDR blocks to check for overlap.",
				ElementType:         basetypes.StringType{},
				AllowNullValue:      true,
				AllowUnknownValues:  true,
			},
		},
	}
}

func (cidrOverlapFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var list types.List
	if err := req.Arguments.GetArgument(ctx, 0, &list); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	if list.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}
	if list.IsNull() {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(path.Root("cidrs"), "Missing List", "CIDR list must be provided.")
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	var elements []basetypes.StringValue
	diags := list.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		diags.AddAttributeError(path.Root("cidrs"), "Invalid Elements", "CIDR list must contain only strings.")
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	cidrs := make([]string, 0, len(elements))
	for _, el := range elements {
		if el.IsUnknown() {
			resp.Result = function.NewResultData(types.BoolUnknown())
			return
		}
		if el.IsNull() {
			diags := diag.Diagnostics{}
			diags.AddAttributeError(path.Root("cidrs"), "Null Element", "CIDR list must not contain null values.")
			resp.Error = function.FuncErrorFromDiags(ctx, diags)
			return
		}
		cidrs = append(cidrs, el.ValueString())
	}

	v := validators.NewCIDROverlap()
	if err := v.Validate(cidrs); err != nil {
		diags := diag.Diagnostics{}
		diags.AddAttributeError(path.Root("cidrs"), "CIDR Overlap", err.Error())
		resp.Error = function.FuncErrorFromDiags(ctx, diags)
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
