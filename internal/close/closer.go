package close

import "github.com/xlab/closer"

type ShutdownApp struct {
	Bind func()
}

func (s *ShutdownApp) Binder() {
	closer.Bind(s.Bind)
}

func (s *ShutdownApp) SetBind(bind func()) {
	s.Bind = bind
}
