package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jacobintern/GoChat/service"
	"golang.org/x/net/websocket"
)

// RegisterchatHandler is
func RegisterchatHandler(c *gin.Engine) {
	go service.Broadcaster.Start()

	c.GET("/ws", gin.WrapH(websocket.Handler(Echo)))
}

// Echo is
func Echo(conn *websocket.Conn) {
	// 建立使用者
	user := service.User{
		Conn: conn,
	}
	user.NewUser()
	// 建立傳送訊息通道 goroutine監聽
	go user.SendMessage()

	// 使用者進入
	enterMsg := user.NewUserEnterMessage()
	service.Broadcaster.UserEntering(&user)
	service.Broadcaster.Broadcast(enterMsg)

	// 訊息接收並傳送給其他使用者
	err := user.ReceiveMessage()

	// 使用者離開
	leaveMsg := user.NewUserLeaveMessage()
	service.Broadcaster.UserLeaving(&user)
	service.Broadcaster.Broadcast(leaveMsg)

	if err == nil {
		conn.Close()
	} else {
		log.Println("read from client error:", err)
		conn.Close()
	}
}
