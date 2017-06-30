package logs

const (
	DefaultFormat     = "${Level}\t${Log}\t${Message}\t${Context}"
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

type Config struct {
	Formats []FormatConfig
	Writers []WriterConfig
}

type FormatConfig struct {
	Name    string            // Unique format name, "" is the default format.
	Type    FormatType        // Format type.
	Time    string            // Time layout for time.Time.Format().
	Message string            // Message layout for strs.Formatter.
	Context map[string]string // Context keys to names.
}

type WriterConfig struct {
	Type   WriterType
	Level  Level
	Format string

	// File writer
	File           string // File path.
	FileMaxSize    int    // Maximum size in megabytes of a log file.
	FileMaxAge     int    // Maximum number of days to retain old log files.
	FileMaxBackups int    // Maximum number of old log files to retain.
}

func NewConfig() Config {
	return Config{
		Writers: []WriterConfig{
			{
				Type:  WriterConsole,
				Level: LevelInfo,
			},
		},
	}
}
