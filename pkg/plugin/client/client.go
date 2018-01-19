package client

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xuebing1110/notify-inspect/pkg/log"
	"github.com/xuebing1110/notify-inspect/pkg/plugin"
)

type registerClient struct {
	addr  string
	mutex sync.Mutex
}

var (
	DefaultRegisterClient *registerClient
	RECONN_TIME           = 15 * time.Second
)

func init() {
	addr := os.Getenv("WS_SERVER_URL")
	if addr == "" {
		addr = "wss://m.bingbaba.com/api/v2/notify/plugins/register"
	}
	DefaultRegisterClient = NewRegisterClient(addr)
}

func NewRegisterClient(addr string) *registerClient {
	return &registerClient{addr: addr}
}

func (c *registerClient) Register(p *plugin.Plugin) error {
	return c.register(p)
}

func (c *registerClient) register(p *plugin.Plugin) error {
	log.GlobalLogger.Infof("now to register the plugin...")
	conn, _, err := websocket.DefaultDialer.Dial(c.addr, nil)
	if err != nil {
		return err
	}

	// send message
	if err := conn.WriteMessage(websocket.TextMessage, p.ToJson()); err != nil {
		return err
	}

	// read message
	_, resp_bytes, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	// decode
	p_resp, err := plugin.NewPlugin(resp_bytes)
	if err != nil {
		return fmt.Errorf("registe plugin failed: %v", err)
	}
	if p_resp.RegisterTime <= 0 {
		return fmt.Errorf("registe plugin failed: %s", resp_bytes)
	}

	// read loop
	go func(p *plugin.Plugin, conn *websocket.Conn) {
		conn.SetCloseHandler(func(code int, text string) error {
			log.GlobalLogger.Error("the connection with register server has been disconnected")
			time.Sleep(RECONN_TIME)
			return c.register(p)
		})

		for {
			msgtype, resp_bytes, err := conn.ReadMessage()
			if err != nil {
				log.GlobalLogger.Errorf("read message error:%v", err)
				time.Sleep(RECONN_TIME)
				err = c.register(p)
				if err != nil {
					log.GlobalLogger.Errorf("register failed:%v", err)
				} else {
					return
				}
			}

			if msgtype == websocket.CloseMessage {
				log.GlobalLogger.Errorf("get close message from register server: %s", resp_bytes)
				time.Sleep(RECONN_TIME)
				err = c.register(p)
				if err != nil {
					log.GlobalLogger.Errorf("register failed:%v", err)
				} else {
					return
				}
			}
		}
	}(p, conn)

	return nil
}
