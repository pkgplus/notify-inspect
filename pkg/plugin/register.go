package plugin

import (
	"errors"
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

func (r *RegisterServer) Register(p *Plugin) error {
	if p.RegisterTime == 0 {
		p.SetRegisterTime()
	}
	p.ResetLost()

	value, exist := r.Load(p.Id)
	if exist && value.(*Plugin).LostTime == 0 {
		return errors.New("pluginRegisterRepeat")
	}
	r.Store(p.Id, p)

	return nil
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
