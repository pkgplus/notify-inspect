package client

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

	c.conn.SetCloseHandler(func(code int, text string) error {
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