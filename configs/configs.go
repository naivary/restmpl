package configs

import "embed"

// Fs includes the instace.yaml
// and will be used for test purposes
// by koanf.
//
//go:embed default.yaml
var Fs embed.FS
