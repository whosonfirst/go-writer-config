package writer

import (
	"encoding/json"
	"os"
	"testing"
)

func TestWriterConfig(t *testing.T) {

	r, err := os.Open("fixtures/config.json")

	if err != nil {
		t.Fatalf("Failed to open example config, %v", err)
	}

	defer r.Close()

	var cfg *WriterConfig

	dec := json.NewDecoder(r)
	err = dec.Decode(&cfg)

	if err != nil {
		t.Fatalf("Failed to decode config, %v", err)
	}

	target, ok := cfg.Target("dev")

	if !ok {
		t.Fatalf("Failed to resolve dev target")
	}

	runtimevar_cfg, ok := target.RuntimevarConfigs("test")

	if !ok {
		t.Fatalf("Failed to resolve test runtimevar configs")
	}

	if len(runtimevar_cfg) != 2 {
		t.Fatalf("Unexpected config count")
	}

}
