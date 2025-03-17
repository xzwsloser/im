package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

/**
测试 Golang 中使用 WebSocket 连接
*/

func server(w http.ResponseWriter, r *http.Request) {
	// 1. 建立 WebSocket 连接
	var upgarder = websocket.Upgrader{}
	conn, err := upgarder.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("建立 WebSocket 连接成功 ...")
	// 2. 接受消息并且处理消息
	defer conn.Close()

	for {
		// 2.1 接受消息
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read: %v\n", err)
			break
		}

		log.Printf("server received message: %v\n", message)

		// 2.2 发送消息
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Printf("write: %v\n", err)
			break
		}
	}

}

func test() {
	http.HandleFunc("/ws", server)
	err := http.ListenAndServe("localhost:1234", nil)
	if err != nil {
		log.Fatal(err)
	}
}
