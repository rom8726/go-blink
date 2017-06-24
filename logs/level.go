package logs

import (
	"strings"
)

type Level int

const (
	LevelUndefined Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

var (
	levelToName map[Level]string = map[Level]string{
		LevelUndefined: "",
		LevelTrace:     "TRACE",
		LevelDebug:     "DEBUG",
		LevelInfo:      "INFO",
		LevelWarn:      "WARN",
		LevelError:     "ERROR",
	}

	nameToLevel map[string]Level = map[string]Level{
		"":      LevelUndefined,
		"TRACE": LevelTrace,
		"DEBUG": LevelDebug,
		"INFO":  LevelInfo,
		"WARN":  LevelWarn,
		"ERROR": LevelError,
	}
)

func (level Level) String() string {
	return levelToName[level]
}

// UnmarshalJSON implements the json.Marshaler interface.
func (level Level) MarshalJSON() ([]byte, error) {
	return []byte(`"` + level.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (level *Level) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.TrimSuffix(s, `"`)
	s = strings.TrimPrefix(s, `"`)
	s = strings.ToUpper(s)
	*level = nameToLevel[s]
	return nil
}
