package writer

// type RuntimevarConfig defines a struct for configuration data used to build a gocloud.dev/runtimevar URI which are
// used to resolve final `whosonfirst/go-writer/v2` URIs.
type RuntimevarConfig struct {
	// The scheme of the gocloud.dev/runtimevar URI to build
	Runtimevar string `json:"runtimevar"`
	// The value of the gocloud.dev/runtimevar URI to build
	Value string `json:"value"`
	// An optional gocloud.dev/runtimevar URI used to replace the string "{credentials}" in `Value`.
	Credentials string `json:"credentials,omitempty"`
	// An optional boolean flag used to flag the config as currently disabled.
	Disabled bool `json:"disabled,omitempty"`
	// An optional string label used to match any "?exclude={LABEL}" parameters in a `whosonfirst/go-writer/v2` URI constuctor; this allows individual configs to be disabled at runtime
	Label string `json:"label,omitempty"`
}

// type TargetConfig is a map where the keys are arbitrary labels (string) mapped to a list `RuntimevarConfig` instances.
type TargetConfig map[string][]*RuntimevarConfig

// RuntimevarConfigs returns the list of `RuntimevarConfig` instances associated with the key 'target'.
func (cfg *TargetConfig) RuntimevarConfigs(target string) ([]*RuntimevarConfig, bool) {
	c, ok := (*cfg)[target]
	return c, ok
}

// type WriterConfig is a map where the keys are arbitrary labels (strings) mapped to `TargetConfig` instances.
type WriterConfig map[string]*TargetConfig

// Target returns the `TargetConfig` instance associated with the key 'target'.
func (cfg *WriterConfig) Target(environment string) (*TargetConfig, bool) {
	t, ok := (*cfg)[environment]
	return t, ok
}
