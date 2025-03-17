package websocket

/**
@Author: loser
@Description: handle the request by the method
*/

type Route struct {
	Method  string
	Handler HandlerFunc
}

type HandlerFunc func(srv *Server, conn *Conn, msg *Message)
