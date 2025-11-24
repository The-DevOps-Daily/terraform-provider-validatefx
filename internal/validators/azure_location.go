package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AzureLocation validates that a string is a valid Azure location.
func AzureLocation() validator.String { return azureLocationValidator{} }

type azureLocationValidator struct{}

var _ validator.String = (*azureLocationValidator)(nil)

// Valid Azure locations as of 2024
var validAzureLocations = map[string]bool{
	// Americas
	"eastus":          true,
	"eastus2":         true,
	"southcentralus":  true,
	"westus2":         true,
	"westus3":         true,
	"centralus":       true,
	"northcentralus":  true,
	"westus":          true,
	"westcentralus":   true,
	"canadacentral":   true,
	"canadaeast":      true,
	"brazilsouth":     true,
	"brazilsoutheast": true,
	// Europe
	"northeurope":        true,
	"westeurope":         true,
	"uksouth":            true,
	"ukwest":             true,
	"francecentral":      true,
	"francesouth":        true,
	"germanywestcentral": true,
	"germanynorth":       true,
	"norwayeast":         true,
	"norwaywest":         true,
	"switzerlandnorth":   true,
	"switzerlandwest":    true,
	"swedencentral":      true,
	"polandcentral":      true,
	// Asia Pacific
	"eastasia":           true,
	"southeastasia":      true,
	"australiaeast":      true,
	"australiasoutheast": true,
	"australiacentral":   true,
	"australiacentral2":  true,
	"japaneast":          true,
	"japanwest":          true,
	"koreacentral":       true,
	"koreasouth":         true,
	"centralindia":       true,
	"southindia":         true,
	"westindia":          true,
	// Middle East
	"uaenorth":     true,
	"uaecentral":   true,
	"qatarcentral": true,
	// Africa
	"southafricanorth": true,
	"southafricawest":  true,
	// Special regions
	"chinaeast":     true,
	"chinaeast2":    true,
	"chinaeast3":    true,
	"chinanorth":    true,
	"chinanorth2":   true,
	"chinanorth3":   true,
	"usgovvirginia": true,
	"usgovarizona":  true,
	"usgovtexas":    true,
}

func (azureLocationValidator) Description(_ context.Context) string {
	return "value must be a valid Azure location"
}

func (v azureLocationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (azureLocationValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if diag := validateStringInMap(value, validAzureLocations, req.Path, "Invalid Azure Location", "Azure location"); diag != nil {
		resp.Diagnostics.Append(diag)
	}
}
