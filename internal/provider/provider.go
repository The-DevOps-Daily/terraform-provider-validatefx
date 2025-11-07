package provider

import (
	"context"
	"fmt"
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
			"datetime_layouts": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Optional list of Go time layouts applied by the `provider::validatefx::datetime` function when no layouts are provided.",
			},
			"timezone": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional default timezone (IANA identifier such as `UTC` or `America/New_York`) used when layout parsing requires an explicit location.",
			},
		},
	}
}

func (p *validateFXProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data struct {
		DatetimeLayouts types.List   `tfsdk:"datetime_layouts"`
		Timezone        types.String `tfsdk:"timezone"`
	}

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	layouts, layoutDiags := listToStrings(ctx, path.Root("datetime_layouts"), data.DatetimeLayouts)
	resp.Diagnostics.Append(layoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var location *time.Location
	if tzName, ok := optionalString(data.Timezone); ok {
		var err error
		location, err = time.LoadLocation(tzName)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("timezone"),
				"Invalid Timezone",
				fmt.Sprintf("Failed to load timezone %q: %s", tzName, err.Error()),
			)
			return
		}
	}

	functions.SetProviderConfiguration(functions.ProviderConfiguration{
		DatetimeLayouts: layouts,
		Timezone:        location,
	})
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
