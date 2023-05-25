package configs

import "embed"

//go:embed instance.yaml test_config.yaml
var Fs embed.FS
