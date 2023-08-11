package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jacobintern/GoChat/controllers"
)

func main() {
	r := gin.Default()
	// page
	r.GET("/login", controllers.GetLogin)
	r.POST("/login", controllers.Login)
	r.GET("/register", controllers.GetRegister)
	r.GET("/chatroom", controllers.GetRoom)

	// api
	// r.HandleFunc("/api/GetUserList", controllers.GetUsers).Methods(http.MethodGet)
	// r.HandleFunc("/api/GetCookies", controllers.GetUsrCookies).Methods(http.MethodGet)
	// websocket
	controllers.RegisterchatHandler()

	// static
	r.LoadHTMLGlob("views/*")
	r.Static("/static", "./static")

	// if any err log
	err := r.Run(":8888")
	if err != nil {
		fmt.Println(err)
	}
}
