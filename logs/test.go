package logs

import "github.com/ivankorobkov/di"

func TestModule(m *di.Module) {
	m.Import(Module)
	m.AddConstructor(NewConfig)
}
