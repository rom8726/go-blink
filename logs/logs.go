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
	writers []writer
	formats map[string]format
}

func New(config Config) Logs {
	formats := map[string]format{}
	for _, fconf := range config.Formats {
		if _, ok := formats[fconf.Name]; ok {
			panic("logs: Duplicate format \"" + fconf.Name + "\"")
		}

		formats[fconf.Name] = newFormat(fconf)
	}
	if _, ok := formats[""]; !ok {
		formats[""] = newDefaultFormat()
	}

	writers := []writer{}
	for _, wconf := range config.Writers {
		fmt := formats[wconf.Format]
		if fmt == nil {
			panic("logs: Unknown format \"" + wconf.Format + "\"")
		}

		writer := newWriter(wconf, fmt)
		writers = append(writers, writer)
	}

	return &logs{
		logs:    make(map[string]Log),
		formats: formats,
		writers: writers,
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
