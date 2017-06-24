package logs

import "time"

type message struct {
	Log    string        // Log name.
	Time   time.Time     // Creation time.
	Level  Level         // Message level.
	Format string        // Optional format string.
	Args   []interface{} // Args to print or format args.
}
