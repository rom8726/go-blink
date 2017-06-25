package logs

import (
	"github.com/ivankorobkov/di"
	"sync"
)

func Module(m *di.Module) {
	m.MarkPackageDep(Config{})
	m.AddConstructor(New)
}

type Logs interface {
	Log(name string) Log
}

func New(config Config) Logs {
	formats := map[string]format{}
	for _, fc := range config.Formats {
		if _, ok := formats[fc.Name]; ok {
			panic("logs: Duplicate format \"" + fc.Name + "\"")
		}
		formats[fc.Name] = newFormat(fc)
	}
	if _, ok := formats[""]; !ok {
		formats[""] = newDefaultFormat()
	}

	loggers := []logger{}
	for _, lc := range config.Loggers {
		f := formats[lc.Format]
		if f == nil {
			panic("logs: Undefined format \"" + lc.Format + "\"")
		}

		l := newLogger(lc, f)
		loggers = append(loggers, l)
	}

	return &logs{
		logs:    make(map[string]Log),
		formats: formats,
		loggers: loggers,
	}
}

type logs struct {
	mu      sync.Mutex
	logs    map[string]Log
	loggers []logger
	formats map[string]format
}

func (logs *logs) Log(name string) Log {
	logs.mu.Lock()
	defer logs.mu.Unlock()

	log, ok := logs.logs[name]
	if ok {
		return log
	}

	log = newLog(logs, name)
	logs.logs[name] = log
	return log
}
