package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

/**
@Author: loser
@Description: the websocket connection
*/

type Conn struct {
	idleMu sync.Mutex

	Uid string
	*websocket.Conn
	s *Server

	idle              time.Time     // the idle time
	maxConnectionIdle time.Duration // the max idle time

	done chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err %v , ", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		done:              make(chan struct{}), // 没有缓冲管道
	}

	go conn.keepalive()
	return conn
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()

	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{} // 有消息的时候设置为繁忙
	return
}

func (c *Conn) WriteMessage(messageType int, p []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()

	err := c.Conn.WriteMessage(messageType, p)
	c.idle = time.Now()
	return err
}

func (c *Conn) Close() error {
	// 可以从已经的通道中读取到无数个值
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return c.Conn.Close()
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			// 表示连接繁忙,延长检测时间
			if idle.IsZero() {
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			// 连接不繁忙,检测剩余时间
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				//c.s.Close(c)
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}
