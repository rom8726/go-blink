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

type logs struct {
	mu      sync.Mutex
	logs    map[string]Log
	loggers []logger
}

func New(config Config) Logs {
	loggers := []logger{}
	for _, lconf := range config {
		logger := newLogger(lconf)
		loggers = append(loggers, logger)
	}

	return &logs{
		logs:    make(map[string]Log),
		loggers: loggers,
	}
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
