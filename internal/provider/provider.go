package provider

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/functions"
)

var (
	_ provider.Provider              = &validateFXProvider{}
	_ provider.ProviderWithFunctions = &validateFXProvider{}
)

// validateFXProvider defines the ValidateFX Terraform provider implementation.
type validateFXProvider struct {
	version string
}

// New returns a new instance of the ValidateFX provider factory function.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		functions.SetProviderVersion(version)
		return &validateFXProvider{
			version: version,
		}
	}
}

func (p *validateFXProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "validatefx"
	resp.Version = p.version
}

func (p *validateFXProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The validatefx provider exposes a suite of reusable validation functions that can be invoked from Terraform expressions using the `provider::validatefx::<name>` syntax.",
		Attributes: map[string]schema.Attribute{
			"default_datetime_layouts": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Optional default datetime layouts applied by `provider::validatefx::datetime` when call-site layouts are null/empty.",
			},
			"default_timezone": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional default timezone (IANA identifier such as `UTC` or `America/New_York`) for datetime parsing when relevant.",
			},
		},
	}
}

func (p *validateFXProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg struct {
		DefaultLayouts types.List   `tfsdk:"default_datetime_layouts"`
		DefaultTZ      types.String `tfsdk:"default_timezone"`
	}

	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	layouts, ldiags := listToStrings(ctx, path.Root("default_datetime_layouts"), cfg.DefaultLayouts)
	resp.Diagnostics.Append(ldiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// timezone is accepted but currently unused in function wrappers
	if s, ok := optionalString(cfg.DefaultTZ); ok {
		if _, err := time.LoadLocation(s); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("default_timezone"),
				"Invalid Timezone",
				"Failed to load timezone "+s+": "+err.Error(),
			)
			return
		}
	}

	functions.SetProviderConfiguration(functions.ProviderConfiguration{DatetimeLayouts: layouts})
}

func (p *validateFXProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *validateFXProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *validateFXProvider) Functions(ctx context.Context) []func() function.Function {
	return functions.ProviderFunctionFactories()
}

func listToStrings(ctx context.Context, attributePath path.Path, value types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if value.IsNull() || value.IsUnknown() {
		return nil, diags
	}

	var elements []basetypes.StringValue
	if d := value.ElementsAs(ctx, &elements, false); d.HasError() {
		d.AddAttributeError(
			attributePath,
			"Invalid Layouts",
			"Layouts must be provided as string values.",
		)
		diags.Append(d...)
		return nil, diags
	}

	result := make([]string, 0, len(elements))
	for _, elem := range elements {
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}
		trimmed := strings.TrimSpace(elem.ValueString())
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	return result, diags
}

func optionalString(value types.String) (string, bool) {
	if value.IsNull() || value.IsUnknown() {
		return "", false
	}
	trimmed := strings.TrimSpace(value.ValueString())
	if trimmed == "" {
		return "", false
	}
	return trimmed, true
}
