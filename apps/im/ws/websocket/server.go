package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"sync"
	"time"
)

/**
@Author: loser
@Description: define a websocket server
*/

// 消息确认方式
type AckType int

const (
	// 不需要 ACK 应答
	NoAck AckType = iota

	// 只需要服务器端回复 ACK 即可
	OnlyAck

	// 服务器端和客户端严格回复 ACK
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}

	return "NoAck"
}

type Server struct {
	sync.RWMutex
	routes         map[string]HandlerFunc
	addr           string
	upgrader       websocket.Upgrader
	connToUser     sync.Map
	userToConn     sync.Map
	authentication Authentication
	pattern        string
	opt            *serverOption
	*threading.TaskRunner
	logx.Logger
}

// @Description: create a websocket server
func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routes:         make(map[string]HandlerFunc),
		addr:           addr,
		upgrader:       websocket.Upgrader{},
		Logger:         logx.WithContext(context.Background()),
		authentication: opt.Authentication,
		opt:            &opt,
		pattern:        opt.pattern,
		TaskRunner:     threading.NewTaskRunner(opt.concurrency),
	}
}

// @Description: the method that server do(like epoll)
func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	s.Info("begin to server ws ...")
	defer func() {
		// 防止抛出异常导致 main 退出
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err %v", r)
		}
	}()

	/*conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("server get the websocket connection failed err %v", err)
		return
	}*/

	conn := NewConn(s, w, r)

	s.Info("begin to authentication ...")
	if !s.authentication.Auth(w, r) {
		s.Send(&Message{
			FrameType: FrameData,
			Data:      fmt.Sprint("不具备访问权限"),
		}, conn)
		s.Close(conn)
		return
	}

	s.addConn(conn, r)

	go s.handleConn(conn)
}

// @Description: dispatch the request to different method(like epoll?)
func (s *Server) handleConn(conn *Conn) {
	uid := s.GetUsers(conn)
	conn.Uid = uid[0]

	s.Info("begin to handle conn ...")

	// 统一开启协程
	go s.handlerWrite(conn)

	// 开启了 ACK 机制
	if s.opt.ack != NoAck {
		go s.readAck(conn)
	}

	for {
		// 1. get the message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message failed err %v", err)
			s.Close(conn)
			return
		}

		// 2. parse the message
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v , msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		if s.isAck(&message) {
			s.Infof("conn message read ack msg %v ", message)
			conn.appendMsg(&message)
		} else {
			conn.message <- &message
		}
	}
}

func (s *Server) isAck(message *Message) bool {
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck
}

// @Description: 读取 ACK 确认信息
func (s *Server) readAck(conn *Conn) {

	// 超时重试机制
	send := func(msg *Message, conn *Conn) error {
		err := s.Send(msg, conn)
		if err == nil {
			return nil
		}

		s.Errorf("message ack OnlyAck send err %v message", msg.Id)
		conn.readMessage[0].errCount++
		conn.messageMu.Unlock()

		tempDelay := time.Duration(200*conn.readMessage[0].errCount) * time.Millisecond
		if max := time.Second * 1; tempDelay > max {
			tempDelay = max
		}

		time.Sleep(tempDelay)
		return err
	}

	for {
		select {
		case <-conn.done:
			s.Infof("close message ack uid: %v , ", conn.Uid)
			return
		default:
		}

		conn.messageMu.Lock()
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			time.Sleep(3 * time.Second)
			continue
		}

		// 读取第一条消息
		message := conn.readMessage[0]
		switch s.opt.ack {
		case OnlyAck:
			// 一次应答
			send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			// 进行业务处理即可
			// 移除消息
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()
			conn.message <- message
		case RigorAck:
			// 首先发送 Ack
			if message.AckSeq == 0 {
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()

				send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)

				s.Infof("message ack RigroAck send mid %v , time %v ", message.Id, message.AckSeq, message.ackTime)
				conn.messageMu.Unlock()
				continue
			}

			// 验证
			// 1. 客户端已经确认
			// 2. 客户端没有确认: 是否超过 ACK 应答时间
			// 没有超过 -> 重新发送
			// 超过最大时间 -> 结束

			// 1. 客户端已经确认
			msgSeq := conn.readMessageSeq[message.Id]
			// 已经确认
			if msgSeq.AckSeq > message.AckSeq {
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack rigorAck: %v", message.Id)
				continue
			}

			val := s.opt.ackTimeout - time.Since(message.ackTime)
			// 超时
			if !message.ackTime.IsZero() && val <= 0 {
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}

			// 没有超时
			conn.messageMu.Unlock()
			send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)

			time.Sleep(3 * time.Second)
		}
	}
}

// @Description: 处理连接消息
func (s *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			return
		case message := <-conn.message:
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{FrameType: FramePing}, conn)
			case FrameData:
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行方法!")},
						conn)
				}
			}

			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

// @Description: send message by the user id
func (s *Server) SendByUserId(msg any, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	// 注意可变参数的使用方式
	return s.Send(msg, s.GetConns(sendIds...)...)
}

// @Description: send message through the connections
func (s *Server) Send(msg any, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		err = conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

// @Description: add connection obj
func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)
	// 验证用户是否登录过
	if _, ok := s.userToConn.Load(uid); ok {
		s.Close(conn)
	}

	s.connToUser.Store(conn, uid)
	s.userToConn.Store(uid, conn)
}

// @Description: get the connection by the uid
func (s *Server) GetConn(uid string) *Conn {
	conn, ok := s.userToConn.Load(uid)
	if ok {
		return conn.(*Conn)
	}
	return nil
}

// @Description: get the user id by the connection
func (s *Server) GetUser(conn *Conn) string {
	uid, ok := s.connToUser.Load(conn)
	if ok {
		return uid.(string)
	}
	return ""
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		conn, _ := s.userToConn.Load(uid)
		if conn == nil {
			res = append(res, nil)
		} else {
			res = append(res, conn.(*Conn))
		}
	}

	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {
	var res []string
	if len(conns) == 0 {
		// get all uids
		s.connToUser.Range(func(key, value any) bool {
			res = append(res, value.(string))
			return true
		})
	} else {
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			uid, _ := s.connToUser.Load(conn)
			res = append(res, uid.(string))
		}
	}

	return res
}

// @Detail: 如果这里不加上锁,那么依然可能产生并发安全问题,这是由于这里读和删除的操作不是原子性的
func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	uid, ok := s.connToUser.Load(conn)
	if !ok {
		return
	}

	s.connToUser.Delete(conn)
	s.userToConn.Delete(uid)

	conn.Close()
}

// @Description: the method that the server invoke
func (s *Server) Start() {
	s.Info("into start method")
	http.HandleFunc(s.pattern, s.ServerWs)
	s.Info("end start method", s.addr)
	//s.Info(http.ListenAndServe(s.addr, nil))
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		s.Errorf("failed to start websocket server: %v", err)
		return
	}
	s.Info("websocket server start at: ", s.addr)
}

func (s *Server) Stop() {
	s.Info("stop the websocket server")
}
