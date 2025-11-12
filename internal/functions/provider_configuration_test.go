package functions

import "testing"

func TestProviderConfigurationSetGet(t *testing.T) {
	// Save current and restore after
	orig := GetProviderConfiguration()
	t.Cleanup(func() { SetProviderConfiguration(orig) })

	cfg := ProviderConfiguration{DatetimeLayouts: []string{"2006-01-02", "2006-01-02 15:04"}}
	SetProviderConfiguration(cfg)

	got := GetProviderConfiguration()
	if len(got.DatetimeLayouts) != 2 || got.DatetimeLayouts[0] != "2006-01-02" || got.DatetimeLayouts[1] != "2006-01-02 15:04" {
		t.Fatalf("unexpected provider configuration: %#v", got)
	}
}
