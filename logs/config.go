package logs

type Config []LoggerConfig

type LoggerConfig struct {
	Type    LoggerType `yaml:"type"`
	Level   Level      `yaml:"level"`
	Message string     `yaml:"format"`
	Time    string     `yaml:"time_format"`
	Context []string   `yaml:"context"`

	// File logger
	File           string `yaml:"file"`             // File path.
	FileMaxSize    int    `yaml:"file_max_size"`    // Maximum size in megabytes of a log file.
	FileMaxAge     int    `yaml:"file_max_age"`     // Maximum number of days to retain old log files.
	FileMaxBackups int    `yaml:"file_max_backups"` // Maximum number of old log files to retain.
}

func NewConfig() Config {
	return Config{
		{
			Type:  LoggerTypeConsole,
			Level: LevelInfo,
		},
	}
}
