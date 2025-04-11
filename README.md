# go-writer-config

Go package implementing the `whosonfirst/go-writer/v3` interfaces to provide methods for creating a new `whosonfirst/go-writer/v3.MultiWriter` instance derived from a JSON-encoded config file.

## Documentation (godoc)

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-writer-config.svg)](https://pkg.go.dev/github.com/whosonfirst/go-writer-config)

## Documentation

This package implements the `whosonfirst/go-writer/v3` interfaces to provide methods for creating a new `whosonfirst/go-writer/v3.MultiWriter` instance derived from a JSON-encoded config file. These config files are organized (grouped) by named "environments" and "targets". A config file may have multiple environments and an environment may have multiple targets. Each target may contain one or more "configs". For example:

```
+ environment
  + target
    - config
    - config
```

This is what the type definition for configs looks like:

```
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
```

An example config file is described below and included in [fixtures/config.json](fixtures/config.json).

### Runtimevars

Under the hood this package uses the [gocloud.dev/runtimevar](https://gocloud.dev/howto/runtimevar) package and the [sfomuseum/runtimevar.StringVar](https://github.com/sfomuseum/runtimevar) method to deference `Value` property of any given config. Please consult the [sfomuseum/runtimevar.StringVar](https://github.com/sfomuseum/runtimevar) documentation for details.

## Example

Here is how you might create a writer from a config file located at `/fixtures/config.json`. This writer will load configs (and create `writer.Writer` instances) from the `dev` environment and `test` target set:

```
import (
	"context"

	_ "github.com/whosonfirst/go-writer-config/v3"
	
	"github.com/whosonfirst/go-writer/v3"	
)

func  main(){

	ctx := context.Background()
	
	wr_uri := "config://dev/test?config=file:///fixtures/config.json&exclude=stdout"
	wr, _ := writer.NewWriter(ctx, wr_uri)

	r := bytes.NewReader([]byte("testing"))
	wr.Write(ctx, "test", r)
```

_Error handling removed for the sake of brevity._

"Config" writers are instantiated by passing a `config://` URI to the `writer.NewWriter` method. Config writer URIs take the form of:

```
config://{ENVIRONMENT}/{TARGET}?{QUERY_PARAMETERS}
```

Where:

* `{ENVIRONMENT}` is the name of the top-level environment to load configs for.
* `{TARGET}` is the name of the second-level target to load configs for.
* `{QUERY_PARAMETERS}` are one or more query parameters:

| Name | Description | Type | Required | Notes |
| --- | --- | --- | --- | --- |
| `config` | A valid `gocloud.dev/runtimevar` URI used to load the config file to parse | string | yes | As mentioned this package uses the `sfomuseum/runtimevar.StringVar` wrapper method for resolving runtimevar URIs so consult the documentation for details. |
| `exclude` | Zero or more query parameters referencing the "label" properties of individual configs to exclude when parsing a config file | string | no | |
| `async` | A boolean flag indicating whether writers (defined by the config file) should be written to asynchronously | boolean | no | |

## Config files

Config files are parsed using the [tidwall/jsonc](https://github.com/tidwall/jsonc) package in order to allow for descriptive comments.

### Example

```
{
   "dev": {
        "indexing": [
            {
	    	/* Write data using the whosonfirst/go-whosonfirst-mysql package */
		/* https://github.com/whosonfirst/go-whosonfirst-mysql/blob/main/writer/writer.go */
		
                "label": "mysqln",
                "runtimevar": "constant",
                "value": "mysql://?dsn={credentials}@tcp(localhost:3306)/your_database",
                "credentials": "file:///etc/mysql/credentials"
            },
            {
	    	/* Write data using the whosonfirst/go-whosonfirst-opensearch package */		
		/* https://github.com/whosonfirst/go-whosonfirst-opensearch/blob/main/writer/writer_opensearch2.go */
	    
             	"label": "opensearch",
		"runtimevar": "constant",
		"value": "opensearch://localhost:9200/your_index"
            },
            {
		/* Write data using the whosonfirst/go-writer null writer (basically /dev/null) */
		/* https://github.com/whosonfirst/go-writer/blob/main/null.go */
	    
             	"label": "stdout",
		"runtimevar": "constant",
		"value": "null://",
		"disabled": true
            }
	],
        "test": [
            {
    		/* Write data using the whosonfirst/go-writer STDOUT writer */
		/* https://github.com/whosonfirst/go-writer/blob/main/stdout.go */
		
             	"runtimevar": "constant",
		"value": "stdout://"
            },
            {
		/* Write data using the whosonfirst/go-writer null writer (basically /dev/null) */
		/* https://github.com/whosonfirst/go-writer/blob/main/null.go */
	    
             	"runtimevar": "constant",
		"value": "null://"
            }
	]
    }
}
```

## See also

* https://github.com/whosonfirst/go-writer
* https://github.com/sfomuseum/runtimevar
* https://gocloud.dev/howto/runtimevar