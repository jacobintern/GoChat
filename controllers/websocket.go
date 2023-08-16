package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jacobintern/GoChat/service"
)

// // Echo is
// func Echo(conn *websocket.Conn) {
// 	// 建立使用者
// 	user := service.User{
// 		Conn: conn,
// 	}
// 	user.NewUser()
// 	// 建立傳送訊息通道 goroutine監聽
// 	go user.SendMessage()

// 	Enter(&user)

// 	// 訊息接收並傳送給其他使用者
// 	err := user.ReceiveMessage()

// 	Leave(&user)

// 	if err == nil {
// 		conn.Close()
// 	} else {
// 		log.Println("read from client error:", err)
// 		conn.Close()
// 	}
// }

func Leave(user *service.User) {
	// 使用者離開
	leaveMsg := user.NewUserLeaveMessage()
	service.Hub.UserLeaving(user)
	service.Hub.Broadcast(leaveMsg)
}

func Enter(user *service.User) {
	// 使用者進入
	enterMsg := user.NewUserEnterMessage()
	service.Hub.UserEntering(user)
	service.Hub.Broadcast(enterMsg)
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandShake(c *gin.Context) {
	ws, _ := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	userID := c.Query("clientId")
	defer ws.Close()

	// 建立使用者
	user := service.User{
		Conn: ws,
	}
	user.NewUser(userID)
	// 建立傳送訊息通道 goroutine監聽
	go user.SendMessage()

	Enter(&user)

	// 訊息接收並傳送給其他使用者
	err := user.ReceiveMessage()

	Leave(&user)

	if err == nil {
		ws.Close()
	} else {
		log.Println("read from client error:", err)
		ws.Close()
	}

}
