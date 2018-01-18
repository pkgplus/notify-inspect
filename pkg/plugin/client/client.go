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
	conn  *websocket.Conn
	mutex sync.Mutex
}

var (
	DefaultRegisterClient *registerClient
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
	log.GlobalLogger.Infof("now to register the plugin...")
	if err := c.dial(); err != nil {
		return err
	}

	// send message
	if err := c.conn.WriteMessage(websocket.TextMessage, p.ToJson()); err != nil {
		return err
	}

	// read message
	_, resp_bytes, err := c.conn.ReadMessage()
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
	go func(p *plugin.Plugin) {
		for {
			msgtype, resp_bytes, err := c.conn.ReadMessage()
			if err != nil {
				log.GlobalLogger.Errorf("read message error:%v", err)
				time.Sleep(time.Minute)
				err = c.Register(p)
				if err != nil {
					log.GlobalLogger.Errorf("register failed:%v", err)
				}
				return
			}

			if msgtype == websocket.CloseMessage {
				log.GlobalLogger.Errorf("get close message from register server: %s", resp_bytes)
				time.Sleep(time.Minute)
				err = c.Register(p)
				if err != nil {
					log.GlobalLogger.Errorf("register failed:%v", err)
				}
				return
			}
		}
	}(p)

	c.conn.SetCloseHandler(func(code int, text string) error {
		log.GlobalLogger.Error("the connection with register server has been disconnected")
		time.Sleep(time.Minute)
		return c.Register(p)
	})

	return nil
}

func (c *registerClient) dial() (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn == nil {
		c.conn, _, err = websocket.DefaultDialer.Dial(c.addr, nil)
		if err != nil {
			return
		}
	}

	return
}
