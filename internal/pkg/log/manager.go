package log

import "github.com/naivary/instance/internal/pkg/log/builder"

// Should all log to the same file or seperate?
// If so: Every service should provide its own file to which it logs
// in the schema <name>_<id>.log
type Manager interface {
	Log(builder.Recorder)
	Shutdown()
}
