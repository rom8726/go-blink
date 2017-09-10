package logs

import "github.com/ivankorobkov/di"

func Module(m *di.Module) {
	m.MarkPackageDep(Config{})
	m.AddConstructor(New)
}

func TestModule(m *di.Module) {
	m.Import(Module)
	m.AddConstructor(NewConfig)
}
