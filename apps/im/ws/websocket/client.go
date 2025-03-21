package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

type Client interface {
	Close() error // 相当于已经有实现了,内嵌了 *websocket.Conn
	Send(v any) error
	Read(v any) error
}

type client struct {
	*websocket.Conn
	host string
	opt  dailOption // 连接选项
}

func NewClient(host string, opts ...DailOptions) *client {
	opt := newDailOptions(opts...)

	c := &client{
		Conn: nil,
		host: host,
		opt:  opt,
	}

	conn, err := c.dail()
	if err != nil {
		panic(err)
	}

	c.Conn = conn
	return c
}

// @Description: connect to the websocket server
func (c *client) dail() (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   c.host,
		Path:   c.opt.pattern,
	}

	// 头部信息用于存储 Jwt Token
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

func (c *client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err == nil {
		return err
	}

	// TODO: 重新建立连接
	conn, err := c.dail()
	if err != nil {
		return err
	}

	c.Conn = conn
	return c.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Read(v any) error {
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg, v)
}
