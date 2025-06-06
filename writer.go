package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sfomuseum/runtimevar"
	"github.com/tidwall/jsonc"
	wof_writer "github.com/whosonfirst/go-writer/v3"	
)

func init() {
	ctx := context.Background()
	wof_writer.RegisterWriter(ctx, "config", NewConfigWriter)
}

// ConfigWriterOptions is a struct containing configuration options for create a new `go-writer/v2.MultiWriter` instance
// derived from a JSON configuration file
type ConfigWriterOptions struct {
	// Config is a `WriterConfig` instance containing configuration data for instantiating one or more `go-writer/v2.Writer` instances.
	Config *WriterConfig
	// Exclude is an optional list of string labels to compare against individual `RuntimevarConfig.Label` values; if there is a match that `RuntimevarConfig` instance will be excluded
	Exclude []string
	// Target is the string label mapped to the list of `RuntimevarConfig` instances used to create a new `go-writer/v2.Writer` instance.
	Target string
	// Environment is the string label mapped to the `TargetConfig` instance used to create a new `go-writer/v2.Writer` instance.
	Environment string
	// Async is an optional boolean value to signal that a new asynchronous `go-writer/v2.MultiWriter` instance should be created.
	Async   bool
	// Verbose is an optional boolean value to signal to the underlying `go-writer/v2.MultiWriter` instance that it should be verbose in logging events.
}

// NewConfigWriter return a new `go-writer/v2.Writer` instance derived from 'uri' which is expected to take the form of:
//
//	config://{ENVIRONMENT}/{TARGET}?config={VALID_GOCLOUD_DEV_RUNTIMEVAR_URI}
//
// For example:
//
//	config://dev/test?config=file:///usr/local/config.json&async=true
func NewConfigWriter(ctx context.Context, uri string) (wof_writer.Writer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	env := u.Host

	target := u.Path
	target = strings.TrimLeft(target, "/")

	q := u.Query()

	config_uri := q.Get("config")

	str_config, err := runtimevar.StringVar(ctx, config_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive config from ?config= parameter, %w", err)
	}

	var writerConfig *WriterConfig

	config_body := jsonc.ToJSON([]byte(str_config))

	err = json.Unmarshal(config_body, &writerConfig)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config provided by ?config= parameter, %w", err)
	}

	to_exclude := q["exclude"]

	async := false

	async_param := q.Get("async")

	if async_param != "" {

		b, err := strconv.ParseBool(async_param)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?async= parameter, %w", err)
		}

		async = b
	}

	opts := &ConfigWriterOptions{
		Config:      writerConfig,
		Exclude:     to_exclude,
		Environment: env,
		Target:      target,
		Async:       async,
	}

	return NewConfigWriterFromOptions(ctx, opts)
}

// NewConfigWriterFromOptions return a new `go-writer/v2.Writer` instance derived from 'opts'.
func NewConfigWriterFromOptions(ctx context.Context, opts *ConfigWriterOptions) (wof_writer.Writer, error) {

	logger := slog.Default()
	logger = logger.With("environment", opts.Environment)
	logger = logger.With("target", opts.Target)	
	
	cfg, ok := opts.Config.Target(opts.Environment)

	if !ok {
		logger.Error("Failed to derive target for enviromment.")
		return nil, fmt.Errorf("Invalid environment")
	}

	runtimevar_configs, ok := cfg.RuntimevarConfigs(opts.Target)

	if !ok {
		logger.Error("Failed to derive config(s) from target.")
		return nil, fmt.Errorf("Invalid target")
	}

	logger.Debug("Create writers.")
	
	writers, err := createWriters(ctx, runtimevar_configs, opts.Exclude)

	if err != nil {
		logger.Error("Failed to create writers.", "error", err)
		return nil, fmt.Errorf("Failed to create writers, %w", err)
	}

	logger.Debug("Create multi writer.", "count", len(writers))
	
	mw_opts := &wof_writer.MultiWriterOptions{
		Writers: writers,
		Async:   opts.Async,
	}

	mw, err := wof_writer.NewMultiWriterWithOptions(ctx, mw_opts)

	if err != nil {
		logger.Error("Failed to create multi writer.", "error", err)
		return nil, fmt.Errorf("Failed to create multi writer, %w", err)
	}

	logger.Debug("Return multi writer.")
	return mw, nil
}

func createWriters(ctx context.Context, runtimevar_configs []*RuntimevarConfig, to_exclude []string) ([]wof_writer.Writer, error) {

	writers := make([]wof_writer.Writer, 0)

	runtime_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for _, cfg := range runtimevar_configs {

		logger := slog.Default()
		logger = logger.With("label", cfg.Label)
		
		if cfg.Disabled {
			logger.Debug("Config is disabled, skipping.")
			continue
		}

		if cfg.Label != "" && len(to_exclude) > 0 {

			exclude := false

			for _, label := range to_exclude {

				if label == cfg.Label {
					logger.Debug("Config is flagged to exclude.")
					exclude = true
					break
				}
			}

			if exclude {
				logger.Debug("Config is flagged to exclude, skipping.")			
				continue
			}
		}

		var runtimevar_uri string

		rt_var := cfg.Runtimevar
		rt_value := cfg.Value
		rt_creds := cfg.Credentials

		if strings.Contains(rt_value, "{credentials}") && rt_creds != "" {

			logger.Debug("Derive credentials from runtimevar.")
			
			creds, err := runtimevar.StringVar(runtime_ctx, rt_creds)

			if err != nil {
				logger.Error("Failed to derive credentials from runtimevar.", "error", err)
				return nil, fmt.Errorf("Unable to resolve runtime var for %s, %w", rt_value, err)
			}

			creds = strings.TrimSpace(creds)

			rt_value = strings.Replace(rt_value, "{credentials}", creds, 1)
		}

		switch rt_var {
		case "file":
			logger.Debug("Derive writer URI from file:// runtimevar.")
			runtimevar_uri = fmt.Sprintf("file://%s", rt_value)
		default:
			logger.Debug("Derive writer URI from runtimevar URI.", "scheme", rt_var)		
			runtimevar_uri = fmt.Sprintf("%s://?val=%s", rt_var, url.QueryEscape(rt_value))
		}

		wr_uri, err := runtimevar.StringVar(runtime_ctx, runtimevar_uri)

		// See the way we're referencing 'cfg.Value' rather than 'rt_value' or 'runtimevar_uri'
		// in the errors below? That's so we don't accidently leak credentials that may have been
		// interpolated above.

		if err != nil {
			logger.Error("Failed to derive writer URI from runtimevar.", "error", err)
			return nil, fmt.Errorf("Failed to derive writer URI from '%s', %w", cfg.Value, err)
		}

		wr_uri = strings.TrimSpace(wr_uri)

		logger.Debug("Create new writer.")
		wr, err := wof_writer.NewWriter(ctx, wr_uri)

		if err != nil {
			logger.Error("Failed to create writer.")
			return nil, fmt.Errorf("Failed to create writer for '%s', %w", cfg.Value, err)
		}

		logger.Debug("Append writer.")
		writers = append(writers, wr)
	}

	return writers, nil
}
