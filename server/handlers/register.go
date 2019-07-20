package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/xuebing1110/notify-inspect/pkg/plugin"
)

var (
	regServer *plugin.RegisterServer
)

func init() {
	regServer = plugin.DefaultRegisterServer
}

func RegisterPlugin(ctx *gin.Context) {
	registerPluginWs(ctx.Writer, ctx.Request)
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func registerPluginWs(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	log.Printf("get a plugin register request from %s ...", conn.RemoteAddr().String())

	// register timeout check
	disconn_chan := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(time.Minute):
			log.Printf("read register plugin request from %s timeout!!!", conn.RemoteAddr().String())
			conn.Close()
		case disconn := <-disconn_chan:
			if disconn {
				log.Printf("get register plugin status(failed) from %s, now disconnect it", conn.RemoteAddr().String())
				conn.Close()
			}
		}
	}()

	pluginIds := make([]string, 0)
	for {
		t, data, err := conn.ReadMessage()
		if err != nil {
			for _, pid := range pluginIds {
				regServer.Deregister(pid)
			}
			break
		}

		switch t {
		case websocket.TextMessage:
			// plugin register request
			plugin_req, err := plugin.NewPlugin(data)
			if err != nil {
				log.Printf("get bad register plugin request: %s, error: %v", string(data), err)
				disconn_chan <- true
				return
			}

			// registerr plugin
			plugin_req.ResetLost()
			plugin_req.SetRegisterTime()
			err = regServer.Register(plugin_req)
			if err != nil {
				log.Printf("get bad register plugin request: %s, error: %v", string(data), err)
				plugin_req.RegisterTime = -1
				plugin_req.RegisterTimeStr = err.Error()
				conn.WriteMessage(websocket.TextMessage, plugin_req.ToJson())
				disconn_chan <- true
				return
			} else {
				disconn_chan <- false
			}

			log.Printf("the plugin %s register suc from %s !", plugin_req.Id, conn.RemoteAddr().String())
			pluginIds = append(pluginIds, plugin_req.Id)

			// reply
			conn.WriteMessage(websocket.TextMessage, plugin_req.ToJson())

		case websocket.CloseMessage:
			log.Printf("the client %s has been disconnected!", conn.RemoteAddr().String())

			for _, pid := range pluginIds {
				regServer.Deregister(pid)
			}
		}
	}
}
