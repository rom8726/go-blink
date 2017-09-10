package logs

import "time"

// Record is a log request.
type Record struct {
	Log     string        // Write name.
	Time    time.Time     // Creation time.
	Level   Level         // Message level.
	Message string        // Optional format string.
	Args    []interface{} // Args to print or format args.
}
