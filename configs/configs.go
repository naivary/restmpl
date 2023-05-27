package configs

import "embed"

// Fs includes the instace.yaml
// and will be used for test purposes
// by koanf.
//
//go:embed instance.yaml
var Fs embed.FS
