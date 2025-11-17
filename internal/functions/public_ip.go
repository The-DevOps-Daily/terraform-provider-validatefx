package functions

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewPublicIPFunction exposes the public IP validator as a Terraform function with options.
// Signature: public_ip(value, exclude_link_local, exclude_reserved)
func NewPublicIPFunction() function.Function { return &publicIPFunction{} }

type publicIPFunction struct{}

var _ function.Function = (*publicIPFunction)(nil)

func (publicIPFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "public_ip"
}

func (publicIPFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that an IP address is public (not private).",
		MarkdownDescription: "Returns true when the input IP address is not in private ranges. Optional flags exclude link-local and reserved/documentation/multicast ranges.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:               "value",
				AllowNullValue:     true,
				AllowUnknownValues: true,
			},
			function.BoolParameter{
				Name:               "exclude_link_local",
				AllowNullValue:     true,
				AllowUnknownValues: true,
			},
			function.BoolParameter{
				Name:               "exclude_reserved",
				AllowNullValue:     true,
				AllowUnknownValues: true,
			},
		},
	}
}

func (publicIPFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// value
	value, vState, ok := stringArgument(ctx, req, resp, 0)
	if !ok {
		return
	}
	if vState == valueUnknown {
		resp.Result = function.NewResultData(basetypes.NewBoolUnknown())
		return
	}

	// flags
	exclLL, sLL, ok := ignoreCaseFlag(ctx, req, resp, 1) // reuse: returns bool and state
	if !ok {
		return
	}
	exclRes, sRes, ok := ignoreCaseFlag(ctx, req, resp, 2)
	if !ok {
		return
	}
	if unknownIf(resp, vState, sLL, sRes) {
		return
	}

	// Build a tiny wrapper validator honoring flags
	validator := publicIPWithOptions(exclLL, exclRes)
	vr := frameworkvalidator.StringResponse{}
	validator.ValidateString(ctx, frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: value}, &vr)
	if vr.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, vr.Diagnostics)
		return
	}
	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

func publicIPWithOptions(excludeLinkLocal, excludeReserved bool) frameworkvalidator.String {
	base := validators.PublicIP()
	if !excludeLinkLocal && !excludeReserved {
		return base
	}
	return &publicIPOptValidator{base: base, excludeLinkLocal: excludeLinkLocal, excludeReserved: excludeReserved}
}

type publicIPOptValidator struct {
	base             frameworkvalidator.String
	excludeLinkLocal bool
	excludeReserved  bool
}

func (v *publicIPOptValidator) Description(_ context.Context) string {
	return "public ip validation (with options)"
}

func (v *publicIPOptValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *publicIPOptValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	// Run base first (checks format and private ranges)
	v.base.ValidateString(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	parsed := net.ParseIP(req.ConfigValue.ValueString())
	if parsed == nil {
		return
	}

	if v.excludeLinkLocal && validators.IsLinkLocalIP(parsed) {
		resp.Diagnostics.AddAttributeError(req.Path, "Not a Public IP", "Value is link-local and excluded by options.")
		return
	}
	if v.excludeReserved && validators.IsReservedIP(parsed) {
		resp.Diagnostics.AddAttributeError(req.Path, "Not a Public IP", "Value is from reserved/documentation/multicast ranges and excluded by options.")
		return
	}
}
