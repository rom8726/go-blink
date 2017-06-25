package logs

type Config struct {
	Formats []FormatConfig
	Loggers []LoggerConfig
}

type FormatConfig struct {
	Name    string
	Type    FormatType
	Message string
	Time    string
	Context map[string]string // Map of context keys to param names
}

type LoggerConfig struct {
	Type   LoggerType
	Level  Level
	Format string

	// File logger
	File           string
	FileMaxSize    int // Maximum size in megabytes of a log file.
	FileMaxAge     int // Maximum number of days to retain old log files.
	FileMaxBackups int // Maximum number of old log files to retain.
}

func NewConfig() Config {
	return Config{
		Loggers: []LoggerConfig{
			{
				Type:  LoggerConsole,
				Level: LevelInfo,
			},
		},
	}
}
