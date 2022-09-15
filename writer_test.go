package writer

import (
	"bytes"
	"context"
	"fmt"
	wof_writer "github.com/whosonfirst/go-writer/v2"
	"path/filepath"
	"testing"
)

func TestConfigWriter(t *testing.T) {

	rel_path := "fixtures/config.json"
	abs_path, err := filepath.Abs(rel_path)

	if err != nil {
		t.Fatalf("Failed to derive absolute path for %s, %v", rel_path, err)
	}

	ctx := context.Background()

	wr_uri := fmt.Sprintf("config://dev/test?config=file://%s&exclude=stdout", abs_path)

	wr, err := wof_writer.NewWriter(ctx, wr_uri)

	if err != nil {
		t.Fatalf("Failed to create new writer for %s, %v", wr_uri, err)
	}

	r := bytes.NewReader([]byte("testing"))

	_, err = wr.Write(ctx, "test", r)

	if err != nil {
		t.Fatalf("Failed to write data, %v", err)
	}
}
