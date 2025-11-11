package functions

import "sync"

var (
	cfgMu sync.RWMutex
	cfg   ProviderConfiguration
)

// ProviderConfiguration holds provider-level defaults used by function wrappers.
type ProviderConfiguration struct {
	DatetimeLayouts []string
}

// SetProviderConfiguration updates provider-level defaults.
func SetProviderConfiguration(c ProviderConfiguration) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	cfg = c
}

// GetProviderConfiguration returns a snapshot of provider-level defaults.
func GetProviderConfiguration() ProviderConfiguration {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg
}
