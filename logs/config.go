package logs

const (
	DefaultFormat     = "${Level}\t${Log}\t${Message}\t${Context}"
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

type Config []LoggerConfig

type LoggerConfig struct {
	Type   LoggerType    `yaml:"type"`
	Level  Level         `yaml:"level"`
	Format *FormatConfig `yaml:"format"`

	// File logger
	File           string `yaml:"file"`             // File path.
	FileMaxSize    int    `yaml:"file_max_size"`    // Maximum size in megabytes of a log file.
	FileMaxAge     int    `yaml:"file_max_age"`     // Maximum number of days to retain old log files.
	FileMaxBackups int    `yaml:"file_max_backups"` // Maximum number of old log files to retain.
}

type FormatConfig struct {
	Time    string            `yaml:"time"`    // Time layout for time.Time.Format().
	Message string            `yaml:"message"` // Message layout for strs.Formatter.
	Context map[string]string `yaml:"context"` // Context keys to names.
}

func newDefaultFormatConfig() *FormatConfig {
	return &FormatConfig{
		Message: DefaultFormat,
	}
}

func NewConfig() Config {
	return Config{
		{
			Type:  LoggerTypeConsole,
			Level: LevelInfo,
		},
	}
}
