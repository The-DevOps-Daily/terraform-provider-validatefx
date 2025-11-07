package functions

import (
	"sync"
	"time"
)

var (
	providerConfigMu sync.RWMutex
	providerConfig   ProviderConfiguration
)

// ProviderConfiguration exposes provider defaults for function implementations.
type ProviderConfiguration struct {
	DatetimeLayouts []string
	Timezone        *time.Location
}

// SetProviderConfiguration stores provider defaults for function usage.
func SetProviderConfiguration(cfg ProviderConfiguration) {
	providerConfigMu.Lock()
	defer providerConfigMu.Unlock()

	providerConfig = cfg
}

// GetProviderConfiguration returns a snapshot of the current provider defaults.
func GetProviderConfiguration() ProviderConfiguration {
	providerConfigMu.RLock()
	defer providerConfigMu.RUnlock()

	return providerConfig
}
