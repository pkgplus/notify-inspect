package handlers

import (
	"sync"
	"time"

	"github.com/kataras/iris/websocket"

	"github.com/xuebing1110/notify-inspect/pkg/plugin"
)

var (
	regServer *plugin.RegisterServer
)

func init() {
	regServer = plugin.DefaultRegisterServer
}

func RegistePlugin(c websocket.Connection) {
	logger := c.Context().Application().Logger()
	logger.Infof("get a register plugin request from %s ...", c.ID())

	// register timeout check
	disconn_chan := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(time.Minute):
			logger.Errorf("read register plugin request from %s timeout!!!", c.ID())
			c.Disconnect()
		case disconn := <-disconn_chan:
			if disconn {
				logger.Infof("get register plugin status(failed) from %s, now disconnect it", c.ID())
				c.Disconnect()
			}
		}
	}()

	// read register request
	var l sync.Mutex
	pluginIds := make([]string, 0)
	c.OnMessage(func(data []byte) {
		// heartbeat
		if string(data) == "@heart" {
			return
		}

		// plugin register request
		plugin_req, err := plugin.NewPlugin(data)
		if err != nil {
			logger.Errorf("get bad register plugin request: %s, error: %v", string(data), err)
			disconn_chan <- true
			return
		}

		// registerr plugin
		plugin_req.ResetLost()
		plugin_req.SetRegisterTime()
		regServer.Register(plugin_req)
		disconn_chan <- false

		l.Lock()
		pluginIds = append(pluginIds, plugin_req.Id)
		l.Unlock()

		// reply
		c.EmitMessage(plugin_req.ToJson())

		logger.Infof("register plugin suc: %+v", plugin_req)
	})

	// disconnect
	c.OnDisconnect(func() {
		logger.Warnf("the client %s has been disconnected!", c.ID())

		for _, pid := range pluginIds {
			regServer.Deregister(pid)
		}
	})
}
