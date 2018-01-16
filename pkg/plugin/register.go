package plugin

import (
	"sync"
)

type RegisterServer struct {
	*sync.Map
}

var (
	DefaultRegisterServer *RegisterServer
)

func init() {
	DefaultRegisterServer = NewRegisterServer()
}

func NewRegisterServer() *RegisterServer {
	return &RegisterServer{
		Map: new(sync.Map),
	}
}

func (r *RegisterServer) Register(p *Plugin) {
	if p.RegisterTime == 0 {
		p.SetRegisterTime()
	}
	p.ResetLost()

	r.Store(p.Id, p)
}

func (r *RegisterServer) Deregister(pluginid string) {
	p, found := r.GetPlugin(pluginid)
	if !found {
		return
	}

	p.SetLost()
	r.Store(p.Id, p)
}

func (r *RegisterServer) GetPlugin(pluginid string) (*Plugin, bool) {
	p, ok := r.Load(pluginid)
	if !ok {
		return nil, ok
	}

	plugin, ok := p.(*Plugin)
	return plugin, ok
}

func (r *RegisterServer) ListPlugins() []*Plugin {
	plugins := make([]*Plugin, 0)
	r.Range(func(_, v interface{}) bool {
		plugins = append(plugins, v.(*Plugin))
		return true
	})
	return plugins
}
