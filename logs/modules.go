package logs

import "github.com/ivankorobkov/go-di"

func Module(m *di.Module) {
	m.Dep(Config{})
	m.Add(New)
}

func TestModule(m *di.Module) {
	m.Import(Module)
	m.Add(NewConfig)
}
