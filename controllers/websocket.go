package controllers

import (
	"log"
	"net/http"

	"github.com/jacobintern/GoChat/service"
	"golang.org/x/net/websocket"
)

// RegisterchatHandler is
func RegisterchatHandler() {
	go service.Broadcaster.Start()

	http.Handle("/ws", websocket.Handler(Echo))
}

// Echo is
func Echo(conn *websocket.Conn) {
	// 建立使用者
	user := service.NewUser(conn)
	// 建立傳送訊息通道 goroutine監聽
	go user.SendMessage()

	// 使用者進入
	msg := service.NewUserEnterMessage(user)
	service.Broadcaster.UserEntering(user)
	service.Broadcaster.Broadcast(msg)

	// 訊息接收並傳送給其他使用者
	err := user.ReceiveMessage()

	// 使用者離開
	msg = service.NewUserLeaveMessage(user)
	service.Broadcaster.UserLeaving(user)
	service.Broadcaster.Broadcast(msg)

	if err == nil {
		conn.Close()
	} else {
		log.Println("read from client error:", err)
		conn.Close()
	}
}
